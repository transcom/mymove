package models

import (
	"github.com/gofrs/uuid"
)

type ServiceItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	MoveTaskOrderID uuid.UUID     `json:"move_task_order_id" db:"move_task_order_id"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_order"`
}

type ServiceItems []ServiceItem
