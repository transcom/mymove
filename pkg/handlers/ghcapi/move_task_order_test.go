//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package ghcapi

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveTaskOrderHandlerIntegration() {
	order := testdatagen.MakeDefaultOrder(suite.DB())
	moveTaskOrder := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: order,
	})
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: "MS",
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := move_task_order.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveTaskOrderHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveTaskOrderResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	moveTaskOrderPayload := moveTaskOrderResponse.Payload

	suite.Assertions.IsType(&move_task_order.GetMoveTaskOrderOK{}, response)
	suite.Equal(strfmt.UUID(moveTaskOrder.ID.String()), moveTaskOrderPayload.ID)
	suite.Nil(moveTaskOrderPayload.AvailableToPrimeAt)
	suite.False(*moveTaskOrderPayload.IsCanceled)
	suite.Equal(strfmt.UUID(moveTaskOrder.OrdersID.String()), moveTaskOrderPayload.OrderID)
	suite.NotNil(moveTaskOrderPayload.ReferenceID)
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationSuccess() {
	moveTaskOrder := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusSUBMITTED}})

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	request = suite.AuthenticateUserRequest(request, requestUser)

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: "MS",
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})

	serviceItemCodes := ghcmessages.MTOApprovalServiceItemCodes{
		ServiceCodeMS: true,
		ServiceCodeCS: true,
	}
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:      request,
		MoveTaskOrderID:  moveTaskOrder.ID.String(),
		IfMatch:          etag.GenerateEtag(moveTaskOrder.UpdatedAt),
		ServiceItemCodes: &serviceItemCodes,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

	// setup the handler
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
	}
	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	handler.SetTraceID(traceID)

	// make the request
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	updatedMove := models.Move{}
	suite.DB().Find(&updatedMove, moveTaskOrdersPayload.ID)
	suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.NotNil(moveTaskOrdersPayload.AvailableToPrimeAt)
	suite.HasWebhookNotification(moveTaskOrder.ID, traceID) // this action always creates a notification for the Prime

	// also check MTO level service items are properly created
	var serviceItems models.MTOServiceItems
	suite.DB().Eager("ReService").Where("move_id = ?", moveTaskOrder.ID).All(&serviceItems)
	suite.Len(serviceItems, 2, "Expected to find at most 2 service items")

	containsServiceCode := func(items models.MTOServiceItems, target models.ReServiceCode) bool {
		for _, si := range items {
			if si.ReService.Code == target {
				return true
			}
		}

		return false
	}

	suite.True(containsServiceCode(serviceItems, models.ReServiceCodeMS), "Expected to find reServiceCode, MS, in array.")
	suite.True(containsServiceCode(serviceItems, models.ReServiceCodeCS), "Expected to find reServiceCode, CS, in array.")
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationWithStaleEtag() {
	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	request = suite.AuthenticateUserRequest(request, requestUser)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
		IfMatch:         etag.GenerateEtag(time.Now()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
	}
	response := handler.Handle(params)
	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusPreconditionFailed{}, response)
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationWithIncompleteOrder() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	order.TAC = nil
	suite.MustSave(&order)
	err := move.Submit()
	if err != nil {
		suite.T().Fatal("Should transition.")
	}
	suite.MustSave(&move)

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	request = suite.AuthenticateUserRequest(request, requestUser)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: move.ID.String(),
		IfMatch:         etag.GenerateEtag(move.UpdatedAt),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
	}
	response := handler.Handle(params)

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusUnprocessableEntity{}, response)
	invalidResponse := response.(*move_task_order.UpdateMoveTaskOrderStatusUnprocessableEntity).Payload
	errorDetail := invalidResponse.Detail

	suite.Contains(*errorDetail, "TransportationAccountingCode cannot be blank.")
}

func (suite *HandlerSuite) TestUpdateMTOStatusServiceCounselingCompletedHandler() {
	order := testdatagen.MakeDefaultOrder(suite.DB())
	moveTaskOrder := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusNeedsServiceCounseling,
		},
		Order: order,
	})

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status/service-counseling-completed", nil)
	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	request = suite.AuthenticateUserRequest(request, requestUser)
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
	handler := UpdateMTOStatusServiceCounselingCompletedHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, siCreator),
	}

	params := move_task_order.UpdateMTOStatusServiceCounselingCompletedParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
		IfMatch:         etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}

	suite.T().Run("Successful move status update to Service Counseling Completed - Integration", func(t *testing.T) {
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		moveTaskOrderResponse := response.(*movetaskorderops.UpdateMTOStatusServiceCounselingCompletedOK)
		moveTaskOrderPayload := moveTaskOrderResponse.Payload

		suite.Assertions.IsType(&move_task_order.UpdateMTOStatusServiceCounselingCompletedOK{}, response)
		suite.Equal(strfmt.UUID(moveTaskOrder.ID.String()), moveTaskOrderPayload.ID)
		suite.Nil(moveTaskOrderPayload.ServiceCounselingCompletedAt)
		suite.EqualValues(models.MoveStatusServiceCounselingCompleted, moveTaskOrderPayload.Status)
	})

	suite.T().Run("Unsuccessful move status update to Service Counseling Completed, not found - Integration", func(t *testing.T) {
		params = move_task_order.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: uuid.FromStringOrNil("").String(),
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&move_task_order.UpdateMTOStatusServiceCounselingCompletedNotFound{}, response)
	})

	suite.T().Run("Unsuccessful move status update to Service Counseling Completed, eTag does not match - Integration", func(t *testing.T) {
		moveTaskOrder = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: order,
		})
		params = move_task_order.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: moveTaskOrder.ID.String(),
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&move_task_order.UpdateMTOStatusServiceCounselingCompletedPreconditionFailed{}, response)
	})

	suite.T().Run("Unsuccessful move status update to Service Counseling Completed, state conflict - Integration", func(t *testing.T) {
		moveTaskOrder = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusDRAFT,
			},
			Order: order,
		})

		params = move_task_order.UpdateMTOStatusServiceCounselingCompletedParams{
			HTTPRequest:     request,
			MoveTaskOrderID: moveTaskOrder.ID.String(),
			IfMatch:         etag.GenerateEtag(moveTaskOrder.UpdatedAt),
		}
		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&move_task_order.UpdateMTOStatusServiceCounselingCompletedConflict{}, response)
	})
}
