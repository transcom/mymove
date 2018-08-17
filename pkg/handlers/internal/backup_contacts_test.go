package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	contactop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateBackupContactHandler() {
	t := suite.T()

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.db)

	newContactPayload := internalmessages.CreateServiceMemberBackupContactPayload{
		Email:      swag.String("email@example.com"),
		Name:       swag.String("name"),
		Permission: internalmessages.BackupContactPermissionEDIT,
		Telephone:  swag.String("5555555555"),
	}
	req := httptest.NewRequest("POST", fmt.Sprintf("/service_member/%v/backup_contacts", serviceMember.ID.String()), nil)

	params := contactop.CreateServiceMemberBackupContactParams{
		CreateBackupContactPayload: &newContactPayload,
		ServiceMemberID:            *utils.FmtUUID(serviceMember.ID),
	}

	params.HTTPRequest = suite.AuthenticateRequest(req, serviceMember)

	handler := CreateBackupContactHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*contactop.CreateServiceMemberBackupContactCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	contacts := models.BackupContacts{}
	suite.db.Q().Eager().All(&contacts)

	if len(contacts) != 1 {
		t.Errorf("Expected to find 1 result but found %v", len(contacts))
	}

	if !uuid.Equal(contacts[0].ServiceMember.ID, serviceMember.ID) {
		t.Errorf("Expected to find a backup contact for service member")
	}

}

func (suite *HandlerSuite) TestIndexBackupContactsHandler() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *utils.FmtUUID(contact.ServiceMember.ID),
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := IndexBackupContactsHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse := response.(*contactop.IndexServiceMemberBackupContactsOK)
	contacts := okResponse.Payload

	if len(contacts) != 1 {
		t.Errorf("Expected to find 1 result but found %v", len(contacts))
	}

	if contacts[0].ID.String() != contact.ID.String() {
		t.Errorf("Expected to find a particular backup contact, found something else.")
	}
}

func (suite *HandlerSuite) TestIndexBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)
	otherServiceMember := testdatagen.MakeDefaultServiceMember(suite.db)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *utils.FmtUUID(contact.ServiceMember.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := IndexBackupContactsHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*utils.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandler() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *utils.FmtUUID(contact.ID),
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := ShowBackupContactHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse := response.(*contactop.ShowServiceMemberBackupContactOK)
	payload := okResponse.Payload

	if payload.ID.String() != contact.ID.String() {
		t.Errorf("Expected to find a particular backup contact, found something else.")
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)
	otherServiceMember := testdatagen.MakeDefaultServiceMember(suite.db)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *utils.FmtUUID(contact.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := ShowBackupContactHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*utils.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandler() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      swag.String("otheremail@example.com"),
		Name:       swag.String("other name"),
		Permission: internalmessages.BackupContactPermissionNONE,
		Telephone:  swag.String("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *utils.FmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := UpdateBackupContactHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse := response.(*contactop.UpdateServiceMemberBackupContactCreated)
	payload := okResponse.Payload

	if *payload.Name != "other name" {
		t.Errorf("Expected backup contact to be updated, but it wasn't.")
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact := testdatagen.MakeDefaultBackupContact(suite.db)
	otherServiceMember := testdatagen.MakeDefaultServiceMember(suite.db)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      swag.String("otheremail@example.com"),
		Name:       swag.String("other name"),
		Permission: internalmessages.BackupContactPermissionNONE,
		Telephone:  swag.String("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *utils.FmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := UpdateBackupContactHandler(utils.NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*utils.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}
