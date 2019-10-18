package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// CustomerMoveItem represents a single move queue item within a queue.
type CustomerMoveItem struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	CustomerName          string     `json:"customer_name" db:"customer_name"`
	ConfirmationNumber    string     `json:"locator" db:"locator"`
	SubmittedDate         *time.Time `json:"submitted_date" db:"submitted_date"`
	LastModifiedDate      time.Time  `json:"last_modified_date" db:"last_modified_date"`
	OriginDutyStationName string     `json:"origin_duty_station_name" db:"origin_duty_station_name"`
	BranchOfService       string     `json:"branch_of_service" db:"branch_of_service"`
}

// GetCustomerMoveItems gets all CustomerMoveItems
func GetCustomerMoveItems(db *pop.Connection) ([]CustomerMoveItem, error) {
	var CustomerMoveItems []CustomerMoveItem

	err := db.RawQuery(`
	SELECT moves.ID,
		CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
		sm.affiliation AS branch_of_service,
		moves.locator AS locator,
		moves.created_at AS created_at,
				moves.updated_at AS last_modified_date,
				origin_duty_station.name AS origin_duty_station_name
			FROM moves
			JOIN orders AS ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_stations AS origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			JOIN duty_stations AS destination_duty_station ON ord.new_duty_station_id = destination_duty_station.id
		WHERE moves.show IS true
	`).All(&CustomerMoveItems)
	return CustomerMoveItems, err
}
