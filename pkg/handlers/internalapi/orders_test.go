package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/gofrs/uuid"

	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateOrder() {
	sm := testdatagen.MakeDefaultServiceMember(suite.DB())
	station := testdatagen.FetchOrMakeDefaultDutyStation(suite.DB())

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
		OrdersType:          ordersType,
		NewDutyStationID:    handlers.FmtUUID(station.ID),
		ServiceMemberID:     handlers.FmtUUID(sm.ID),
		OrdersNumber:        handlers.FmtString("123456"),
		ParagraphNumber:     handlers.FmtString("123"),
		OrdersIssuingAgency: handlers.FmtString("Test Agency"),
		Tac:                 handlers.FmtString("TacNumber"),
		Sac:                 handlers.FmtString("SacNumber"),
		DepartmentIndicator: &deptIndicator,
	}

	params := ordersop.CreateOrdersParams{
		HTTPRequest:  req,
		CreateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetFileStorer(fakeS3)
	createHandler := CreateOrdersHandler{context}

	response := createHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
	okResponse := response.(*ordersop.CreateOrdersCreated)

	suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Len(okResponse.Payload.Moves, 1)
	suite.Assertions.Equal(ordersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(handlers.FmtString("123"), okResponse.Payload.ParagraphNumber)
	suite.Assertions.Equal(handlers.FmtString("Test Agency"), okResponse.Payload.OrdersIssuingAgency)
	suite.Assertions.Equal(handlers.FmtString("TacNumber"), okResponse.Payload.Tac)
	suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
	suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
}

func (suite *HandlerSuite) TestShowOrder() {
	order := testdatagen.MakeDefaultOrder(suite.DB())

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.ShowOrdersParams{
		HTTPRequest: req,
		OrdersID:    *handlers.FmtUUID(order.ID),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetFileStorer(fakeS3)
	showHandler := ShowOrdersHandler{context}

	response := showHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.ShowOrdersOK{}, response)
	okResponse := response.(*ordersop.ShowOrdersOK)

	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Equal(order.OrdersType, okResponse.Payload.OrdersType)
}

func (suite *HandlerSuite) TestUpdateOrder() {
	order := testdatagen.MakeDefaultOrder(suite.DB())

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
		NewDutyStationID:    handlers.FmtUUID(order.NewDutyStationID),
		Tac:                 order.TAC,
		Sac:                 handlers.FmtString("N3TEST"),
		OrdersIssuingAgency: handlers.FmtString("TEST AGENCY"),
		ParagraphNumber:     handlers.FmtString("123456"),
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
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetFileStorer(fakeS3)
	updateHandler := UpdateOrdersHandler{context}

	response := updateHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.UpdateOrdersOK{}, response)
	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String(), "service member id should not change")
	suite.Assertions.Equal(newOrdersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(newOrdersTypeDetail, *okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(handlers.FmtString("N3TEST"), okResponse.Payload.Sac)
	suite.Assertions.Equal(handlers.FmtString("TEST AGENCY"), okResponse.Payload.OrdersIssuingAgency)
	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.ParagraphNumber)
}
