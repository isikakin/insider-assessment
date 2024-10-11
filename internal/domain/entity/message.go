package entity

import (
	"assesment/internal/domain/enum"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	MessageId uuid.UUID
	Recipient string
	Content   string
	Status    enum.MessageStatus
	SentDate  *time.Time
	CreatedAt time.Time
}

func NewMessage(messageId uuid.UUID, recipient, content string) *Message {
	return &Message{
		MessageId: messageId,
		Recipient: recipient,
		Content:   content,
		Status:    enum.Pending,
		SentDate:  nil,
		CreatedAt: time.Now(),
	}
}
