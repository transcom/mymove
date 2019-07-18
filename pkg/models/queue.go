package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// MoveQueueItem represents a single move queue item within a queue.
type MoveQueueItem struct {
	ID                         uuid.UUID                           `json:"id" db:"id"`
	CreatedAt                  time.Time                           `json:"created_at" db:"created_at"`
	Edipi                      string                              `json:"edipi" db:"edipi"`
	Rank                       *internalmessages.ServiceMemberRank `json:"rank" db:"rank"`
	CustomerName               string                              `json:"customer_name" db:"customer_name"`
	Locator                    string                              `json:"locator" db:"locator"`
	GBLNumber                  *string                             `json:"gbl_number" db:"gbl_number"`
	Status                     string                              `json:"status" db:"status"`
	PpmStatus                  *string                             `json:"ppm_status" db:"ppm_status"`
	HhgStatus                  *string                             `json:"hhg_status" db:"hhg_status"`
	OrdersType                 string                              `json:"orders_type" db:"orders_type"`
	MoveDate                   *time.Time                          `json:"move_date" db:"move_date"`
	SubmittedDate              *time.Time                          `json:"submitted_date" db:"submitted_date"`
	LastModifiedDate           time.Time                           `json:"last_modified_date" db:"last_modified_date"`
	ShipmentID                 uuid.UUID                           `json:"shipment_id" db:"shipment_id"`
	OriginDutyStationName      string                              `json:"origin_duty_station_name" db:"origin_duty_station_name"`
	DestinationDutyStationName string                              `json:"destination_duty_station_name" db:"destination_duty_station_name"`
	SitArray                   string                              `json:"sit_array" db:"sit_array"`
	SliArray                   string                              `json:"sli_array" db:"sli_array"`
	PmSurveyConductedDate      *time.Time                          `json:"pm_survey_conducted_date" db:"pm_survey_conducted_date"`
	OriginGBLOC                *string                             `json:"origin_gbloc" db:"origin_gbloc"`
	DestinationGBLOC           *string                             `json:"destination_gbloc" db:"destination_gbloc"`
	DeliveredDate              *time.Time                          `json:"delivered_date" db:"delivered_date"`
	InvoiceApprovedDate        *time.Time                          `json:"invoice_approved_date" db:"invoice_approved_date"`
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
				moves.updated_at as last_modified_date,
				moves.status as status,
				ppm.status as ppm_status,
				shipment.status as hhg_status,
				shipment.gbl_number as gbl_number,
				shipment.pm_survey_conducted_date as pm_survey_conducted_date,
				json_agg(json_build_object('id', sits.id , 'location', sits.location, 'status', sits.status, 'actual_start_date', sits.actual_start_date, 'out_date', sits.out_date)) as sit_array,
				json_agg(slis.status) as sli_array,
				origin_duty_station.name as origin_duty_station_name
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN duty_stations as origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			LEFT JOIN storage_in_transits as sits ON sits.shipment_id = shipment.id
			LEFT JOIN shipment_line_items as slis ON slis.shipment_id = shipment.id
			WHERE (moves.status = 'SUBMITTED'
			OR ((shipment.status in ('SUBMITTED', 'AWARDED', 'ACCEPTED') OR ppm.status = 'SUBMITTED')
				AND (NOT moves.status in ('CANCELED', 'DRAFT'))))
			AND moves.show is true
			GROUP BY moves.ID, rank, customer_name, edipi, locator, orders_type, move_date, moves.created_at, last_modified_date, moves.status, shipment.id, ppm.submit_date, ppm_status, origin_duty_station.name
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
				shipment.gbl_number as gbl_number,
				origin_duty_station.name as origin_duty_station_name,
				destination_duty_station.name as destination_duty_station_name
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_stations as origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			JOIN duty_stations as destination_duty_station ON ord.new_duty_station_id = destination_duty_station.id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			WHERE moves.show is true
			and ppm.status in ('APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED')
		`
	} else if lifecycleState == "hhg_active" {
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
				shipment.gbl_number as gbl_number,
				origin_duty_station.name as origin_duty_station_name,
				destination_duty_station.name as destination_duty_station_name,
				shipment.id as shipment_id,
				shipment.pm_survey_conducted_date as pm_survey_conducted_date,
				json_agg(json_build_object('id', sits.id , 'location', sits.location, 'status', sits.status, 'actual_start_date', sits.actual_start_date, 'out_date', sits.out_date)) as sit_array,
				json_agg(slis.status) as sli_array
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			JOIN duty_stations as origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			JOIN duty_stations as destination_duty_station ON ord.new_duty_station_id = destination_duty_station.id
			LEFT JOIN shipments as shipment ON moves.id = shipment.move_id
			LEFT JOIN storage_in_transits as sits ON sits.shipment_id = shipment.id
			LEFT JOIN shipment_line_items as slis ON slis.shipment_id = shipment.id
			WHERE ((shipment.status IN ('IN_TRANSIT', 'APPROVED')) OR (shipment.status = 'ACCEPTED' AND shipment.pm_survey_conducted_date IS NOT NULL))
			AND moves.show is true AND moves.status != 'CANCELED'
			GROUP BY moves.ID, rank, customer_name, edipi, locator, orders_type, move_date, moves.created_at, last_modified_date, moves.status, origin_duty_station_name, destination_duty_station_name, shipment.id
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
				shipment.gbl_number as gbl_number,
				shipment.pm_survey_conducted_date as pm_survey_conducted_date,
				json_agg(json_build_object('id', sits.id , 'location', sits.location, 'status', sits.status, 'actual_start_date', sits.actual_start_date, 'out_date', sits.out_date)) as sit_array,
				json_agg(slis.status) as sli_array,
				shipment.source_gbloc as origin_gbloc,
				shipment.destination_gbloc as destination_gbloc,
				shipment.actual_delivery_date as delivered_date,
				(case when invoice.status = 'SUBMITTED' then invoice.invoiced_date end) as invoice_approved_date
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments as shipment ON moves.id = shipment.move_id
			LEFT JOIN storage_in_transits as sits ON sits.shipment_id = shipment.id
			LEFT JOIN shipment_line_items as slis ON slis.shipment_id = shipment.id
			LEFT JOIN invoices as invoice on invoice.shipment_id = shipment.id
			WHERE shipment.status = 'DELIVERED'
			and moves.show is true
			GROUP BY moves.ID, rank, customer_name, edipi, locator, orders_type, move_date, moves.created_at, last_modified_date, moves.status,
			shipment.id, origin_gbloc, destination_gbloc, invoice_approved_date, delivered_date

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
				shipment.gbl_number as gbl_number,
				origin_duty_station.name as origin_duty_station_name,
				destination_duty_station.name as destination_duty_station_name
			FROM moves
			JOIN orders as ord ON moves.orders_id = ord.id
			JOIN service_members AS sm ON ord.service_member_id = sm.id
			LEFT JOIN shipments AS shipment ON moves.id = shipment.move_id
			LEFT JOIN personally_procured_moves AS ppm ON moves.id = ppm.move_id
			JOIN duty_stations as origin_duty_station ON sm.duty_station_id = origin_duty_station.id
			JOIN duty_stations as destination_duty_station ON ord.new_duty_station_id = destination_duty_station.id
			WHERE moves.show is true
		`
	} else {
		return moveQueueItems, ErrFetchNotFound
	}

	err := db.RawQuery(query).All(&moveQueueItems)
	return moveQueueItems, err
}
