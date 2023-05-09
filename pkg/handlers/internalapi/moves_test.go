// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
// RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
// RA: in which this would be considered a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
)

func (suite *HandlerSuite) TestPatchMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				ProvidesCloseout: true,
			},
		},
	}, nil)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
	patchPayload := internalmessages.PatchMovePayload{
		CloseoutOfficeID: &closeoutOfficeID,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		IfMatch:          etag.GenerateEtag(move.UpdatedAt),
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	closeoutOfficeUpdater := moverouter.NewCloseoutOfficeUpdater(moverouter.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())
	// And: a move is patched
	handler := PatchMoveHandler{suite.HandlerConfig(), closeoutOfficeUpdater}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.PatchMoveOK{}, response)
	okResponse := response.(*moveop.PatchMoveOK)
	suite.NoError(okResponse.Payload.Validate(strfmt.Default))

	suite.Equal(transportationOffice.ID.String(), okResponse.Payload.CloseoutOffice.ID.String())
	suite.Equal(transportationOffice.Name, *okResponse.Payload.CloseoutOffice.Name)

	suite.Equal(transportationOffice.Address.ID.String(), okResponse.Payload.CloseoutOffice.Address.ID.String())
}

func (suite *HandlerSuite) TestPatchMoveHandlerWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	// And: another logged in user
	anotherUser := factory.BuildServiceMember(suite.DB(), nil, nil)

	// And: the context contains a different user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	closeoutOfficeID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	patchPayload := internalmessages.PatchMovePayload{
		CloseoutOfficeID: &closeoutOfficeID,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		IfMatch:          etag.GenerateEtag(move.UpdatedAt),
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	closeoutOfficeUpdater := moverouter.NewCloseoutOfficeUpdater(moverouter.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())
	handler := PatchMoveHandler{suite.HandlerConfig(), closeoutOfficeUpdater}
	response := handler.Handle(params)

	suite.IsType(&moveop.PatchMoveForbidden{}, response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoMove() {
	// Given: a logged in user and no Move
	user := factory.BuildServiceMember(suite.DB(), nil, nil)

	moveUUID := uuid.Must(uuid.NewV4())

	// And: the context contains a logged in user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, user)

	closeoutOfficeID := strfmt.UUID(uuid.Must(uuid.NewV4()).String())
	patchPayload := internalmessages.PatchMovePayload{
		CloseoutOfficeID: &closeoutOfficeID,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		IfMatch:          "fake",
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	closeoutOfficeUpdater := moverouter.NewCloseoutOfficeUpdater(moverouter.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())
	handler := PatchMoveHandler{suite.HandlerConfig(), closeoutOfficeUpdater}
	response := handler.Handle(params)

	suite.IsType(&moveop.PatchMoveNotFound{}, response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerCloseoutOfficeNotFound() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	// TransportationOffice doesn't provide PPM closeout so should not be found
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
	patchPayload := internalmessages.PatchMovePayload{
		CloseoutOfficeID: &closeoutOfficeID,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		IfMatch:          etag.GenerateEtag(move.UpdatedAt),
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	closeoutOfficeUpdater := moverouter.NewCloseoutOfficeUpdater(moverouter.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())
	// And: a move is patched
	handler := PatchMoveHandler{suite.HandlerConfig(), closeoutOfficeUpdater}
	response := handler.Handle(params)

	// Then: expect a 404 status code
	suite.Assertions.IsType(&moveop.PatchMoveNotFound{}, response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerETagPreconditionFailure() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				ProvidesCloseout: true,
			},
		},
	}, nil)

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	closeoutOfficeID := strfmt.UUID(transportationOffice.ID.String())
	patchPayload := internalmessages.PatchMovePayload{
		CloseoutOfficeID: &closeoutOfficeID,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		IfMatch:          etag.GenerateEtag(time.Now()), // should not match move updatedAt value
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	closeoutOfficeUpdater := moverouter.NewCloseoutOfficeUpdater(moverouter.NewMoveFetcher(), transportationoffice.NewTransportationOfficesFetcher())
	// And: a move is patched
	handler := PatchMoveHandler{suite.HandlerConfig(), closeoutOfficeUpdater}
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.PatchMovePreconditionFailed{}, response)
}

func (suite *HandlerSuite) TestShowMoveHandler() {

	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler{suite.HandlerConfig()}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&moveop.ShowMoveOK{}, showResponse)
	okResponse := showResponse.(*moveop.ShowMoveOK)

	// And: Returned query to include our added move
	suite.Assertions.Equal(move.OrdersID.String(), okResponse.Payload.OrdersID.String())

}

func (suite *HandlerSuite) TestShowMoveWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	// And: another logged in user
	anotherUser := factory.BuildServiceMember(suite.DB(), nil, nil)

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler{suite.HandlerConfig()}
	showResponse := showHandler.Handle(showMoveParams)
	// Then: expect a forbidden response
	suite.CheckResponseForbidden(showResponse)

}

func (suite *HandlerSuite) TestSubmitMoveForApprovalHandler() {
	suite.Run("Submits ppm success", func() {
		// Given: a set of orders, a move, user and servicemember
		move := factory.BuildMove(suite.DB(), nil, nil)
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText: models.StringPointer("This is your legal message"),
			CertificationType: &certType,
			Date:              &signingDate,
			Signature:         models.StringPointer("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// When: a move is submitted
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal"))
		handler := SubmitMoveHandler{handlerConfig, moverouter.NewMoveRouter()}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)
		updatedMove, err := models.FetchMoveByMoveID(suite.DB(), move.ID)
		suite.NoError(err)

		// And: Returned query to have a submitted status
		suite.Assertions.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)
		suite.Assertions.NotNil(okResponse.Payload.SubmittedAt)

		// And: SignedCertification was created
		signedCertification := models.SignedCertification{}
		err = suite.DB().Where("move_id = $1", move.ID).First(&signedCertification)
		suite.NoError(err)
		suite.NotNil(signedCertification)

		// Test that the move was submitted within a few seconds of the current time.
		// This is better than asserting that it's not Nil, and avoids trying to mock
		// time.Now() or having to pass in a date to MoveRouter.Submit just to be able to test it.
		actualSubmittedAt := updatedMove.SubmittedAt
		suite.WithinDuration(time.Now(), *actualSubmittedAt, 2*time.Second)
	})
	suite.Run("Submits hhg shipment success", func() {
		// Given: a set of orders, a move, user and servicemember
		hhg := factory.BuildMTOShipment(suite.DB(), nil, nil)
		move := hhg.MoveTaskOrder

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText: models.StringPointer("This is your legal message"),
			CertificationType: &certType,
			Date:              &signingDate,
			Signature:         models.StringPointer("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// And: a move is submitted
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal"))
		handler := SubmitMoveHandler{handlerConfig, moverouter.NewMoveRouter()}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)

		// And: Returned query to have a submitted status
		suite.Assertions.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)

		// And: SignedCertification was created
		signedCertification := models.SignedCertification{}
		err := suite.DB().Where("move_id = $1", move.ID).First(&signedCertification)
		suite.NoError(err)
		suite.NotNil(signedCertification)
	})

}

func (suite *HandlerSuite) TestSubmitMoveForServiceCounselingHandler() {
	suite.Run("Routes to service counseling when feature flag is true", func() {
		// Given: a set of orders with an origin duty location that provides services counseling,
		// a move, user and servicemember
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					ProvidesServicesCounseling: true,
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText: models.StringPointer("This is your legal message"),
			CertificationType: &certType,
			Date:              &signingDate,
			Signature:         models.StringPointer("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// When: a move is submitted
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal"))
		handler := SubmitMoveHandler{handlerConfig, moverouter.NewMoveRouter()}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)

		updatedMove, err := models.FetchMoveByMoveID(suite.DB(), move.ID)
		suite.NoError(err)

		// Test that the move was submitted within a few seconds of the current time.
		// This is better than asserting that it's not Nil, and avoids trying to mock
		// time.Now() or having to pass in a date to sendToServiceCounseling just
		// to be able to test it.
		actualSubmittedAt := updatedMove.SubmittedAt
		suite.WithinDuration(time.Now(), *actualSubmittedAt, 2*time.Second)

		suite.Equal(models.MoveStatusNeedsServiceCounseling, updatedMove.Status)
		// And: Returned query to have a needs service counseling status
		suite.Equal(internalmessages.MoveStatusNEEDSSERVICECOUNSELING, okResponse.Payload.Status)
		suite.NotNil(okResponse.Payload.SubmittedAt)

		// And: SignedCertification was created
		signedCertification := models.SignedCertification{}
		err = suite.DB().Where("move_id = $1", move.ID).First(&signedCertification)
		suite.NoError(err)
		suite.NotNil(signedCertification)
	})
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryHandler() {
	dutyLocationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "Fort Gordon",
				City:           "Augusta",
				State:          "GA",
				PostalCode:     "30813",
				Country:        models.StringPointer("United States"),
			},
		},
	}, nil)

	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				Name: "Fort Sam Houston",
			},
		},
		{
			Model:    dutyLocationAddress,
			LinkOnly: true,
		},
	}, nil)

	rank := models.ServiceMemberRankE4
	serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Rank: &rank,
			},
		},
		{
			Model:    dutyLocation,
			LinkOnly: true,
		},
	}, nil)

	newDutyLocationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "n/a",
				City:           "San Antonio",
				State:          "TX",
				PostalCode:     "78234",
				Country:        models.StringPointer("United States"),
			},
		},
	}, nil)

	newDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				Name: "Fort Gordon",
			},
		},
		{
			Model:    newDutyLocationAddress,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				ReportByDate:     time.Date(2018, 10, 31, 0, 0, 0, 0, time.UTC),
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
		{
			Model:    dutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	path := fmt.Sprintf("/moves/%s/move_dates", move.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	moveID := strfmt.UUID(move.ID.String())
	moveDate := strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC))
	params := moveop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      moveID,
		MoveDate:    moveDate,
	}
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1125, nil)

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetPlanner(planner)

	showHandler := ShowMoveDatesSummaryHandler{handlerConfig}
	response := showHandler.Handle(params)

	suite.IsType(&moveop.ShowMoveDatesSummaryOK{}, response)
	okResponse := response.(*moveop.ShowMoveDatesSummaryOK)

	id := move.ID.String() + ":" + moveDate.String()
	suite.Equal(id, *okResponse.Payload.ID)
	suite.Equal(moveID, *okResponse.Payload.MoveID)
	suite.Equal(moveDate, *okResponse.Payload.MoveDate)

	pack := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 9, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pack, okResponse.Payload.Pack)

	pickup := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pickup, okResponse.Payload.Pickup)

	transit := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(transit, okResponse.Payload.Transit)

	delivery := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 22, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(delivery, okResponse.Payload.Delivery)

	report := []strfmt.Date{
		strfmt.Date(move.Orders.ReportByDate),
	}
	suite.Equal(report, okResponse.Payload.Report)
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryForbiddenUser() {
	// Given: a set of orders, a move, user and servicemember
	move := factory.BuildMove(suite.DB(), nil, nil)
	// And: another logged in user
	anotherUser := factory.BuildServiceMember(suite.DB(), nil, nil)

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id/", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	moveDate := strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC))
	params := moveop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		MoveDate:    moveDate,
	}
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1125, nil)

	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetPlanner(planner)

	showHandler := ShowMoveDatesSummaryHandler{handlerConfig}
	response := showHandler.Handle(params)

	// Then: expect a forbidden response
	suite.CheckResponseForbidden(response)

}

func (suite *HandlerSuite) TestSubmitAmendedOrdersHandler() {
	suite.Run("Submits move with amended orders for review", func() {
		// Given: a set of orders, a move, user and service member
		move := factory.BuildSubmittedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Document{},
				Type:  &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)
		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit_amended_orders", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := moveop.SubmitAmendedOrdersParams{
			HTTPRequest: req,
			MoveID:      strfmt.UUID(move.ID.String()),
		}
		// And: a move is submitted
		handlerConfig := suite.HandlerConfig()

		handler := SubmitAmendedOrdersHandler{handlerConfig, moverouter.NewMoveRouter()}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitAmendedOrdersOK{}, response)
		okResponse := response.(*moveop.SubmitAmendedOrdersOK)

		// And: Returned query to have a submitted status
		suite.Assertions.Equal(internalmessages.MoveStatusAPPROVALSREQUESTED, okResponse.Payload.Status)

		// And: Check status in database
		move, err := models.FetchMoveByMoveID(suite.DB(), move.ID)
		suite.NoError(err)
		suite.Assertions.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
	})
}
