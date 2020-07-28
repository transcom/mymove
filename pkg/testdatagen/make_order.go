package testdatagen

import (
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeGrade makes a service member grade
func MakeGrade() string {
	grades := [28]string{"E_1",
		"E_2",
		"E_3",
		"E_4",
		"E_5",
		"E_6",
		"E_7",
		"E_8",
		"E_9",
		"O_1_ACADEMY_GRADUATE",
		"O_2",
		"O_3",
		"O_4",
		"O_5",
		"O_6",
		"O_7",
		"O_8",
		"O_9",
		"O_10",
		"W_1",
		"W_2",
		"W_3",
		"W_4",
		"W_5",
		"AVIATION_CADET",
		"CIVILIAN_EMPLOYEE",
		"ACADEMY_CADET",
		"MIDSHIPMAN"}
	return grades[rand.Intn(len(grades))]
}

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

	grade := assertions.Order.Grade
	if grade == nil || *grade == "" {
		grade = stringPointer(MakeGrade())
	}

	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		assertions.Order.Grade = grade
		entitlement = MakeEntitlement(db, assertions)
	}

	originDutyStation := assertions.OriginDutyStation
	if isZeroUUID(originDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, assertions)
	}

	orderTypeDetail := assertions.Order.OrdersTypeDetail
	tbdString := internalmessages.OrdersTypeDetail("TBD")
	if orderTypeDetail == nil || *orderTypeDetail == "" {
		orderTypeDetail = &tbdString
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
		Grade:               grade,
		Entitlement:         &entitlement,
		EntitlementID:       &entitlement.ID,
		OriginDutyStation:   &originDutyStation,
		OriginDutyStationID: &originDutyStation.ID,
		OrdersTypeDetail:    orderTypeDetail,
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
