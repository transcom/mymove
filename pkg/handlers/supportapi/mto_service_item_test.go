package supportapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/supportmessages"

	"github.com/transcom/mymove/pkg/models"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerApproveSuccess() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to approve it.
	// Expected outcome:
	//             Success, MTO service item is approved

	// SETUP
	// Create a service item on a move
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVALSREQUESTED,
		},
	})

	// Create a request to the endpoint
	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)

	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status:          supportmessages.MTOServiceItemStatusAPPROVED,
		RejectionReason: swag.String("Should not update the reason"),
	}
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload
	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)

	// Check the status is APPROVED
	suite.Equal(supportmessages.MTOServiceItemStatusAPPROVED, mtoServiceItemPayload.Status())
	// Check that reason was NOT set
	suite.NotEqual(requestPayload.RejectionReason, mtoServiceItemPayload.RejectionReason())
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectSuccess() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to reject it.
	// Expected outcome:
	//             Success, MTO service item is rejected, rejectionReason is populated

	// SETUP
	// Create a service item on a move
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVALSREQUESTED,
		},
	})

	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)
	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status:          supportmessages.MTOServiceItemStatusREJECTED,
		RejectionReason: swag.String("Should definitely update the reason"),
	}
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
	suite.Equal(supportmessages.MTOServiceItemStatusREJECTED, mtoServiceItemPayload.Status())
	suite.Equal(requestPayload.RejectionReason, mtoServiceItemPayload.RejectionReason())
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectionFailedNoReason() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to reject it, but fail
	//             to send a rejectionReason
	// Expected outcome:
	//             Fail, RejectionReason must be provided to service item is not updated.

	// SETUP
	// Create a service item on a move
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVALSREQUESTED,
		},
	})

	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)
	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status:          supportmessages.MTOServiceItemStatusREJECTED,
		RejectionReason: nil,
	}

	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusConflict)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusConflict{}, mtoServiceItemResponse)
	suite.Assertions.IsType(&supportmessages.ClientError{}, mtoServiceItemPayload)

	// Check that the status in DB is still SUBMITTED, not APPROVED or REJECTED
	serviceItemInDB := models.MTOServiceItem{}
	suite.DB().Find(&serviceItemInDB, mtoServiceItem.ID)
	suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItemInDB.Status)

}
