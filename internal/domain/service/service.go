package service

import (
	"assesment/internal/domain/entity"
	"assesment/internal/domain/repository"
	model2 "assesment/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"math"
	"net/http"
)

type MessageService interface {
	AddMessage(messageId uuid.UUID, request model2.SendMessageRequest) fiber.Map
	RetrieveSentMessages(page int) model2.GetMessagesResponse
	RetrieveUnSentMessagesByLimit(limit int) []entity.Message
	MarkAsSent(messageId uuid.UUID)
}

type messageService struct {
	ctx               *fiber.Ctx
	messageRepository repository.MessageRepository
}

func (s *messageService) RetrieveSentMessages(page int) (response model2.GetMessagesResponse) {

	var (
		messages   []entity.Message
		totalCount int
	)

	if page <= 0 {
		page = 1
	}

	messages = s.messageRepository.RetrieveSentMessages(page)

	totalCount = s.messageRepository.RetrieveSentMessagesTotalCount()

	totalPageCount := int(math.Ceil(float64(totalCount) / float64(10)))

	response = model2.GetMessagesResponse{
		Messages:       messages,
		CurrentPage:    page,
		TotalPageCount: totalPageCount}

	return response
}

func (s *messageService) AddMessage(messageId uuid.UUID, request model2.SendMessageRequest) (response fiber.Map) {

	message := entity.NewMessage(messageId,
		request.To,
		request.Content,
	)

	s.messageRepository.Insert(message)

	response = fiber.Map{
		"message":   http.StatusText(http.StatusAccepted),
		"messageId": messageId,
	}

	return response
}

func (s *messageService) RetrieveUnSentMessagesByLimit(limit int) (messages []entity.Message) {
	return s.messageRepository.RetrieveUnSentMessagesByLimit(limit)
}

func (s *messageService) MarkAsSent(messageId uuid.UUID) {
	s.messageRepository.MarkAsSent(messageId)
}

func NewMessageService(messageRepository repository.MessageRepository) MessageService {
	return &messageService{messageRepository: messageRepository}
}
