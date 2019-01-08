package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// ShipmentSummaryWorksheetExtractor is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetExtractor struct {
	ServiceMemberName string `db:"service_member_full_name"`
}

// FetchShipmentSummaryWorksheetExtractor fetches a single ShipmentSummaryWorksheetExtractor for a given Shipment ID
func FetchShipmentSummaryWorksheetExtractor(db *pop.Connection, shipmentID uuid.UUID) (ShipmentSummaryWorksheetExtractor, error) {
	var ssw ShipmentSummaryWorksheetExtractor
	sql := ` SELECT
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name) AS service_member_full_name
				FROM shipments s
				INNER JOIN service_members sm
					ON s.service_member_id = sm.id
				WHERE s.id = $1
				`
	err := db.RawQuery(sql, shipmentID).Eager().First(&ssw)
	if err != nil {
		return ssw, err
	}

	return ssw, nil
}
