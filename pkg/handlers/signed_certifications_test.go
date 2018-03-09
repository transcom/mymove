package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gorilla/context"
	"github.com/satori/go.uuid"

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
		SelectedMoveType: swag.String("HHG"),
	}
	suite.mustSave(&move)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: *fmtUUID(move.ID),
	}

	context.Set(params.HTTPRequest, "user_id", user.ID.String())

	handler := NewCreateSignedCertificationHandler(suite.db, suite.logger)
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
		SelectedMoveType: swag.String("HHG"),
	}
	suite.mustSave(&move)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: *fmtUUID(move.ID),
	}

	handler := NewCreateSignedCertificationHandler(suite.db, suite.logger)
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
		SelectedMoveType: swag.String("HHG"),
	}
	suite.mustSave(&move)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: *fmtUUID(move.ID),
	}

	// Uses a different user than is on the move object
	context.Set(params.HTTPRequest, "user_id", user2.ID.String())

	handler := NewCreateSignedCertificationHandler(suite.db, suite.logger)
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
		SelectedMoveType: swag.String("HHG"),
	}
	suite.mustSave(&move)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}

	badMoveID := strfmt.UUID("3511d4d6-019d-4031-9c27-8a553e055543")
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: badMoveID,
	}

	// Uses a different user than is on the move object
	context.Set(params.HTTPRequest, "user_id", user.ID.String())

	handler := NewCreateSignedCertificationHandler(suite.db, suite.logger)
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
