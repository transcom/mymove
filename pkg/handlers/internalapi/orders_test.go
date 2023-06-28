package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
)

func (suite *HandlerSuite) TestCreateOrder() {
	sm := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)
	dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), dutyLocation.Address.PostalCode, "KKFA")
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	req := httptest.NewRequest("POST", "/orders", nil)
	req = suite.AuthenticateRequest(req, sm)

	hasDependents := true
	spouseHasProGear := true
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	deptIndicator := internalmessages.DeptIndicatorAIRFORCE
	payload := &internalmessages.CreateUpdateOrders{
		HasDependents:       handlers.FmtBool(hasDependents),
		SpouseHasProGear:    handlers.FmtBool(spouseHasProGear),
		IssueDate:           handlers.FmtDate(issueDate),
		ReportByDate:        handlers.FmtDate(reportByDate),
		OrdersType:          internalmessages.NewOrdersType(ordersType),
		NewDutyLocationID:   handlers.FmtUUID(dutyLocation.ID),
		ServiceMemberID:     handlers.FmtUUID(sm.ID),
		OrdersNumber:        handlers.FmtString("123456"),
		Tac:                 handlers.FmtString("E19A"),
		Sac:                 handlers.FmtString("SacNumber"),
		DepartmentIndicator: internalmessages.NewDeptIndicator(deptIndicator),
	}

	params := ordersop.CreateOrdersParams{
		HTTPRequest:  req,
		CreateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	createHandler := CreateOrdersHandler{handlerConfig}

	response := createHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
	okResponse := response.(*ordersop.CreateOrdersCreated)
	orderID := okResponse.Payload.ID.String()
	createdOrder, _ := models.FetchOrder(suite.DB(), uuid.FromStringOrNil(orderID))

	suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Len(okResponse.Payload.Moves, 1)
	suite.Assertions.Equal(ordersType, *okResponse.Payload.OrdersType)
	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(handlers.FmtString("E19A"), okResponse.Payload.Tac)
	suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
	suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
	suite.Equal(sm.DutyLocationID, createdOrder.OriginDutyLocationID)
	suite.Equal((*string)(sm.Rank), createdOrder.Grade)
	suite.Assertions.Equal(*models.Int64Pointer(8000), *okResponse.Payload.AuthorizedWeight)
	suite.NotNil(&createdOrder.Entitlement)
	suite.NotEmpty(createdOrder.SupplyAndServicesCostEstimate)
	suite.NotEmpty(createdOrder.PackingAndShippingInstructions)
	suite.NotEmpty(createdOrder.MethodOfPayment)
	suite.NotEmpty(createdOrder.NAICS)
}

func (suite *HandlerSuite) TestShowOrder() {
	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
		},
	}, nil)
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)
	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.ShowOrdersParams{
		HTTPRequest: req,
		OrdersID:    *handlers.FmtUUID(order.ID),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	showHandler := ShowOrdersHandler{handlerConfig}

	response := showHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.ShowOrdersOK{}, response)
	okResponse := response.(*ordersop.ShowOrdersOK)

	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Equal(order.OrdersType, *okResponse.Payload.OrdersType)
	suite.Assertions.Equal(order.OrdersTypeDetail, okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(*order.Grade, *okResponse.Payload.Grade)
	suite.Assertions.Equal(*order.TAC, *okResponse.Payload.Tac)
	suite.Assertions.Equal(*order.DepartmentIndicator, string(*okResponse.Payload.DepartmentIndicator))
	//suite.Assertions.Equal(order.IssueDate.String(), okResponse.Payload.IssueDate.String()) // TODO: get date formats aligned
	//suite.Assertions.Equal(order.ReportByDate.String(), okResponse.Payload.ReportByDate.String())
	suite.Assertions.Equal(order.HasDependents, *okResponse.Payload.HasDependents)
	suite.Assertions.Equal(order.SpouseHasProGear, *okResponse.Payload.SpouseHasProGear)
}

func (suite *HandlerSuite) TestUploadAmendedOrder() {
	var moves models.Moves
	mto := factory.BuildMove(suite.DB(), nil, nil)
	order := mto.Orders
	order.Moves = append(moves, mto)
	path := fmt.Sprintf("/orders/%v/upload_amended_orders", order.ID.String())
	req := httptest.NewRequest("PATCH", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.UploadAmendedOrdersParams{
		HTTPRequest: req,
		File:        suite.Fixture("test.pdf"),
		OrdersID:    *handlers.FmtUUID(order.ID),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	uploadAmendedHandler := UploadAmendedOrdersHandler{
		HandlerConfig: handlerConfig,
		OrderUpdater:  orderservice.NewOrderUpdater(move.NewMoveRouter()),
	}
	response := uploadAmendedHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.UploadAmendedOrdersCreated{}, response)
	okResponse := response.(*ordersop.UploadAmendedOrdersCreated)
	suite.Assertions.NotNil(okResponse.Payload.ID.String()) // UploadPayload
	suite.Assertions.Equal("test.pdf", okResponse.Payload.Filename)
}

// TODO: Fix now that we capture transaction error. May be a data setup problem
/*
func (suite *HandlerSuite) TestUpdateOrder() {
	order := factory.BuildOrder(suite.DB(), nil, nil)

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("PUT", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	newOrdersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
	departmentIndicator := internalmessages.DeptIndicatorAIRFORCE
	otherServiceMemberUUID := uuid.Must(uuid.NewV4())

	payload := &internalmessages.CreateUpdateOrders{
		OrdersNumber:        handlers.FmtString("123456"),
		HasDependents:       handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:    handlers.FmtBool(order.SpouseHasProGear),
		IssueDate:           handlers.FmtDate(order.IssueDate),
		ReportByDate:        handlers.FmtDate(order.ReportByDate),
		OrdersType:          newOrdersType,
		OrdersTypeDetail:    &newOrdersTypeDetail,
		NewDutyLocationID:    handlers.FmtUUID(order.NewDutyLocationID),
		Tac:                 order.TAC,
		Sac:                 handlers.FmtString("N3TEST"),
		DepartmentIndicator: &departmentIndicator,
		// Attempt to assign to another service member
		ServiceMemberID: handlers.FmtUUID(otherServiceMemberUUID),
	}

	params := ordersop.UpdateOrdersParams{
		HTTPRequest:  req,
		OrdersID:     *handlers.FmtUUID(order.ID),
		UpdateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := handlers.NewHandlerCOnfig(suite.DB(), suite.TestLogger())
	handlerConfig.SetFileStorer(fakeS3)
	updateHandler := UpdateOrdersHandler{context}

	response := updateHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.UpdateOrdersOK{}, response)
	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String(), "service member id should not change")
	suite.Assertions.Equal(newOrdersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(newOrdersTypeDetail, *okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(handlers.FmtString("N3TEST"), okResponse.Payload.Sac)
}
*/
