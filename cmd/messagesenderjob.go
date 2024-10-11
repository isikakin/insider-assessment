package cmd

import (
	_const "assesment/internal/const"
	"assesment/internal/domain/entity"
	"assesment/internal/domain/repository"
	"assesment/internal/domain/service"
	"assesment/internal/messagesender"
	"assesment/internal/model"
	"assesment/pkg/cache"
	"assesment/pkg/config"
	customHttp "assesment/pkg/http"
	"assesment/pkg/log"
	"assesment/pkg/sqlite"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var messageSenderJobCmd = &cobra.Command{
	Use:  "messagesender-job",
	RunE: runMessageSenderJob,
}

func init() {
	RootCmd.AddCommand(messageSenderJobCmd)
}

func runMessageSenderJob(_ *cobra.Command, _ []string) error {

	var configuration config.Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		panic("configuration is invalid!")
	}

	var logger = log.NewLogger()

	sqlClient := sqlite.NewClient(false)

	messageRepository := repository.NewMessageRepository(sqlClient)
	messageService := service.NewMessageService(messageRepository)
	redisClient := cache.NewDistributed(configuration)
	clientDecorator := customHttp.NewHttpClientDecorator(logger, 10)
	scheduler := messagesender.NewScheduler(configuration, clientDecorator, logger)

	jobStatusExist := redisClient.Exists(_const.JobStatusRedisKey)

	if !jobStatusExist {
		redisClient.Set(_const.JobStatusRedisKey, true, 1*time.Hour)
	}

	sendMessage(logger, messageService, redisClient, scheduler)

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sendMessage(logger, messageService, redisClient, scheduler)
		}
	}

	return nil
}

func sendMessage(logger log.Logger, messageService service.MessageService, redisClient cache.Cache, scheduler messagesender.Scheduler) {

	redisData, _ := redisClient.Retrieve(_const.JobStatusRedisKey)

	jobStatus, _ := strconv.ParseBool(redisData.(string))

	if !jobStatus {
		logger.Info(context.Background(), "job status is passive")
		return
	}

	fmt.Println("job started at: ", time.Now())

	var (
		messages []entity.Message
	)

	messages = messageService.RetrieveUnSentMessagesByLimit(2)

	var isSent bool

	for _, message := range messages {
		body := model.SendMessageRequest{
			To:      message.Recipient,
			Content: message.Content,
		}

		if isSent = scheduler.SendMessage(context.Background(), message.MessageId, body); !isSent {
			continue
		}

		redisClient.Set(message.MessageId.String(), time.Now().UTC(), 1*time.Hour)
	}
}
