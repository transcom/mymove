package models_test

import (
	"time"

	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicOrderInstantiation() {
	order := &Order{}

	expErrors := map[string][]string{
		"orders_type":         {"OrdersType can not be blank."},
		"issue_date":          {"IssueDate can not be blank."},
		"report_by_date":      {"ReportByDate can not be blank."},
		"service_member_id":   {"ServiceMemberID can not be blank."},
		"new_duty_station_id": {"NewDutyStationID can not be blank."},
		"status":              {"Status can not be blank."},
	}

	suite.verifyValidationErrors(order, expErrors)
}

func (suite *ModelSuite) TestFetchOrder() {

	serviceMember1, _ := testdatagen.MakeServiceMember(suite.db)
	serviceMember2, _ := testdatagen.MakeServiceMember(suite.db)

	dutyStation := testdatagen.MakeAnyDutyStation(suite.db)
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
		Name:            UploadedOrdersDocumentName,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.mustSave(&uploadedOrder)
	order := Order{
		ServiceMemberID:     serviceMember1.ID,
		ServiceMember:       serviceMember1,
		IssueDate:           issueDate,
		ReportByDate:        reportByDate,
		OrdersType:          ordersType,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		NewDutyStationID:    dutyStation.ID,
		NewDutyStation:      dutyStation,
		UploadedOrdersID:    uploadedOrder.ID,
		UploadedOrders:      uploadedOrder,
		Status:              OrderStatusSUBMITTED,
		TAC:                 &TAC,
		DepartmentIndicator: &deptIndicator,
	}
	suite.mustSave(&order)

	// User is authorized to fetch order
	session := &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          serviceMember1.UserID,
		ServiceMemberID: serviceMember1.ID,
	}
	goodOrder, err := FetchOrder(suite.db, session, order.ID)
	if suite.NoError(err) {
		suite.True(order.IssueDate.Equal(goodOrder.IssueDate))
		suite.True(order.ReportByDate.Equal(goodOrder.ReportByDate))
		suite.Equal(order.OrdersType, goodOrder.OrdersType)
		suite.Equal(order.HasDependents, goodOrder.HasDependents)
		suite.Equal(order.SpouseHasProGear, goodOrder.SpouseHasProGear)
		suite.Equal(order.NewDutyStation.ID, goodOrder.NewDutyStation.ID)
	}

	// Wrong Order ID
	wrongID, _ := uuid.NewV4()
	_, err = FetchOrder(suite.db, session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}
	// User is forbidden from fetching order
	session.UserID = serviceMember2.UserID
	session.ServiceMemberID = serviceMember2.ID
	_, err = FetchOrder(suite.db, session, order.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}

}

func (suite *ModelSuite) TestOrderStateMachine() {
	serviceMember1, _ := testdatagen.MakeServiceMember(suite.db)

	dutyStation := testdatagen.MakeAnyDutyStation(suite.db)
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
		Name:            UploadedOrdersDocumentName,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.mustSave(&uploadedOrder)
	order := Order{
		ServiceMemberID:     serviceMember1.ID,
		ServiceMember:       serviceMember1,
		IssueDate:           issueDate,
		ReportByDate:        reportByDate,
		OrdersType:          ordersType,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		NewDutyStationID:    dutyStation.ID,
		NewDutyStation:      dutyStation,
		UploadedOrdersID:    uploadedOrder.ID,
		UploadedOrders:      uploadedOrder,
		Status:              OrderStatusDRAFT,
		TAC:                 &TAC,
		DepartmentIndicator: &deptIndicator,
	}
	suite.mustSave(&order)

	// Can't cancel Orders with DRAFT status
	err := order.Cancel()
	suite.Equal(ErrInvalidTransition, errors.Cause(err))

	// Submit Orders
	err = order.Submit()
	suite.Nil(err)
	suite.Equal(OrderStatusSUBMITTED, order.Status, "expected Submitted")

	// Can cancel orders
	err = order.Cancel()
	suite.Nil(err)
	suite.Equal(OrderStatusCANCELED, order.Status, "expected Canceled")
}

func (suite *ModelSuite) TestCanceledMoveCancelsOrder() {
	serviceMember1, _ := testdatagen.MakeServiceMember(suite.db)

	dutyStation := testdatagen.MakeAnyDutyStation(suite.db)
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
		Name:            UploadedOrdersDocumentName,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	selectedType := internalmessages.SelectedMoveTypeCOMBO
	suite.mustSave(&uploadedOrder)
	orders := Order{
		ServiceMemberID:     serviceMember1.ID,
		ServiceMember:       serviceMember1,
		IssueDate:           issueDate,
		ReportByDate:        reportByDate,
		OrdersType:          ordersType,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		NewDutyStationID:    dutyStation.ID,
		NewDutyStation:      dutyStation,
		UploadedOrdersID:    uploadedOrder.ID,
		UploadedOrders:      uploadedOrder,
		Status:              OrderStatusSUBMITTED,
		TAC:                 &TAC,
		DepartmentIndicator: &deptIndicator,
	}
	suite.mustSave(&orders)

	move, verrs, err := orders.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders
	suite.mustSave(move)

	err = move.Submit()
	suite.Nil(err)

	reason := "Mistaken identity"
	err = move.Cancel(reason)
	suite.Nil(err)
	suite.Equal(MoveStatusCANCELED, move.GetStatus(), "expected Canceled")
	suite.Equal(OrderStatusCANCELED, move.Orders.Status, "expected Canceled")

}
