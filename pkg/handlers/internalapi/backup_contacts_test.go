package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	contactop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/backup_contacts"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestCreateBackupContactHandler() {
	t := suite.T()

	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

	newContactPayload := internalmessages.CreateServiceMemberBackupContactPayload{
		Email:      models.StringPointer("email@example.com"),
		FirstName:  models.StringPointer("firstName"),
		LastName:   models.StringPointer("lastName"),
		Permission: internalmessages.NewBackupContactPermission(internalmessages.BackupContactPermissionEDIT),
		Telephone:  models.StringPointer("5555555555"),
	}
	req := httptest.NewRequest("POST", fmt.Sprintf("/service_member/%v/backup_contacts", serviceMember.ID.String()), nil)

	params := contactop.CreateServiceMemberBackupContactParams{
		CreateBackupContactPayload: &newContactPayload,
		ServiceMemberID:            *handlers.FmtUUID(serviceMember.ID),
	}

	params.HTTPRequest = suite.AuthenticateRequest(req, serviceMember)

	handler := CreateBackupContactHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	_, ok := response.(*contactop.CreateServiceMemberBackupContactCreated)
	if !ok {
		t.Fatalf("Request failed: %#v", response)
	}

	contacts := models.BackupContacts{}
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
	//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
	//RA: in a unit test, then there is no risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	suite.DB().Q().Eager().All(&contacts)
	if len(contacts) != 1 {
		t.Errorf("Expected to find 1 result but found %v", len(contacts))
	}

	if contacts[0].ServiceMember.ID != serviceMember.ID {
		t.Errorf("Expected to find a backup contact for service member")
	}

}

func (suite *HandlerSuite) TestIndexBackupContactsHandler() {
	t := suite.T()

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *handlers.FmtUUID(contact.ServiceMember.ID),
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := IndexBackupContactsHandler{suite.HandlerConfig()}
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

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)
	otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

	indexPath := fmt.Sprintf("/service_member/%v/backup_contacts", contact.ServiceMember.ID.String())
	req := httptest.NewRequest("GET", indexPath, nil)

	params := contactop.IndexServiceMemberBackupContactsParams{
		ServiceMemberID: *handlers.FmtUUID(contact.ServiceMember.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := IndexBackupContactsHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	errResponse := response.(*handlers.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandler() {
	t := suite.T()

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *handlers.FmtUUID(contact.ID),
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := ShowBackupContactHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	okResponse := response.(*contactop.ShowServiceMemberBackupContactOK)
	payload := okResponse.Payload

	if payload.ID.String() != contact.ID.String() {
		t.Errorf("Expected to find a particular backup contact, found something else.")
	}
}

func (suite *HandlerSuite) TestShowBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)
	otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

	showPath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("GET", showPath, nil)

	params := contactop.ShowServiceMemberBackupContactParams{
		BackupContactID: *handlers.FmtUUID(contact.ID),
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := ShowBackupContactHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	errResponse := response.(*handlers.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandler() {
	t := suite.T()

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      models.StringPointer("otheremail@example.com"),
		FirstName:  models.StringPointer("other"),
		LastName:   models.StringPointer("name"),
		Permission: internalmessages.NewBackupContactPermission(internalmessages.BackupContactPermissionNONE),
		Telephone:  models.StringPointer("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *handlers.FmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	params.HTTPRequest = suite.AuthenticateRequest(req, contact.ServiceMember)

	handler := UpdateBackupContactHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	okResponse := response.(*contactop.UpdateServiceMemberBackupContactCreated)
	payload := okResponse.Payload

	if *payload.FirstName != "other" && *payload.LastName != "name" {
		t.Errorf("Expected backup contact to be updated, but it wasn't.")
	}
}

func (suite *HandlerSuite) TestUpdateBackupContactsHandlerWrongUser() {
	t := suite.T()

	contact := factory.BuildBackupContact(suite.DB(), nil, nil)
	otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

	updatePath := fmt.Sprintf("/service_member/%v/backup_contacts/%v", contact.ServiceMember.ID.String(), contact.ID.String())
	req := httptest.NewRequest("PUT", updatePath, nil)

	updateContactPayload := internalmessages.UpdateServiceMemberBackupContactPayload{
		Email:      models.StringPointer("otheremail@example.com"),
		FirstName:  models.StringPointer("other"),
		LastName:   models.StringPointer("name"),
		Permission: internalmessages.NewBackupContactPermission(internalmessages.BackupContactPermissionNONE),
		Telephone:  models.StringPointer("4444444444"),
	}

	params := contactop.UpdateServiceMemberBackupContactParams{
		BackupContactID:                         *handlers.FmtUUID(contact.ID),
		UpdateServiceMemberBackupContactPayload: &updateContactPayload,
	}
	// Logged in as other user
	params.HTTPRequest = suite.AuthenticateRequest(req, otherServiceMember)

	handler := UpdateBackupContactHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	errResponse := response.(*handlers.ErrResponse)
	code := errResponse.Code

	if code != http.StatusForbidden {
		t.Errorf("Expected to receive a forbidden HTTP code, got %v", code)
	}
}
