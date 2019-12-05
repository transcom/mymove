package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CustomerInfo struct {
	ID                         uuid.UUID `json:"id" db:"id"`
	CustomerName               string    `json:"customer_name" db:"customer_name"`
	OriginDutyStationName      string    `json:"origin_duty_station_name" db:"origin_duty_station_name"`
	DestinationDutyStationName string    `json:"destination_duty_station_name" db:"destination_duty_station_name"`
	Agency                     string    `json:"agency" db:"agency"`
	DependentsAuthorized       bool      `json:"dependents_authorized" db:"dependents_authorized"`
	Grade                      string    `json:"grade" db:"grade"`
	Email                      string    `json:"email" db:"email"`
	Telephone                  string    `json:"telephone" db:"telephone"`
}

// CustomerMoveItem represents a single move queue item within a queue.
type CustomerMoveItem struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	CustomerID            uuid.UUID  `json:"customer_id" db:"customer_id"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	CustomerName          string     `json:"customer_name" db:"customer_name"`
	ConfirmationNumber    string     `json:"locator" db:"locator"`
	SubmittedDate         *time.Time `json:"submitted_date" db:"submitted_date"`
	LastModifiedDate      time.Time  `json:"last_modified_date" db:"last_modified_date"`
	OriginDutyStationName string     `json:"origin_duty_station_name" db:"origin_duty_station_name"`
	BranchOfService       string     `json:"branch_of_service" db:"branch_of_service"`
	ReferenceID           *string    `json:"reference_id" db:"reference_id"`
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
		origin_duty_station.name AS origin_duty_station_name,
		sm.id AS customer_id,
		mto.reference_id
			FROM moves
			LEFT JOIN orders AS ord ON moves.orders_id = ord.id
			LEFT JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN duty_stations AS origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			LEFT JOIN move_task_orders AS mto ON mto.move_id = moves.id
		WHERE moves.show IS TRUE
	`).All(&CustomerMoveItems)
	return CustomerMoveItems, err
}

func GetCustomerInfo(db *pop.Connection, customerID uuid.UUID) (CustomerInfo, error) {
	var customer CustomerInfo
	err := db.RawQuery(`
	SELECT sm.ID,
	   CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
       sm.affiliation AS agency,
       sm.rank AS grade,
       sm.telephone AS telephone,
       sm.personal_email AS email,
       origin_duty_station.name AS origin_duty_station_name,
       destination_duty_station.name AS destination_duty_station_name,
       ord.has_dependents AS dependents_authorized
			FROM moves
			JOIN orders AS ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_stations AS origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			JOIN duty_stations AS destination_duty_station ON ord.new_duty_station_id = destination_duty_station.id
           WHERE sm.id = $1
	`, customerID).First(&customer)
	return customer, err
}
