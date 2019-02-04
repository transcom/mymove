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
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Order.ServiceMemberID) {
		sm = MakeExtendedServiceMember(db, assertions)
	}

	station := assertions.Order.NewDutyStation
	// Note above
	if isZeroUUID(assertions.Order.NewDutyStationID) {
		station = FetchOrMakeDefaultNewOrdersDutyStation(db)
	}

	document := assertions.Order.UploadedOrders
	// Note above
	if isZeroUUID(assertions.Order.UploadedOrdersID) {
		document = MakeDocument(db, Assertions{
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
		u := MakeUpload(db, Assertions{
			Upload: models.Upload{
				DocumentID: &document.ID,
				Document:   document,
				UploaderID: sm.UserID,
			},
			Uploader: assertions.Uploader,
		})
		document.Uploads = append(document.Uploads, u)
	}

	ordersNumber := "ORDER3"
	TAC := "F8E1"
	SAC := "SAC"
	departmentIndicator := "AIR_FORCE"
	hasDependents := assertions.Order.HasDependents || false
	spouseHasProGear := assertions.Order.SpouseHasProGear || false

	order := models.Order{
		ServiceMember:       sm,
		ServiceMemberID:     sm.ID,
		NewDutyStation:      station,
		NewDutyStationID:    station.ID,
		UploadedOrders:      document,
		UploadedOrdersID:    document.ID,
		IssueDate:           time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:        time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:          internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:        &ordersNumber,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		Status:              models.OrderStatusDRAFT,
		TAC:                 &TAC,
		SAC:                 &SAC,
		DepartmentIndicator: &departmentIndicator,
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
