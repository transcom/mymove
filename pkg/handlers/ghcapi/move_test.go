package ghcapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	moveops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	movelocker "github.com/transcom/mymove/pkg/services/lock_move"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservice "github.com/transcom/mymove/pkg/services/move"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestGetMoveHandler() {
	swaggerTimeFormat := "2006-01-02T15:04:05.99Z07:00"
	availableToPrimeAt := time.Now()
	submittedAt := availableToPrimeAt.Add(-1 * time.Hour)

	var move models.Move
	var requestUser models.OfficeUser
	setupTestData := func() {
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:             models.MoveStatusAPPROVED,
					AvailableToPrimeAt: &availableToPrimeAt,
					ApprovedAt:         &availableToPrimeAt,
					SubmittedAt:        &submittedAt,
				},
			},
		}, nil)
		requestUser = factory.BuildOfficeUser(nil, nil, nil)
	}

	suite.Run("Successful move fetch", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}
		mockLocker := movelocker.NewMoveLocker()

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
			MoveLocker:    mockLocker,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&move, nil)

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser.User)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)
		payload := response.(*moveops.GetMoveOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), payload.ID.String())
		suite.Equal(move.AvailableToPrimeAt.Format(swaggerTimeFormat), time.Time(*payload.AvailableToPrimeAt).Format(swaggerTimeFormat))
		suite.Equal(move.ApprovedAt.Format(swaggerTimeFormat), time.Time(*payload.AvailableToPrimeAt).Format(swaggerTimeFormat))
		suite.Equal(move.ContractorID.String(), payload.ContractorID.String())
		suite.Equal(move.Locator, payload.Locator)
		suite.Equal(move.OrdersID.String(), payload.OrdersID.String())
		suite.Equal(move.ReferenceID, payload.ReferenceID)
		suite.Equal(string(move.Status), string(payload.Status))
		suite.Equal(move.CreatedAt.Format(swaggerTimeFormat), time.Time(payload.CreatedAt).Format(swaggerTimeFormat))
		suite.Equal(move.SubmittedAt.Format(swaggerTimeFormat), time.Time(*payload.SubmittedAt).Format(swaggerTimeFormat))
		suite.Equal(move.UpdatedAt.Format(swaggerTimeFormat), time.Time(payload.UpdatedAt).Format(swaggerTimeFormat))
		suite.Nil(payload.CloseoutOffice)
	})

	suite.Run("Successful move with a saved transportation office", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: true,
				},
			},
		}, nil)

		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:      models.MoveStatusSUBMITTED,
					SubmittedAt: &submittedAt,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		}, nil)
		mockFetcher := mocks.MoveFetcher{}
		mockLocker := movelocker.NewMoveLocker()
		requestOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateOfficeRequest(req, requestOfficeUser)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
			MoveLocker:    mockLocker,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&move, nil)

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveOK{}, response)
		payload := response.(*moveops.GetMoveOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(transportationOffice.ID.String(), payload.CloseoutOfficeID.String())
		suite.Equal(transportationOffice.ID.String(), payload.CloseoutOffice.ID.String())
		suite.Equal(transportationOffice.AddressID.String(), payload.CloseoutOffice.Address.ID.String())

	})

	suite.Run("Unsuccessful move fetch - empty string bad request", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
		}
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser.User)

		// Validate incoming payload: no body to validate

		response := handler.Handle(moveops.GetMoveParams{HTTPRequest: req, Locator: ""})
		suite.IsType(&moveops.GetMoveBadRequest{}, response)
		payload := response.(*moveops.GetMoveBadRequest).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move fetch - locator not found", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}
		mockLocker := movelocker.NewMoveLocker()

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
			MoveLocker:    mockLocker,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.NotFoundError{})
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser.User)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
		payload := response.(*moveops.GetMoveNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move fetch - internal server error", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}
		mockLocker := movelocker.NewMoveLocker()

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
			MoveLocker:    mockLocker,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.QueryError{})

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser.User)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveInternalServerError{}, response)
		payload := response.(*moveops.GetMoveInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful move fetch - invalid privileges", func() {
		setupTestData()
		mockFetcher := mocks.MoveFetcher{}
		mockLocker := movelocker.NewMoveLocker()

		handler := GetMoveHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveFetcher:   &mockFetcher,
			MoveLocker:    mockLocker,
		}

		mockFetcher.On("FetchMove",
			mock.AnythingOfType("*appcontext.appContext"),
			move.Locator,
			mock.Anything,
		).Return(&models.Move{}, apperror.NotFoundError{})

		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser.User)
		params := moveops.GetMoveParams{
			HTTPRequest: req,
			Locator:     move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&moveops.GetMoveNotFound{}, response)
		payload := response.(*moveops.GetMoveNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestSearchMovesHandler() {

	var requestUser models.User
	setupTestData := func() *http.Request {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		return req

	}

	/* setupGblocTestData is a helper function to set up test data for the search moves handler specifically for testing GBLOCs.
	 * returns a non-PPM move and a PPM move with different destination postal codes and GBLOCs. */
	setupGblocTestData := func() (*models.Move, *models.Move, *models.Move) {
		// ZIPs takes a GBLOC and returns a ZIP
		ZIPs := map[string]string{
			"AGFM": "62225",
			"KKFA": "90210",
			"BGNC": "47712",
			"CLPK": "33009",
		}

		for k, v := range ZIPs {
			factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), v, k)
		}

		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		defaultPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		addressAGFM := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     ZIPs["AGFM"],
				},
			},
		}, nil)
		addressKKFA := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     ZIPs["KKFA"],
				},
			},
		}, nil)
		addressBGNC := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     ZIPs["BGNC"],
				},
			},
		}, nil)
		addressCLPK := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main St",
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     ZIPs["CLPK"],
				},
			},
		}, nil)

		destDutyLocationAGFM := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:      "Test AGFM",
					AddressID: addressAGFM.ID,
				},
			},
		}, nil)
		destDutyLocationCLPK := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:      "Test CLPK",
					AddressID: addressCLPK.ID,
				},
			},
		}, nil)

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					ServiceMemberID:   serviceMember.ID,
					NewDutyLocationID: destDutyLocationAGFM.ID,
					DestinationGBLOC:  handlers.FmtString("AGFM"),
					HasDependents:     false,
					SpouseHasProGear:  false,
					OrdersType:        "PERMANENT_CHANGE_OF_STATION",
					OrdersTypeDetail:  nil,
				},
			},
		}, nil)

		orderWithShipment := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					ServiceMemberID:   serviceMember.ID,
					NewDutyLocationID: destDutyLocationAGFM.ID,
					DestinationGBLOC:  handlers.FmtString("AGFM"),
					HasDependents:     false,
					SpouseHasProGear:  false,
					OrdersType:        "PERMANENT_CHANGE_OF_STATION",
					OrdersTypeDetail:  nil,
				},
			},
		}, nil)

		orderWithShipmentPPM := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					ServiceMemberID:   serviceMember.ID,
					NewDutyLocationID: destDutyLocationCLPK.ID,
					DestinationGBLOC:  handlers.FmtString("CLPK"),
					HasDependents:     false,
					SpouseHasProGear:  false,
					OrdersType:        "PERMANENT_CHANGE_OF_STATION",
					OrdersTypeDetail:  nil,
				},
			},
		}, nil)

		moveWithoutShipment := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					OrdersID: order.ID,
				},
			},
		}, nil)

		moveWithShipment := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					DestinationAddressID: &addressKKFA.ID,
					PickupAddressID:      &defaultPickupAddress.ID,
					Status:               models.MTOShipmentStatusSubmitted,
					ShipmentType:         models.MTOShipmentTypeHHG,
				},
			},
			{
				Model: models.Move{
					OrdersID: orderWithShipment.ID,
				},
			},
		}, nil)

		moveWithShipmentPPM := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status:   models.MoveStatusSUBMITTED,
					OrdersID: orderWithShipmentPPM.ID,
				},
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusSubmitted,
					DestinationAddressID: &addressBGNC.ID,
					PickupAddressID:      &defaultPickupAddress.ID,
				},
			},
			{
				Model: models.PPMShipment{
					Status:               models.PPMShipmentStatusSubmitted,
					DestinationAddressID: &addressBGNC.ID,
					PickupAddressID:      &defaultPickupAddress.ID,
				},
			},
		}, nil)

		return &moveWithoutShipment, &moveWithShipment, &moveWithShipmentPPM
	}

	suite.Run("Successful move search by locator", func() {
		req := setupTestData()
		move := factory.BuildMove(suite.DB(), nil, nil)
		moves := models.Moves{move}

		mockSearcher := mocks.MoveSearcher{}

		mockUnlocker := movelocker.NewMoveUnlocker()
		handler := SearchMovesHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveSearcher:  &mockSearcher,
			MoveUnlocker:  mockUnlocker,
		}
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.Locator == move.Locator
			}),
		).Return(moves, 1, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &move.Locator,
				Edipi:   nil,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload := response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		payloadMove := *(*payload).SearchMoves[0]
		suite.Equal(move.ID.String(), payloadMove.ID.String())
		suite.Equal(*move.Orders.ServiceMember.Edipi, *payloadMove.Edipi)
		suite.Equal(move.Orders.NewDutyLocation.Address.PostalCode, payloadMove.DestinationPostalCode)
		suite.Equal(move.Orders.OriginDutyLocation.Address.PostalCode, payloadMove.OriginDutyLocationPostalCode)
		suite.Equal(ghcmessages.MoveStatusDRAFT, payloadMove.Status)
		suite.Equal("ARMY", payloadMove.Branch)
		suite.Equal((*string)(nil), payloadMove.Emplid)
		suite.Equal(int64(0), payloadMove.ShipmentsCount)
		suite.NotEmpty(payloadMove.FirstName)
		suite.NotEmpty(payloadMove.LastName)
		suite.NotEmpty(payloadMove.OriginGBLOC)
		suite.NotEmpty(payloadMove.DestinationGBLOC)
	})

	suite.Run("Successful move search by DoD ID", func() {
		req := setupTestData()
		move := factory.BuildMove(suite.DB(), nil, nil)
		moves := models.Moves{move}

		mockSearcher := mocks.MoveSearcher{}

		handler := SearchMovesHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveSearcher:  &mockSearcher,
		}
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.DodID == *move.Orders.ServiceMember.Edipi &&
					params.Locator == nil &&
					params.CustomerName == nil
			}),
		).Return(moves, 1, nil)

		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: nil,
				Edipi:   move.Orders.ServiceMember.Edipi,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload := response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), (*payload).SearchMoves[0].ID.String())
	})

	suite.Run("Destination Postal Code and GBLOC is correct for different shipment types", func() {
		req := setupTestData()
		move, moveWithShipment, moveWithShipmentPPM := setupGblocTestData()

		moves := models.Moves{*move}
		movesWithShipment := models.Moves{*moveWithShipment}
		movesWithShipmentPPM := models.Moves{*moveWithShipmentPPM}

		// Mocks
		mockSearcher := mocks.MoveSearcher{}
		mockUnlocker := movelocker.NewMoveUnlocker()
		handler := SearchMovesHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveSearcher:  &mockSearcher,
			MoveUnlocker:  mockUnlocker,
		}

		// Set Mock Search settings for move without Shipment
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.Locator == move.Locator
			}),
		).Return(moves, 1, nil)

		// Move search params without Shipment
		params := moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &move.Locator,
				Edipi:   nil,
			},
		}

		// Validate incoming payload non-PPM
		suite.NoError(params.Body.Validate(strfmt.Default))

		// set and validate response and payload
		response := handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload := response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload without shipment
		suite.NoError(payload.Validate(strfmt.Default))

		var moveDestinationAddress *models.Address
		var moveDestinationGBLOC string
		var err error

		// Get destination postal code and GBLOC based on business logic
		moveDestinationAddress, err = move.GetDestinationAddress(suite.DB())
		suite.NoError(err)
		moveDestinationGBLOC, err = move.GetDestinationGBLOC(suite.DB())
		suite.NoError(err)

		suite.Equal(moveDestinationAddress.PostalCode, "62225")
		suite.Equal(ghcmessages.GBLOC(moveDestinationGBLOC), ghcmessages.GBLOC("AGFM"))

		// Set Mock Search settings for move with MTO Shipment
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.Locator == moveWithShipment.Locator
			}),
		).Return(movesWithShipment, 1, nil)

		// Move search params with MTO Shipment
		params = moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &moveWithShipment.Locator,
				Edipi:   nil,
			},
		}

		// Validate incoming payload with shipment
		suite.NoError(params.Body.Validate(strfmt.Default))

		// reset and validate response and payload
		response = handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload = response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload with shipment
		suite.NoError(payload.Validate(strfmt.Default))

		// Get destination postal code and GBLOC based on business logic
		moveDestinationAddress, err = moveWithShipment.GetDestinationAddress(suite.DB())
		suite.NoError(err)
		moveDestinationGBLOC, err = moveWithShipment.GetDestinationGBLOC(suite.DB())
		suite.NoError(err)

		suite.Equal(moveDestinationAddress.PostalCode, "90210")
		suite.Equal(ghcmessages.GBLOC(moveDestinationGBLOC), ghcmessages.GBLOC("KKFA"))

		// Set Mock Search settings for move with PPM Shipment
		mockSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(params *services.SearchMovesParams) bool {
				return *params.Locator == moveWithShipmentPPM.Locator
			}),
		).Return(movesWithShipmentPPM, 1, nil)

		// Move search params with PPM Shipment
		params = moveops.SearchMovesParams{
			HTTPRequest: req,
			Body: moveops.SearchMovesBody{
				Locator: &moveWithShipmentPPM.Locator,
				Edipi:   nil,
			},
		}

		// Validate incoming payload with PPM shipment
		suite.NoError(params.Body.Validate(strfmt.Default))

		// reset and validate response and payload
		response = handler.Handle(params)
		suite.IsType(&moveops.SearchMovesOK{}, response)
		payload = response.(*moveops.SearchMovesOK).Payload

		// Validate outgoing payload non-PPM
		suite.NoError(payload.Validate(strfmt.Default))

		// Get destination postal code and GBLOC based on business logic
		moveDestinationAddress, err = moveWithShipmentPPM.GetDestinationAddress(suite.DB())
		suite.NoError(err)
		moveDestinationGBLOC, err = moveWithShipmentPPM.GetDestinationGBLOC(suite.DB())
		suite.NoError(err)

		suite.Equal(moveDestinationAddress.PostalCode, payload.SearchMoves[0].DestinationPostalCode)
		suite.Equal(ghcmessages.GBLOC(moveDestinationGBLOC), payload.SearchMoves[0].DestinationGBLOC)
	})
}

func (suite *HandlerSuite) TestSetFinancialReviewFlagHandler() {
	var move models.Move
	var requestUser models.User
	setupTestData := func() (*http.Request, models.Move) {
		move = factory.BuildMove(suite.DB(), nil, nil)
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", "/move/#{move.locator}", nil)
		req = suite.AuthenticateUserRequest(req, requestUser)
		return req, move
	}
	defaultRemarks := "destination address is on the moon"
	fakeEtag := ""

	suite.Run("Successful flag setting to true", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			move.ID,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&move, nil)

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: models.BoolPointer(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagOK{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful flag - missing remarks", func() {
		req, move := setupTestData()
		paramsNilRemarks := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       nil,
				FlagForReview: models.BoolPointer(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}

		// Validate incoming payload
		suite.NoError(paramsNilRemarks.Body.Validate(strfmt.Default))

		response := handler.Handle(paramsNilRemarks)
		suite.IsType(&moveops.SetFinancialReviewFlagUnprocessableEntity{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Unsuccessful flag - move not found", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.NotFoundError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: models.BoolPointer(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagNotFound{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful flag - internal server error", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.QueryError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: models.BoolPointer(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagInternalServerError{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Unsuccessful flag - bad etag", func() {
		req, move := setupTestData()
		mockFlagSetter := mocks.MoveFinancialReviewFlagSetter{}
		handler := SetFinancialReviewFlagHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveFinancialReviewFlagSetter: &mockFlagSetter,
		}
		mockFlagSetter.On("SetFinancialReviewFlag",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.AnythingOfType("bool"),
			&defaultRemarks,
		).Return(&models.Move{}, apperror.PreconditionFailedError{})

		params := moveops.SetFinancialReviewFlagParams{
			HTTPRequest: req,
			IfMatch:     &fakeEtag,
			Body: moveops.SetFinancialReviewFlagBody{
				Remarks:       &defaultRemarks,
				FlagForReview: models.BoolPointer(true),
			},
			MoveID: *handlers.FmtUUID(move.ID),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.SetFinancialReviewFlagPreconditionFailed{}, response)
		payload := response.(*moveops.SetFinancialReviewFlagPreconditionFailed).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestUpdateMoveCloseoutOfficeHandler() {
	var move models.Move
	var requestUser models.OfficeUser
	var transportationOffice models.TransportationOffice

	closeoutOfficeUpdater := moveservice.NewCloseoutOfficeUpdater(moveservice.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())

	setupTestData := func() (*http.Request, models.Move, models.TransportationOffice) {
		move = factory.BuildMove(suite.DB(), nil, nil)
		requestUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		transportationOffice = factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: true,
				},
			},
		}, nil)

		req := httptest.NewRequest("GET", "/move/#{move.locator}/closeout-office", nil)
		req = suite.AuthenticateOfficeRequest(req, requestUser)
		return req, move, transportationOffice
	}

	suite.Run("Successful update of closeout office", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeOK{}, response)
		payload := response.(*moveops.UpdateCloseoutOfficeOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(closeoutOfficeID, *payload.CloseoutOfficeID)
		suite.Equal(closeoutOfficeID, *payload.CloseoutOffice.ID)
		suite.Equal(transportationOffice.AddressID.String(), payload.CloseoutOffice.Address.ID.String())
	})

	suite.Run("Unsuccessful move not found", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: "ABC123",
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeNotFound{}, response)
	})

	suite.Run("Unsuccessful closeout office not found", func() {
		transportationOfficeNonCloseout := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: false,
				},
			},
		}, nil)

		req, move, _ := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOfficeNonCloseout.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficeNotFound{}, response)
	})

	suite.Run("Unsuccessful eTag does not match", func() {
		req, move, transportationOffice := setupTestData()
		handler := UpdateMoveCloseoutOfficeHandler{
			HandlerConfig:             suite.HandlerConfig(),
			MoveCloseoutOfficeUpdater: closeoutOfficeUpdater,
		}

		closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
		params := moveops.UpdateCloseoutOfficeParams{
			HTTPRequest: req,
			IfMatch:     etag.GenerateEtag(time.Now()),
			Body: moveops.UpdateCloseoutOfficeBody{
				CloseoutOfficeID: &closeoutOfficeID,
			},
			Locator: move.Locator,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateCloseoutOfficePreconditionFailed{}, response)
	})
}

func (suite *HandlerSuite) TestUploadAdditionalDocumentsHander() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeS3)
	additionalDocumentsUploader := moveservice.NewMoveAdditionalDocumentsUploader(uploadCreator)

	setupRequestAndParams := func(move models.Move) *moveops.UploadAdditionalDocumentsParams {
		endpoint := fmt.Sprintf("/moves/%v/upload_additional_documents", move.ID)
		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := moveops.UploadAdditionalDocumentsParams{
			HTTPRequest: req,
			File:        suite.Fixture("filled-out-orders.pdf"),
			MoveID:      *handlers.FmtUUID(move.ID),
		}

		return &params
	}

	setupHandler := func() UploadAdditionalDocumentsHandler {
		return UploadAdditionalDocumentsHandler{
			suite.createS3HandlerConfig(),
			additionalDocumentsUploader,
		}
	}

	suite.Run("Returns 201 if the additional documents uploaded successfully", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		params := setupRequestAndParams(move)
		handler := setupHandler()
		response := handler.Handle(*params)

		if suite.IsType(&moveops.UploadAdditionalDocumentsCreated{}, response) {
			payload := response.(*moveops.UploadAdditionalDocumentsCreated).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.NotEqual("", string(payload.ID))
			suite.Equal("filled-out-orders.pdf", payload.Filename)
			suite.Equal(uploader.FileTypePDF, payload.ContentType)
			suite.NotEqual("", string(payload.URL))
		}
	})

	suite.Run("Returns 400 - Bad Request if there is an issue with the file being uploaded", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		params := setupRequestAndParams(move)
		params.File = factory.FixtureOpen("empty.pdf")

		handler := setupHandler()
		response := handler.Handle(*params)

		suite.IsType(&moveops.UploadAdditionalDocumentsInternalServerError{}, response)

	})
}

func (suite *HandlerSuite) TestUpdateAssignedOfficeUserHandler() {
	var move models.Move
	var assignedUser models.OfficeUser

	assignedOfficeUserUpdater := moveservice.NewAssignedOfficeUserUpdater(moveservice.NewMoveFetcher())
	officeUserFetcher := officeuser.NewOfficeUserFetcherPop()

	setupTestData := func() (*http.Request, UpdateAssignedOfficeUserHandler, models.Move, models.OfficeUser) {
		move = factory.BuildMove(suite.DB(), nil, nil)
		assignedUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		req := httptest.NewRequest("GET", "/moves/{moveID}/assignOfficeUser", nil)
		req = suite.AuthenticateOfficeRequest(req, assignedUser)

		handler := UpdateAssignedOfficeUserHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveAssignedOfficeUserUpdater: assignedOfficeUserUpdater,
			officeUserFetcherPop:          officeUserFetcher,
		}
		return req, handler, move, assignedUser
	}

	suite.Run("Successful update of a move's SC", func() {
		req, handler, move, officeUser := setupTestData()

		officeUserID := strfmt.UUID(officeUser.ID.String())
		moveID := strfmt.UUID(move.ID.String())
		roleType := string(roles.RoleTypeServicesCounselor)
		params := moveops.UpdateAssignedOfficeUserParams{
			HTTPRequest: req,
			Body: &ghcmessages.AssignOfficeUserBody{
				OfficeUserID: &officeUserID,
				RoleType:     &roleType,
			},
			MoveID: moveID,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateAssignedOfficeUserOK{}, response)
		payload := response.(*moveops.UpdateAssignedOfficeUserOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(officeUserID, payload.SCAssignedUser.OfficeUserID)
	})
	suite.Run("Successful update of a move's TOO", func() {
		req, handler, move, officeUser := setupTestData()

		officeUserID := strfmt.UUID(officeUser.ID.String())
		moveID := strfmt.UUID(move.ID.String())
		roleType := string(roles.RoleTypeTOO)
		params := moveops.UpdateAssignedOfficeUserParams{
			HTTPRequest: req,
			Body: &ghcmessages.AssignOfficeUserBody{
				OfficeUserID: &officeUserID,
				RoleType:     &roleType,
			},
			MoveID: moveID,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateAssignedOfficeUserOK{}, response)
		payload := response.(*moveops.UpdateAssignedOfficeUserOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(officeUserID, payload.TOOAssignedUser.OfficeUserID)
	})
	suite.Run("Successful update of a move's TIO", func() {
		req, handler, move, officeUser := setupTestData()

		officeUserID := strfmt.UUID(officeUser.ID.String())
		moveID := strfmt.UUID(move.ID.String())
		roleType := string(roles.RoleTypeTIO)
		params := moveops.UpdateAssignedOfficeUserParams{
			HTTPRequest: req,
			Body: &ghcmessages.AssignOfficeUserBody{
				OfficeUserID: &officeUserID,
				RoleType:     &roleType,
			},
			MoveID: moveID,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&moveops.UpdateAssignedOfficeUserOK{}, response)
		payload := response.(*moveops.UpdateAssignedOfficeUserOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(officeUserID, payload.TIOAssignedUser.OfficeUserID)
	})
	suite.Run("Successful unassign of an office user", func() {
		move = factory.BuildMove(suite.DB(), nil, nil)

		req := httptest.NewRequest("GET", "/moves/{moveID}/unassignOfficeUser", nil)
		req = suite.AuthenticateOfficeRequest(req, assignedUser)

		handler := DeleteAssignedOfficeUserHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			MoveAssignedOfficeUserUpdater: assignedOfficeUserUpdater,
		}

		moveID := strfmt.UUID(move.ID.String())
		roleType := string(roles.RoleTypeTIO)
		params := moveops.DeleteAssignedOfficeUserParams{
			HTTPRequest: req,
			Body: moveops.DeleteAssignedOfficeUserBody{
				RoleType: &roleType,
			},
			MoveID: moveID,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&moveops.DeleteAssignedOfficeUserOK{}, response)
		payload := response.(*moveops.DeleteAssignedOfficeUserOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Nil(payload.TIOAssignedUser)
	})
}

func (suite *HandlerSuite) TestCheckForLockedMovesAndUnlockHandler() {
	var validOfficeUser models.OfficeUser
	var move models.Move

	mockLocker := movelocker.NewMoveLocker()
	setupLockedMove := func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           validOfficeUser.User.Roles,
			OfficeUserID:    validOfficeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
			UserID:          validOfficeUser.ID,
		})

		validOfficeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					ID: validOfficeUser.ID,
				},
			},
		}, nil)

		move, err := mockLocker.LockMove(appCtx, &move, validOfficeUser.ID)

		suite.NoError(err)
		suite.NotNil(move.LockedByOfficeUserID)
	}

	setupTestData := func() (*http.Request, CheckForLockedMovesAndUnlockHandler) {
		req := httptest.NewRequest("GET", "/moves/{officeUserID}/CheckForLockedMovesAndUnlock", nil)

		handler := CheckForLockedMovesAndUnlockHandler{
			HandlerConfig: suite.HandlerConfig(),
			MoveUnlocker:  movelocker.NewMoveUnlocker(),
		}

		return req, handler
	}
	suite.PreloadData(setupLockedMove)

	suite.Run("Successful unlocking of move", func() {
		req, handler := setupTestData()

		expectedPayloadMessage := "Successfully unlocked all move(s) for current office user"

		officeUserID := strfmt.UUID(validOfficeUser.ID.String())
		params := moveops.CheckForLockedMovesAndUnlockParams{
			HTTPRequest:  req,
			OfficeUserID: officeUserID,
		}

		handler.Handle(params)
		suite.NotNil(move)

		response := handler.Handle(params)
		suite.IsType(&moveops.CheckForLockedMovesAndUnlockOK{}, response)
		payload := response.(*moveops.CheckForLockedMovesAndUnlockOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))

		actualMessage := payload.SuccessMessage
		suite.Equal(expectedPayloadMessage, actualMessage)
	})

	suite.Run("Unsucceful unlocking of move - nil officerUserId", func() {
		req, handler := setupTestData()

		invalidOfficeUserID := strfmt.UUID(uuid.Nil.String())
		params := moveops.CheckForLockedMovesAndUnlockParams{
			HTTPRequest:  req,
			OfficeUserID: invalidOfficeUserID,
		}

		handler.Handle(params)
		response := handler.Handle(params)
		suite.IsType(&moveops.CheckForLockedMovesAndUnlockInternalServerError{}, response)
		payload := response.(*moveops.CheckForLockedMovesAndUnlockInternalServerError).Payload
		suite.Nil(payload)
	})
}
