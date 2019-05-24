package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

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
	GBLNumber        *string                             `json:"gbl_number" db:"gbl_number"`
	Status           string                              `json:"status" db:"status"`
	PpmStatus        *string                             `json:"ppm_status" db:"ppm_status"`
	HhgStatus        *string                             `json:"hhg_status" db:"hhg_status"`
	OrdersType       string                              `json:"orders_type" db:"orders_type"`
	MoveDate         *time.Time                          `json:"move_date" db:"move_date"`
	SubmittedDate    *time.Time                          `json:"submitted_date" db:"submitted_date"`
	LastModifiedDate time.Time                           `json:"last_modified_date" db:"last_modified_date"`
	LastModifiedName string                              `json:"last_modified_name" db:"last_modified_name"`
}

// GetMoveQueueItems gets all moveQueueItems for a specific lifecycleState
func GetMoveQueueItems(db *pop.Connection, lifecycleState string) ([]MoveQueueItem, error) {
	var moveQueueItems []MoveQueueItem
	var query string

	if lifecycleState == "new" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				COALESCE(
					shipment.actual_pickup_date,
					shipment.pm_survey_planned_pickup_date,
					shipment.requested_pickup_date,
					ppm.actual_move_date,
					ppm.original_move_date
				) as move_date,
				COALESCE(
					shipment.submit_date,
					ppm.submit_date
				) as submitted_date,
				moves.created_at as created_at,
				moves.status as status,
				ppm.status as ppm_status,
				shipment.status as hhg_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			WHERE moves.status = 'SUBMITTED'
			and moves.show is true
		`
	} else if lifecycleState == "ppm" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				COALESCE(ppm.actual_move_date, ppm.original_move_date) as move_date,
				moves.created_at as created_at,
				ppm.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			WHERE moves.status = 'APPROVED'
			and moves.show is true
		`
	} else if lifecycleState == "hhg_accepted" {
		// Move date is the Requested Pickup Date because accepted shipments haven't yet gone through the
		// premove survey to set the actual Pickup Date.
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				shipment.requested_pickup_date as move_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				shipment.status as hhg_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments as shipment ON moves.id = shipment.move_id
			WHERE shipment.status = 'ACCEPTED'
			and moves.show is true
		`
	} else if lifecycleState == "hhg_in_transit" {
		// Move date is the Actual Pickup Date.
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				shipment.actual_pickup_date as move_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				shipment.status as hhg_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments as shipment ON moves.id = shipment.move_id
			WHERE shipment.status = 'IN_TRANSIT'
			and moves.show is true
		`
	} else if lifecycleState == "hhg_delivered" {
		// Move date is the Actual Pickup Date.
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				shipment.actual_pickup_date as move_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				shipment.status as hhg_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments as shipment ON moves.id = shipment.move_id
			WHERE shipment.status = 'DELIVERED'
			and moves.show is true
		`
	} else if lifecycleState == "all" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				ord.orders_type as orders_type,
				COALESCE(
					shipment.actual_pickup_date,
					shipment.pm_survey_planned_pickup_date,
					shipment.requested_pickup_date,
					ppm.actual_move_date,
					ppm.original_move_date
				) as move_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				shipment.gbl_number as gbl_number
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			WHERE moves.show is true
		`
	}

	err := db.RawQuery(query).All(&moveQueueItems)
	return moveQueueItems, err
}
