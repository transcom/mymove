package primeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
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

		mto := testdatagen.MakeAvailableMove(suite.DB())
		subtestData.mtoShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
		})
		testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			},
		})
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		reason := "lorem ipsum"
		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		// Customer gets new pickup address for SIT Origin Pickup (DOPSIT) which gets added when
		// creating DOFSIT (SIT origin first day).
		//
		// Do not create Address in the database (Assertions.Stub = true), because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		actualPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID:           mto.ID,
			MTOShipmentID:             &subtestData.mtoShipment.ID,
			ReService:                 models.ReService{Code: models.ReServiceCodeDOFSIT},
			Reason:                    &reason,
			SITEntryDate:              &sitEntryDate,
			SITPostalCode:             &sitPostalCode,
			SITOriginHHGActualAddress: &actualPickupAddress,
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
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("POST failure - 500", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}
		err := fmt.Errorf("ServerError")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemInternalServerError{}, response)

		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")

	})

	suite.Run("POST failure - 422 Unprocessable Entity Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}
		// InvalidInputError should generate an UnprocessableEntity response
		err := apperror.InvalidInputError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 409 Conflict Error", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}
		// ConflictError should generate a Conflict response
		err := apperror.ConflictError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemConflict{}, response)
	})

	suite.Run("POST failure - 404", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}
		err := apperror.NotFoundError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
	})

	suite.Run("POST failure - 404 - MTO is not available to Prime", func() {
		subtestData := makeSubtestData()
		mtoNotAvailable := testdatagen.MakeDefaultMove(suite.DB())
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		body := payloads.MTOServiceItem(&subtestData.mtoServiceItem)
		body.SetMoveTaskOrderID(handlers.FmtUUID(mtoNotAvailable.ID))

		paramsNotAvailable := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        body,
		}

		response := handler.Handle(paramsNotAvailable)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)

		typedResponse := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound)
		suite.Contains(*typedResponse.Payload.Detail, mtoNotAvailable.ID.String())
	})

	suite.Run("POST failure - 404 - Integration - ShipmentID not linked by MoveTaskOrderID", func() {
		subtestData := makeSubtestData()
		mto2 := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto2,
		})
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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

		response := handler.Handle(newParams)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
	})

	suite.Run("POST failure - 422 - Model validation errors", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			&mockCreator,
			mtoChecker,
		}
		verrs := validate.NewErrors()
		verrs.Add("test", "testing")

		mockCreator.On("CreateMTOServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, verrs, nil)

		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.Run("POST failure - 422 - modelType() not supported", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
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

		mto := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
		})
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDCRT,
			},
		})
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDUCRT,
			},
		})
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
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("Successful POST - Integration Test - Domestic Uncrating", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDUCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.req,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("POST failure - 422", func() {
		subtestData := makeSubtestData()
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
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

		subtestData.mto = testdatagen.MakeAvailableMove(suite.DB())
		subtestData.mtoShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.mto,
		})
		testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{})

		reason := "lorem ipsum"
		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{},
			Reason:          &reason,
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
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOPSIT
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.Error(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)

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
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)

	})

	suite.Run("Successful POST - Create DOASIT with DOFSIT", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a standalone DOASIT MTOServiceItem
		// Expected outcome:
		//             Receive a 404 - Not Found
		// SETUP
		// Create the payload
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
			Move:        subtestData.mto,
			MTOShipment: subtestData.mtoShipment,
		})

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOASIT
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

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
		mto := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
		})
		testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{})
		reason := "lorem ipsum"
		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: mto.ID,
			MTOShipmentID:   &mtoShipment.ID,
			ReService:       models.ReService{},
			Reason:          &reason,
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

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOFSIT
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
		unprocessableEntity := response.(*mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity)
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
		mto := testdatagen.MakeAvailableMove(suite.DB())
		subtestData.mtoShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
		})
		testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{})
		reason := "lorem ipsum"
		sitEntryDate := time.Now()
		sitPostalCode := "00000"

		// Original customer pickup address
		subtestData.originalPickupAddress = subtestData.mtoShipment.PickupAddress
		subtestData.originalPickupAddressID = subtestData.mtoShipment.PickupAddressID

		// Customer gets new pickup address

		// Do not create the Address in the database (Assertions.Stub = true), because if the information is coming from the Prime
		// via the Prime API, the address will not have a valid database ID. And tests need to ensure
		// that we properly create the address coming in from the API.
		subtestData.actualPickupAddress = testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true})

		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID:           mto.ID,
			MTOShipmentID:             &subtestData.mtoShipment.ID,
			ReService:                 models.ReService{},
			Reason:                    &reason,
			SITEntryDate:              &sitEntryDate,
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

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDOFSIT
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

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
		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
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
	sitEntryDate := time.Now()

	type localSubtestData struct {
		mto            models.Move
		mtoShipment    models.MTOShipment
		mtoServiceItem models.MTOServiceItem
		params         mtoserviceitemops.CreateMTOServiceItemParams
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		subtestData.mto = testdatagen.MakeAvailableMove(suite.DB())
		subtestData.mtoShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: subtestData.mto,
		})
		testdatagen.MakeDDFSITReService(suite.DB())

		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		subtestData.mtoServiceItem = models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Description:     handlers.FmtString("description"),
			CustomerContacts: models.MTOServiceItemCustomerContacts{
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeFirst,
					TimeMilitary:               "0400Z",
					FirstAvailableDeliveryDate: time.Now(),
				},
				models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeSecond,
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

	suite.Run("POST failure - 422 Cannot create DDFSIT with missing fields", func() {
		subtestData := makeSubtestData()
		// Under test: createMTOServiceItemHandler function
		// Set up:     We hit the endpoint with a DDFSIT MTOServiceItem missing Customer Contact fields
		// Expected outcome:
		//             Receive a 422 - Unprocessable Entity
		// SETUP
		// Create the payload

		mtoServiceItemDDFSIT := models.MTOServiceItem{
			MoveTaskOrderID: subtestData.mto.ID,
			MTOShipmentID:   &subtestData.mtoShipment.ID,
			ReService:       models.ReService{Code: models.ReServiceCodeDDFSIT},
			Description:     handlers.FmtString("description"),
			SITEntryDate:    &sitEntryDate,
		}
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		// CALL FUNCTION UNDER TEST
		req := httptest.NewRequest("POST", "/mto-service-items", nil)
		paramsDDFSIT := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItemDDFSIT),
		}

		// Run swagger validations
		suite.NoError(paramsDDFSIT.Body.Validate(strfmt.Default))

		// CHECK RESULTS
		response := handler.Handle(paramsDDFSIT)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)

	})

	suite.Run("Successful POST - Integration Test", func() {
		subtestData := makeSubtestData()
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := handler.Handle(subtestData.params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.Run("Successful POST - Create DDASIT standalone", func() {
		subtestData := makeSubtestData()
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		// now that the mto service item has been created, create a standalone
		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDDASIT
		params = mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		suite.NoError(params.Body.Validate(strfmt.Default))
		response = handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())

	})

	suite.Run("POST Failure - Cannot create DDASIT without DDFSIT", func() {
		subtestData := makeSubtestData()
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{})

		subtestData.mtoServiceItem.ReService.Code = models.ReServiceCodeDDASIT
		subtestData.mtoServiceItem.MTOShipment = mtoShipment
		subtestData.mtoServiceItem.MTOShipmentID = &mtoShipment.ID

		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: subtestData.params.HTTPRequest,
			Body:        payloads.MTOServiceItem(&subtestData.mtoServiceItem),
		}
		moveRouter := moverouter.NewMoveRouter()
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			creator,
			mtoChecker,
		}

		// CHECK RESULTS
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)

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
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
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
		suite.Error(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)

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

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		timeNow := time.Now()
		subtestData.dddsit = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &timeNow,
			},
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: swag.Time(time.Now()),
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		})

		destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
		addr := primemessages.Address{
			StreetAddress1: &destinationAddress.StreetAddress1,
			City:           &destinationAddress.City,
			State:          &destinationAddress.State,
			PostalCode:     &destinationAddress.PostalCode,
			Country:        destinationAddress.Country,
		}

		// Create the payload with the desired update
		subtestData.reqPayload = &primemessages.UpdateMTOServiceItemSIT{
			ReServiceCode:              models.ReServiceCodeDDDSIT.String(),
			SitDepartureDate:           *handlers.FmtDate(time.Now().AddDate(0, 0, 5)),
			SitDestinationFinalAddress: &addr,
		}
		subtestData.reqPayload.SetID(strfmt.UUID(subtestData.dddsit.ID.String()))

		// Create the handler
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		subtestData.handler = UpdateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter),
		}

		// create the params struct
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service_items/%s", subtestData.dddsit.ID), nil)
		eTag := etag.GenerateEtag(subtestData.dddsit.UpdatedAt)
		subtestData.params = mtoserviceitemops.UpdateMTOServiceItemParams{
			HTTPRequest:      req,
			Body:             subtestData.reqPayload,
			MtoServiceItemID: subtestData.dddsit.ID.String(),
			IfMatch:          eTag,
		}
		return subtestData
	}

	suite.Run("Successful PATCH - Updated SITDepartureDate on DDDSIT", func() {
		subtestData := makeSubtestData()
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We create an mto service item using DDDSIT (which was created above)
		//             And send an update to the sit entry date
		// Expected outcome:
		//             Receive a success response with the SitDepartureDate updated

		// CALL FUNCTION UNDER TEST
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemOK{}, response)
		r := response.(*mtoserviceitemops.UpdateMTOServiceItemOK)
		resp1 := r.Payload

		respPayload := resp1.(*primemessages.MTOServiceItemDestSIT)
		suite.Equal(subtestData.reqPayload.ID(), respPayload.ID())
		suite.Equal(subtestData.reqPayload.SitDepartureDate.String(), respPayload.SitDepartureDate.String())
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.StreetAddress1, respPayload.SitDestinationFinalAddress.StreetAddress1)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.City, respPayload.SitDestinationFinalAddress.City)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.PostalCode, respPayload.SitDestinationFinalAddress.PostalCode)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.State, respPayload.SitDestinationFinalAddress.State)
		suite.Equal(subtestData.reqPayload.SitDestinationFinalAddress.Country, respPayload.SitDestinationFinalAddress.Country)

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
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemNotFound{}, response)
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
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity{}, response)
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
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		cost := unit.Cents(20000)
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
			},
			PaymentRequest: paymentRequest,
			MTOServiceItem: subtestData.dddsit,
		})

		// CALL FUNCTION UNDER TEST
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemConflict{}, response)
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

	type localSubtestData struct {
		dopsit     models.MTOServiceItem
		handler    UpdateMTOServiceItemHandler
		reqPayload *primemessages.UpdateMTOServiceItemSIT
		params     mtoserviceitemops.UpdateMTOServiceItemParams
	}

	makeSubtestData := func() (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}
		timeNow := time.Now()
		subtestData.dopsit = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &timeNow,
			},
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: swag.Time(time.Now()),
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		// Create the payload with the desired update
		subtestData.reqPayload = &primemessages.UpdateMTOServiceItemSIT{
			ReServiceCode:    models.ReServiceCodeDOPSIT.String(),
			SitDepartureDate: *handlers.FmtDate(time.Now().AddDate(0, 0, 5)),
		}
		subtestData.reqPayload.SetID(strfmt.UUID(subtestData.dopsit.ID.String()))

		subtestData.handler = UpdateMTOServiceItemHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter),
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
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemOK{}, response)
		r := response.(*mtoserviceitemops.UpdateMTOServiceItemOK)
		resp1 := r.Payload

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
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemNotFound{}, response)

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
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemUnprocessableEntity{}, response)

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
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		cost := unit.Cents(20000)
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
			},
			PaymentRequest: paymentRequest,
			MTOServiceItem: subtestData.dopsit,
		})

		// CALL FUNCTION UNDER TEST
		suite.NoError(subtestData.params.Body.Validate(strfmt.Default))
		response := subtestData.handler.Handle(subtestData.params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemConflict{}, response)
	})

}
