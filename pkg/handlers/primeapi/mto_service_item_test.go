package primeapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMTOServiceItemHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})
	builder := query.NewQueryBuilder(suite.DB())

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s/mto_service_items", mto.ID.String(), mtoShipment.ID.String()), nil)

	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrderID:  mto.ID,
		MTOShipmentID:    &mtoShipment.ID,
		ReService:        models.ReService{Code: models.ReServiceCodeCS},
		Reason:           nil,
		PickupPostalCode: nil,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	params := mtoserviceitemops.CreateMTOServiceItemParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
		MtoShipmentID:   *handlers.FmtUUID(mtoShipment.ID),
		Body:            payloads.MTOServiceItem(&mtoServiceItem),
	}

	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload.ID())
	})

	//suite.T().Run("PUT failure - 500", func(t *testing.T) {
	//	mockUpdater := mocks.MTOShipmentUpdater{}
	//	handler := UpdateMTOShipmentHandler{
	//		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	//		&mockUpdater,
	//	}
	//	internalServerErr := errors.New("ServerError")
	//
	//	mockUpdater.On("UpdateMTOShipment",
	//		mock.Anything,
	//		mock.Anything,
	//	).Return(nil, internalServerErr)
	//
	//	response := handler.Handle(params)
	//	suite.IsType(&mtoserviceitemops.UpdateMTOShipmentInternalServerError{}, response)
	//})
	//
	//suite.T().Run("PUT failure - 400", func(t *testing.T) {
	//	mockUpdater := mocks.MTOShipmentUpdater{}
	//	handler := UpdateMTOShipmentHandler{
	//		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	//		&mockUpdater,
	//	}
	//
	//	mockUpdater.On("UpdateMTOShipment",
	//		mock.Anything,
	//		mock.Anything,
	//	).Return(nil, services.NewInvalidInputError(mtoShipment.ID, nil, nil, "invalid input"))
	//
	//	response := handler.Handle(params)
	//	suite.IsType(&mtoserviceitemops.UpdateMTOShipmentBadRequest{}, response)
	//})
	//
	//suite.T().Run("PUT failure - 404", func(t *testing.T) {
	//	mockUpdater := mocks.MTOShipmentUpdater{}
	//	handler := UpdateMTOShipmentHandler{
	//		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	//		&mockUpdater,
	//	}
	//
	//	mockUpdater.On("UpdateMTOShipment",
	//		mock.Anything,
	//		mock.Anything,
	//	).Return(nil, services.NotFoundError{})
	//
	//	response := handler.Handle(params)
	//	suite.IsType(&mtoserviceitemops.UpdateMTOShipmentNotFound{}, response)
	//})
	//
	//suite.T().Run("PUT failure - 412", func(t *testing.T) {
	//	mockUpdater := mocks.MTOShipmentUpdater{}
	//	handler := UpdateMTOShipmentHandler{
	//		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	//		&mockUpdater,
	//	}
	//
	//	mockUpdater.On("UpdateMTOShipment",
	//		mock.Anything,
	//		mock.Anything,
	//	).Return(nil, services.PreconditionFailedError{})
	//
	//	response := handler.Handle(params)
	//	suite.IsType(&mtoserviceitemops.UpdateMTOShipmentPreconditionFailed{}, response)
	//})
	//
	//mto2 := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	//mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
	//	MoveTaskOrder: mto,
	//})
	//
	//payload := primemessages.MTOShipment{
	//	ID:              strfmt.UUID(mtoShipment2.ID.String()),
	//	MoveTaskOrderID: strfmt.UUID(mtoShipment2.MoveTaskOrderID.String()),
	//}
	//
	//req2 := httptest.NewRequest("PUT", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", mto2.ID.String(), mtoShipment2.ID.String()), nil)
	//
	//eTag = etag.GenerateEtag(mtoShipment2.UpdatedAt)
	//params = mtoserviceitemops.UpdateMTOShipmentParams{
	//	HTTPRequest:     req2,
	//	MoveTaskOrderID: *handlers.FmtUUID(mtoShipment2.MoveTaskOrderID),
	//	MtoShipmentID:   *handlers.FmtUUID(mtoShipment2.ID),
	//	Body:            &payload,
	//	IfMatch:         eTag,
	//}
	//
	//suite.T().Run("Successful PUT - Integration Test with Only Required Fields in Payload", func(t *testing.T) {
	//	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher)
	//	handler := UpdateMTOShipmentHandler{
	//		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	//		updater,
	//	}
	//
	//	response := handler.Handle(params)
	//	suite.IsType(&mtoserviceitemops.UpdateMTOShipmentOK{}, response)
	//
	//	okResponse := response.(*mtoserviceitemops.UpdateMTOShipmentOK)
	//	suite.Equal(mtoShipment2.ID.String(), okResponse.Payload.ID.String())
	//})
}
