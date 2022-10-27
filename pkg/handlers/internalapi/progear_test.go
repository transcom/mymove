package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	progear "github.com/transcom/mymove/pkg/services/progear"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateProgearHandler() {
	// Reusable objects
	progearCreator := progear.NewCustomerProgearCreator()

	type progearCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      progearops.CreateProGearWeightTicketParams
		handler     CreateProgearHandler
	}
	makeCreateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData progearCreateSubtestData) {
		db := appCtx.DB()

		subtestData.ppmShipment = testdatagen.MakePPMShipment(db, testdatagen.Assertions{})
		endpoint := fmt.Sprintf("/ppm-shipments/%s/pro-gear-weight-tickets", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = progearops.CreateProGearWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateProgearHandler{
			suite.HandlerConfig(),
			progearCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketOK{}, response)

		createdProgear := response.(*progearops.CreateProGearWeightTicketOK).Payload

		suite.NotEmpty(createdProgear.ID.String())
		suite.NotNil(createdProgear.EmptyDocumentID.String())
		suite.NotNil(createdProgear.FullDocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&progearops.CreateProGearWeightTicketUnauthorized{}, response)
	})

	suite.Run("POST failure - 403- permission denied - wrong application", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&progearops.CreateProGearWeightTicketForbidden{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.ProgearCreator{}
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateProgear",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateProgearHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&progearops.CreateProGearWeightTicketInternalServerError{}, response)
	})
}

//
// UPDATE test
//

// func (suite *HandlerSuite) TestUpdateProgearHandler() {
// 	// Reusable objects
// 	progearUpdater := progear.NewCustomerProgearUpdater()

// 	type progearUpdateSubtestData struct {
// 		ppmShipment models.PPMShipment
// 		progear     models.Progear
// 		params      progearops.UpdateProgearParams
// 		handler     UpdateProgearHandler
// 	}
// 	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData progearUpdateSubtestData) {
// 		db := appCtx.DB()

// 		// Use fake data:
// 		subtestData.progear = testdatagen.MakeProgear(db, testdatagen.Assertions{})
// 		subtestData.ppmShipment = subtestData.progear.PPMShipment
// 		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

// 		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.progear.ID.String())
// 		req := httptest.NewRequest("PATCH", endpoint, nil)
// 		if authenticateRequest {
// 			req = suite.AuthenticateRequest(req, serviceMember)
// 		}
// 		eTag := etag.GenerateEtag(subtestData.progear.UpdatedAt)
// 		subtestData.params = progearops.UpdateProgearParams{
// 			HTTPRequest:   req,
// 			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
// 			ProgearID:     *handlers.FmtUUID(subtestData.progear.ID),
// 			IfMatch:       eTag,
// 		}

// 		subtestData.handler = UpdateProgearHandler{
// 			suite.createS3HandlerConfig(),
// 			progearUpdater,
// 		}

// 		return subtestData
// 	}

// 	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)

// 		params := subtestData.params

// 		// An upload must exist if trailer is owned and qualifies to be claimed
// 		testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
// 			UserUpload: models.UserUpload{
// 				DocumentID: &subtestData.progear.ProofOfTrailerOwnershipDocumentID,
// 				Document:   subtestData.progear.ProofOfTrailerOwnershipDocument,
// 			},
// 		})

// 		// Add vehicleDescription
// 		params.UpdateProgearPayload = &internalmessages.UpdateProgear{
// 			VehicleDescription:   "Subaru",
// 			EmptyWeight:          handlers.FmtInt64(1),
// 			MissingEmptyProgear:  false,
// 			FullWeight:           handlers.FmtInt64(4000),
// 			MissingFullProgear:   false,
// 			OwnsTrailer:          true,
// 			TrailerMeetsCriteria: true,
// 		}

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&progearops.UpdateProgearOK{}, response)

// 		updatedProgear := response.(*progearops.UpdateProgearOK).Payload
// 		suite.Equal(subtestData.progear.ID.String(), updatedProgear.ID.String())
// 		suite.Equal(params.UpdateProgearPayload.VehicleDescription, *updatedProgear.VehicleDescription)
// 	})

// 	suite.Run("PATCH failure -400 - nil body", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		subtestData.params.UpdateProgearPayload = nil
// 		response := subtestData.handler.Handle(subtestData.params)

// 		suite.IsType(&progearops.UpdateProgearBadRequest{}, response)
// 	})

// 	suite.Run("PATCH failure -422 - Invalid Input", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateProgearPayload = &internalmessages.UpdateProgear{
// 			VehicleDescription:   "Subaru",
// 			EmptyWeight:          handlers.FmtInt64(0),
// 			MissingEmptyProgear:  false,
// 			FullWeight:           handlers.FmtInt64(0),
// 			MissingFullProgear:   false,
// 			OwnsTrailer:          true,
// 			TrailerMeetsCriteria: true,
// 		}

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&progearops.UpdateProgearUnprocessableEntity{}, response)
// 	})

// 	suite.Run("PATCH failure - 404- not found", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateProgearPayload = &internalmessages.UpdateProgear{}
// 		// This test should fail because of the wrong ID
// 		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
// 		params.ProgearID = *uuidString

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&progearops.UpdateProgearNotFound{}, response)
// 	})

// 	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateProgearPayload = &internalmessages.UpdateProgear{}
// 		params.IfMatch = "intentionally-bad-if-match-header-value"

// 		response := subtestData.handler.Handle(params)

// 		suite.IsType(&progearops.UpdateProgearPreconditionFailed{}, response)
// 	})

// 	suite.Run("PATCH failure - 500", func() {
// 		mockUpdater := mocks.ProgearUpdater{}
// 		appCtx := suite.AppContextForTest()

// 		subtestData := makeUpdateSubtestData(appCtx, true)
// 		params := subtestData.params
// 		params.UpdateProgearPayload = &internalmessages.UpdateProgear{
// 			VehicleDescription:   "Subaru",
// 			EmptyWeight:          handlers.FmtInt64(1),
// 			MissingEmptyProgear:  false,
// 			FullWeight:           handlers.FmtInt64(4000),
// 			MissingFullProgear:   false,
// 			OwnsTrailer:          true,
// 			TrailerMeetsCriteria: true,
// 		}

// 		err := errors.New("ServerError")

// 		mockUpdater.On("UpdateProgear",
// 			mock.AnythingOfType("*appcontext.appContext"),
// 			mock.AnythingOfType("models.Progear"),
// 			mock.AnythingOfType("string"),
// 		).Return(nil, err)

// 		handler := UpdateProgearHandler{
// 			suite.HandlerConfig(),
// 			&mockUpdater,
// 		}

// 		response := handler.Handle(params)

// 		suite.IsType(&progearops.UpdateProgearInternalServerError{}, response)
// 		errResponse := response.(*progearops.UpdateProgearInternalServerError)
// 		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
// 	})
// }
