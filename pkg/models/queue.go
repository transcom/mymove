package models

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

// MoveQueueItem represents a single move queue item within a queue.
type MoveQueueItem struct {
	ID                          uuid.UUID          `json:"id" db:"id"`
	CreatedAt                   time.Time          `json:"created_at" db:"created_at"`
	Edipi                       string             `json:"edipi" db:"edipi"`
	Rank                        *ServiceMemberRank `json:"rank" db:"rank"`
	CustomerName                string             `json:"customer_name" db:"customer_name"`
	Locator                     string             `json:"locator" db:"locator"`
	Status                      string             `json:"status" db:"status"`
	PpmStatus                   *string            `json:"ppm_status" db:"ppm_status"`
	OrdersType                  string             `json:"orders_type" db:"orders_type"`
	MoveDate                    *time.Time         `json:"move_date" db:"move_date"`
	SubmittedDate               *time.Time         `json:"submitted_date" db:"submitted_date"`
	LastModifiedDate            time.Time          `json:"last_modified_date" db:"last_modified_date"`
	OriginDutyLocationName      string             `json:"origin_duty_location_name" db:"origin_duty_location_name"`
	DestinationDutyLocationName string             `json:"destination_duty_location_name" db:"destination_duty_location_name"`
	PmSurveyConductedDate       *time.Time         `json:"pm_survey_conducted_date" db:"pm_survey_conducted_date"`
	OriginGBLOC                 *string            `json:"origin_gbloc" db:"origin_gbloc"`
	DestinationGBLOC            *string            `json:"destination_gbloc" db:"destination_gbloc"`
	DeliveredDate               *time.Time         `json:"delivered_date" db:"delivered_date"`
	InvoiceApprovedDate         *time.Time         `json:"invoice_approved_date" db:"invoice_approved_date"`
	BranchOfService             string             `json:"branch_of_service" db:"branch_of_service"`
	ActualMoveDate              *time.Time         `json:"actual_move_date" db:"actual_move_date"`
	OriginalMoveDate            *time.Time         `json:"original_move_date" db:"original_move_date"`
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
				sm.affiliation as branch_of_service,
				ord.orders_type as orders_type,
				COALESCE(
					ppm.actual_move_date,
					ppm.original_move_date
				) as move_date,
				COALESCE(
					ppm.submit_date
				) as submitted_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				origin_duty_location.name as origin_duty_location_name,
                ppm.actual_move_date,
                ppm.original_move_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN duty_locations as origin_duty_location ON sm.duty_station_id = origin_duty_location.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			WHERE (moves.status = 'SUBMITTED'
			OR (ppm.status = 'SUBMITTED'
				AND (NOT moves.status in ('CANCELED', 'DRAFT'))))
			AND moves.show is true
			GROUP BY moves.ID, rank, customer_name, edipi, locator, orders_type, move_date, moves.created_at, last_modified_date, moves.status, ppm.submit_date, ppm_status, origin_duty_location.name, sm.affiliation, ppm.actual_move_date, ppm.original_move_date
		`
	} else if lifecycleState == "ppm_payment_requested" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				sm.affiliation as branch_of_service,
				ord.orders_type as orders_type,
				COALESCE(ppm.actual_move_date, ppm.original_move_date) as move_date,
				moves.created_at as created_at,
				ppm.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				origin_duty_location.name as origin_duty_location_name,
				destination_duty_location.name as destination_duty_location_name,
                ppm.actual_move_date,
                ppm.original_move_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_locations as origin_duty_location ON sm.duty_station_id = origin_duty_location.id
			JOIN duty_locations as destination_duty_location ON ord.new_duty_location_id = destination_duty_location.id
			WHERE moves.show is true
			and ppm.status = 'PAYMENT_REQUESTED'
		`
	} else if lifecycleState == "ppm_completed" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				sm.affiliation as branch_of_service,
				ord.orders_type as orders_type,
				COALESCE(ppm.actual_move_date, ppm.original_move_date) as move_date,
				moves.created_at as created_at,
				ppm.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				origin_duty_location.name as origin_duty_location_name,
				destination_duty_location.name as destination_duty_location_name,
                ppm.actual_move_date,
                ppm.original_move_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_locations as origin_duty_location ON sm.duty_station_id = origin_duty_location.id
			JOIN duty_locations as destination_duty_location ON ord.new_duty_location_id = destination_duty_location.id
			WHERE moves.show is true
			and ppm.status = 'COMPLETED'
		`
	} else if lifecycleState == "ppm_approved" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				sm.affiliation as branch_of_service,
				ord.orders_type as orders_type,
				COALESCE(ppm.actual_move_date, ppm.original_move_date) as move_date,
				moves.created_at as created_at,
				ppm.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				origin_duty_location.name as origin_duty_location_name,
				destination_duty_location.name as destination_duty_location_name,
                ppm.actual_move_date,
                ppm.original_move_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_locations as origin_duty_location ON sm.duty_station_id = origin_duty_location.id
			JOIN duty_locations as destination_duty_location ON ord.new_duty_location_id = destination_duty_location.id
			WHERE moves.show is true
			and ppm.status = 'APPROVED'
		`
	} else if lifecycleState == "all" {
		query = `
			SELECT moves.ID,
				COALESCE(sm.edipi, '*missing*') as edipi,
				COALESCE(sm.rank, '*missing*') as rank,
				CONCAT(COALESCE(sm.last_name, '*missing*'), ', ', COALESCE(sm.first_name, '*missing*')) AS customer_name,
				moves.locator as locator,
				sm.affiliation as branch_of_service,
				ord.orders_type as orders_type,
				COALESCE(
					ppm.actual_move_date,
					ppm.original_move_date
				) as move_date,
				moves.created_at as created_at,
				moves.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				origin_duty_location.name as origin_duty_location_name,
				destination_duty_location.name as destination_duty_location_name,
                ppm.actual_move_date,
                ppm.original_move_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_locations as origin_duty_location ON sm.duty_station_id = origin_duty_location.id
			JOIN duty_locations as destination_duty_location ON ord.new_duty_location_id = destination_duty_location.id
			WHERE moves.show is true
		`
	} else {
		return moveQueueItems, ErrFetchNotFound
	}

	err := db.RawQuery(query).All(&moveQueueItems)
	return moveQueueItems, err
}
