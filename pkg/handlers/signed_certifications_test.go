package handlers

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.T()
	move, _ := testdatagen.MakeMove(suite.db)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:      *fmtUUID(move.ID),
		HTTPRequest: req,
	}

	ctx := params.HTTPRequest.Context()
	ctx = auth.PopulateAuthContext(ctx, move.Orders.ServiceMember.User.ID, "fake token")
	ctx = auth.PopulateUserModel(ctx, move.Orders.ServiceMember.User)

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	query := suite.db.Where(fmt.Sprintf("submitting_user_id='%v'", move.Orders.ServiceMember.User.ID)).Where(fmt.Sprintf("move_id='%v'", move.ID))
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
	suite.mustSave(&user2)
	move, _ := testdatagen.MakeMove(suite.db)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:      *fmtUUID(move.ID),
		HTTPRequest: req,
	}

	// Uses a different user than is on the move object
	ctx := params.HTTPRequest.Context()
	ctx = auth.PopulateAuthContext(ctx, user2.ID, "fake token")

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseForbidden(response)

	certs := []models.SignedCertification{}
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerBadMoveID() {
	t := suite.T()

	move, _ := testdatagen.MakeMove(suite.db)
	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}

	badMoveID := strfmt.UUID("3511d4d6-019d-4031-9c27-8a553e055543")
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID:      badMoveID,
		HTTPRequest: req,
	}

	// Uses a different user than is on the move object
	ctx := params.HTTPRequest.Context()
	ctx = auth.PopulateAuthContext(ctx, move.Orders.ServiceMember.User.ID, "fake token")

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseNotFound(response)

	var certs []models.SignedCertification
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}
