package model

import (
	"assesment/internal/domain/entity"
	"github.com/google/uuid"
)

type GetMessagesResponse struct {
	Messages       []entity.Message `json:"messages"`
	CurrentPage    int              `json:"currentPage"`
	TotalPageCount int              `json:"totalPageCount"`
}

type SendMessageResponse struct {
	Message   string    `json:"message"`
	MessageId uuid.UUID `json:"messageId"`
}
