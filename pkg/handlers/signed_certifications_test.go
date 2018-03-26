package handlers

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth/context"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.T()

	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
	}
	suite.mustSave(&move)

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
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	query := suite.db.Where(fmt.Sprintf("submitting_user_id='%v'", user.ID)).Where(fmt.Sprintf("move_id='%v'", move.ID))
	certs := []models.SignedCertification{}
	query.All(&certs)

	if len(certs) != 1 {
		t.Errorf("Expected to find 1 signed certification but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerNoUserID() {
	t := suite.T()

	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
	}
	suite.mustSave(&move)

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

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationUnauthorized)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	certs := []models.SignedCertification{}
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerMismatchedUser() {
	t := suite.T()

	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	userUUID2, _ := uuid.FromString("3511d4d6-019d-4031-9c27-8a553e055543")
	user2 := models.User{
		LoginGovUUID:  userUUID2,
		LoginGovEmail: "email2@example.com",
	}
	suite.mustSave(&user2)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
	}
	suite.mustSave(&move)

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
	ctx = context.PopulateAuthContext(ctx, user2.ID, "fake token")

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationForbidden)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	certs := []models.SignedCertification{}
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerBadMoveID() {
	t := suite.T()

	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	move := models.Move{
		UserID:           user.ID,
		SelectedMoveType: "HHG",
	}
	suite.mustSave(&move)

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
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")

	params.HTTPRequest = params.HTTPRequest.WithContext(ctx)

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationNotFound)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	certs := []models.SignedCertification{}
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}
