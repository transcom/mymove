package internalapi

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
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
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment:   models.MTOShipment{},
	})
	mtoShipment.MoveTaskOrderID = mto.ID

	builder := query.NewQueryBuilder(suite.DB())

	req := httptest.NewRequest("POST", "/mto-shipments", nil)
	unauthorizedReq := httptest.NewRequest("POST", "/mto-shipments", nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	params := mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &internalmessages.CreateShipment{
			MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:          internalmessages.MTOAgents{},
			CustomerRemarks: nil,
			DestinationAddress: &internalmessages.Address{
				City:           &destinationAddress.City,
				Country:        destinationAddress.Country,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
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
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					Country:        destinationAddress.Country,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
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
		unauthorizedParams := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: unauthorizedReq,
			Body: &internalmessages.CreateShipment{
				MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
				Agents:          internalmessages.MTOAgents{},
				CustomerRemarks: nil,
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					Country:        destinationAddress.Country,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
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

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment:   models.MTOShipment{},
	})
	mtoShipment.MoveTaskOrderID = mto.ID

	builder := query.NewQueryBuilder(suite.DB())

	req := httptest.NewRequest("PATCH", "/mto-shipments", nil)
	unauthorizedReq := httptest.NewRequest("PATCH", "/mto-shipments", nil)
	req = suite.AuthenticateRequest(req, serviceMember)

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest: req,
		Body: &internalmessages.UpdateShipment{
			MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:          internalmessages.MTOAgents{},
			CustomerRemarks: nil,
			DestinationAddress: &internalmessages.Address{
				City:           &destinationAddress.City,
				Country:        destinationAddress.Country,
				PostalCode:     &destinationAddress.PostalCode,
				State:          &destinationAddress.State,
				StreetAddress1: &destinationAddress.StreetAddress1,
				StreetAddress2: destinationAddress.StreetAddress2,
				StreetAddress3: destinationAddress.StreetAddress3,
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
			RequestedPickupDate:   strfmt.Date(*mtoShipment.RequestedPickupDate),
			RequestedDeliveryDate: strfmt.Date(*mtoShipment.RequestedDeliveryDate),
			ShipmentType:          internalmessages.MTOShipmentTypeHHG,
		},
	}

	suite.T().Run("Successful POST - Integration Test", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}
		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
	})

	suite.T().Run("POST failure - 400 - invalid input, missing pickup address", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		badParams := params
		badParams.Body.PickupAddress = nil

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.T().Run("POST failure - 401- permission denied - not authenticated", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)
		unauthorizedParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest: unauthorizedReq,
			Body: &internalmessages.UpdateShipment{
				MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
				Agents:          internalmessages.MTOAgents{},
				CustomerRemarks: nil,
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					Country:        destinationAddress.Country,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
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
				RequestedPickupDate:   strfmt.Date(*mtoShipment.RequestedPickupDate),
				RequestedDeliveryDate: strfmt.Date(*mtoShipment.RequestedDeliveryDate),
				ShipmentType:          internalmessages.MTOShipmentTypeHHG,
			},
		}

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.T().Run("POST failure - 403- permission denied - wrong application", func(t *testing.T) {
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)
		unauthorizedReq = suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest: unauthorizedReq,
			Body: &internalmessages.UpdateShipment{
				MoveTaskOrderID: handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
				Agents:          internalmessages.MTOAgents{},
				CustomerRemarks: nil,
				DestinationAddress: &internalmessages.Address{
					City:           &destinationAddress.City,
					Country:        destinationAddress.Country,
					PostalCode:     &destinationAddress.PostalCode,
					State:          &destinationAddress.State,
					StreetAddress1: &destinationAddress.StreetAddress1,
					StreetAddress2: destinationAddress.StreetAddress2,
					StreetAddress3: destinationAddress.StreetAddress3,
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
				RequestedPickupDate:   strfmt.Date(*mtoShipment.RequestedPickupDate),
				RequestedDeliveryDate: strfmt.Date(*mtoShipment.RequestedDeliveryDate),
				ShipmentType:          internalmessages.MTOShipmentTypeHHG,
			},
		}

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(unauthorizedParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnauthorized{}, response)
	})

	suite.T().Run("POST failure - 404 -- not found", func(t *testing.T) {

		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := handler.Handle(badParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("POST failure - 400 -- nil body", func(t *testing.T) {
		fetcher := fetch.NewFetcher(builder)
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, nil)

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		otherParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(otherParams)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentBadRequest{}, response)
	})

	suite.T().Run("POST failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}

		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})
}
