package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeOrder creates a single Order and associated data.
func MakeOrder(db *pop.Connection, assertions Assertions) models.Order {
	// Create new relational data if not provided
	sm := assertions.Order.ServiceMember
	if isZeroUUID(assertions.Order.ServiceMemberID) {
		sm = MakeExtendedServiceMember(db, assertions)
	}

	station := assertions.Order.NewDutyStation
	if isZeroUUID(assertions.Order.NewDutyStationID) {
		station = MakeAnyDutyStation(db)
	}

	document := assertions.Order.UploadedOrders
	if isZeroUUID(assertions.Order.UploadedOrdersID) {
		document = MakeDocument(db, Assertions{
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
				Name:            models.UploadedOrdersDocumentName,
			},
		})
	}

	order := models.Order{
		ServiceMember:    sm,
		ServiceMemberID:  sm.ID,
		NewDutyStation:   station,
		NewDutyStationID: station.ID,
		UploadedOrders:   document,
		UploadedOrdersID: document.ID,
		IssueDate:        time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:     time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:       internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		HasDependents:    true,
		SpouseHasProGear: true,
		Status:           models.OrderStatusDRAFT,
	}

	// Overwrite values with those from assertions
	mergeModels(&order, assertions.Order)

	mustCreate(db, &order)

	return order
}

// MakeDefaultOrder return an Order with default values
func MakeDefaultOrder(db *pop.Connection) models.Order {
	return MakeOrder(db, Assertions{})
}
