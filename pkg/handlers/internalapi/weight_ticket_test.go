package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateWeightTicketHandler() {
	// Reusable objects
	weightTicketCreator := weightticket.NewCustomerWeightTicketCreator()

	type weightTicketCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      weightticketops.CreateWeightTicketParams
		handler     CreateWeightTicketHandler
	}
	makeCreateSubtestData := func(authenticateRequest bool) (subtestData weightTicketCreateSubtestData) {
		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight_ticket", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = weightticketops.CreateWeightTicketParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateWeightTicketHandler{
			suite.NewHandlerConfig(),
			weightTicketCreator,
		}

		return subtestData
	}

	suite.Run("Successfully Create Weight Ticket - Integration Test", func() {
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketOK{}, response)

		createdWeightTicket := response.(*weightticketops.CreateWeightTicketOK).Payload

		suite.NotEmpty(createdWeightTicket.ID.String())
		suite.NotNil(createdWeightTicket.EmptyDocumentID.String())
		suite.NotNil(createdWeightTicket.FullDocumentID.String())
		suite.NotNil(createdWeightTicket.ProofOfTrailerOwnershipDocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)
		// Missing PPM Shipment ID
		params := subtestData.params

		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		subtestData := makeCreateSubtestData(false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.CreateWeightTicketUnauthorized{}, response)
	})

	suite.Run("POST failure - 404- not found - wrong service member", func() {
		subtestData := makeCreateSubtestData(false)

		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, unauthorizedUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&weightticketops.CreateWeightTicketNotFound{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.WeightTicketCreator{}

		subtestData := makeCreateSubtestData(true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		mockCreator.On("CreateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.CreateWeightTicketInternalServerError{}, response)
	})
}

//
// UPDATE test
//

func (suite *HandlerSuite) TestUpdateWeightTicketHandler() {
	// Reusable objects
	ppmShipmentUpdater := mocks.PPMShipmentUpdater{}
	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	weightTicketUpdater := weightticket.NewCustomerWeightTicketUpdater(weightTicketFetcher, &ppmShipmentUpdater)

	type weightTicketUpdateSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.UpdateWeightTicketParams
		handler      UpdateWeightTicketHandler
	}
	makeUpdateSubtestData := func(authenticateRequest bool) (subtestData weightTicketUpdateSubtestData) {
		// Use fake data:
		subtestData.weightTicket = factory.BuildWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.weightTicket.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		eTag := etag.GenerateEtag(subtestData.weightTicket.UpdatedAt)
		subtestData.params = weightticketops.UpdateWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
			IfMatch:        eTag,
		}

		subtestData.handler = UpdateWeightTicketHandler{
			suite.createS3HandlerConfig(),
			weightTicketUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Weight Ticket - Integration Test", func() {
		subtestData := makeUpdateSubtestData(true)

		ppmShipmentUpdater.On(
			"UpdatePPMShipmentWithDefaultCheck",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, nil)

		params := subtestData.params

		// Add vehicleDescription
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(1),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(4000),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
			AdjustedNetWeight:        handlers.FmtInt64(3999),
			NetWeightRemarks:         "Adjusted net weight",
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketOK{}, response)

		updatedWeightTicket := response.(*weightticketops.UpdateWeightTicketOK).Payload
		suite.Equal(subtestData.weightTicket.ID.String(), updatedWeightTicket.ID.String())
		suite.Equal(params.UpdateWeightTicketPayload.VehicleDescription, *updatedWeightTicket.VehicleDescription)
		suite.Equal(params.UpdateWeightTicketPayload.AdjustedNetWeight, updatedWeightTicket.AdjustedNetWeight)
		suite.Equal(params.UpdateWeightTicketPayload.NetWeightRemarks, *updatedWeightTicket.NetWeightRemarks)
	})

	suite.Run("PATCH failure -400 - nil body", func() {
		subtestData := makeUpdateSubtestData(true)
		subtestData.params.UpdateWeightTicketPayload = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.UpdateWeightTicketBadRequest{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(0),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(0),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{}
		// This test should fail because of the wrong ID
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("f20d9c9b-2de5-4860-ad31-fd5c10e739f6"))
		params.WeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {

		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.WeightTicketUpdater{}

		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateWeightTicketPayload = &internalmessages.UpdateWeightTicket{
			VehicleDescription:       "Subaru",
			EmptyWeight:              handlers.FmtInt64(1),
			MissingEmptyWeightTicket: false,
			FullWeight:               handlers.FmtInt64(4000),
			MissingFullWeightTicket:  false,
			OwnsTrailer:              true,
			TrailerMeetsCriteria:     true,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.WeightTicket"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		handler := UpdateWeightTicketHandler{
			suite.NewHandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.UpdateWeightTicketInternalServerError{}, response)
		errResponse := response.(*weightticketops.UpdateWeightTicketInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Payload title is wrong")
	})
}

//
// DELETE test
//

func (suite *HandlerSuite) TestDeleteWeightTicketHandler() {
	// Create Reusable objects
	fetcher := weightticket.NewWeightTicketFetcher()
	estimator := mocks.PPMEstimator{}
	weightTicketDeleter := weightticket.NewWeightTicketDeleter(fetcher, &estimator)

	type weightTicketDeleteSubtestData struct {
		ppmShipment  models.PPMShipment
		weightTicket models.WeightTicket
		params       weightticketops.DeleteWeightTicketParams
		handler      DeleteWeightTicketHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData weightTicketDeleteSubtestData) {
		// Fake data:
		subtestData.weightTicket = factory.BuildWeightTicket(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.weightTicket.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/weight-ticket/%s", subtestData.ppmShipment.ID.String(), subtestData.weightTicket.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = weightticketops.
			DeleteWeightTicketParams{
			HTTPRequest:    req,
			PpmShipmentID:  *handlers.FmtUUID(subtestData.ppmShipment.ID),
			WeightTicketID: *handlers.FmtUUID(subtestData.weightTicket.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteWeightTicketHandler{
			suite.createS3HandlerConfig(),
			weightTicketDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete Weight Ticket - Integration Test", func() {
		mockIncentive := unit.Cents(100000)
		estimator.On("FinalIncentiveWithDefaultChecks", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("models.PPMShipment"), mock.AnythingOfType("*models.PPMShipment")).Return(&mockIncentive, nil)

		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&weightticketops.DeleteWeightTicketUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		subtestData := makeDeleteSubtestData(false)

		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{{
			Model: models.User{
				Roles: roles.Roles{
					{
						RoleType: roles.RoleTypeTOO,
					},
				},
			},
		}}, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&weightticketops.DeleteWeightTicketForbidden{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong service member user", func() {
		subtestData := makeDeleteSubtestData(false)

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, otherServiceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&weightticketops.DeleteWeightTicketForbidden{}, response)
	})
	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and weight ticket ID don't match", func() {
		subtestData := makeDeleteSubtestData(false)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, serviceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)
		// otherPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID matches serviceMember.ID
		suite.IsType(&weightticketops.DeleteWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.WeightTicketID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketNotFound{}, response)
	})

	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.WeightTicketDeleter{}

		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteWeightTicket",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := DeleteWeightTicketHandler{
			suite.createS3HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&weightticketops.DeleteWeightTicketInternalServerError{}, response)
	})
}
