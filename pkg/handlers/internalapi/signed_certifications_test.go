package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.T()
	move := testdatagen.MakeDefaultMove(suite.DB())

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:                           *handlers.FmtUUID(move.ID),
	}

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

	userUUID2, _ := uuid.FromString("3511d4d6-019d-4031-9c27-8a553e055543")
	user2 := models.User{
		LoginGovUUID:  userUUID2,
		LoginGovEmail: "email2@example.com",
	}
	suite.MustSave(&user2)
	move := testdatagen.MakeDefaultMove(suite.DB())

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:                           *handlers.FmtUUID(move.ID),
	}

	// Uses a different user than is on the move object
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateUserRequest(req, user2)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

	move := testdatagen.MakeDefaultMove(suite.DB())
	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.DateTime)(&date),
		Signature:         swag.String("Scruff McGruff"),
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

	handler := CreateSignedCertificationHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)

	var certs []models.SignedCertification
	suite.DB().All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandlerBadMoveID() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember

	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmPayment,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
	badMoveID := strfmt.UUID("3511d4d6-019d-4031-9c27-8a553e055543")
	params := certop.IndexSignedCertificationParams{
		MoveID: badMoveID,
	}

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandlerMismatchedUser() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmPayment,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	userUUID2, _ := uuid.FromString("3511d4d6-019d-4031-9c27-8a553e055543")
	unauthorizedUser := models.User{
		LoginGovUUID:  userUUID2,
		LoginGovEmail: "email2@example.com",
	}
	params := certop.IndexSignedCertificationParams{
		MoveID: *handlers.FmtUUID(move.ID),
	}
	suite.MustSave(&unauthorizedUser)

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateUserRequest(req, unauthorizedUser)

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseForbidden(response)
}

func (suite *HandlerSuite) TestIndexSignedCertificationHandler() {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	ppmPayment := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:                   ppm.Move.ID,
			SubmittingUserID:         sm.User.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &ppmPayment,
			CertificationText:        "LEGAL",
			Signature:                "ACCEPT",
			Date:                     testdatagen.NextValidMoveDate,
		},
	})
	params := certop.IndexSignedCertificationParams{
		MoveID: *handlers.FmtUUID(move.ID),
	}

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.AuthenticateRequest(req, sm)

	params.HTTPRequest = req

	handler := IndexSignedCertificationsHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	okResponse, ok := response.(*certop.IndexSignedCertificationOK)
	suite.True(ok)
	suite.Equal(1, len(okResponse.Payload))
	responsePayload := okResponse.Payload[0]
	suite.Equal(move.ID.String(), responsePayload.MoveID.String())
}
