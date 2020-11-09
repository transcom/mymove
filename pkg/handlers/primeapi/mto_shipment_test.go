package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

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
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	mto := testdatagen.MakeAvailableMove(suite.DB())
	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: models.MTOShipment{},
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
			RequestedPickupDate: handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
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

		mtoNotAvailable := testdatagen.MakeDefaultMove(suite.DB())
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
	primeActualWeight := unit.Pound(600)
	move := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipment.PrimeEstimatedWeight = &primeEstimatedWeight
	mtoShipment.PrimeEstimatedWeightRecordedDate = &primeEstimatedWeightDate

	mtoNotAvailable := testdatagen.MakeDefaultMove(suite.DB())
	mtoShipmentNotAvailable := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mtoNotAvailable,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	minimalMtoShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:              models.MTOShipmentStatusDraft,
			ScheduledPickupDate: &time.Time{},
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
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)
	handler := UpdateMTOShipmentHandler{
		context,
		updater,
	}

	now := time.Now()

	suite.T().Run("PUT failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		mockHandler := UpdateMTOShipmentHandler{
			context,
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

	suite.T().Run("Successful PUT - Update weights on minimal shipment estimated weights", func(t *testing.T) {
		minimalMtoShipment.PrimeEstimatedWeight = &primeEstimatedWeight
		minimalMtoShipment.PrimeActualWeight = &primeActualWeight

		req = httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s", minimalMtoShipment.ID.String()), nil)
		eTag = etag.GenerateEtag(minimalMtoShipment.UpdatedAt)
		params = mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(minimalMtoShipment.ID),
			Body:          ClearNonUpdateFields(&minimalMtoShipment),
			IfMatch:       eTag,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(minimalMtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(minimalMtoShipment.PrimeActualWeight.Int64(), okResponse.Payload.PrimeActualWeight)
		suite.Equal(minimalMtoShipment.PrimeEstimatedWeight.Int64(), okResponse.Payload.PrimeEstimatedWeight)
	})

	suite.T().Run("PUT Failure (422) - Cannot update weight without scheduledPickupDate", func(t *testing.T) {
		noScheduledPickupShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              models.MTOShipmentStatusSubmitted,
				ScheduledPickupDate: nil,
			},
		})

		noScheduledPickupShipment.PrimeEstimatedWeight = &primeEstimatedWeight
		noScheduledPickupShipment.PrimeActualWeight = &primeActualWeight

		noPickupReq := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s", noScheduledPickupShipment.ID.String()), nil)
		noPickupETag := etag.GenerateEtag(noScheduledPickupShipment.UpdatedAt)
		noPickupParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   noPickupReq,
			MtoShipmentID: *handlers.FmtUUID(noScheduledPickupShipment.ID),
			Body:          ClearNonUpdateFields(&noScheduledPickupShipment),
			IfMatch:       noPickupETag,
		}

		response := handler.Handle(noPickupParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity)
		suite.NotEmpty(errResponse.Payload.InvalidFields)
		suite.Contains(errResponse.Payload.InvalidFields, "primeEstimatedWeight")
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
		Move: move,
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
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment2.ID.String(), okResponse.Payload.ID.String())
	})
	//}

	// Prime-specific validations tested below
	suite.T().Run("Failed case if not both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		eightDaysFromNow := now.AddDate(0, 0, 8)
		threeDaysBefore := now.AddDate(0, 0, -3)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &eightDaysFromNow,
				ApprovedDate:        &threeDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}
		params = mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req2,
			MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("Successful case if both approved date and estimated weight recorded date is more than ten days prior to scheduled move date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &tenDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(oldShipment.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("Successful case if scheduled pickup is changed. RequiredDeliveryDate should be generated.", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:       "APPROVED",
				ApprovedDate: &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  strfmt.Date(tenDaysFromNow),
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)

		// Let's double check our maths.
		expectedRDD := time.Time(responsePayload.ScheduledPickupDate).AddDate(0, 0, 12)
		actualRDD := time.Time(responsePayload.RequiredDeliveryDate)
		suite.Equal(expectedRDD.Year(), actualRDD.Year())
		suite.Equal(expectedRDD.Month(), actualRDD.Month())
		suite.Equal(expectedRDD.Day(), actualRDD.Day())

	})

	suite.T().Run("Successful case if in Alaska, should add an extra 10 days to required delivery date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		akAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "12345",
				State:      "AK",
			},
		})
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:               "APPROVED",
				ApprovedDate:         &now,
				DestinationAddress:   &akAddress,
				DestinationAddressID: &akAddress.ID,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payloadAKAddress := primemessages.Address{
			City:           &akAddress.City,
			Country:        akAddress.Country,
			ETag:           eTag,
			ID:             strfmt.UUID(akAddress.ID.String()),
			PostalCode:     &akAddress.PostalCode,
			State:          &akAddress.State,
			StreetAddress1: &akAddress.StreetAddress1,
			StreetAddress2: akAddress.StreetAddress2,
			StreetAddress3: akAddress.StreetAddress3,
		}
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  strfmt.Date(tenDaysFromNow),
			DestinationAddress:   &payloadAKAddress,
		}
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)

		// Let's double check our maths.
		expectedRDD := time.Time(responsePayload.ScheduledPickupDate).AddDate(0, 0, 22)
		actualRDD := time.Time(responsePayload.RequiredDeliveryDate)
		suite.Equal(expectedRDD.Year(), actualRDD.Year())
		suite.Equal(expectedRDD.Month(), actualRDD.Month())
		suite.Equal(expectedRDD.Day(), actualRDD.Day())
	})

	suite.T().Run("Successful case in Adak, Alaska, should add 20 days to required delivery date", func(t *testing.T) {
		tenDaysFromNow := now.AddDate(0, 0, 11)
		adakAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "99546",
				State:      "AK",
				City:       "Adak",
			},
		})
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:               "APPROVED",
				ApprovedDate:         &now,
				DestinationAddress:   &adakAddress,
				DestinationAddressID: &adakAddress.ID,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		payloadAdakAddress := primemessages.Address{
			City:           &adakAddress.City,
			Country:        adakAddress.Country,
			ETag:           eTag,
			ID:             strfmt.UUID(adakAddress.ID.String()),
			PostalCode:     &adakAddress.PostalCode,
			State:          &adakAddress.State,
			StreetAddress1: &adakAddress.StreetAddress1,
			StreetAddress2: adakAddress.StreetAddress2,
			StreetAddress3: adakAddress.StreetAddress3,
		}

		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
			ScheduledPickupDate:  strfmt.Date(tenDaysFromNow),
			DestinationAddress:   &payloadAdakAddress,
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)

		// Let's double check our maths.
		expectedRDD := time.Time(responsePayload.ScheduledPickupDate).AddDate(0, 0, 32)
		actualRDD := time.Time(responsePayload.RequiredDeliveryDate)
		suite.Equal(expectedRDD.Year(), actualRDD.Year())
		suite.Equal(expectedRDD.Month(), actualRDD.Month())
		suite.Equal(expectedRDD.Day(), actualRDD.Day())

	})

	suite.T().Run("Failed case if approved date is 3-9 days from scheduled move date but estimated weight recorded date isn't at least 3 days prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("Successful case if approved date is 3-9 days from scheduled move date and estimated weight recorded date is at least 3 days prior to scheduled move date", func(t *testing.T) {
		sixDaysFromNow := now.AddDate(0, 0, 6)
		twoDaysBefore := now.AddDate(0, 0, -2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &sixDaysFromNow,
				ApprovedDate:        &twoDaysBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
	})

	suite.T().Run("Failed case if approved date is less than 3 days from scheduled move date but estimated weight recorded date isn't at least 1 day prior to scheduled move date", func(t *testing.T) {
		oneDayFromNow := now.AddDate(0, 0, 1)
		oneDayBefore := now.AddDate(0, 0, -1)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &oneDayFromNow,
				ApprovedDate:        &oneDayBefore,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("Successful case if approved date is less than 3 days from scheduled move date and estimated weight recorded date is at least 1 day prior to scheduled move date", func(t *testing.T) {
		twoDaysFromNow := now.AddDate(0, 0, 2)
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:              "APPROVED",
				ScheduledPickupDate: &twoDaysFromNow,
				ApprovedDate:        &now,
			},
		})
		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
		payload := primemessages.MTOShipment{
			ID:                   strfmt.UUID(oldShipment.ID.String()),
			PrimeEstimatedWeight: int64(primeEstimatedWeight),
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.NotNil(responsePayload.RequiredDeliveryDate)
	})

	suite.T().Run("Successful case for valid and complete payload including approved date and re service code", func(t *testing.T) {
		reServiceID, _ := uuid.NewV4()
		mto := testdatagen.MakeAvailableMove(suite.DB())

		reService := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				ID:   reServiceID,
				Code: models.ReServiceCodeDDFSIT,
			},
		})

		now := time.Now()
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: mto.ID,
				Status:          "APPROVED",
				ApprovedDate:    &now,
			},
		})

		mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:        mto,
			MTOShipment: oldShipment,
			ReService:   reService,
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: mto.ID,
				ReServiceID:     reService.ID,
				MTOShipmentID:   &oldShipment.ID,
			},
		})
		serviceItems := models.MTOServiceItems{mtoServiceItem1}
		oldShipment.MTOServiceItems = serviceItems

		eTag := etag.GenerateEtag(oldShipment.UpdatedAt)

		payload := primemessages.MTOShipment{
			ID:             strfmt.UUID(oldShipment.ID.String()),
			PointOfContact: "John McRand",
		}

		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto_shipments/%s", oldShipment.ID.String()), nil)

		params := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
			Body:          &payload,
			IfMatch:       eTag,
		}

		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			context,
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		responsePayload := okResponse.Payload

		suite.Equal(1, len(responsePayload.MtoServiceItems()))
		var serviceItemDDFSIT *primemessages.MTOServiceItemDDFSIT
		var serviceItemDDFSITCode string

		for _, item := range responsePayload.MtoServiceItems() {
			if item.ModelType() == primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT {
				serviceItemDDFSIT = item.(*primemessages.MTOServiceItemDDFSIT)
				serviceItemDDFSITCode = *serviceItemDDFSIT.ReServiceCode
				break
			}

		}
		suite.Equal(oldShipment.ID.String(), responsePayload.ID.String())
		suite.Equal(string(reService.Code), serviceItemDDFSITCode)
		suite.Equal(oldShipment.ApprovedDate, &now)
	})
}
