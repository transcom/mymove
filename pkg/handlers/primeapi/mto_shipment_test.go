package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", mto.ID.String(), mtoShipment.ID.String()), nil)

	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
		MtoShipmentID:   *handlers.FmtUUID(mtoShipment.ID),
		Body:            payloads.MTOShipment(&mtoShipment),
		IfMatch:         eTag,
	}

	suite.T().Run("Successful PUT - Integration Test", func(t *testing.T) {
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("PUT failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}
		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
	})

	suite.T().Run("PUT failure - 400", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NewInvalidInputError(mtoShipment.ID, nil, nil, "invalid input"))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("PUT failure - 404", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PUT failure - 412", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, services.PreconditionFailedError{})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	mto2 := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	payload := primemessages.MTOShipment{
		ID:              strfmt.UUID(mtoShipment2.ID.String()),
		MoveTaskOrderID: strfmt.UUID(mtoShipment2.MoveTaskOrderID.String()),
	}

	req2 := httptest.NewRequest("PUT", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", mto2.ID.String(), mtoShipment2.ID.String()), nil)

	eTag = etag.GenerateEtag(mtoShipment2.UpdatedAt)
	params = mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:     req2,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment2.MoveTaskOrderID),
		MtoShipmentID:   *handlers.FmtUUID(mtoShipment2.ID),
		Body:            &payload,
		IfMatch:         eTag,
	}

	suite.T().Run("Successful PUT - Integration Test with Only Required Fields in Payload", func(t *testing.T) {
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment2.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("Successful PUT - Integration Test with updating the releasing/receiving agents", func(t *testing.T) {
		mtoAgents := make(primemessages.MTOAgents, 2)
		newFirstName := "NewTestName"
		newLastName := "NewLastName"
		mtoAgents[0].FirstName = &newFirstName
		mtoAgents[1].LastName = &newLastName

		agentPayload := primemessages.MTOShipment{
			ID:              strfmt.UUID(mtoShipment2.ID.String()),
			MoveTaskOrderID: strfmt.UUID(mtoShipment2.MoveTaskOrderID.String()),
			Agents:          mtoAgents,
		}
		agentParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:     req2,
			MoveTaskOrderID: *handlers.FmtUUID(mtoShipment2.MoveTaskOrderID),
			MtoShipmentID:   *handlers.FmtUUID(mtoShipment2.ID),
			Body:            &agentPayload,
			IfMatch:         eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(agentParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(newFirstName, okResponse.Payload.Agents[0].FirstName)
		suite.Equal(newLastName, okResponse.Payload.Agents[1].LastName)
	})
}
