package messagereceiver

import (
	_const "assesment/internal/const"
	"assesment/internal/domain/service"
	_ "assesment/internal/messagereceiver/docs"
	"assesment/internal/model"
	"assesment/pkg/cache"
	"assesment/pkg/log"
	"assesment/pkg/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// Controller
// @title Message Receiver API
// @version 1.0
// @description This is a message receiver api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-ins-auth-key
// @Security ApiKeyAuth
type Controller struct {
	logger         log.Logger
	messageService service.MessageService
	redisClient    cache.Cache
}

// SendMessage godoc
// @Summary     Send Message
// @Description Used to send messages
// @Tags         Message
// @Accept       json
// @Produce      json
// @Param messageId path string true "uuid formatted id"
// @Param 		 requestBody body model.SendMessageRequest true "request body"
// @Router     /  [post]
// @Success 202 {object} model.SendMessageResponse
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-ins-auth-key
// @Security ApiKeyAuth
func (self *Controller) SendMessage(c *fiber.Ctx) error {

	var (
		requestBody model.SendMessageRequest
		response    model.SendMessageResponse
		err         error
	)

	ctx := middleware.BindFiberContext(c, self.logger)

	messageId, err := uuid.Parse(c.Params("messageId"))

	if err = c.BodyParser(&requestBody); err != nil {
		return err
	}

	self.messageService.MarkAsSent(messageId)

	self.logger.Info(ctx, fmt.Sprintf("%s message mark as sent successfully", messageId.String()))

	response = model.SendMessageResponse{
		Message:   http.StatusText(http.StatusAccepted),
		MessageId: messageId,
	}

	return c.Status(http.StatusAccepted).JSON(response)
}

// RetrieveSentMessages godoc
// @Summary     Retrieve sent messages
// @Description Gets the list of sent messages
// @Tags         Message
// @Accept       json
// @Produce      json
// @Param page query int true "value for pagination"
// @Router     /  [get]
// @Success 200 {object} model.GetMessagesResponse
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-ins-auth-key
// @Security ApiKeyAuth
func (self *Controller) RetrieveSentMessages(c *fiber.Ctx) error {

	var (
		response model.GetMessagesResponse
	)

	page := c.QueryInt("page")

	response = self.messageService.RetrieveSentMessages(page)
	return c.Status(http.StatusOK).JSON(response)
}

// UpdateJobStatus godoc
// @Summary     Update job status
// @Description Update auto message sender job status
// @Tags         Message
// @Accept       json
// @Produce      json
// @Param 		 requestBody body model.UpdateJobStatus true "request body"
// @Router     /  [put]
// @Success 204
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-ins-auth-key
// @Security ApiKeyAuth
func (self *Controller) UpdateJobStatus(c *fiber.Ctx) error {
	var (
		requestBody model.UpdateJobStatus
		err         error
	)

	ctx := middleware.BindFiberContext(c, self.logger)

	if err = c.BodyParser(&requestBody); err != nil {
		return err
	}

	self.redisClient.Set(_const.JobStatusRedisKey, requestBody.Status, 1*time.Hour)

	self.logger.Info(ctx, fmt.Sprintf("message sender job status updated as %v successfully", requestBody.Status))

	return c.Status(http.StatusNoContent).Send(nil)
}

func NewController(logger log.Logger, messageService service.MessageService, redisClient cache.Cache) *Controller {
	return &Controller{
		logger:         logger,
		messageService: messageService,
		redisClient:    redisClient,
	}
}
