package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// ApproverFor descreve uma categoria em que o usuário atua como aprovador,
// usada na tela de seletor de perfil para dar contexto ao avaliador.
type ApproverFor struct {
	CategoryID   uuid.UUID `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Order        int16     `json:"order"`
}

type UserProfile struct {
	User
	ApproverFor []ApproverFor `json:"approverFor"`
}

// CreateUserInput é usado pela área de gestão de cadastros (só gestor) para
// criar novas pessoas fictícias — ex.: para virarem aprovadoras de uma
// categoria nova.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  Role   `json:"role"`
}
