package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicOrderInstantiation() {
	order := &Order{
		TAC:    StringPointer(""),
		SAC:    StringPointer(""),
		NtsTAC: StringPointer(""),
		NtsSAC: StringPointer(""),
	}

	expErrors := map[string][]string{
		"orders_type":                       {"OrdersType can not be blank."},
		"issue_date":                        {"IssueDate can not be blank."},
		"report_by_date":                    {"ReportByDate can not be blank."},
		"service_member_id":                 {"ServiceMemberID can not be blank."},
		"new_duty_location_id":              {"NewDutyLocationID can not be blank."},
		"status":                            {"Status can not be blank."},
		"uploaded_orders_id":                {"UploadedOrdersID can not be blank."},
		"transportation_accounting_code":    {"TAC must be exactly 4 alphanumeric characters.", "TransportationAccountingCode can not be blank."},
		"sac":                               {"SAC can not be blank."},
		"nts_tac":                           {"NtsTAC can not be blank."},
		"nts_sac":                           {"NtsSAC can not be blank."},
		"supply_and_services_cost_estimate": {"SupplyAndServicesCostEstimate can not be blank."},
		"method_of_payment":                 {"MethodOfPayment can not be blank."},
		"naics":                             {"NAICS can not be blank."},
		"packing_and_shipping_instructions": {"PackingAndShippingInstructions can not be blank."},
	}

	suite.verifyValidationErrors(order, expErrors)
}

func (suite *ModelSuite) TestMiscValidationsAfterSubmission() {
	move := factory.BuildStubbedMoveWithStatus(MoveStatusSUBMITTED)
	order := move.Orders
	order.Moves = append(order.Moves, move)

	suite.Run("test valid UploadedAmendedOrdersID", func() {
		testUUID := uuid.Must(uuid.NewV4())
		order.UploadedAmendedOrdersID = &testUUID

		expErrors := map[string][]string{}

		suite.verifyValidationErrors(&order, expErrors)
	})

	suite.Run("test UploadedAmendedOrdersID is not nil UUID", func() {
		order.UploadedAmendedOrdersID = &uuid.Nil

		expErrors := map[string][]string{
			"uploaded_amended_orders_id": {"UploadedAmendedOrdersID can not be blank."},
		}

		suite.verifyValidationErrors(&order, expErrors)
	})
}

func (suite *ModelSuite) TestTacCanBeNilBeforeSubmissionToTOO() {
	validStatuses := []struct {
		desc  string
		value MoveStatus
	}{
		{"Draft", MoveStatusDRAFT},
		{"NeedsServiceCounseling", MoveStatusNeedsServiceCounseling},
	}
	for _, validStatus := range validStatuses {
		move := factory.BuildStubbedMoveWithStatus(validStatus.value)
		order := move.Orders
		order.TAC = nil
		order.Moves = append(order.Moves, move)

		expErrors := map[string][]string{}

		suite.verifyValidationErrors(&order, expErrors)
	}
}

func (suite *ModelSuite) TestTacFormat() {
	invalidCases := []struct {
		desc string
		tac  string
	}{
		{"TestOneCharacter", "A"},
		{"TestTwoCharacters", "AB"},
		{"TestThreeCharacters", "ABC"},
		{"TestGreaterThanFourChars", "ABCD1"},
		{"TestNonAlphaNumChars", "AB-C"},
	}
	for _, invalidCase := range invalidCases {
		move := factory.BuildStubbedMoveWithStatus(MoveStatusSUBMITTED)
		order := move.Orders
		order.TAC = &invalidCase.tac
		order.Moves = append(order.Moves, move)

		expErrors := map[string][]string{
			"transportation_accounting_code": {"TAC must be exactly 4 alphanumeric characters."},
		}

		suite.verifyValidationErrors(&order, expErrors)
	}
}

func (suite *ModelSuite) TestFetchOrderForUser() {

	suite.Run("successful fetch by authorized user", func() {
		order := factory.BuildOrder(suite.DB(), nil, nil)

		// User is authorized to fetch order
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          order.ServiceMember.UserID,
			ServiceMemberID: order.ServiceMemberID,
		}
		goodOrder, err := FetchOrderForUser(suite.DB(), session, order.ID)

		suite.NoError(err)
		suite.True(order.IssueDate.Equal(goodOrder.IssueDate))
		suite.True(order.ReportByDate.Equal(goodOrder.ReportByDate))
		suite.Equal(order.OrdersType, goodOrder.OrdersType)
		suite.Equal(order.HasDependents, goodOrder.HasDependents)
		suite.Equal(order.SpouseHasProGear, goodOrder.SpouseHasProGear)
		suite.Equal(order.OriginDutyLocation.ID, goodOrder.OriginDutyLocation.ID)
		suite.Equal(order.NewDutyLocation.ID, goodOrder.NewDutyLocation.ID)
		suite.Equal(order.Grade, goodOrder.Grade)
		suite.Equal(order.UploadedOrdersID, goodOrder.UploadedOrdersID)
	})

	suite.Run("check for closeout office", func() {
		closeoutOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    closeoutOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		}, nil)
		orders := move.Orders
		orders.Moves = append(orders.Moves, move)

		// User is authorized to fetch order
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          orders.ServiceMember.UserID,
			ServiceMemberID: orders.ServiceMemberID,
		}

		goodOrder, err := FetchOrderForUser(suite.DB(), session, orders.ID)

		suite.NoError(err)
		suite.Equal(orders.Moves[0].CloseoutOffice.ID, goodOrder.Moves[0].CloseoutOffice.ID)
		suite.Equal(orders.Moves[0].CloseoutOffice.Name, goodOrder.Moves[0].CloseoutOffice.Name)
		suite.Equal(orders.Moves[0].CloseoutOffice.Address.ID, goodOrder.Moves[0].CloseoutOffice.Address.ID)
		suite.Equal(orders.Moves[0].CloseoutOffice.Gbloc, goodOrder.Moves[0].CloseoutOffice.Gbloc)

	})

	suite.Run("fetch not found due to bad id", func() {
		sm := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          sm.UserID,
			ServiceMemberID: sm.ID,
		}
		// Wrong Order ID
		wrongID, _ := uuid.NewV4()
		_, err := FetchOrderForUser(suite.DB(), session, wrongID)

		suite.Error(err)
		suite.Equal(ErrFetchNotFound, err)
	})

	suite.Run("forbidden user cannot fetch order", func() {
		order := factory.BuildOrder(suite.DB(), nil, nil)
		// User is forbidden from fetching order
		serviceMember2 := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          serviceMember2.UserID,
			ServiceMemberID: serviceMember2.ID,
		}
		_, err := FetchOrderForUser(suite.DB(), session, order.ID)

		suite.Error(err)
		suite.Equal(ErrFetchForbidden, err)
	})

	suite.Run("successfully excludes deleted orders uploads", func() {
		nonDeletedOrdersUpload := factory.BuildUserUpload(suite.DB(), nil, nil)
		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedOrdersUpload.Document,
				LinkOnly: true,
			},
			{
				Model: UserUpload{
					DeletedAt: TimePointer(time.Now()),
				},
			},
		}, nil)

		nonDeletedAmendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model: UserUpload{
					UploaderID: nonDeletedOrdersUpload.Document.ServiceMember.UserID,
				},
			},
		}, nil)
		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedAmendedUpload.Document,
				LinkOnly: true,
			},
			{
				Model: UserUpload{
					DeletedAt: TimePointer(time.Now()),
				},
			},
		}, nil)

		expectedOrder := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedOrdersUpload.Document.ServiceMember,
				LinkOnly: true,
			},
			{
				Model:    nonDeletedOrdersUpload.Document,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedOrders,
			},
			{
				Model:    nonDeletedAmendedUpload.Document,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		userSession := auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          expectedOrder.ServiceMember.ID,
			ServiceMemberID: expectedOrder.ServiceMemberID,
		}

		actualOrder, err := FetchOrderForUser(suite.DB(), &userSession, expectedOrder.ID)

		suite.NoError(err)
		suite.Len(actualOrder.UploadedOrders.UserUploads, 1)
		suite.Equal(actualOrder.UploadedOrders.UserUploads[0].ID, nonDeletedOrdersUpload.ID)
		suite.Len(actualOrder.UploadedAmendedOrders.UserUploads, 1)
		suite.Equal(actualOrder.UploadedAmendedOrders.UserUploads[0].ID, nonDeletedAmendedUpload.ID)
	})
}

func (suite *ModelSuite) TestFetchOrderNotForUser() {
	serviceMember1 := factory.BuildServiceMember(suite.DB(), nil, nil)

	dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
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
	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	packingAndShippingInstructions := InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + InstructionsAfterContractNumber
	order := Order{
		ServiceMemberID:                serviceMember1.ID,
		ServiceMember:                  serviceMember1,
		IssueDate:                      issueDate,
		ReportByDate:                   reportByDate,
		OrdersType:                     ordersType,
		HasDependents:                  hasDependents,
		SpouseHasProGear:               spouseHasProGear,
		NewDutyLocationID:              dutyLocation.ID,
		NewDutyLocation:                dutyLocation,
		UploadedOrdersID:               uploadedOrder.ID,
		UploadedOrders:                 uploadedOrder,
		Status:                         OrderStatusSUBMITTED,
		TAC:                            &TAC,
		DepartmentIndicator:            &deptIndicator,
		SupplyAndServicesCostEstimate:  SupplyAndServicesCostEstimate,
		MethodOfPayment:                MethodOfPayment,
		NAICS:                          NAICS,
		PackingAndShippingInstructions: packingAndShippingInstructions,
	}
	suite.MustSave(&order)

	// No session
	goodOrder, err := FetchOrder(suite.DB(), order.ID)
	suite.NoError(err)
	suite.True(order.IssueDate.Equal(goodOrder.IssueDate))
	suite.True(order.ReportByDate.Equal(goodOrder.ReportByDate))
	suite.Equal(order.OrdersType, goodOrder.OrdersType)
	suite.Equal(order.HasDependents, goodOrder.HasDependents)
	suite.Equal(order.SpouseHasProGear, goodOrder.SpouseHasProGear)
	suite.Equal(order.NewDutyLocationID, goodOrder.NewDutyLocationID)
}

func (suite *ModelSuite) TestOrderStateMachine() {
	serviceMember1 := factory.BuildServiceMember(suite.DB(), nil, nil)

	dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
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
	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
	packingAndShippingInstructions := InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + InstructionsAfterContractNumber
	suite.MustSave(&uploadedOrder)
	order := Order{
		ServiceMemberID:                serviceMember1.ID,
		ServiceMember:                  serviceMember1,
		IssueDate:                      issueDate,
		ReportByDate:                   reportByDate,
		OrdersType:                     ordersType,
		HasDependents:                  hasDependents,
		SpouseHasProGear:               spouseHasProGear,
		NewDutyLocationID:              dutyLocation.ID,
		NewDutyLocation:                dutyLocation,
		UploadedOrdersID:               uploadedOrder.ID,
		UploadedOrders:                 uploadedOrder,
		Status:                         OrderStatusDRAFT,
		TAC:                            &TAC,
		DepartmentIndicator:            &deptIndicator,
		SupplyAndServicesCostEstimate:  SupplyAndServicesCostEstimate,
		MethodOfPayment:                MethodOfPayment,
		NAICS:                          NAICS,
		PackingAndShippingInstructions: packingAndShippingInstructions,
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

func (suite *ModelSuite) TestSaveOrder() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: Move{
				ID: moveID,
			},
		},
		{
			Model: Order{
				ID: orderID,
			},
		},
	}, nil)

	order := move.Orders

	postalCode := "30813"
	newPostalCode := "12345"
	address := Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     newPostalCode,
	}
	suite.MustSave(&address)

	dutyLocationName := "New Duty Location"
	location := DutyLocation{
		Name:      dutyLocationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	suite.Equal(postalCode, order.NewDutyLocation.Address.PostalCode, "Wrong orig postal code")
	order.NewDutyLocationID = location.ID
	order.NewDutyLocation = location
	verrs, err := SaveOrder(suite.DB(), &order)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	orderUpdated, err := FetchOrder(suite.DB(), orderID)
	suite.NoError(err)
	suite.Equal(location.ID, orderUpdated.NewDutyLocationID, "Wrong order new_duty_location_id")
	suite.Equal(newPostalCode, order.NewDutyLocation.Address.PostalCode, "Wrong orig postal code")

}

func (suite *ModelSuite) TestSaveOrderWithoutPPM() {
	orderID := uuid.Must(uuid.NewV4())
	moveID, _ := uuid.FromString("7112b18b-7e03-4b28-adde-532b541bba8d")

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: Move{
				ID: moveID,
			},
		},
		{
			Model: Order{
				ID: orderID,
			},
		},
	}, nil)

	order := move.Orders

	postalCode := "30813"
	newPostalCode := "12345"
	address := Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     newPostalCode,
	}
	suite.MustSave(&address)

	dutyLocationName := "New Duty Location"
	location := DutyLocation{
		Name:      dutyLocationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	suite.Equal(postalCode, order.NewDutyLocation.Address.PostalCode, "Wrong orig postal code")

	order.NewDutyLocationID = location.ID
	order.NewDutyLocation = location

	verrs, err := SaveOrder(suite.DB(), &order)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	orderUpdated, err := FetchOrder(suite.DB(), orderID)
	suite.NoError(err)
	suite.Equal(location.ID, orderUpdated.NewDutyLocationID, "Wrong order new_duty_location_id")
	suite.Equal(newPostalCode, order.NewDutyLocation.Address.PostalCode, "Wrong orig postal code")
}
