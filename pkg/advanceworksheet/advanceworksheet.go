package advanceworksheet

import (
	"github.com/gobuffalo/pop"
	"github.com/satori/go.uuid"
	"github.com/transcom/mymove/pkg/models"
)

// FetchAdvanceInfo returns the the populated struct with the info for a PPM Advance Worksheet
func FetchAdvanceInfo(dbConnection *pop.Connection, moveID uuid.UUID) (models.AdvanceWorksheet, error) {
	var advanceInfo models.AdvanceWorksheet

	sql := `SELECT
    sm.first_name,
    sm.middle_name,
    sm.last_name,
    sm.telephone,
    sm.edipi,
    sm.affiliation,
    sm.rank,
    sm.personal_email,
    ord.issue_date,
    ord.orders_type,
    COALESCE(ord.orders_number, 'Not entered') AS orders_number,
    COALESCE(ord.department_indicator, 'Not entered') AS department_indicator,
    duty_stations.name,
    ppm.pickup_postal_code,
    ppm.destination_postal_code,
    ppm.planned_move_date,
    ppm.weight_estimate,
    ppm.status,
    COALESCE(ppm.days_in_storage, 0) as days_in_storage,
    reimbursements.id AS reimbursement_id,
    reimbursements.method_of_receipt AS reimbursement_payment_method,
    reimbursements.requested_amount AS requested_reimbursement_amount,
    bu.name AS backup_contact_name,
    bu.created_at AS backup_contact_authorization_date,
    bu.email AS backup_contact_email,
    COALESCE(bu.phone, 'Not provided') AS backup_contact_phone
  FROM service_members AS sm
  LEFT JOIN orders AS ord ON
    ord.service_member_id = sm.id
  LEFT JOIN duty_stations ON
    ord.new_duty_station_id = duty_stations.id
  LEFT JOIN backup_contacts as bu ON
    sm.id = bu.service_member_id
  LEFT JOIN moves ON
    ord.id = moves.orders_id
  LEFT JOIN personally_procured_moves AS ppm ON
    moves.id = ppm.move_id
  LEFT JOIN reimbursements ON
    ppm.advance_id = reimbursements.id;`

	err := dbConnection.RawQuery(sql).All(&advanceInfo)
	return advanceInfo, err
}
