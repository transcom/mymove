package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// MoveQueueItem represents a single move queue item within a queue.
type MoveQueueItem struct {
	ID               uuid.UUID                           `json:"id" db:"id"`
	CreatedAt        time.Time                           `json:"created_at" db:"created_at"`
	Edipi            string                              `json:"edipi" db:"edipi"`
	Rank             *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	CustomerName     string                              `json:"customer_name" db:"customer_name"`
	Locator          string                              `json:"locator" db:"locator"`
	Status           string                              `json:"status" db:"status"`
	OrdersType       string                              `json:"orders_type" db:"orders_type"`
	MoveDate         time.Time                           `json:"move_date" db:"move_date"`
	CustomerDeadline time.Time                           `json:"customer_deadline" db:"customer_deadline"`
	LastModifiedDate time.Time                           `json:"last_modified_date" db:"last_modified_date"`
	LastModifiedName string                              `json:"last_modified_name" db:"last_modified_name"`
}

// GetMoveQueueItems gets all moveQueueItems for a specific lifecycleState
func GetMoveQueueItems(db *pop.Connection, lifecycleState string) ([]MoveQueueItem, error) {
	var moveQueueItems []MoveQueueItem
	// TODO: add clause `JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id`
	var query string

	if lifecycleState == "new" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				current_date as move_date,
				moves.created_at as created_at,
				current_time as last_modified_date,
				moves.status as status
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			WHERE moves.status = 'SUBMITTED'
		`
	} else if lifecycleState == "ppm" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				current_date as move_date,
				moves.created_at as created_at,
				current_time as last_modified_date,
				moves.status as status
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			WHERE moves.status = 'APPROVED'
		`
	} else {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				current_date as move_date,
				moves.created_at as created_at,
				current_time as last_modified_date,
				moves.status as status
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
		`
	}

	err := db.RawQuery(query).All(&moveQueueItems)
	return moveQueueItems, err
}
