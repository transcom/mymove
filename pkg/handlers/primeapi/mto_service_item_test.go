package primeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestCreateMTOServiceItemHandler() {
	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

	type localSubtestData struct {
		params         mtoserviceitemops.CreateMTOServiceItemParams
		mtoShipment    models.MTOShipment
		mtoServiceItem models.MTOServiceItem
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}

		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		subtestData.mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		sitEntryDate := time.Now()
		sitPostalCode := "00000"
		requestApprovalRequestedStatus := false

		// Customer gets new pickup address for SIT Origin Pickup (DOPSIT) which gets added when
		// creating DOFSIT (SIT origin first day).
		//
		// Do not create Address in the database (Assertions.Stub = true), because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID:                   mto.ID,
			MTOShipmentID:                     &subtestData.mtoShipment.ID,
			ReService:                         models.ReService{Code: models.ReServiceCodeDOFSIT},
			Reason:                            models.StringPointer("lorem ipsum"),
			SITEntryDate:                      &sitEntryDate,
			SITPostalCode:                     &sitPostalCode,
			SITOriginHHGActualAddress:         &actualPickupAddress,
			RequestedApprovalsRequestedStatus: &requestApprovalRequestedStatus,
		}

		subtestData.params = mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		return subtestData
	}

	suite.Run("Successful POST - Integration Test", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		// TODO: This is failing because DOPSIT and DDDSIT are being sent back in the response
		//   but those are not listed in the enum in the swagger file.  They aren't allowed for
		//   incoming payloads, but are allowed for outgoing payloads, but the same payload spec
		//   is used for both.  Need to figure out best way to resolve.
		// Validate outgoing payload (each element of slice)
		// for _, mtoServiceItem := range okResponse.Payload {
		// 	suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		// }

		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("Successful POST for Creating Shuttling without PrimeEstimatedWeight set - Integration Test", func() {
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		mtoShipment.PrimeEstimatedWeight = nil
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOSHUT)
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		reason := "lorem ipsum"

		mtoServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: mto.ID,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDOSHUT},
			Reason:          &reason,
		}

		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("POST failure - 500", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		err := fmt.Errorf("ServerError")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemInternalServerError{}, response)
		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemInternalServerError)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))

		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})

	suite.Run("POST failure - 422 Unprocessable Entity Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		// InvalidInputError should generate an UnprocessableEntity response
		// Need verrs incorporated to satisfy swagger validation
		verrs := validate.NewErrors()
		verrs.Add("some key", "some value")
		err := apperror.NewInvalidInputError(subtestData.mtoServiceItem.ID, nil, verrs, "some error")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 409 Conflict Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		// ConflictError should generate a Conflict response
		err := apperror.ConflictError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemConflict{}, response)
		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemConflict)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		err := apperror.NotFoundError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound)

		// Validate outgoing payload
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 404 - MTO is not available to Prime", func() {
		subtestData := makeSubtestData()
		mtoNotAvailable := factory.BuildMove(suite.DB(), nil, nil)
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		body := payloads.MTOServiceItem(&subtestData.mtoServiceItem)
		body.SetMoveTaskOrderID(handlers.FmtUUID(mtoNotAvailable.ID))

		paramsNotAvailable := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        body,
		}

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(paramsNotAvailable)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
		typedResponse := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.Contains(*typedResponse.Payload.Detail, mtoNotAvailable.ID.String())
	})

	suite.Run("POST failure - 404 - Integration - ShipmentID not linked by MoveTaskOrderID", func() {
		subtestData := makeSubtestData()
		mto2 := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto2,
				LinkOnly: true,
			},
		}, nil)
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		body := payloads.MTOServiceItem(&subtestData.mtoServiceItem)
		body.SetMoveTaskOrderID(handlers.FmtUUID(subtestData.mtoShipment.MoveTaskOrderID))
		body.SetMtoShipmentID(strfmt.UUID(mtoShipment2.ID.String()))

		newParams := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        body,
		}

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(newParams)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 422 - Model validation errors", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		verrs := validate.NewErrors()
		verrs.Add("test", "testing")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, verrs, nil)

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 422 - modelType() not supported", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		err := apperror.NotFoundError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		mtoServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mtoShipment.MoveTaskOrder.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeMS},
			Reason:          nil,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestCreateMTOServiceItemDomesticCratingHandler() {
	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

	type localSubtestData struct {
		req            *http.Request
		mtoServiceItem models.MTOServiceItem
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}

		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDCRT)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDUCRT)
		subtestData.req = httptest.NewRequest("POST", "/mto-service-items", nil)

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: mto.ID,
			MTOShipmentID:   &mtoShipment.ID,
			Description:     handlers.FmtString("description"),
			Dimensions: models.MTOServiceItemDimensions{
				models.MTOServiceItemDimension{
					Type:   models.DimensionTypeItem,
					Length: 1000,
					Height: 1000,
					Width:  1000,
				},
				models.MTOServiceItemDimension{
					Type:   models.DimensionTypeCrate,
					Length: 10000,
					Height: 10000,
					Width:  10000,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Reason:    handlers.FmtString("reason"),
		}
		return subtestData
	}

	suite.Run("Successful POST - Integration Test - Domestic Crating", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		// Validate outgoing payload (each element of slice)
		for _, mtoServiceItem := range okResponse.Payload {
			suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		}

		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("Successful POST - Integration Test - Domestic Uncrating", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDUCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		// Validate outgoing payload (each element of slice)
		for _, mtoServiceItem := range okResponse.Payload {
			suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		}

		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("POST failure - 422", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			&mockCreator,
			mtoChecker,
		}
		err := fmt.Errorf("ServerError")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDUCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		var height int32
		params.Body.(*primemessages.MTOServiceItemDomesticCrating).Crate.Height = &height

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestCreateMTOServiceItemOriginSITHandler() {
	// Under test: createMTOServiceItemHandler function,
	// - no DOPSIT standalone
	// -  DOASIT standalone with DOFSIT

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

	type localSubtestData struct {
		mto            models.Move
		mtoShipment    models.MTOShipment
		mtoServiceItem models.MTOServiceItem
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}

		subtestData.mto = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		subtestData.mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{},
			Reason:          models.StringPointer("lorem ipsum"),
			SITEntryDate:    &sitEntryDate,
			SITPostalCode:   &sitPostalCode,
		}

		return subtestData
	}

	suite.Run("POST failure - 422 Cannot create DOPSIT standalone", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a DOPSIT MTOServiceItem
		// Expected outcome:
		//             Receive a 422 - Unprocessable Entity
		// SETUP
		// Create the payload
		requestApprovalRequestedStatus := false
		subtestData.mtoServiceItem.RequestedApprovalsRequestedStatus = &requestApprovalRequestedStatus
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOPSIT
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.Error(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST Failure - Cannot create DOASIT without DOFSIT", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a standalone DOASIT MTOServiceItem, no DOFSIT
		// Expected outcome:
		//             Receive a 404 - Not Found
		// SETUP
		// Create the payload
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOASIT
		requestApprovalRequestedStatus := false
		subtestData.mtoServiceItem.RequestedApprovalsRequestedStatus = &requestApprovalRequestedStatus

		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Successful POST - Create DOASIT with DOFSIT", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a standalone DOASIT MTOServiceItem
		// Expected outcome:
		//             Receive a 404 - Not Found
		// SETUP
		// Create the payload
		requestedApprovalsRequestedStatus := false
		subtestData.mtoServiceItem.RequestedApprovalsRequestedStatus = &requestedApprovalsRequestedStatus
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model:    subtestData.mto,
				LinkOnly: true,
			},
			{
				Model:    subtestData.mtoShipment,
				LinkOnly: true,
			},
			// These get copied over to the DOASIT as part of creation and are needed for the response to validate
			{
				Model: models.MTOServiceItem{
					Reason:        models.StringPointer("lorem ipsum"),
					SITEntryDate:  models.TimePointer(time.Now()),
					SITPostalCode: models.StringPointer("00000"),
				},
			},
		}, nil)

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOASIT
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemOK).Payload

		// Validate outgoing payload (each element of slice)
		for _, mtoServiceItem := range responsePayload {
			suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		}
	})
}

func (suite *HandlerSuite) TestCreateMTOServiceItemOriginSITHandlerWithDOFSITNoAddress() {
	// Under test: createMTOServiceItemHandler function,
	// - fail to create DOFSIT because of missing sitHHGActualAddress

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

	type localSubtestData struct {
		mtoServiceItem models.MTOServiceItem
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: mto.ID,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       models.ReService{},
			Reason:          models.StringPointer("lorem ipsum"),
			SITEntryDate:    &sitEntryDate,
			SITPostalCode:   &sitPostalCode,
		}
		return subtestData
	}

	suite.Run("Failed POST - Does not DOFSIT with missing SitHHGActualOrigin", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a standalone DOFSIT MTOServiceItem
		// Expected outcome:
		//             CreateMTOServiceItemUnprocessableEntity
		// SETUP
		// Create the payload

		requstedApprovalsRequestedStatus := false
		subtestData.mtoServiceItem.RequestedApprovalsRequestedStatus = &requstedApprovalsRequestedStatus
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOFSIT
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		unprocessableEntity := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(unprocessableEntity.Payload.Validate(strfmt.Default))

		suite.Contains(*unprocessableEntity.Payload.Detail, "must have the sitHHGActualOrigin")
	})

}

func (suite *HandlerSuite) TestCreateMTOServiceItemOriginSITHandlerWithDOFSITWithAddress() {
	// Under test: createMTOServiceItemHandler function,
	// - no DOPSIT standalone
	// -  DOASIT standalone with DOFSIT

	type localSubtestData struct {
		mtoShipment             models.MTOShipment
		mtoServiceItem          models.MTOServiceItem
		actualPickupAddress     models.Address
		originalPickupAddress   *models.Address
		originalPickupAddressID *uuid.UUID
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		subtestData.mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)
		sitEntryDate := time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC)
		sitDepartureDate := time.Date(2024, time.February, 27, 0, 0, 0, 0, time.UTC)
		sitPostalCode := "00000"

		// Original customer pickup address
		subtestData.originalPickupAddress = subtestData.mtoShipment.PickupAddress
		subtestData.originalPickupAddressID = subtestData.mtoShipment.PickupAddressID

		// Customer gets new pickup address

		// Do not create the Address in the database (factory.BuildAddress(nil, nil, nil)), because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		subtestData.actualPickupAddress = factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID:           mto.ID,
			MTOShipmentID:             &subtestData.mtoShipment.ID,
			ReService:                 models.ReService{},
			Reason:                    models.StringPointer("lorem ipsum"),
			SITEntryDate:              &sitEntryDate,
			SITDepartureDate:          &sitDepartureDate,
			SITPostalCode:             &sitPostalCode,
			SITOriginHHGActualAddress: &subtestData.actualPickupAddress,
		}

		// Verify the addresses for original pickup and new pickup are not the same
		suite.NotEqual(subtestData.originalPickupAddressID, subtestData.mtoServiceItem.SITOriginHHGActualAddressID, "address ID is not the same")
		suite.NotEqual(subtestData.originalPickupAddress.StreetAddress1, subtestData.mtoServiceItem.SITOriginHHGActualAddress.StreetAddress1, "street address is not the same")
		suite.NotEqual(subtestData.originalPickupAddress.City, subtestData.mtoServiceItem.SITOriginHHGActualAddress.City, "city is not the same")
		suite.NotEqual(subtestData.originalPickupAddress.PostalCode, subtestData.mtoServiceItem.SITOriginHHGActualAddress.PostalCode, "zip is not the same")

		return subtestData
	}
	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

	suite.Run("Successful POST - Create DOFSIT", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a standalone DOFSIT MTOServiceItem
		// Expected outcome:
		//             Successful creation of DOFSIT with DOPSIT added
		// SETUP
		// Create the payload

		requestedApprovalsRequestedStatus := false
		subtestData.mtoServiceItem.RequestedApprovalsRequestedStatus = &requestedApprovalsRequestedStatus
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOFSIT
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		// TODO: This is failing because DOPSIT and DDDSIT are being sent back in the response
		//   but those are not listed in the enum in the swagger file.  They aren't allowed for
		//   incoming payloads, but are allowed for outgoing payloads, but the same payload spec
		//   is used for both.  Need to figure out best way to resolve.
		// Validate outgoing payload (each element of slice)
		// for _, mtoServiceItem := range okResponse.Payload {
		// 	suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		// }

		// Verify address was updated on MTO Shipment
		var updatedMTOShipment models.MTOShipment
		suite.NoError(suite.DB().Eager("PickupAddress").Find(&updatedMTOShipment, subtestData.mtoShipment.ID))

		// Verify the HHG pickup address is the actual address on the shipment
		suite.Equal(*subtestData.mtoShipment.PickupAddressID, *updatedMTOShipment.PickupAddressID, "hhg actual address id is the same")
		suite.Equal(subtestData.actualPickupAddress.StreetAddress1, updatedMTOShipment.PickupAddress.StreetAddress1, "hhg actual street address is the same")
		suite.Equal(subtestData.actualPickupAddress.City, updatedMTOShipment.PickupAddress.City, "hhg actual city is the same")
		suite.Equal(subtestData.actualPickupAddress.State, updatedMTOShipment.PickupAddress.State, "hhg actual state is the same")
		suite.Equal(subtestData.actualPickupAddress.PostalCode, updatedMTOShipment.PickupAddress.PostalCode, "hhg actual zip is the same")

		// Verify address on SIT service item
		suite.NotZero(okResponse.Payload[0].ID())

		foundDOFSIT := false
		foundDOPSIT := false
		foundDOASIT := false

		for _, serviceItem := range okResponse.Payload {

			// Find the matching MTO Service Item from the DB for the returned payload
			var mtosi models.MTOServiceItem
			id := serviceItem.ID()
			findServiceItemErr := suite.DB().Eager("ReService", "SITOriginHHGOriginalAddress", "SITOriginHHGActualAddress").Find(&mtosi, &id)
			suite.NoError(findServiceItemErr)

			if mtosi.ReService.Code == models.ReServiceCodeDOPSIT || mtosi.ReService.Code == models.ReServiceCodeDOFSIT || mtosi.ReService.Code == models.ReServiceCodeDOASIT {
				suite.IsType(&primemessages.MTOServiceItemOriginSIT{}, serviceItem)
				sitItem := serviceItem.(*primemessages.MTOServiceItemOriginSIT)

				if mtosi.ReService.Code == models.ReServiceCodeDOPSIT {
					foundDOPSIT = true
				} else if mtosi.ReService.Code == models.ReServiceCodeDOFSIT {
					foundDOFSIT = true
				} else if mtosi.ReService.Code == models.ReServiceCodeDOASIT {
					foundDOASIT = true
				}

				// Verify the return primemessages payload has the correct addresses
				suite.NotNil(sitItem.SitHHGActualOrigin, "primemessages SitHHGActualOrigin is not Nil")
				suite.NotEqual(uuid.Nil, sitItem.SitHHGActualOrigin.ID, "primemessages actual address ID is not nil")
				suite.Equal(updatedMTOShipment.PickupAddress.StreetAddress1, *sitItem.SitHHGActualOrigin.StreetAddress1, "primemessages actual street address is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.City, *sitItem.SitHHGActualOrigin.City, "primemessages actual city is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.State, *sitItem.SitHHGActualOrigin.State, "primemessages actual state is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.PostalCode, *sitItem.SitHHGActualOrigin.PostalCode, "primemessages actual zip is the same")

				// Verify the HHG original pickup address is the original address on the service item
				suite.NotNil(mtosi.SITOriginHHGOriginalAddressID, "original address ID is not nil")
				suite.NotEqual(uuid.Nil, *mtosi.SITOriginHHGOriginalAddressID)
				suite.Equal(subtestData.originalPickupAddress.StreetAddress1, mtosi.SITOriginHHGOriginalAddress.StreetAddress1, "original street address is the same")
				suite.Equal(subtestData.originalPickupAddress.City, mtosi.SITOriginHHGOriginalAddress.City, "original city is the same")
				suite.Equal(subtestData.originalPickupAddress.State, mtosi.SITOriginHHGOriginalAddress.State, "original state is the same")
				suite.Equal(subtestData.originalPickupAddress.PostalCode, mtosi.SITOriginHHGOriginalAddress.PostalCode, "original zip is the same")

				// Verify the HHG pickup address is the actual address on the service item
				suite.NotNil(mtosi.SITOriginHHGActualAddressID, "actual address ID is not nil")
				suite.NotEqual(uuid.Nil, *mtosi.SITOriginHHGActualAddressID)
				suite.Equal(updatedMTOShipment.PickupAddress.StreetAddress1, mtosi.SITOriginHHGActualAddress.StreetAddress1, "shipment actual street address is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.City, mtosi.SITOriginHHGActualAddress.City, "shipment actual city is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.State, mtosi.SITOriginHHGActualAddress.State, "shipment actual state is the same")
				suite.Equal(updatedMTOShipment.PickupAddress.PostalCode, mtosi.SITOriginHHGActualAddress.PostalCode, "shipment actual zip is the same")

				// Verify SIT entry date and SIT departure date are updated on the service item
				sitEntryDate := mtosi.SITEntryDate.UTC()
				sitDepartureDate := mtosi.SITDepartureDate.UTC()
				suite.Equal(subtestData.mtoServiceItem.SITEntryDate, &sitEntryDate, "sit entry date is the same")
				suite.Equal(subtestData.mtoServiceItem.SITDepartureDate, &sitDepartureDate, "sit departure date is the same")
			}
		}
		suite.Equal(true, foundDOFSIT, "Found expected ReServiceCodeDOFSIT")
		suite.Equal(true, foundDOPSIT, "Found expected ReServiceCodeDOPSIT")
		suite.Equal(true, foundDOASIT, "Found expected ReServiceCodeDOASIT")
	})

}

func (suite *HandlerSuite) TestCreateMTOServiceItemDestSITHandler() {

	builder := query.NewQueryBuilder()
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	sitEntryDate := time.Now().Add(time.Hour * 24)

	type localSubtestData struct {
		mto            models.Move
		mtoShipment    models.MTOShipment
		mtoServiceItem models.MTOServiceItem
		params         mtoserviceitemops.CreateMTOServiceItemParams
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		subtestData.mto = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		subtestData.mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.mto,
				LinkOnly: true,
			},
		}, nil)
		factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDFSIT)

		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Reason:          models.StringPointer("lorem ipsum"),
			Description:     handlers.FmtString("description"),
			CustomerContacts: models.MTOServiceItemCustomerContacts{
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeFirst,
					DateOfContact:              time.Now().Add(time.Hour * 24),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeSecond,
					DateOfContact:              time.Now(),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
			},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			SITEntryDate: &sitEntryDate,
		}
		subtestData.params = mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		return subtestData
	}

	suite.Run("Successful POST - Integration Test", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		mtoServiceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Description:     handlers.FmtString("description"),
			SITEntryDate:    &sitEntryDate,
			Reason:          models.StringPointer("lorem ipsum"),
			CustomerContacts: models.MTOServiceItemCustomerContacts{
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeFirst,
					DateOfContact:              time.Now().Add(time.Hour * 24),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeSecond,
					DateOfContact:              time.Now().Add(time.Hour * 24),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
			},
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		paramsDDFSIT := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItemDDFSIT),
		}

		// Validate incoming payload
		suite.NoError(paramsDDFSIT.Body.Validate(strfmt.Default))

		// CHECK RESULTS
		response := handler.Handle(paramsDDFSIT)

		//Validate incoming payload
		suite.NoError(paramsDDFSIT.Body.Validate(strfmt.Default))

		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemOK).Payload
		suite.NotZero(responsePayload[0].ID())
	})

	suite.Run("Successful POST - create DDFSIT without customer contact fields", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a DDFSIT MTOServiceItem missing Customer Contact fields
		// Expected outcome:
		//             Successful creation of Destination SIT service items
		// SETUP
		// Create the payload
		mtoServiceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Description:     handlers.FmtString("description"),
			SITEntryDate:    &sitEntryDate,
			Reason:          models.StringPointer("lorem ipsum"),
		}
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		paramsDDFSIT := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItemDDFSIT),
		}

		// Validate incoming payload
		suite.NoError(paramsDDFSIT.Body.Validate(strfmt.Default))

		// CHECK RESULTS
		response := handler.Handle(paramsDDFSIT)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemOK).Payload
		suite.NotZero(responsePayload[0].ID())
	})

	suite.Run("Failure POST - Integration Test - Missing reason", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		mtoServiceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Description:     handlers.FmtString("description"),
			SITEntryDate:    &sitEntryDate,
			Reason:          nil,
			CustomerContacts: models.MTOServiceItemCustomerContacts{
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeFirst,
					DateOfContact:              time.Now().Add(time.Hour * 24),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeSecond,
					DateOfContact:              time.Now().Add(time.Hour * 24),
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
			},
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		paramsDDFSIT := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItemDDFSIT),
		}

		// CHECK RESULTS
		response := handler.Handle(paramsDDFSIT)

		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.Run("Successful POST - Create DDASIT standalone", func() {
		subtestData := makeSubtestData()
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		// now that the mto service item has been created, create a standalone
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDDASIT
		params = mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response = handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)

		// Validate outgoing payload (each element of slice)
		for _, mtoServiceItem := range okResponse.Payload {
			suite.NoError(mtoServiceItem.Validate(strfmt.Default))
		}

		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("POST Failure - Cannot create DDASIT without DDFSIT", func() {
		subtestData := makeSubtestData()
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDDASIT
		subtestData.mtoServiceItem.MTOShipment = mtoShipment
		subtestData.mtoServiceItem.MTOShipmentID = &mtoShipment.ID

		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 422 Cannot create DDDSIT standalone", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a DDDSIT MTOServiceItem
		// Expected outcome:
		//             Receive a 422 - Unprocessable Entity
		// SETUP
		// Create the payload
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDDDSIT
		moveRouter := moverouter.NewMoveRouter()
		planner := &routemocks.Planner{}
		creator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		handler := CreateMTOServiceItemHandler{
			suite.HandlerConfig(),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		// CHECK RESULTS

		// Validate incoming payload
		suite.Error(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemDDDSIT() {

	// Under test: updateMTOServiceItemHandler.Handle function
	//             MTOServiceItemUpdater.Update service object function
	// SETUP
	// Create the service item in the db for dddsit
	type localSubtestData struct {
		dddsit     models.MTOServiceItem
		handler    UpdateMTOServiceItemHandler
		reqPayload *primemessages.UpdateMTOServiceItemSIT
		params     mtoserviceitemops.UpdateMTOServiceItemParams
	}

	sitStatusService := sitstatus.NewShipmentSITStatus()

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		timeNow := time.Now()
		requestApproavalsRequestedStatus := false
		// Number of days of grace period after customer contacts prime for delivery out of SIT
		const GracePeriodDays = 5

		year, month, day := timeNow.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := timeNow.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &timeNow,
					ApprovedAt:         &timeNow,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// We need to create a destination first day sit in order to properly calculate authorized end date
		ddfsitServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &timeNow,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "APPROVED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		subtestData.dddsit = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     models.TimePointer(time.Now()),
					SITDepartureDate: &timeNow,
					Status:           "REJECTED",
					Reason:           models.StringPointer("reason"),
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, subtestData.dddsit, ddfsitServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		addr := primemessages.Address{
			StreetAddress1: &destinationAddress.StreetAddress1,
			City:           &destinationAddress.City,
			State:          &destinationAddress.State,
			PostalCode:     &destinationAddress.PostalCode,
			Country:        models.StringPointer("US"),
		}

		milTime := "1400Z"
		requestApprovalRequestedStatus := false

		// Create the payload with the desired update
		subtestData.reqPayload = &primemessages.UpdateMTOServiceItemSIT{
			ReServiceCode:                   models.ReServiceCodeDDDSIT.String(),
			SitDepartureDate:                *handlers.FmtDate(time.Now().AddDate(0, 0, 5)),
			SitDestinationFinalAddress:      &addr,
			DateOfContact1:                  handlers.FmtDate(time.Date(2020, time.December, 04, 0, 0, 0, 0, time.UTC)),
			TimeMilitary1:                   handlers.FmtStringPtrNonEmpty(&milTime),
			FirstAvailableDeliveryDate1:     handlers.FmtDate(time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)),
			SitCustomerContacted:            handlers.FmtDate(time.Now()),
			SitRequestedDelivery:            handlers.FmtDate(time.Now().AddDate(0, 0, 3)),
			RequestApprovalsRequestedStatus: &requestApprovalRequestedStatus,
			UpdateReason:                    models.StringPointer("reason for updating"),
		}
		subtestData.reqPayload.SetID(strfmt.UUID(subtestData.dddsit.ID.String()))

		// Create the handler
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
		addressCreator := address.NewAddressCreator()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		subtestData.handler = UpdateMTOServiceItemHandler{
			suite.HandlerConfig(),
			mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		}

		// create the params struct
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service-items/%s", subtestData.dddsit.ID), nil)
		eTag := etag.GenerateEtag(subtestData.dddsit.UpdatedAt)
		subtestData.params = mtoserviceitemops.UpdateMTOServiceItemParams{
			HTTPRequest:      req,
			Body:             subtestData.reqPayload,
			MtoServiceItemID: subtestData.dddsit.ID.String(),
			IfMatch:          eTag,
		}
		return subtestData
	}

	suite.Run("Successful PATCH - Updated SITDepartureDate and FADD Fields and other date fields on DDDSIT", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We create an mto service item using DDDSIT (which was created above)
		//             And send an update to the sit entry date and customer contact fields
		// Expected outcome:
		//             Receive a success response with the SitDepartureDate, TimeMilitary and FirstAvailableDeliveryDate1 updated

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemOK{}, response)
		r := response.(*mtoserviceitemops.UpdateMTOServiceItemOK)
		resp1 := r.Payload

		// TODO: This is failing because DOPSIT and DDDSIT are being sent back in the response
		//   but those are not listed in the enum in the swagger file for an MTO service item.
		//   Need to figure out best way to resolve.
		// Validate outgoing payload
		// suite.NoError(resp1.Validate(strfmt.Default))

		respPayload := resp1.(*primemessages.MTOServiceItemDestSIT)

		suite.Equal(subtestData.reqPayload.ID(), respPayload.ID())
		suite.Equal(subtestData.reqPayload.SitDepartureDate.String(), respPayload.SitDepartureDate.String())
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.StreetAddress1, respPayload.SitDestinationFinalAddress.StreetAddress1)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.City, respPayload.SitDestinationFinalAddress.City)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.PostalCode, respPayload.SitDestinationFinalAddress.PostalCode)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.State, respPayload.SitDestinationFinalAddress.State)
		suite.Equal(subtestData.reqPayload.DateOfContact1, respPayload.DateOfContact1)
		suite.Equal(subtestData.reqPayload.TimeMilitary1, respPayload.TimeMilitary1)
		suite.Equal(subtestData.reqPayload.FirstAvailableDeliveryDate1, respPayload.FirstAvailableDeliveryDate1)
		suite.Equal(*subtestData.reqPayload.SitCustomerContacted, *respPayload.SitCustomerContacted)
		suite.Equal(*subtestData.reqPayload.SitRequestedDelivery, *respPayload.SitRequestedDelivery)
	})

	suite.Run("Failed PATCH - No DDDSIT found", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a non existent DDDSIT item
		//             And send an update to DOPSIT to the SitDepartureDate
		// Expected outcome:
		//             Receive a NotFound error response

		// SETUP
		// Replace the request path with a bad id that won't be found
		badUUID := uuid.Must(uuid.NewV4())
		badReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service_items/%s", badUUID), nil)
		subtestData.params.HTTPRequest = badReq
		subtestData.params.MtoServiceItemID = badUUID.String()
		subtestData.reqPayload.SetID(strfmt.UUID(badUUID.String()))

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemNotFound).Payload

		// CHECK RESULTS

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Failure 422 - Unprocessable Entity", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a non existent DDDSIT item ID in the param body
		//             And send an update to DDDSIT to the SitDepartureDate
		// Expected outcome:
		//             Receive an unprocessable entity error response

		// SETUP
		// Replace the payload ID with one that does not match request param
		badUUID := uuid.Must(uuid.NewV4())
		subtestData.reqPayload.SetID(strfmt.UUID(badUUID.String()))

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Failed PATCH - Payment request created", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a DDDSIT that already has a payment request associated
		//             Then try to update the SitDepartureDate on that
		// Expected outcome:
		//             Receive a ConflictError response

		// SETUP
		// Make a payment request and link to the dddsit service item
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		cost := unit.Cents(20000)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    subtestData.dddsit,
				LinkOnly: true,
			},
		}, nil)

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemConflict{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemConflict).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemDOPSIT() {

	// Under test: updateMTOServiceItemHandler.Handle function
	//             MTOServiceItemUpdater.Update service object function
	// SETUP
	// Create the service item in the db for dofsit and DOPSIT
	// Create the handler
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	sitStatusService := sitstatus.NewShipmentSITStatus()

	type localSubtestData struct {
		dopsit     models.MTOServiceItem
		handler    UpdateMTOServiceItemHandler
		reqPayload *primemessages.UpdateMTOServiceItemSIT
		params     mtoserviceitemops.UpdateMTOServiceItemParams
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		requestApprovalRequestedStatus := false
		subtestData = &localSubtestData{}
		timeNow := time.Now()
		requestApproavalsRequestedStatus := false
		// Number of days of grace period after customer contacts prime for delivery out of SIT
		const GracePeriodDays = 5

		year, month, day := timeNow.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := timeNow.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &timeNow,
					ApprovedAt:         &timeNow,
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// We need to create an origin first day sit in order to properly calculate authorized end date
		dofsitServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &timeNow,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "APPROVED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		subtestData.dopsit = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: models.TimePointer(time.Now()),
					Status:       models.MTOServiceItemStatusRejected,
					Reason:       models.StringPointer("reason"),
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOPSIT,
				},
			},
		}, nil)

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, subtestData.dopsit, dofsitServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Create the payload with the desired update
		subtestData.reqPayload = &primemessages.UpdateMTOServiceItemSIT{
			ReServiceCode:                   models.ReServiceCodeDOPSIT.String(),
			SitDepartureDate:                *handlers.FmtDate(time.Now().AddDate(0, 0, 5)),
			RequestApprovalsRequestedStatus: &requestApprovalRequestedStatus,
			UpdateReason:                    models.StringPointer("a new reason"),
		}
		subtestData.reqPayload.SetID(strfmt.UUID(subtestData.dopsit.ID.String()))
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		subtestData.handler = UpdateMTOServiceItemHandler{
			suite.HandlerConfig(),
			mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		}

		// create the params struct
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service_items/%s", subtestData.dopsit.ID), nil)
		eTag := etag.GenerateEtag(subtestData.dopsit.UpdatedAt)
		subtestData.params = mtoserviceitemops.UpdateMTOServiceItemParams{
			HTTPRequest:      req,
			Body:             subtestData.reqPayload,
			MtoServiceItemID: subtestData.dopsit.ID.String(),
			IfMatch:          eTag,
		}
		return subtestData
	}

	suite.Run("Successful PATCH - Updated SITDepartureDate on DOPSIT", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We create an mto service item using DOFSIT (which was created above)
		//             And send an update to the sit entry date
		// Expected outcome:
		//             Receive a success response with the SitDepartureDate updated

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemOK{}, response)
		r := response.(*mtoserviceitemops.UpdateMTOServiceItemOK)
		resp1 := r.Payload

		// TODO: This is failing because DOPSIT and DDDSIT are being sent back in the response
		//   but those are not listed in the enum in the swagger file for an MTO service item.
		// Validate outgoing payload
		// suite.NoError(resp1.Validate(strfmt.Default))

		respPayload := resp1.(*primemessages.MTOServiceItemOriginSIT)
		suite.Equal(subtestData.reqPayload.ID(), respPayload.ID())
		suite.Equal(subtestData.reqPayload.SitDepartureDate.String(), respPayload.SitDepartureDate.String())
	})

	suite.Run("Failed PATCH - No DOPSIT found", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a non existent DOPSIT item
		//             And send an update to DOPSIT to the SitDepartureDate
		// Expected outcome:
		//             Receive a NotFound error response

		// SETUP
		// Replace the request path with a bad id that won't be found
		badUUID := uuid.Must(uuid.NewV4())
		badReq := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service_items/%s", badUUID), nil)
		subtestData.params.HTTPRequest = badReq
		subtestData.params.MtoServiceItemID = badUUID.String()
		subtestData.reqPayload.SetID(strfmt.UUID(badUUID.String()))

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemNotFound{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("Failure 422 - Unprocessable Entity", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a non existent DOPSIT item ID in the param body
		//             And send an update to DOPSIT to the SitDepartureDate
		// Expected outcome:
		//             Receive an unprocessable entity error response

		// SETUP
		// Replace the payload ID with one that does not match request param
		badUUID := uuid.Must(uuid.NewV4())
		subtestData.reqPayload.SetID(strfmt.UUID(badUUID.String()))

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))

		// return to good state for next test
		subtestData.reqPayload.SetID(strfmt.UUID(subtestData.dopsit.ID.String()))
	})

	suite.Run("Failed PATCH - Payment request created", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We use a DOPSIT that already has a payment request associated
		//             Then try to update the SitDepartureDate on that
		// Expected outcome:
		//             Receive a ConflictError response

		// SETUP
		// Make a payment request and link to the DOPSIT service item
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		cost := unit.Cents(20000)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    subtestData.dopsit,
				LinkOnly: true,
			},
		}, nil)

		// CALL FUNCTION UNDER TEST

		// Validate incoming payload
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))

		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemConflict{}, response)
		responsePayload := response.(*mtoserviceitemops.UpdateMTOServiceItemConflict).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}
