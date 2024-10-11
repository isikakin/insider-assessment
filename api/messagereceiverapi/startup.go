package messagereceiverapi

import (
	"assesment/internal/domain/repository"
	"assesment/internal/domain/service"
	"assesment/internal/messagereceiver"
	_ "assesment/internal/messagereceiver/docs"
	"assesment/internal/model"
	"assesment/pkg/cache"
	"assesment/pkg/config"
	"assesment/pkg/log"
	"assesment/pkg/middleware"
	"assesment/pkg/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Init(_ *cobra.Command, _ []string) error {

	var configuration config.Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		panic("configuration is invalid!")
	}

	app := fiber.New()

	middleware.Recover(app)
	middleware.AddCorrelationId(app)
	middleware.BasicAuth(app, configuration.ApiKey)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Service Up")
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	logger := log.NewLogger()

	sqlClient := sqlite.NewClient(true)

	messageRepository := repository.NewMessageRepository(sqlClient)
	messageService := service.NewMessageService(messageRepository)
	redisClient := cache.NewDistributed(configuration)

	controller := messagereceiver.NewController(logger, messageService, redisClient)

	app.Post("/:messageId", model.ValidateSendMessage, controller.SendMessage)
	app.Get("/", controller.RetrieveSentMessages)
	app.Put("/", controller.UpdateJobStatus)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}

	return nil
}
