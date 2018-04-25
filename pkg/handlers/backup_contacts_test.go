package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	contactop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateBackupContactHandler() {
	t := suite.T()

	serviceMember, _ := testdatagen.MakeServiceMember(suite.db)

	newContactPayload := internalmessages.CreateServiceMemberBackupContactPayload{
		Email:      fmtEmail("email@example.com"),
		Name:       swag.String("name"),
		Permission: internalmessages.BackupContactPermissionEDIT,
		Telephone:  swag.String("5555555555"),
	}
	req := httptest.NewRequest("POST", fmt.Sprintf("/service_member/%v/backup_contacts", serviceMember.ID.String()), nil)

	params := contactop.CreateServiceMemberBackupContactParams{
		CreateBackupContactPayload: &newContactPayload,
		ServiceMemberID:            *fmtUUID(serviceMember.ID),
	}

	params.HTTPRequest = suite.authenticateRequest(req, serviceMember.User)

	handler := CreateBackupContactHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	_, ok := response.(*contactop.CreateServiceMemberBackupContactCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	contacts := models.BackupContacts{}
	eagerConn := *suite.db
	eagerConn.Eager().All(&contacts)

	if len(contacts) != 1 {
		t.Errorf("Expected to find 1 result but found %v", len(contacts))
	}

	if !uuid.Equal(contacts[0].ServiceMember.ID, serviceMember.ID) {
		t.Errorf("Expected to find a backup contact for service member")
	}

}

func (suite *HandlerSuite) TestIndexBackupContactsHandler() {
	t := suite.T()

	contact, _ := testdatagen.MakeBackupContact(suite.db)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *fmtUUID(contact.ServiceMember.ID),
	}
	params.HTTPRequest = suite.authenticateRequest(req, contact.ServiceMember.User)

	handler := IndexBackupContactsHandler(NewHandlerContext(suite.db, suite.logger))
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

	contact, _ := testdatagen.MakeBackupContact(suite.db)
	otherServiceMember, _ := testdatagen.MakeServiceMember(suite.db)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *fmtUUID(contact.ServiceMember.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.authenticateRequest(req, otherServiceMember.User)

	handler := IndexBackupContactsHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*errResponse)
	code := errResponse.code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandler() {
	t := suite.T()

	contact, _ := testdatagen.MakeBackupContact(suite.db)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *fmtUUID(contact.ID),
	}
	params.HTTPRequest = suite.authenticateRequest(req, contact.ServiceMember.User)

	handler := ShowBackupContactHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse := response.(*contactop.ShowServiceMemberBackupContactOK)
	payload := okResponse.Payload

	if payload.ID.String() != contact.ID.String() {
		t.Errorf("Expected to find a particular backup contact, found something else.")
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact, _ := testdatagen.MakeBackupContact(suite.db)
	otherServiceMember, _ := testdatagen.MakeServiceMember(suite.db)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *fmtUUID(contact.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.authenticateRequest(req, otherServiceMember.User)

	handler := ShowBackupContactHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*errResponse)
	code := errResponse.code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandler() {
	t := suite.T()

	contact, _ := testdatagen.MakeBackupContact(suite.db)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      fmtEmail("otheremail@example.com"),
		Name:       swag.String("other name"),
		Permission: internalmessages.BackupContactPermissionNONE,
		Telephone:  swag.String("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *fmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	params.HTTPRequest = suite.authenticateRequest(req, contact.ServiceMember.User)

	handler := UpdateBackupContactHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	okResponse := response.(*contactop.UpdateServiceMemberBackupContactCreated)
	payload := okResponse.Payload

	if *payload.Name != "other name" {
		t.Errorf("Expected backup contact to be updated, but it wasn't.")
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact, _ := testdatagen.MakeBackupContact(suite.db)
	otherServiceMember, _ := testdatagen.MakeServiceMember(suite.db)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      fmtEmail("otheremail@example.com"),
		Name:       swag.String("other name"),
		Permission: internalmessages.BackupContactPermissionNONE,
		Telephone:  swag.String("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *fmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	// Logged in as other user
	params.HTTPRequest = suite.authenticateRequest(req, otherServiceMember.User)

	handler := UpdateBackupContactHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(params)

	errResponse := response.(*errResponse)
	code := errResponse.code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}
