package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"

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
		u := MakeUserUpload(db, Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &document.ID,
				Document:   document,
				UploaderID: sm.UserID,
			},
			UserUploader: assertions.UserUploader,
		})
		document.UserUploads = append(document.UserUploads, u)
	}

	ordersNumber := "ORDER3"
	TAC := "F8E1"
	departmentIndicator := "AIR_FORCE"
	hasDependents := assertions.Order.HasDependents || false
	spouseHasProGear := assertions.Order.SpouseHasProGear || false
	grade := "E_1"

	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		assertions.Order.Grade = &grade
		entitlement = MakeEntitlement(db, assertions)
	}

	originDutyStation := assertions.OriginDutyStation
	if isZeroUUID(originDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, assertions)
	}

	orderTypeDetail := assertions.Order.OrdersTypeDetail
	hhgPermittedString := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	if orderTypeDetail == nil || *orderTypeDetail == "" {
		orderTypeDetail = &hhgPermittedString
	}

	order := models.Order{
		ServiceMember:       sm,
		ServiceMemberID:     sm.ID,
		NewDutyStation:      station,
		NewDutyStationID:    station.ID,
		UploadedOrders:      document,
		UploadedOrdersID:    document.ID,
		IssueDate:           time.Date(TestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:        time.Date(TestYear, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:          internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		OrdersNumber:        &ordersNumber,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		Status:              models.OrderStatusDRAFT,
		TAC:                 &TAC,
		DepartmentIndicator: &departmentIndicator,
		Grade:               &grade,
		Entitlement:         &entitlement,
		EntitlementID:       &entitlement.ID,
		OriginDutyStation:   &originDutyStation,
		OriginDutyStationID: &originDutyStation.ID,
		OrdersTypeDetail:    orderTypeDetail,
	}

	// Overwrite values with those from assertions
	mergeModels(&order, assertions.Order)

	mustCreate(db, &order, assertions.Stub)

	return order
}

// MakeDefaultOrder return an Order with default values
func MakeDefaultOrder(db *pop.Connection) models.Order {
	return MakeOrder(db, Assertions{})
}
