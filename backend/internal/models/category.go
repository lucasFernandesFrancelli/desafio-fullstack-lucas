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

// CreateCategoryInput exige os dois aprovadores já na criação — uma
// categoria sem os dois aprovadores definidos nunca deveria existir,
// já que o envio de uma solicitação depende do aprovador de ordem 1.
type CreateCategoryInput struct {
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	FirstApproverID   uuid.UUID `json:"firstApproverId"`
	SecondApproverID  uuid.UUID `json:"secondApproverId"`
}

// UpdateCategoryInput edita só nome/descrição — trocar aprovadores é uma
// operação separada (SetApproversInput), com sua própria validação.
type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type SetApproversInput struct {
	FirstApproverID  uuid.UUID `json:"firstApproverId"`
	SecondApproverID uuid.UUID `json:"secondApproverId"`
}
