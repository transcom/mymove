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

// MakeMoveOrder creates a single MoveOrder and associated set relationships
func MakeMoveOrder(db *pop.Connection, assertions Assertions) models.Order {
	grade := assertions.MoveOrder.Grade
	if grade == nil || *grade == "" {
		grade = stringPointer(MakeGrade())
	}
	customer := assertions.Customer
	if isZeroUUID(customer.ID) {
		customer = MakeCustomer(db, assertions)
	}
	entitlement := assertions.Entitlement
	if isZeroUUID(entitlement.ID) {
		assertions.MoveOrder.Grade = grade
		entitlement = MakeEntitlement(db, assertions)
	}
	originDutyStation := assertions.OriginDutyStation
	if isZeroUUID(originDutyStation.ID) {
		originDutyStation = MakeDutyStation(db, assertions)
	}
	destinationDutyStation := assertions.DestinationDutyStation
	if isZeroUUID(destinationDutyStation.ID) {
		destinationDutyStation = MakeDutyStation(db, assertions)
	}

	orderNumber := assertions.MoveOrder.OrdersNumber
	if orderNumber == nil || *orderNumber == "" {
		orderNumber = stringPointer("ORDER123")
	}

	orderType := assertions.MoveOrder.OrdersType
	if &orderType == nil || orderType == "" {
		orderType = "GHC"
	}

	orderTypeDetail := assertions.MoveOrder.OrdersTypeDetail
	tbdString := internalmessages.OrdersTypeDetail("TBD")
	if orderTypeDetail == nil || *orderTypeDetail == "" {
		orderTypeDetail = &tbdString
	}

	reportByDate := &assertions.MoveOrder.ReportByDate

	if &reportByDate == nil || time.Time.IsZero(*reportByDate) {
		reportByDate = models.TimePointer(time.Date(2020, time.February, 15, 0, 0, 0, 0, time.UTC))
	}

	dateIssued := &assertions.MoveOrder.IssueDate

	if &dateIssued == nil || time.Time.IsZero(*dateIssued) {
		dateIssued = models.TimePointer(time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC))
	}

	linesOfAccounting := "F8E1"

	document := assertions.Order.UploadedOrders
	// Note above
	if isZeroUUID(assertions.Order.UploadedOrdersID) {
		document = MakeDocument(db, Assertions{
			Document: models.Document{
				ServiceMemberID: customer.ID,
				ServiceMember:   customer,
			},
		})
		u := MakeUserUpload(db, Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &document.ID,
				Document:   document,
				UploaderID: customer.UserID,
			},
			UserUploader: assertions.UserUploader,
		})
		document.UserUploads = append(document.UserUploads, u)
	}

	moveOrder := models.Order{
		ServiceMember:       customer,
		ServiceMemberID:     customer.ID,
		ConfirmationNumber:  stringPointer(models.GenerateLocator()),
		IssueDate:           *dateIssued,
		Entitlement:         &entitlement,
		EntitlementID:       &entitlement.ID,
		NewDutyStation:      destinationDutyStation,
		NewDutyStationID:    destinationDutyStation.ID,
		Grade:               grade,
		OriginDutyStation:   &originDutyStation,
		OriginDutyStationID: &originDutyStation.ID,
		OrdersNumber:        orderNumber,
		OrdersType:          orderType,
		OrdersTypeDetail:    orderTypeDetail,
		ReportByDate:        *reportByDate,
		TAC:                 &linesOfAccounting,
		Status:              models.OrderStatusDRAFT,
		UploadedOrders:      document,
		UploadedOrdersID:    document.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&moveOrder, assertions.MoveOrder)

	mustCreate(db, &moveOrder)

	return moveOrder
}

// MakeDefaultMoveOrder makes a MoveOrder with default values
func MakeDefaultMoveOrder(db *pop.Connection) models.Order {
	return MakeMoveOrder(db, Assertions{})
}
