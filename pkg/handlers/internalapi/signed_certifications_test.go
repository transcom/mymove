// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.T()
	move := factory.BuildMove(suite.DB(), nil, nil)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: models.StringPointer("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         models.StringPointer("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:                           *handlers.FmtUUID(move.ID),
	}

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	query := suite.DB().Where("submitting_user_id = ?", move.Orders.ServiceMember.User.ID).
		Where("move_id = ?", move.ID)
	certs := []models.SignedCertification{}
	query.All(&certs)

	if len(certs) != 1 {
		t.Errorf("Expected to find 1 signed certification but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerMismatchedUser() {
	t := suite.T()

	userUUID2 := "3511d4d6-019d-4031-9c27-8a553e055543"
	user2 := models.User{
		OktaID:    userUUID2,
		OktaEmail: "email2@example.com",
	}
	suite.MustSave(&user2)
	move := factory.BuildMove(suite.DB(), nil, nil)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: models.StringPointer("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         models.StringPointer("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:                           *handlers.FmtUUID(move.ID),
	}

	// Uses a different user than is on the move object
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateUserRequest(req, user2)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.CheckResponseForbidden(response)

	certs := []models.SignedCertification{}
	suite.DB().All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerBadMoveID() {
	t := suite.T()

	move := factory.BuildMove(suite.DB(), nil, nil)
	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: models.StringPointer("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         models.StringPointer("Scruff McGruff"),
	}

	badMoveID := strfmt.UUID("3511d4d6-019d-4031-9c27-8a553e055543")
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:                           badMoveID,
	}

	// Uses a different user than is on the move object
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)

	var certs []models.SignedCertification
	suite.DB().All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandlerBadMoveID() {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move.MTOShipments = append(move.MTOShipments, mtoShipment)

	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmPayment,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
	badMoveID := strfmt.UUID("3511d4d6-019d-4031-9c27-8a553e055543")
	params := certop.IndexSignedCertificationParams{
		MoveID: badMoveID,
	}

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandlerMismatchedUser() {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move.MTOShipments = append(move.MTOShipments, mtoShipment)
	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmPayment,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)
	userUUID2 := "3511d4d6-019d-4031-9c27-8a553e055543"
	unauthorizedUser := models.User{
		OktaID:    userUUID2,
		OktaEmail: "email2@example.com",
	}
	params := certop.IndexSignedCertificationParams{
		MoveID: *handlers.FmtUUID(move.ID),
	}
	suite.MustSave(&unauthorizedUser)

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateUserRequest(req, unauthorizedUser)

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.CheckResponseForbidden(response)
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandler() {
	ppm := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: ppm,
		},
	}, nil)
	move.MTOShipments = append(move.MTOShipments, mtoShipment)
	sm := move.Orders.ServiceMember
	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	factory.BuildSignedCertification(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				CertificationType: &ppmPayment,
				CertificationText: "LEGAL",
				Signature:         "ACCEPT",
				Date:              testdatagen.NextValidMoveDate,
			},
		},
	}, nil)
	params := certop.IndexSignedCertificationParams{
		MoveID: *handlers.FmtUUID(move.ID),
	}

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, sm)

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	okResponse, ok := response.(*certop.IndexSignedCertificationOK)
	suite.True(ok)
	suite.Equal(1, len(okResponse.Payload))
	responsePayload := okResponse.Payload[0]
	suite.Equal(move.ID.String(), responsePayload.MoveID.String())
}
