package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/swag"

	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"

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
	mto := MakeAvailableMoveTaskOrder(suite.DB())
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment:   models.MTOShipment{},
	})

	mtoShipment.MoveTaskOrderID = mto.ID

	builder := query.NewQueryBuilder(suite.DB())
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

	req := httptest.NewRequest("POST", "/mto-shipments", nil)

	params := mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &primemessages.CreateMTOShipment{
			MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:          nil,
			CustomerRemarks: nil,
			DestinationAddress: &primemessages.Address{
				City:           &destinationAddress.City,
				Country:        destinationAddress.Country,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
			},
			PickupAddress: &primemessages.Address{
				City:           &pickupAddress.City,
				Country:        pickupAddress.Country,
				PostalCode:     &pickupAddress.PostalCode,
				State:          &pickupAddress.State,
				StreetAddress1: &pickupAddress.StreetAddress1,
				StreetAddress2: pickupAddress.StreetAddress2,
				StreetAddress3: pickupAddress.StreetAddress3,
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
			mtoChecker,
		}
		response := handler.Handle(params)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		// check that the mto shipment status is Submitted
		suite.Require().Equal(createMTOShipmentPayload.Status, primemessages.MTOShipmentStatusSUBMITTED, "MTO Shipment should have been submitted")
	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {
		mockCreator := mocks.MTOShipmentCreator{}

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}

		err := errors.New("ServerError")

		mockCreator.On("CreateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.T().Run("POST failure - 422 -- Bad agent IDs set on shipment", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		badID := params.Body.MoveTaskOrderID
		agent := &primemessages.MTOAgent{
			ID:            *badID,
			MtoShipmentID: *badID,
			FirstName:     handlers.FmtString("Mary"),
		}

		paramsBadIDs := params
		paramsBadIDs.Body.Agents = primemessages.MTOAgents{agent}

		response := handler.Handle(paramsBadIDs)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)
		suite.NotEmpty(typedResponse.Payload.InvalidFields)
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:id")
		suite.Contains(typedResponse.Payload.InvalidFields, "agents:mtoShipmentID")
	})

	suite.T().Run("POST failure - 422 - invalid input, missing pickup address", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		badParams := params
		badParams.Body.PickupAddress = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 404 -- not found", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("POST failure - 400 -- nil body", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(paramsNilBody)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("POST failure - 404 -- MTO is not available to Prime", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		mtoNotAvailable := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
		mtoIDNotAvailable := strfmt.UUID(mtoNotAvailable.ID.String())

		paramsNotAvailable := params
		paramsNotAvailable.Body.MoveTaskOrderID = &mtoIDNotAvailable

		response := handler.Handle(paramsNotAvailable)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)

		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentNotFound)
		suite.Contains(*typedResponse.Payload.Detail, mtoNotAvailable.ID.String())
	})

	suite.T().Run("POST failure - 422 - modelType() not supported", func(t *testing.T) {
		mockCreator := mocks.MTOShipmentCreator{}

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}
		err := services.NotFoundError{}

		mockCreator.On("CreateMTOShipment",
			mock.Anything,
		).Return(nil, nil, err)

		mtoServiceItems := models.MTOServiceItems{
			models.MTOServiceItem{
				MoveTaskOrderID:  mto.ID,
				MTOShipmentID:    &mtoShipment.ID,
				ReService:        models.ReService{Code: models.ReServiceCodeMS},
				Reason:           nil,
				PickupPostalCode: nil,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}

		badModelTypeParams := params
		badModelTypeParams.Body.SetMtoServiceItems(*payloads.MTOServiceItems(&mtoServiceItems))
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})
}

// ClearNonUpdateFields clears out the MTOShipment payload fields that CANNOT be sent in for a successful update
func ClearNonUpdateFields(mtoShipment *models.MTOShipment) *primemessages.MTOShipment {
	mtoShipment.MoveTaskOrderID = uuid.FromStringOrNil("")
	mtoShipment.CreatedAt = time.Time{}
	mtoShipment.UpdatedAt = time.Time{}
	mtoShipment.PrimeEstimatedWeightRecordedDate = &time.Time{}
	mtoShipment.RequiredDeliveryDate = &time.Time{}
	mtoShipment.ApprovedDate = &time.Time{}
	mtoShipment.Status = ""
	mtoShipment.RejectionReason = nil
	mtoShipment.CustomerRemarks = nil
	mtoShipment.MTOAgents = nil

	return payloads.MTOShipment(mtoShipment)
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	primeEstimatedWeight := unit.Pound(500)
	primeEstimatedWeightDate := testdatagen.DateInsidePeakRateCycle
	mto := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.Move{
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	mtoShipment.PrimeEstimatedWeight = &primeEstimatedWeight
	mtoShipment.PrimeEstimatedWeightRecordedDate = &primeEstimatedWeightDate

	mtoNotAvailable := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipmentNotAvailable := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mtoNotAvailable,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

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

	req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s", mtoShipment.ID.String()), nil)
	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		Body:          ClearNonUpdateFields(&mtoShipment),
		IfMatch:       eTag,
	}
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	// used for all tests except the 500 server error:
	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
	handler := UpdateMTOShipmentHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		updater,
	}

	suite.T().Run("Successful PUT - Integration Test", func(t *testing.T) {
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("PUT failure - Shipment is not part of MTO available to prime", func(t *testing.T) {
		notAvailableShipment := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: *handlers.FmtUUID(mtoShipmentNotAvailable.ID),
			Body:          ClearNonUpdateFields(&mtoShipmentNotAvailable),
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(notAvailableShipment)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PUT failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		mockHandler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}
		internalServerErr := errors.New("ServerError")

		mockUpdater.On("MTOShipmentsMTOAvailableToPrime",
			mock.Anything,
		).Return(true, nil)

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")

	})

	suite.T().Run("PUT failure - 404", func(t *testing.T) {
		notFoundParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()),
			Body:          &primemessages.MTOShipment{},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(notFoundParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PUT failure - 422 (extra fields)", func(t *testing.T) {
		remarks := fmt.Sprintf("test conflict %s", time.Now())
		conflictParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &primemessages.MTOShipment{CustomerRemarks: &remarks},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(conflictParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("PUT failure - 412", func(t *testing.T) {
		staleParams := params
		// no need to update IfMatch or eTag value because it's still the old value from before the successful PUT
		staleParams.Body.PrimeEstimatedWeight = 0 // causes this test to fail because any updates to this are invalid
		// (and input validation happens first before this check)
		response := handler.Handle(staleParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.T().Run("PUT failure - 422", func(t *testing.T) {
		params.Body.PrimeEstimatedWeight = 1 // cannot update once initial value has been set
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
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
		ID: strfmt.UUID(mtoShipment2.ID.String()),
	}

	req2 := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", mtoShipment2.ID.String()), nil)

	eTag = etag.GenerateEtag(mtoShipment2.UpdatedAt)
	params = mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req2,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
		Body:          &payload,
		IfMatch:       eTag,
	}

	suite.T().Run("Successful PUT - Integration Test with Only Required Fields in Payload", func(t *testing.T) {
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
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
