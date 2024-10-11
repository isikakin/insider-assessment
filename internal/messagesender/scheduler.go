package messagesender

import (
	_const "assesment/internal/const"
	"assesment/internal/model"
	"assesment/pkg/circuitbreaker"
	"assesment/pkg/config"
	customHttpClient "assesment/pkg/http"
	"assesment/pkg/log"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type Scheduler interface {
	SendMessage(ctx context.Context, messageId uuid.UUID, requestBody model.SendMessageRequest) bool
}

type scheduler struct {
	configuration config.Configuration
	httpClient    customHttpClient.ClientDecorator
	logger        log.Logger
}

func (s *scheduler) SendMessage(ctx context.Context, messageId uuid.UUID, requestBody model.SendMessageRequest) bool {
	var (
		request  *http.Request
		response *http.Response
		body     []byte
		err      error
	)

	if body, err = json.Marshal(requestBody); err != nil {
		s.logger.Exception(ctx, "Error while parse body to json", err)
		return false
	}

	decorator := circuitbreaker.NewClientCircuitBreakerProxy(ctx, "message-receiver", s.logger, s.httpClient)

	if request, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://messagereceiver-api:3000/%s", messageId), bytes.NewBuffer(body)); err != nil {
		s.logger.Exception(ctx, "Error while create http request", err)
		return false
	}

	request.Header.Set(_const.AuthorizationKey, s.configuration.ApiKey)

	if response, err = decorator.Do(ctx, request); err != nil {
		s.logger.Exception(ctx, "Error while send http request", err)
		return false
	}

	defer response.Body.Close()

	return response.StatusCode == http.StatusAccepted
}

func NewScheduler(configuration config.Configuration, httpClient customHttpClient.ClientDecorator, logger log.Logger) Scheduler {
	return &scheduler{
		configuration: configuration,
		httpClient:    httpClient,
		logger:        logger,
	}
}
