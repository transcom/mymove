package internalapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"

	"net/http/httptest"
	"testing"

	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

//
// CREATE
//

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	mto := testdatagen.MakeDefaultMove(suite.DB())

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	secondaryPickupAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})

	destinationAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})

	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: models.MTOShipment{},
	})
	mtoShipment.MoveTaskOrderID = mto.ID

	mtoAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
	}}

	customerRemarks := "I have some grandfather clocks."

	builder := query.NewQueryBuilder(suite.DB())

	req := httptest.NewRequest("POST", "/mto_shipments", nil)
	unauthorizedReq := httptest.NewRequest("POST", "/mto_shipments", nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	params := mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:          agents,
			CustomerRemarks: &customerRemarks,
			PickupAddress: &internalmessages.Address{
				City:           &pickupAddress.City,
				Country:        pickupAddress.Country,
				PostalCode:     &pickupAddress.PostalCode,
				State:          &pickupAddress.State,
				StreetAddress1: &pickupAddress.StreetAddress1,
				StreetAddress2: pickupAddress.StreetAddress2,
				StreetAddress3: pickupAddress.StreetAddress3,
			},
			SecondaryPickupAddress: &internalmessages.Address{
				City:           &secondaryPickupAddress.City,
				Country:        secondaryPickupAddress.Country,
				PostalCode:     &secondaryPickupAddress.PostalCode,
				State:          &secondaryPickupAddress.State,
				StreetAddress1: &secondaryPickupAddress.StreetAddress1,
				StreetAddress2: secondaryPickupAddress.StreetAddress2,
				StreetAddress3: secondaryPickupAddress.StreetAddress3,
			},
			DestinationAddress: &internalmessages.Address{
				City:           &destinationAddress.City,
				Country:        destinationAddress.Country,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
			},
			SecondaryDeliveryAddress: &internalmessages.Address{
				City:           &secondaryDeliveryAddress.City,
				Country:        secondaryDeliveryAddress.Country,
				PostalCode:     &secondaryDeliveryAddress.PostalCode,
				State:          &secondaryDeliveryAddress.State,
				StreetAddress1: &secondaryDeliveryAddress.StreetAddress1,
				StreetAddress2: secondaryDeliveryAddress.StreetAddress2,
				StreetAddress3: secondaryDeliveryAddress.StreetAddress3,
			},
			RequestedPickupDate:   strfmt.Date(*mtoShipment.RequestedPickupDate),
			RequestedDeliveryDate: strfmt.Date(*mtoShipment.RequestedDeliveryDate),
			ShipmentType:          internalmessages.MTOShipmentTypeHHG,
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

		createdShipment := response.(*mtoshipmentops.CreateMTOShipmentOK).Payload

		suite.NotEmpty(createdShipment.ID.String())

		suite.Equal(*params.Body.CustomerRemarks, *createdShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *createdShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *createdShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *createdShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *createdShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), createdShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.RequestedDeliveryDate.String(), createdShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, createdShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, createdShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, createdShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, createdShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, createdShipment.Agents[0].AgentType)
		suite.Equal(createdShipment.ID.String(), string(createdShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(createdShipment.Agents[0].ID)
	})

	suite.T().Run("POST failure - 400 - invalid input, missing pickup address", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		badParams := params
		badParams.Body.PickupAddress = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 401- permission denied - not authenticated", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)
		unauthorizedParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: unauthorizedReq,
			Body: &internalmessages.CreateShipment{
				MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
				Agents:          internalmessages.MTOAgents{},
				CustomerRemarks: nil,
				PickupAddress: &internalmessages.Address{
					City:           &pickupAddress.City,
					Country:        pickupAddress.Country,
					PostalCode:     &pickupAddress.PostalCode,
					State:          &pickupAddress.State,
					StreetAddress1: &pickupAddress.StreetAddress1,
					StreetAddress2: pickupAddress.StreetAddress2,
					StreetAddress3: pickupAddress.StreetAddress3,
				},
				RequestedPickupDate:   strfmt.Date(*mtoShipment.RequestedPickupDate),
				RequestedDeliveryDate: strfmt.Date(*mtoShipment.RequestedDeliveryDate),
				ShipmentType:          internalmessages.MTOShipmentTypeHHG,
			},
		}

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
	})

	suite.T().Run("POST failure - 403- permission denied - wrong application", func(t *testing.T) {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(suite.DB(), builder, fetcher)
		unauthorizedReq = suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			creator,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnauthorized{}, response)
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
		}

		otherParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(otherParams)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {
		mockCreator := mocks.MTOShipmentCreator{}

		handler := CreateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockCreator,
		}

		err := errors.New("ServerError")

		mockCreator.On("CreateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})
}

//
// UPDATE
//

func (suite *HandlerSuite) getUpdateMTOShipmentParams(originalShipment models.MTOShipment) mtoshipmentops.UpdateMTOShipmentParams {
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	pickupAddress.StreetAddress1 = "123 Fake Test St NW"

	secondaryPickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	secondaryPickupAddress.StreetAddress1 = "89999 Other Test St NW"

	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"

	secondaryDeliveryAddress := testdatagen.MakeDefaultAddress(suite.DB())
	secondaryDeliveryAddress.StreetAddress1 = "9999 Test Fake Rd SE"

	mtoAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	agents := internalmessages.MTOAgents{&internalmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
	}}

	customerRemarks := ""

	req := httptest.NewRequest("PATCH", "/mto-shipments/"+originalShipment.ID.String(), nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	payload := internalmessages.UpdateShipment{
		Agents:          agents,
		CustomerRemarks: &customerRemarks,
		DestinationAddress: &internalmessages.Address{
			City:           &destinationAddress.City,
			Country:        destinationAddress.Country,
			PostalCode:     &destinationAddress.PostalCode,
			State:          &destinationAddress.State,
			StreetAddress1: &destinationAddress.StreetAddress1,
			StreetAddress2: destinationAddress.StreetAddress2,
			StreetAddress3: destinationAddress.StreetAddress3,
		},
		SecondaryDeliveryAddress: &internalmessages.Address{
			City:           &secondaryDeliveryAddress.City,
			Country:        secondaryDeliveryAddress.Country,
			PostalCode:     &secondaryDeliveryAddress.PostalCode,
			State:          &secondaryDeliveryAddress.State,
			StreetAddress1: &secondaryDeliveryAddress.StreetAddress1,
			StreetAddress2: secondaryDeliveryAddress.StreetAddress2,
			StreetAddress3: secondaryDeliveryAddress.StreetAddress3,
		},
		PickupAddress: &internalmessages.Address{
			City:           &pickupAddress.City,
			Country:        pickupAddress.Country,
			PostalCode:     &pickupAddress.PostalCode,
			State:          &pickupAddress.State,
			StreetAddress1: &pickupAddress.StreetAddress1,
			StreetAddress2: pickupAddress.StreetAddress2,
			StreetAddress3: pickupAddress.StreetAddress3,
		},
		SecondaryPickupAddress: &internalmessages.Address{
			City:           &secondaryPickupAddress.City,
			Country:        secondaryPickupAddress.Country,
			PostalCode:     &secondaryPickupAddress.PostalCode,
			State:          &secondaryPickupAddress.State,
			StreetAddress1: &secondaryPickupAddress.StreetAddress1,
			StreetAddress2: secondaryPickupAddress.StreetAddress2,
			StreetAddress3: secondaryPickupAddress.StreetAddress3,
		},
		RequestedPickupDate:   strfmt.Date(*originalShipment.RequestedPickupDate),
		RequestedDeliveryDate: strfmt.Date(*originalShipment.RequestedDeliveryDate),
		ShipmentType:          internalmessages.MTOShipmentTypeHHG,
	}

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(originalShipment.ID),
		Body:          &payload,
		IfMatch:       eTag,
	}

	return params
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	builder := query.NewQueryBuilder(suite.DB())
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	suite.T().Run("Successful PATCH - Integration Test", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		suite.Equal(oldShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(*params.Body.CustomerRemarks, *updatedShipment.CustomerRemarks)
		suite.Equal(*params.Body.PickupAddress.StreetAddress1, *updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryPickupAddress.StreetAddress1, *updatedShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(*params.Body.DestinationAddress.StreetAddress1, *updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(*params.Body.SecondaryDeliveryAddress.StreetAddress1, *updatedShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), updatedShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.RequestedDeliveryDate.String(), updatedShipment.RequestedDeliveryDate.String())

		suite.Equal(params.Body.Agents[0].FirstName, updatedShipment.Agents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, updatedShipment.Agents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, updatedShipment.Agents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, updatedShipment.Agents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, updatedShipment.Agents[0].AgentType)
		suite.Equal(oldShipment.ID.String(), string(updatedShipment.Agents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.Agents[0].ID)
	})

	suite.T().Run("Successful PATCH - Can update shipment status", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		expectedStatus := internalmessages.MTOShipmentStatusSUBMITTED

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		params.Body.Status = expectedStatus

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)

		suite.Equal(expectedStatus, updatedResponse.Payload.Status)
	})

	suite.T().Run("PATCH failure - 400 -- nil body", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		params.Body = nil

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("PATCH failure - 400 -- invalid requested status update", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		params.Body.Status = internalmessages.MTOShipmentStatusREJECTED

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("PATCH failure - 401- permission denied - not authenticated", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		updateURI := "/mto-shipments/" + oldShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		params.HTTPRequest = unauthorizedReq

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.T().Run("PATCH failure - 403- permission denied - wrong application / user", func(t *testing.T) {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		updateURI := "/mto-shipments/" + oldShipment.ID.String()

		unauthorizedReq := httptest.NewRequest("PATCH", updateURI, nil)
		unauthorizedReq = suite.AuthenticateOfficeRequest(unauthorizedReq, officeUser)
		params.HTTPRequest = unauthorizedReq

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
	})

	suite.T().Run("PATCH failure - 404 -- not found", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		params.MtoShipmentID = *uuidString

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PATCH failure - 412 -- etag mismatch", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	//TODO: when the address bug (MB-3691) gets fixed, this test should pass
	// Update: This test is not passing due swagger validation failing and no server-side validation
	// happening. These changes weren't covered in MB-3691, so we'll need to do addt'l
	// work to fix. Since we have refactoring slated for addresses, we can do then.
	// suite.T().Run("PATCH failure - 422 -- invalid input", func(t *testing.T) {
	// 	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	// 	fetcher := fetch.NewFetcher(builder)
	// 	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, planner)
	// 	handler := UpdateMTOShipmentHandler{
	// 		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
	// 		updater,
	// 	}

	// 	oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
	// 	oldShipment.RequestedPickupDate = nil

	// 	req := httptest.NewRequest("PATCH", "/mto-shipments/"+oldShipment.ID.String(), nil)
	// 	req = suite.AuthenticateRequest(req, serviceMember)
	// 	// invalid zip
	// 	payloadDestinationAddress := &internalmessages.Address{
	// 		City:           swag.String("Stumptown"),
	// 		Country:        swag.String("USA"),
	// 		ID:             "6e07a670-a072-4014-be9f-4926c1389f9a",
	// 		State:          swag.String("CA"),
	// 		StreetAddress1: swag.String("321 Main St."),
	// 		PostalCode:     swag.String("123"),
	// 	}

	// 	payload := internalmessages.UpdateShipment{
	// 		DestinationAddress: payloadDestinationAddress,
	// 	}
	// 	eTag := etag.GenerateEtag(oldShipment.UpdatedAt)
	// 	params := mtoshipmentops.UpdateMTOShipmentParams{
	// 		HTTPRequest:   req,
	// 		MtoShipmentID: *handlers.FmtUUID(oldShipment.ID),
	// 		Body:          &payload,
	// 		IfMatch:       eTag,
	// 	}

	// 	response := handler.Handle(params)
	// 	suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)

	// })

	suite.T().Run("PATCH failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipmentCustomer",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		oldShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		params := suite.getUpdateMTOShipmentParams(oldShipment)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")
	})
}

//
// GET ALL
//

func (suite *HandlerSuite) TestListMTOShipmentsHandler() {
	mto := testdatagen.MakeDefaultMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})

	requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)

	pickupAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
	secondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "123 Nowhere",
			StreetAddress2: swag.String("P.O. Box 5555"),
			StreetAddress3: swag.String("c/o Some Other Person"),
			City:           "El Paso",
			State:          "TX",
			PostalCode:     "79916",
			Country:        swag.String("US"),
		},
	})

	deliveryAddress := testdatagen.MakeAddress4(suite.DB(), testdatagen.Assertions{})
	secondaryDeliveryAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "5432 Everywhere",
			StreetAddress2: swag.String("P.O. Box 111"),
			StreetAddress3: swag.String("c/o Some Other Person"),
			City:           "Portsmouth",
			State:          "NH",
			PostalCode:     "03801",
			Country:        swag.String("US"),
		},
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:                   models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:      &requestedPickupDate,
			PickupAddress:            &pickupAddress,
			SecondaryPickupAddress:   &secondaryPickupAddress,
			DestinationAddress:       &deliveryAddress,
			SecondaryDeliveryAddress: &secondaryDeliveryAddress,
		},
	})

	shipments := models.MTOShipments{mtoShipment, mtoShipment2}
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoshipmentops.ListMTOShipmentsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
	}

	suite.T().Run("Successful list fetch - 200 - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			listFetcher,
			fetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsOK{}, response)

		okResponse := response.(*mtoshipmentops.ListMTOShipmentsOK)

		suite.Len(okResponse.Payload, 2)

		firstShipmentReturned := okResponse.Payload[0]
		secondShipmentReturned := okResponse.Payload[1]

		// we expect the shipment that was created first to come first in the response
		suite.Equal(mtoShipment.ID.String(), firstShipmentReturned.ID.String())
		suite.Equal(mtoShipment2.ID.String(), secondShipmentReturned.ID.String())

		for i, returnedShipment := range okResponse.Payload {
			expectedShipment := shipments[i]

			suite.Equal(expectedShipment.Status, models.MTOShipmentStatus(returnedShipment.Status))

			suite.EqualDatePtr(expectedShipment.RequestedPickupDate, returnedShipment.RequestedPickupDate)

			suite.Equal(expectedShipment.PickupAddress.StreetAddress1, *returnedShipment.PickupAddress.StreetAddress1)
			suite.Equal(*expectedShipment.PickupAddress.StreetAddress2, *returnedShipment.PickupAddress.StreetAddress2)
			suite.Equal(*expectedShipment.PickupAddress.StreetAddress3, *returnedShipment.PickupAddress.StreetAddress3)
			suite.Equal(expectedShipment.PickupAddress.City, *returnedShipment.PickupAddress.City)
			suite.Equal(expectedShipment.PickupAddress.State, *returnedShipment.PickupAddress.State)
			suite.Equal(expectedShipment.PickupAddress.PostalCode, *returnedShipment.PickupAddress.PostalCode)

			suite.Equal(expectedShipment.SecondaryPickupAddress.StreetAddress1, *returnedShipment.SecondaryPickupAddress.StreetAddress1)
			suite.Equal(*expectedShipment.SecondaryPickupAddress.StreetAddress2, *returnedShipment.SecondaryPickupAddress.StreetAddress2)
			suite.Equal(*expectedShipment.SecondaryPickupAddress.StreetAddress3, *returnedShipment.SecondaryPickupAddress.StreetAddress3)
			suite.Equal(expectedShipment.SecondaryPickupAddress.City, *returnedShipment.SecondaryPickupAddress.City)
			suite.Equal(expectedShipment.SecondaryPickupAddress.State, *returnedShipment.SecondaryPickupAddress.State)
			suite.Equal(expectedShipment.SecondaryPickupAddress.PostalCode, *returnedShipment.SecondaryPickupAddress.PostalCode)

			suite.Equal(expectedShipment.DestinationAddress.StreetAddress1, *returnedShipment.DestinationAddress.StreetAddress1)
			suite.Equal(*expectedShipment.DestinationAddress.StreetAddress2, *returnedShipment.DestinationAddress.StreetAddress2)
			suite.Equal(*expectedShipment.DestinationAddress.StreetAddress3, *returnedShipment.DestinationAddress.StreetAddress3)
			suite.Equal(expectedShipment.DestinationAddress.City, *returnedShipment.DestinationAddress.City)
			suite.Equal(expectedShipment.DestinationAddress.State, *returnedShipment.DestinationAddress.State)
			suite.Equal(expectedShipment.DestinationAddress.PostalCode, *returnedShipment.DestinationAddress.PostalCode)

			suite.Equal(expectedShipment.SecondaryDeliveryAddress.StreetAddress1, *returnedShipment.SecondaryDeliveryAddress.StreetAddress1)
			suite.Equal(*expectedShipment.SecondaryDeliveryAddress.StreetAddress2, *returnedShipment.SecondaryDeliveryAddress.StreetAddress2)
			suite.Equal(*expectedShipment.SecondaryDeliveryAddress.StreetAddress3, *returnedShipment.SecondaryDeliveryAddress.StreetAddress3)
			suite.Equal(expectedShipment.SecondaryDeliveryAddress.City, *returnedShipment.SecondaryDeliveryAddress.City)
			suite.Equal(expectedShipment.SecondaryDeliveryAddress.State, *returnedShipment.SecondaryDeliveryAddress.State)
			suite.Equal(expectedShipment.SecondaryDeliveryAddress.PostalCode, *returnedShipment.SecondaryDeliveryAddress.PostalCode)
		}
	})

	suite.T().Run("POST failure - 400 - Bad Request", func(t *testing.T) {
		emtpyMTOID := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: "",
		}
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		response := handler.Handle(emtpyMTOID)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsBadRequest{}, response)
	})

	suite.T().Run("POST failure - 401 - permission denied - not authenticated", func(t *testing.T) {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := mtoshipmentops.ListMTOShipmentsParams{
			HTTPRequest:     unauthorizedReq,
			MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
		}
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.ListMTOShipmentsUnauthorized{}, response)
	})

	suite.T().Run("Failure list fetch - 404 Not Found - Move Task Order ID", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
	})

	suite.T().Run("Failure list fetch - 500 Internal Server Error", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
	})
}
