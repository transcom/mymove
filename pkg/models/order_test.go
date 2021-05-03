package models_test

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

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

func (suite *ModelSuite) TestTacNotNilAfterSubmission() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	order.TAC = nil
	err := move.Submit()
	if err != nil {
		suite.T().Fatal("Should transition.")
	}
	suite.MustSave(&move)
	err = suite.DB().Load(&order, "Moves")
	suite.NoError(err)

	expErrors := map[string][]string{
		"transportation_accounting_code": {"TransportationAccountingCode cannot be blank."},
	}

	suite.verifyValidationErrors(&order, expErrors)
}

func (suite *ModelSuite) TestOrdersNumberPresenceAfterSubmission() {
	invalidCases := []struct {
		desc  string
		value *string
	}{
		{"EmptyString", swag.String("")},
		{"Nil", nil},
	}
	for _, invalidCase := range invalidCases {
		move := testdatagen.MakeDefaultMove(suite.DB())
		order := move.Orders
		order.OrdersNumber = invalidCase.value
		err := move.Submit()
		if err != nil {
			suite.T().Fatal("Should transition.")
		}
		suite.MustSave(&move)
		err = suite.DB().Load(&order, "Moves")
		suite.NoError(err)

		expErrors := map[string][]string{
			"orders_number": {"OrdersNumber cannot be blank."},
		}

		suite.verifyValidationErrors(&order, expErrors)
	}
}

func (suite *ModelSuite) TestOrdersTypeDetailPresenceAfterSubmission() {
	emptyString := internalmessages.OrdersTypeDetail("")

	invalidCases := []struct {
		desc  string
		value *internalmessages.OrdersTypeDetail
	}{
		{"EmptyString", &emptyString},
		{"Nil", nil},
	}
	for _, invalidCase := range invalidCases {
		move := testdatagen.MakeDefaultMove(suite.DB())
		order := move.Orders

		order.OrdersTypeDetail = invalidCase.value
		err := move.Submit()
		if err != nil {
			suite.T().Fatal("Should transition.")
		}
		suite.MustSave(&move)
		err = suite.DB().Load(&order, "Moves")
		suite.NoError(err)

		expErrors := map[string][]string{
			"orders_type_detail": {"OrdersTypeDetail cannot be blank."},
		}

		suite.verifyValidationErrors(&order, expErrors)
	}
}

func (suite *ModelSuite) TestDepartmentIndicatorNotNilAfterSubmission() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	order.DepartmentIndicator = nil
	err := move.Submit()
	if err != nil {
		suite.T().Fatal("Should transition.")
	}
	suite.MustSave(&move)
	err = suite.DB().Load(&order, "Moves")
	suite.NoError(err)

	expErrors := map[string][]string{
		"department_indicator": {"DepartmentIndicator cannot be blank."},
	}

	suite.verifyValidationErrors(&order, expErrors)
}

func (suite *ModelSuite) TestFetchOrderForUser() {

	serviceMember1 := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMember2 := testdatagen.MakeDefaultServiceMember(suite.DB())

	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	dutyStation2 := testdatagen.FetchOrMakeDefaultNewOrdersDutyStation(suite.DB())
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.MustSave(&uploadedOrder)

	SAC := "N002214CSW32Y9"
	ordersNumber := "FD4534JFJ"

	order := Order{
		ServiceMemberID:     serviceMember1.ID,
		ServiceMember:       serviceMember1,
		IssueDate:           issueDate,
		ReportByDate:        reportByDate,
		OrdersType:          ordersType,
		HasDependents:       hasDependents,
		SpouseHasProGear:    spouseHasProGear,
		OriginDutyStationID: &dutyStation.ID,
		OriginDutyStation:   &dutyStation,
		NewDutyStationID:    dutyStation2.ID,
		NewDutyStation:      dutyStation2,
		UploadedOrdersID:    uploadedOrder.ID,
		UploadedOrders:      uploadedOrder,
		Status:              OrderStatusSUBMITTED,
		OrdersNumber:        &ordersNumber,
		TAC:                 &TAC,
		SAC:                 &SAC,
		DepartmentIndicator: &deptIndicator,
		Grade:               swag.String("E-3"),
	}
	suite.MustSave(&order)

	// User is authorized to fetch order
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          serviceMember1.UserID,
		ServiceMemberID: serviceMember1.ID,
	}
	goodOrder, err := FetchOrderForUser(suite.DB(), session, order.ID)

	if suite.NoError(err) {
		suite.True(order.IssueDate.Equal(goodOrder.IssueDate))
		suite.True(order.ReportByDate.Equal(goodOrder.ReportByDate))
		suite.Equal(order.OrdersType, goodOrder.OrdersType)
		suite.Equal(order.HasDependents, goodOrder.HasDependents)
		suite.Equal(order.SpouseHasProGear, goodOrder.SpouseHasProGear)
		suite.Equal(order.OriginDutyStation.ID, goodOrder.OriginDutyStation.ID)
		suite.Equal(order.NewDutyStation.ID, goodOrder.NewDutyStation.ID)
		suite.Equal(order.Grade, goodOrder.Grade)
		suite.Equal(order.UploadedOrdersID, goodOrder.UploadedOrdersID)
	}

	// Wrong Order ID
	wrongID, _ := uuid.NewV4()
	_, err = FetchOrderForUser(suite.DB(), session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}
	// User is forbidden from fetching order
	session.UserID = serviceMember2.UserID
	session.ServiceMemberID = serviceMember2.ID
	_, err = FetchOrderForUser(suite.DB(), session, order.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}
}

func (suite *ModelSuite) TestFetchOrderNotForUser() {

	serviceMember1 := testdatagen.MakeDefaultServiceMember(suite.DB())

	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.MustSave(&uploadedOrder)
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
	suite.MustSave(&order)

	// No session
	goodOrder, err := FetchOrder(suite.DB(), order.ID)
	if suite.NoError(err) {
		suite.True(order.IssueDate.Equal(goodOrder.IssueDate))
		suite.True(order.ReportByDate.Equal(goodOrder.ReportByDate))
		suite.Equal(order.OrdersType, goodOrder.OrdersType)
		suite.Equal(order.HasDependents, goodOrder.HasDependents)
		suite.Equal(order.SpouseHasProGear, goodOrder.SpouseHasProGear)
		suite.Equal(order.NewDutyStationID, goodOrder.NewDutyStationID)
	}
}

func (suite *ModelSuite) TestOrderStateMachine() {
	serviceMember1 := testdatagen.MakeDefaultServiceMember(suite.DB())

	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.MustSave(&uploadedOrder)
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
	suite.MustSave(&order)

	// Submit Orders
	err := order.Submit()
	suite.NoError(err)
	suite.Equal(OrderStatusSUBMITTED, order.Status, "expected Submitted")

	// Can cancel orders
	err = order.Cancel()
	suite.NoError(err)
	suite.Equal(OrderStatusCANCELED, order.Status, "expected Canceled")
}

func (suite *ModelSuite) TestCanceledMoveCancelsOrder() {
	serviceMember1 := testdatagen.MakeDefaultServiceMember(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

	dutyStation := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember1,
		ServiceMemberID: serviceMember1.ID,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.MustSave(&uploadedOrder)
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
	suite.MustSave(&orders)

	selectedMoveType := SelectedMoveTypeHHGPPM
	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	move, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")
	move.Orders = orders
	suite.MustSave(move)

	err = move.Submit()
	suite.NoError(err)

	reason := "Mistaken identity"
	err = move.Cancel(reason)
	suite.NoError(err)
	suite.Equal(MoveStatusCANCELED, move.Status, "expected Canceled")
	suite.Equal(OrderStatusCANCELED, move.Orders.Status, "expected Canceled")

}

func (suite *ModelSuite) TestSaveOrder() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")

	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: Order{
			ID: orderID,
		},
	})
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: Move{
			ID:       moveID,
			OrdersID: orderID,
			Orders:   order,
		},
	})

	postalCode := "30813"
	newPostalCode := "12345"
	address := Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     newPostalCode,
	}
	suite.MustSave(&address)

	stationName := "New Duty Station"
	station := DutyStation{
		Name:        stationName,
		Affiliation: internalmessages.AffiliationAIRFORCE,
		AddressID:   address.ID,
		Address:     address,
	}
	suite.MustSave(&station)

	advance := BuildDraftReimbursement(1000, MethodOfReceiptMILPAY)
	_, verrs, err := move.CreatePPM(suite.DB(), nil, nil, nil, nil, nil, swag.String("55555"), nil, nil, nil, true, &advance)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	suite.Equal(postalCode, order.NewDutyStation.Address.PostalCode, "Wrong orig postal code")
	order.NewDutyStationID = station.ID
	order.NewDutyStation = station
	verrs, err = SaveOrder(suite.DB(), &order)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	orderUpdated, err := FetchOrder(suite.DB(), orderID)
	suite.NoError(err)
	suite.Equal(station.ID, orderUpdated.NewDutyStationID, "Wrong order new_duty_station_id")
	suite.Equal(newPostalCode, order.NewDutyStation.Address.PostalCode, "Wrong orig postal code")

	ppm, err := FetchPersonallyProcuredMoveByOrderID(suite.DB(), orderUpdated.ID)
	suite.NoError(err)
	suite.Equal(newPostalCode, *ppm.DestinationPostalCode, "Wrong ppm postal code")
}
