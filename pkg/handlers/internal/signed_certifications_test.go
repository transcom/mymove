package internal

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateSignedCertificationHandler() {
	t := suite.parent.T()
	move := testdatagen.MakeDefaultMove(suite.parent.Db)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: *utils.FmtUUID(move.ID),
	}

	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.parent.AuthenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	_, ok := response.(*certop.CreateSignedCertificationCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	query := suite.parent.Db.Where(fmt.Sprintf("submitting_user_id='%v'", move.Orders.ServiceMember.User.ID)).Where(fmt.Sprintf("move_id='%v'", move.ID))
	certs := []models.SignedCertification{}
	query.All(&certs)

	if len(certs) != 1 {
		t.Errorf("Expected to find 1 signed certification but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerMismatchedUser() {
	t := suite.parent.T()

	userUUID2, _ := uuid.FromString("3511d4d6-019d-4031-9c27-8a553e055543")
	user2 := models.User{
		LoginGovUUID:  userUUID2,
		LoginGovEmail: "email2@example.com",
	}
	suite.parent.MustSave(&user2)
	move := testdatagen.MakeDefaultMove(suite.parent.Db)

	date := time.Now()
	certPayload := internalmessages.CreateSignedCertificationPayload{
		CertificationText: swag.String("lorem ipsum"),
		Date:              (*strfmt.Date)(&date),
		Signature:         swag.String("Scruff McGruff"),
	}
	params := certop.CreateSignedCertificationParams{
		CreateSignedCertificationPayload: &certPayload,
		MoveID: *utils.FmtUUID(move.ID),
	}

	// Uses a different user than is on the move object
	req := httptest.NewRequest("GET", "/move/id/thing", nil)
	req = suite.parent.AuthenticateUserRequest(req, user2)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.CheckResponseForbidden(response)

	certs := []models.SignedCertification{}
	suite.parent.Db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestCreateSignedCertificationHandlerBadMoveID() {
	t := suite.parent.T()

	move := testdatagen.MakeDefaultMove(suite.parent.Db)
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
	req = suite.parent.AuthenticateRequest(req, move.Orders.ServiceMember)

	params.HTTPRequest = req

	handler := CreateSignedCertificationHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.CheckResponseNotFound(response)

	var certs []models.SignedCertification
	suite.parent.Db.All(&certs)

	if len(certs) > 0 {
		t.Errorf("Expected to find no signed certifications but found %v", len(certs))
	}
}

func (suite *HandlerSuite) TestIndexSignedCertificationsHandler() {
	move := testdatagen.MakeDefaultMove(suite.parent.Db)

	time1 := time.Date(2018, time.January, 1, 1, 1, 1, 1, time.UTC)
	cert1 := models.SignedCertification{
		SubmittingUserID:  move.Orders.ServiceMember.UserID,
		MoveID:            move.ID,
		CertificationText: "You agree, yes?",
		Signature:         "name",
		Date:              time1,
	}
	suite.parent.MustSave(&cert1)

	time2 := time.Date(2018, time.February, 1, 1, 1, 1, 1, time.UTC)
	cert2 := models.SignedCertification{
		SubmittingUserID:  move.Orders.ServiceMember.UserID,
		MoveID:            move.ID,
		CertificationText: "You agree, yes?",
		Signature:         "name",
		Date:              time2,
	}
	suite.parent.MustSave(&cert2)

	req := httptest.NewRequest("GET", "/moves/id/signed_certifications", nil)
	params := certop.IndexSignedCertificationsParams{
		HTTPRequest: suite.parent.AuthenticateRequest(req, move.Orders.ServiceMember),
		MoveID:      *utils.FmtUUID(move.ID),
	}

	handler := IndexSignedCertificationsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response := handler.Handle(params)

	suite.parent.Assertions.IsType(&certop.IndexSignedCertificationsOK{}, response)
	okResponse := response.(*certop.IndexSignedCertificationsOK)

	suite.parent.Require().Len(okResponse.Payload, 2)
	suite.parent.Require().Equal(time2.Month(), (time.Time)(*okResponse.Payload[0].Date).Month())

	// Now test that a limit works
	params.Limit = utils.FmtInt64(1)

	handler = IndexSignedCertificationsHandler(utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger))
	response = handler.Handle(params)

	suite.parent.Assertions.IsType(&certop.IndexSignedCertificationsOK{}, response)
	okResponse = response.(*certop.IndexSignedCertificationsOK)

	suite.parent.Len(okResponse.Payload, 1)
	suite.parent.Equal(time2.Month(), (time.Time)(*okResponse.Payload[0].Date).Month())
}
