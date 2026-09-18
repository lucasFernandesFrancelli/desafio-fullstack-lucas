package models

import "github.com/google/uuid"

// UserPendingCount é usado no dashboard do gestor para responder "quem tem
// mais solicitações esperando a própria ação agora".
type UserPendingCount struct {
	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	Count    int       `json:"count"`
}

// OldestPendingItem resume uma solicitação parada há mais tempo, para o
// gestor enxergar gargalos sem precisar abrir cada item.
type OldestPendingItem struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	Status            Status    `json:"status"`
	CategoryName      string    `json:"categoryName"`
	PendingActorName  string    `json:"pendingActorName"`
	DaysSinceLastMove int       `json:"daysSinceLastMove"`
}

type DashboardSummary struct {
	CountsByStatus map[Status]int      `json:"countsByStatus"`
	OldestPending  []OldestPendingItem `json:"oldestPending"`
	AwaitingByUser []UserPendingCount  `json:"awaitingByUser"`
}
