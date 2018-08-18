package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/gobuffalo/uuid"

	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateOrder() {
	sm := testdatagen.MakeDefaultServiceMember(suite.db)
	station := testdatagen.MakeDefaultDutyStation(suite.db)

	req := httptest.NewRequest("POST", "/orders", nil)
	req = suite.AuthenticateRequest(req, sm)

	hasDependents := true
	spouseHasProGear := true
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	payload := &internalmessages.CreateUpdateOrders{
		HasDependents:    utils.FmtBool(hasDependents),
		SpouseHasProGear: utils.FmtBool(spouseHasProGear),
		IssueDate:        utils.FmtDate(issueDate),
		ReportByDate:     utils.FmtDate(reportByDate),
		OrdersType:       ordersType,
		NewDutyStationID: utils.FmtUUID(station.ID),
		ServiceMemberID:  utils.FmtUUID(sm.ID),
	}

	params := ordersop.CreateOrdersParams{
		HTTPRequest:  req,
		CreateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	createHandler := CreateOrdersHandler(context)

	response := createHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
	okResponse := response.(*ordersop.CreateOrdersCreated)

	suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Len(okResponse.Payload.Moves, 1)
	suite.Assertions.Equal(ordersType, okResponse.Payload.OrdersType)
}

func (suite *HandlerSuite) TestShowOrder() {
	order := testdatagen.MakeDefaultOrder(suite.db)

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.ShowOrdersParams{
		HTTPRequest: req,
		OrdersID:    *utils.FmtUUID(order.ID),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	showHandler := ShowOrdersHandler(context)

	response := showHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.ShowOrdersOK{}, response)
	okResponse := response.(*ordersop.ShowOrdersOK)

	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Equal(order.OrdersType, okResponse.Payload.OrdersType)
}

func (suite *HandlerSuite) TestUpdateOrder() {
	order := testdatagen.MakeDefaultOrder(suite.db)

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("PUT", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	newOrdersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
	departmentIndicator := internalmessages.DeptIndicatorAIRFORCE
	otherServiceMemberUUID := uuid.Must(uuid.NewV4())

	payload := &internalmessages.CreateUpdateOrders{
		OrdersNumber:        utils.FmtString("123456"),
		HasDependents:       utils.FmtBool(order.HasDependents),
		SpouseHasProGear:    utils.FmtBool(order.SpouseHasProGear),
		IssueDate:           utils.FmtDate(order.IssueDate),
		ReportByDate:        utils.FmtDate(order.ReportByDate),
		OrdersType:          newOrdersType,
		OrdersTypeDetail:    &newOrdersTypeDetail,
		NewDutyStationID:    utils.FmtUUID(order.NewDutyStationID),
		Tac:                 order.TAC,
		DepartmentIndicator: &departmentIndicator,
		// Attempt to assign to another service member
		ServiceMemberID: utils.FmtUUID(otherServiceMemberUUID),
	}

	params := ordersop.UpdateOrdersParams{
		HTTPRequest:  req,
		OrdersID:     *utils.FmtUUID(order.ID),
		UpdateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetFileStorer(fakeS3)
	updateHandler := UpdateOrdersHandler(context)

	response := updateHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.UpdateOrdersOK{}, response)
	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.Assertions.Equal(utils.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String(), "service member id should not change")
	suite.Assertions.Equal(newOrdersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(newOrdersTypeDetail, *okResponse.Payload.OrdersTypeDetail)
}
