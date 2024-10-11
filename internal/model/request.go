package model

type SendMessageRequest struct {
	To      string `json:"to" validate:"required"`
	Content string `json:"content" validate:"required,min=3,max=20"`
}

type UpdateJobStatus struct {
	Status bool `json:"status"`
}
