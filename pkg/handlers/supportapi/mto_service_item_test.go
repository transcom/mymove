package supportapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/gen/supportmessages"

	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerApproveSuccess() {
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusSUBMITTED}})

	request := httptest.NewRequest("PATCH", "/service-items/{mtoServiceItemID}/status", nil)
	reason := "should not update reason"
	mtoServiceItem.Status = models.MTOServiceItemStatusApproved
	mtoServiceItem.Reason = &reason
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             payloads.MTOServiceItem(&mtoServiceItem),
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())

	// make the request
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
	suite.Equal(mtoServiceItemPayload.Status, supportmessages.MTOServiceItemStatusAPPROVED)
	suite.NotEqual(mtoServiceItemPayload.RejectionReason, reason)
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectSuccess() {
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusSUBMITTED}})

	request := httptest.NewRequest("PATCH", "/service-items/{mtoServiceItemID}/status", nil)
	reason := "item too heavy"
	mtoServiceItem.Status = models.MTOServiceItemStatusRejected
	mtoServiceItem.RejectionReason = &reason
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             payloads.MTOServiceItem(&mtoServiceItem),
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())

	// make the request
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
	suite.Equal(mtoServiceItemPayload.Status, supportmessages.MTOServiceItemStatusREJECTED)
	suite.Equal(*mtoServiceItemPayload.RejectionReason, reason)
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectionFailedNoReason() {
	mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

	request := httptest.NewRequest("PATCH", "/service-items/{mtoServiceItemID}/status", nil)
	mtoServiceItem.Status = models.MTOServiceItemStatusRejected
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             payloads.MTOServiceItem(&mtoServiceItem),
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())

	// make the request
	handler := UpdateMTOServiceItemStatusHandler{context,
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}
	response := handler.Handle(params)

	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusConflict)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusConflict{}, mtoServiceItemResponse)
	suite.Assertions.IsType(mtoServiceItemPayload, &supportmessages.ClientError{})
}
