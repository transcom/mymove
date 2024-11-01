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

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	move "github.com/transcom/mymove/pkg/services/move"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
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
		suite.NotNil(updatedMove.SubmittedAt)
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

func (suite *HandlerSuite) TestSubmitGetAllMovesHandler() {
	suite.Run("Gets all moves belonging to a service member", func() {

		time := time.Now()
		laterTime := time.AddDate(0, 0, 1)
		// Given: A servicemember and a user
		user := factory.BuildDefaultUser(suite.DB())

		newServiceMember := factory.BuildExtendedServiceMember(suite.DB(), []factory.Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, nil)
		suite.MustSave(&newServiceMember)

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    newServiceMember,
				LinkOnly: true,
				Type:     &factory.ServiceMember,
			},
		}, nil)

		// Given: a set of orders, a move, user and service member
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			},
			{
				Model:    newServiceMember,
				LinkOnly: true,
				Type:     &factory.ServiceMember,
			},
			{
				Model: models.Move{
					CreatedAt: time,
				},
			},
		}, nil)

		move2 := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			},
			{
				Model:    newServiceMember,
				LinkOnly: true,
				Type:     &factory.ServiceMember,
			},
			{
				Model: models.Move{
					CreatedAt: laterTime,
				},
			},
		}, nil)

		// // And: the context contains the auth values
		req := httptest.NewRequest("GET", "/moves/allmoves", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := moveop.GetAllMovesParams{
			HTTPRequest:     req,
			ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
		}

		// And: a move is submitted
		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := GetAllMovesHandler{handlerConfig}
		response := handler.Handle(params)

		// // Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.GetAllMovesOK{}, response)
		okResponse := response.(*moveop.GetAllMovesOK)

		suite.Greater(len(okResponse.Payload.CurrentMove), 0)
		suite.Greater(len(okResponse.Payload.PreviousMoves), 0)
		suite.Equal(okResponse.Payload.CurrentMove[0].ID.String(), move.ID.String())
		suite.Equal(okResponse.Payload.PreviousMoves[0].ID.String(), move2.ID.String())

	})
}

func (suite *HandlerSuite) TestUploadAdditionalDocumentsHander() {
	fakeS3 := storageTest.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeS3)
	additionalDocumentsUploader := move.NewMoveAdditionalDocumentsUploader(uploadCreator)

	setupRequestAndParams := func(move models.Move) *moveop.UploadAdditionalDocumentsParams {
		endpoint := fmt.Sprintf("/moves/%v/upload_additional_documents", move.ID)
		req := httptest.NewRequest("PATCH", endpoint, nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

		params := moveop.UploadAdditionalDocumentsParams{
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

		if suite.IsType(&moveop.UploadAdditionalDocumentsCreated{}, response) {
			payload := response.(*moveop.UploadAdditionalDocumentsCreated).Payload

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

		suite.IsType(&moveop.UploadAdditionalDocumentsInternalServerError{}, response)

	})
}
