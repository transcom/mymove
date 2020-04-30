package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"

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

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			PickupAddress:        &pickupAddress,
			PickupAddressID:      &pickupAddress.ID,
			DestinationAddress:   &destinationAddress,
			DestinationAddressID: &destinationAddress.ID,
		},
	})

	mtoShipment.MoveTaskOrderID = mto.ID

	builder := query.NewQueryBuilder(suite.DB())

	req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/%s/mto_shipments", mto.ID.String()), nil)

	params := mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
		Body: &primemessages.CreateShipmentPayload{
			Agents:          nil,
			CustomerRemarks: mtoShipment.CustomerRemarks,
			DestinationAddress: &primemessages.Address{
				City:    &mtoShipment.DestinationAddress.City,
				Country: mtoShipment.DestinationAddress.Country,
				//ID:             strfmt.UUID(mtoShipment.DestinationAddress.ID.String()),
				PostalCode:     &mtoShipment.DestinationAddress.PostalCode,
				State:          &mtoShipment.DestinationAddress.State,
				StreetAddress1: &mtoShipment.DestinationAddress.StreetAddress1,
				StreetAddress2: mtoShipment.DestinationAddress.StreetAddress2,
				StreetAddress3: mtoShipment.DestinationAddress.StreetAddress3,
			},
			PickupAddress: &primemessages.Address{
				City:    &mtoShipment.PickupAddress.City,
				Country: mtoShipment.PickupAddress.Country,
				//ID:             strfmt.UUID(mtoShipment.PickupAddress.ID.String()),
				PostalCode:     &mtoShipment.PickupAddress.PostalCode,
				State:          &mtoShipment.PickupAddress.State,
				StreetAddress1: &mtoShipment.PickupAddress.StreetAddress1,
				StreetAddress2: mtoShipment.PickupAddress.StreetAddress2,
				StreetAddress3: mtoShipment.PickupAddress.StreetAddress3,
			},
			PointOfContact:      "John Doe",
			RequestedPickupDate: strfmt.Date(*mtoShipment.RequestedPickupDate),
			ShipmentType:        primemessages.MTOShipmentTypeHHG,
		},
	}
	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)
		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		badParams := params
		badParams.Body.PickupAddress = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
	})

	suite.T().Run("POST failure - 404 -- not found", func(t *testing.T) {

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.MoveTaskOrderID = strfmt.UUID(uuidString)

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("POST failure - 400 -- nil body", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		req := httptest.NewRequest("POST", fmt.Sprintf("/move_task_orders/{MoveTaskOrderID}/mto_shipments"), nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	primeEstimatedWeight := unit.Pound(500)
	primeEstimatedWeightDate := testdatagen.DateInsidePeakRateCycle
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	mtoShipment.PrimeEstimatedWeight = &primeEstimatedWeight
	mtoShipment.PrimeEstimatedWeightRecordedDate = &primeEstimatedWeightDate

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

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
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, route.NewTestingPlanner(400))
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
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, route.NewTestingPlanner(400))
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment2.ID.String(), okResponse.Payload.ID.String())
	})
}
