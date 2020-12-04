package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/unit"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/models"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateMTOServiceItemHandler() {
	mto := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})
	testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
		},
	})
	builder := query.NewQueryBuilder(suite.DB())
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

	req := httptest.NewRequest("POST", "/mto-service-items", nil)

	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrderID:  mto.ID,
		MTOShipmentID:    &mtoShipment.ID,
		ReService:        models.ReService{Code: models.ReServiceCodeDOFSIT},
		Reason:           nil,
		PickupPostalCode: nil,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	params := mtoserviceitemops.CreateMTOServiceItemParams{
		HTTPRequest: req,
		Body:        payloads.MTOServiceItem(&mtoServiceItem),
	}

	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		err := errors.New("ServerError")

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemInternalServerError{}, response)

		errResponse := response.(*mtoserviceitemops.CreateMTOServiceItemInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")

	})

	suite.T().Run("POST failure - 422 Unprocessable Entity Error", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		// InvalidInputError should generate an UnprocessableEntity response
		err := services.InvalidInputError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 409 Conflict Error", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		// ConflictError should generate a Conflict response
		err := services.ConflictError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemConflict{}, response)
	})

	suite.T().Run("POST failure - 404", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		err := services.NotFoundError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
	})

	suite.T().Run("POST failure - 404 - MTO is not available to Prime", func(t *testing.T) {
		mtoNotAvailable := testdatagen.MakeDefaultMove(suite.DB())

		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		body := payloads.MTOServiceItem(&mtoServiceItem)
		body.SetMoveTaskOrderID(handlers.FmtUUID(mtoNotAvailable.ID))

		paramsNotAvailable := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        body,
		}

		response := handler.Handle(paramsNotAvailable)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)

		typedResponse := response.(*mtoserviceitemops.CreateMTOServiceItemNotFound)
		suite.Contains(*typedResponse.Payload.Detail, mtoNotAvailable.ID.String())
	})

	suite.T().Run("POST failure - 404 - Integration - ShipmentID not linked by MoveTaskOrderID", func(t *testing.T) {
		mto2 := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto2,
		})
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		body := payloads.MTOServiceItem(&mtoServiceItem)
		body.SetMoveTaskOrderID(handlers.FmtUUID(mtoShipment.MoveTaskOrderID))
		body.SetMtoShipmentID(strfmt.UUID(mtoShipment2.ID.String()))

		newParams := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        body,
		}

		response := handler.Handle(newParams)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemNotFound{}, response)
	})

	suite.T().Run("POST failure - 422 - Model validation errors", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		verrs := validate.NewErrors()
		verrs.Add("test", "testing")

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, verrs, nil)

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 422 - modelType() not supported", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		err := services.NotFoundError{}

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		mtoServiceItem := models.MTOServiceItem{
			MoveTaskOrderID:  mto.ID,
			MTOShipmentID:    &mtoShipment.ID,
			ReService:        models.ReService{Code: models.ReServiceCodeMS},
			Reason:           nil,
			PickupPostalCode: nil,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})
}

func (suite *HandlerSuite) TestCreateMTOServiceItemDomesticCratingHandler() {
	mto := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DCRT",
		},
	})
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DUCRT",
		},
	})
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DCRTSA",
		},
	})
	builder := query.NewQueryBuilder(suite.DB())
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

	req := httptest.NewRequest("POST", "/mto-service-items", nil)

	mtoServiceItem := models.MTOServiceItem{
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
	}

	suite.T().Run("Successful POST - Integration Test - Domestic Crating", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		mtoServiceItem.ReService.Code = models.ReServiceCodeDCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.T().Run("Successful POST - Integration Test - Domestic Uncrating", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		mtoServiceItem.ReService.Code = models.ReServiceCodeDUCRT
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.T().Run("Successful POST - Integration Test - Domestic Crating Standalone", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		mtoServiceItem.ReService.Code = models.ReServiceCodeDCRTSA
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})

	suite.T().Run("POST failure - 422", func(t *testing.T) {
		mockCreator := mocks.MTOServiceItemCreator{}
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
			mtoChecker,
		}
		err := errors.New("ServerError")

		mockCreator.On("CreateMTOServiceItem",
			mock.Anything,
		).Return(nil, nil, err)

		mtoServiceItem.ReService.Code = models.ReServiceCodeDCRTSA
		params := mtoserviceitemops.CreateMTOServiceItemParams{
			HTTPRequest: req,
			Body:        payloads.MTOServiceItem(&mtoServiceItem),
		}

		var height int32 = 0
		params.Body.(*primemessages.MTOServiceItemDomesticCrating).Crate.Height = &height
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemUnprocessableEntity{}, response)
	})
}

func (suite *HandlerSuite) TestCreateMTOServiceItemDDFSITHandler() {
	mto := testdatagen.MakeAvailableMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})
	testdatagen.MakeDDFSITReService(suite.DB())
	builder := query.NewQueryBuilder(suite.DB())
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())

	req := httptest.NewRequest("POST", "/mto-service-items", nil)

	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrderID: mto.ID,
		MTOShipmentID:   &mtoShipment.ID,
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	params := mtoserviceitemops.CreateMTOServiceItemParams{
		HTTPRequest: req,
		Body:        payloads.MTOServiceItem(&mtoServiceItem),
	}

	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		creator := mtoserviceitem.NewMTOServiceItemCreator(builder)
		handler := CreateMTOServiceItemHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
			mtoChecker,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemops.CreateMTOServiceItemOK{}, response)

		okResponse := response.(*mtoserviceitemops.CreateMTOServiceItemOK)
		suite.NotZero(okResponse.Payload[0].ID())
	})
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemDDDSIT() {

	// Under test: updateMTOServiceItemHandler.Handle function
	//             MTOServiceItemUpdater.Update service object function
	// SETUP
	// Create the service item in the db for dofsit and dddsit
	timeNow := time.Now()
	dddsit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &timeNow,
		},
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: swag.Time(time.Now()),
		},
		ReService: models.ReService{
			Code: "DDDSIT",
		},
	})

	// Create the payload with the desired update
	reqPayload := &primemessages.UpdateMTOServiceItemSIT{
		ReServiceCode:    "DDDSIT",
		SitDepartureDate: *handlers.FmtDate(time.Now().AddDate(0, 0, 5)),
	}
	reqPayload.SetID(strfmt.UUID(dddsit.ID.String()))

	// Create the handler
	queryBuilder := query.NewQueryBuilder(suite.DB())
	handler := UpdateMTOServiceItemHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
	}

	// create the params struct
	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-service_items/%s", dddsit.ID), nil)
	eTag := etag.GenerateEtag(dddsit.UpdatedAt)
	params := mtoserviceitemops.UpdateMTOServiceItemParams{
		HTTPRequest: req,
		Body:        reqPayload,
		IfMatch:     eTag,
	}

	suite.T().Run("Successful PATCH - Updated SITDepartureDate on DOPSIT", func(t *testing.T) {
		// Under test: updateMTOServiceItemHandler.Handle function
		//             MTOServiceItemUpdater.Update service object function
		// Set up:     We create an mto service item using DOFSIT (which was created above)
		//             And send an update to the sit entry date
		// Expected outcome:
		//             Receive a success response with the SitDepartureDate updated

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemOK{}, response)
		r := response.(*mtoserviceitemops.UpdateMTOServiceItemOK)
		resp1 := r.Payload

		respPayload := resp1.(*primemessages.MTOServiceItemDDFSIT)
		suite.Equal(reqPayload.ID(), respPayload.ID())
		suite.Equal(reqPayload.SitDepartureDate.String(), respPayload.SitDepartureDate.String())

		// Return to good state for next test
		params.IfMatch = respPayload.ETag()

	})

	suite.T().Run("Failed PATCH - No DDDSIT found", func(t *testing.T) {
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
		params.HTTPRequest = badReq
		reqPayload.SetID(strfmt.UUID(badUUID.String()))

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemNotFound{}, response)

		// return to good state for next test
		params.HTTPRequest = req
		reqPayload.SetID(strfmt.UUID(dddsit.ID.String()))
	})

	suite.T().Run("Failed PATCH - Payment request created", func(t *testing.T) {
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
			MTOServiceItem: dddsit,
		})

		// CALL FUNCTION UNDER TEST
		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)

		// CHECK RESULTS
		suite.IsType(&mtoserviceitemops.UpdateMTOServiceItemConflict{}, response)
	})

}
