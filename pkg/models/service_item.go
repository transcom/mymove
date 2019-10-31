package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ServiceItem struct {
	ID              uuid.UUID     `db:"id"`
	CreatedAt       time.Time     `db:"created_at"`
	MoveTaskOrder   MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID uuid.UUID     `db:"move_task_order_id"`
	UpdatedAt       time.Time     `db:"updated_at"`
}

type ServiceItems []ServiceItem
