package handlers

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.T()
	move := testdatagen.MakeDefaultMove(suite.db)

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

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

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
	move := testdatagen.MakeDefaultMove(suite.db)

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
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.authenticateUserRequest(req, user2)

	params.HTTPRequest = req

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

	move := testdatagen.MakeDefaultMove(suite.db)
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
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.authenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.checkResponseNotFound(response)

	var certs []models.SignedCertification
	suite.db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestIndexSignedCertificationsHandler() {
	move := testdatagen.MakeDefaultMove(suite.db)

	time1 := time.Date(2018, time.January, 1, 1, 1, 1, 1, time.UTC)
	cert1 := models.SignedCertification{
		SubmittingUserID:  move.Orders.ServiceMember.UserID,
		MoveID:            move.ID,
		CertificationText: "You agree, yes?",
		Signature:         "name",
		Date:              time1,
	}
	suite.mustSave(&cert1)

	time2 := time.Date(2018, time.February, 1, 1, 1, 1, 1, time.UTC)
	cert2 := models.SignedCertification{
		SubmittingUserID:  move.Orders.ServiceMember.UserID,
		MoveID:            move.ID,
		CertificationText: "You agree, yes?",
		Signature:         "name",
		Date:              time2,
	}
	suite.mustSave(&cert2)

	req := httptest.NewRequest("GET", "/moves/id/signed_certifications", nil)
	params := certop.IndexSignedCertificationsParams{
		HTTPRequest: suite.authenticateRequest(req, move.Orders.ServiceMember),
		MoveID:      *fmtUUID(move.ID),
	}

	handler := IndexSignedCertificationsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	suite.Assertions.IsType(&certop.IndexSignedCertificationsOK{}, response)
	okResponse := response.(*certop.IndexSignedCertificationsOK)

	suite.Require().Len(okResponse.Payload, 2)
	suite.Require().Equal(time2.Month(), (time.Time)(*okResponse.Payload[0].Date).Month())

	// Now test that a limit works
	params.Limit = fmtInt64(1)

	handler = IndexSignedCertificationsHandler(NewHandlerContext(suite.db, suite.logger))
	response = handler.Handle(params)

	suite.Assertions.IsType(&certop.IndexSignedCertificationsOK{}, response)
	okResponse = response.(*certop.IndexSignedCertificationsOK)

	suite.Len(okResponse.Payload, 1)
	suite.Equal(time2.Month(), (time.Time)(*okResponse.Payload[0].Date).Month())
}
