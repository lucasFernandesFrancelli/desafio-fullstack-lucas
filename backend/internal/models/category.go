package models

import "github.com/google/uuid"

type CategoryApprover struct {
	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	Order    int16     `json:"order"`
}

type Category struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Approvers   []CategoryApprover `json:"approvers"`
}
