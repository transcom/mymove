package supportapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestPatchMTOShipmentHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", mtoShipment.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoshipmentops.PatchMTOShipmentStatusParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		Body:          &supportmessages.PatchMTOShipmentStatus{Status: "APPROVED"},
		IfMatch:       eTag,
	}

	// Used for all tests except 500 error:
	queryBuilder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(queryBuilder)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
	updater := mtoshipment.NewMTOShipmentStatusUpdater(suite.DB(), queryBuilder, siCreator, route.NewTestingPlanner(500))
	handler := PatchMTOShipmentStatusHandlerFunc{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		fetcher,
		updater,
	}

	suite.T().Run("Patch failure - 500", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		mockHandler := PatchMTOShipmentStatusHandlerFunc{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusInternalServerError{}, response)
	})

	suite.T().Run("Patch failure - 404", func(t *testing.T) {
		notFoundParams := mtoshipmentops.PatchMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()),
			Body:          params.Body,
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(notFoundParams)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusNotFound{}, response)
	})

	suite.T().Run("Patch failure - 412", func(t *testing.T) {
		preconditionParams := mtoshipmentops.PatchMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          params.Body,
			IfMatch:       "eTag",
		}
		response := handler.Handle(preconditionParams)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusPreconditionFailed{}, response)
	})

	suite.T().Run("Patch failure - 422", func(t *testing.T) {
		invalidInputParams := mtoshipmentops.PatchMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &supportmessages.PatchMTOShipmentStatus{Status: "X"},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(invalidInputParams)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusUnprocessableEntity{}, response)
	})

	// Second to last because many of the above tests fail because of a conflict error with APPROVED/REJECTED shipments
	// first:
	suite.T().Run("Successful patch - Integration Test", func(t *testing.T) {
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.PatchMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
	})

	// Last because the shipment has to be either APPROVED or REJECTED before triggering this conflict:
	suite.T().Run("Patch failure - 409", func(t *testing.T) {
		conflictParams := mtoshipmentops.PatchMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &supportmessages.PatchMTOShipmentStatus{Status: "SUBMITTED"},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(conflictParams)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusConflict{}, response)
	})
}
