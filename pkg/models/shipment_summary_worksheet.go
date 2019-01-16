package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
)

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
// Convert dates to strings in order to avoid automatic formatting within forms.go
type ShipmentSummaryWorksheetPage1Values struct {
	ServiceMemberName        string `db:"service_member_name"`
	MaxSITStorageEntitlement string
}

// ShipmentSummaryWorksheetPage2Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage2Values struct {
}

// FetchShipmentSummaryWorksheetFormValues fetches a single ShipmentSummaryWorksheetExtractor for a given Shipment ID
func FetchShipmentSummaryWorksheetFormValues(db *pop.Connection, moveID uuid.UUID) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, error) {
	var page1 ShipmentSummaryWorksheetPage1Values
	var page2 ShipmentSummaryWorksheetPage2Values

	sql := ` SELECT
				concat_ws(', ', concat_ws(' ', sm.last_name, sm.suffix), concat_ws(' ', sm.first_name, sm.middle_name)) AS service_member_name
				FROM moves m
				INNER JOIN orders o ON m.orders_id = o.id
				INNER JOIN service_members sm ON o.service_member_id = sm.id
				WHERE m.id = $1
				`
	err := db.RawQuery(sql, moveID).First(&page1)
	if err != nil {
		return page1, page2, err
	}
	page1.MaxSITStorageEntitlement = "90 days per each shipment"

	return page1, page2, nil
}
