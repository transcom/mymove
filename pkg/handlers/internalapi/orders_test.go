package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateOrder() {
	sm := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{})
	station := testdatagen.FetchOrMakeDefaultCurrentDutyStation(suite.DB())
	testdatagen.MakeDefaultContractor(suite.DB())

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
		Tac:                 handlers.FmtString("E19A"),
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
	orderID := okResponse.Payload.ID.String()
	createdOrder, _ := models.FetchOrder(suite.DB(), uuid.FromStringOrNil(orderID))

	suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Len(okResponse.Payload.Moves, 1)
	suite.Assertions.Equal(ordersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(handlers.FmtString("E19A"), okResponse.Payload.Tac)
	suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
	suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
	suite.Equal(sm.DutyStationID, createdOrder.OriginDutyStationID)
	suite.Equal((*string)(sm.Rank), createdOrder.Grade)
	suite.Assertions.Equal(*swag.Int64(8000), *okResponse.Payload.AuthorizedWeight)
	suite.NotNil(&createdOrder.Entitlement)
}

func (suite *HandlerSuite) TestShowOrder() {
	dutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Address: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{}),
		},
	})
	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			OriginDutyStation: &dutyStation,
		},
	})
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
	suite.Assertions.Equal(order.OrdersTypeDetail, okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(*order.Grade, *okResponse.Payload.Grade)
	suite.Assertions.Equal(*order.TAC, *okResponse.Payload.Tac)
	suite.Assertions.Equal(*order.DepartmentIndicator, string(*okResponse.Payload.DepartmentIndicator))
	//suite.Assertions.Equal(order.IssueDate.String(), okResponse.Payload.IssueDate.String()) // TODO: get date formats aligned
	//suite.Assertions.Equal(order.ReportByDate.String(), okResponse.Payload.ReportByDate.String())
	suite.Assertions.Equal(order.HasDependents, *okResponse.Payload.HasDependents)
	suite.Assertions.Equal(order.SpouseHasProGear, *okResponse.Payload.SpouseHasProGear)
}

// TODO: Fix now that we capture transaction error. May be a data setup problem
/*
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
}
*/
