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
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	servicememberop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/service_members"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
)

func (suite *HandlerSuite) TestShowServiceMemberHandler() {
	// Given: A servicemember and a user
	user := factory.BuildDefaultUser(suite.DB())

	newServiceMember := factory.BuildExtendedServiceMember(suite.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
	suite.MustSave(&newServiceMember)

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(newServiceMember.ID.String()),
	}
	// And: show ServiceMember is queried
	showHandler := ShowServiceMemberHandler{suite.HandlerConfig()}
	response := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&servicememberop.ShowServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOK)

	// And: Returned query to include our added servicemember
	suite.Assertions.Equal(user.ID.String(), okResponse.Payload.UserID.String())
}

func (suite *HandlerSuite) TestShowServiceMemberWrongUser() {
	// Given: Servicemember trying to load another
	notLoggedInUser := factory.BuildServiceMember(suite.DB(), nil, nil)
	loggedInUser := factory.BuildServiceMember(suite.DB(), nil, nil)

	req := httptest.NewRequest("GET", fmt.Sprintf("/service_members/%s", notLoggedInUser.ID.String()), nil)
	req = suite.AuthenticateRequest(req, loggedInUser)

	showServiceMemberParams := servicememberop.ShowServiceMemberParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(notLoggedInUser.ID.String()),
	}
	// And: Show servicemember is queried
	showHandler := ShowServiceMemberHandler{suite.HandlerConfig()}
	response := showHandler.Handle(showServiceMemberParams)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerNoValues() {
	// Given: A logged-in user
	user := factory.BuildDefaultUser(suite.DB())

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{}

	req := httptest.NewRequest("POST", "/service_members", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}
	handler := CreateServiceMemberHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.CookieUpdateResponder{}, response)

	unwrappedResponse := response.(*handlers.CookieUpdateResponder).Responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.DB().Where("user_id = ?", user.ID)
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)

	serviceMemberPayload := unwrappedResponse.(*servicememberop.CreateServiceMemberCreated).Payload

	suite.Assertions.NotEqual(*serviceMemberPayload.ID, uuid.Nil)
	suite.Assertions.NotEqual(*serviceMemberPayload.UserID, uuid.Nil)
	suite.Assertions.Equal(*serviceMemberPayload.IsProfileComplete, false)
	suite.Assertions.Equal(len((*serviceMemberPayload).Orders), 0)

	// These shouldn't return any value or Swagger clients will complain during validation
	// because the payloads for these objects are defined to require non-null values for most fields
	// which can't be handled in OpenAPI Spec 2.0. Therefore we don't return them at all.
	suite.Assertions.Equal((*serviceMemberPayload).Grade, (*internalmessages.OrderPayGrade)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).Affiliation, (*internalmessages.Affiliation)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).ResidentialAddress, (*internalmessages.Address)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).BackupMailingAddress, (*internalmessages.Address)(nil))
	suite.Assertions.Equal((*serviceMemberPayload).BackupContacts, internalmessages.IndexServiceMemberBackupContactsPayload{})
}

func (suite *HandlerSuite) TestSubmitServiceMemberHandlerAllValues() {
	// Given: A logged-in user
	user := factory.BuildDefaultUser(suite.DB())

	// When: a new ServiceMember is posted
	newServiceMemberPayload := internalmessages.CreateServiceMemberPayload{
		UserID:               strfmt.UUID(user.ID.String()),
		Edipi:                models.StringPointer("random string bla"),
		FirstName:            models.StringPointer("random string bla"),
		MiddleName:           models.StringPointer("random string bla"),
		LastName:             models.StringPointer("random string bla"),
		Suffix:               models.StringPointer("random string bla"),
		Telephone:            models.StringPointer("random string bla"),
		SecondaryTelephone:   models.StringPointer("random string bla"),
		PersonalEmail:        models.StringPointer("wml@example.com"),
		PhoneIsPreferred:     models.BoolPointer(false),
		EmailIsPreferred:     models.BoolPointer(true),
		ResidentialAddress:   fakeAddressPayload(),
		BackupMailingAddress: fakeAddressPayload(),
	}

	req := httptest.NewRequest("GET", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.CreateServiceMemberParams{
		CreateServiceMemberPayload: &newServiceMemberPayload,
		HTTPRequest:                req,
	}

	handler := CreateServiceMemberHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.CookieUpdateResponder{}, response)
	unwrappedResponse := response.(*handlers.CookieUpdateResponder).Responder
	suite.Assertions.IsType(&servicememberop.CreateServiceMemberCreated{}, unwrappedResponse)

	// Then: we expect a servicemember to have been created for the user
	query := suite.DB().Where("user_id = ?", user.ID)
	var serviceMembers models.ServiceMembers
	query.All(&serviceMembers)

	suite.Assertions.Len(serviceMembers, 1)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandler() {
	// TODO: add more fields to change
	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newEmplid := "1234567"

	origAffiliation := models.AffiliationAIRFORCE
	newAffiliation := internalmessages.AffiliationCOASTGUARD

	origFirstName := models.StringPointer("random string bla")
	newFirstName := models.StringPointer("John")

	origMiddleName := models.StringPointer("random string bla")
	newMiddleName := models.StringPointer("")

	origLastName := models.StringPointer("random string bla")
	newLastName := models.StringPointer("Doe")

	origSuffix := models.StringPointer("random string bla")
	newSuffix := models.StringPointer("Mr.")

	origTelephone := models.StringPointer("random string bla")
	newTelephone := models.StringPointer("555-555-5555")

	origSecondaryTelephone := models.StringPointer("random string bla")
	newSecondaryTelephone := models.StringPointer("555-555-5555")

	origPersonalEmail := models.StringPointer("wml@example.com")
	newPersonalEmail := models.StringPointer("example@email.com")

	origPhoneIsPreferred := models.BoolPointer(false)
	newPhoneIsPreferred := models.BoolPointer(true)

	origEmailIsPreferred := models.BoolPointer(true)
	newEmailIsPreferred := models.BoolPointer(false)

	newServiceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Edipi:              &origEdipi,
				Emplid:             nil,
				Affiliation:        &origAffiliation,
				FirstName:          origFirstName,
				MiddleName:         origMiddleName,
				LastName:           origLastName,
				Suffix:             origSuffix,
				Telephone:          origTelephone,
				SecondaryTelephone: origSecondaryTelephone,
				PersonalEmail:      origPersonalEmail,
				PhoneIsPreferred:   origPhoneIsPreferred,
				EmailIsPreferred:   origEmailIsPreferred,
			},
		},
	}, nil)

	orderGrade := models.ServiceMemberGradeE5
	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				Grade: &orderGrade,
			},
		},
		{
			Model:    newServiceMember,
			LinkOnly: true,
		},
	}, nil)

	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &newEdipi,
		Emplid:               &newEmplid,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &newAffiliation,
		EmailIsPreferred:     newEmailIsPreferred,
		FirstName:            newFirstName,
		LastName:             newLastName,
		MiddleName:           newMiddleName,
		PersonalEmail:        newPersonalEmail,
		PhoneIsPreferred:     newPhoneIsPreferred,
		SecondaryTelephone:   newSecondaryTelephone,
		Suffix:               newSuffix,
		Telephone:            newTelephone,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	handler := PatchServiceMemberHandler{handlerConfig}
	response := handler.Handle(params)

	suite.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	suite.Equal(newEdipi, *serviceMemberPayload.Edipi)
	suite.Equal(newEmplid, *serviceMemberPayload.Emplid)
	suite.Equal(newAffiliation, *serviceMemberPayload.Affiliation)
	suite.Equal(*newFirstName, *serviceMemberPayload.FirstName)
	suite.Equal(*newMiddleName, *serviceMemberPayload.MiddleName)
	suite.Equal(*newLastName, *serviceMemberPayload.LastName)
	suite.Equal(*newSuffix, *serviceMemberPayload.Suffix)
	suite.Equal(*newTelephone, *serviceMemberPayload.Telephone)
	suite.Equal(*newSecondaryTelephone, *serviceMemberPayload.SecondaryTelephone)
	suite.Equal(*newPersonalEmail, *serviceMemberPayload.PersonalEmail)
	suite.Equal(*newPhoneIsPreferred, *serviceMemberPayload.PhoneIsPreferred)
	suite.Equal(*newEmailIsPreferred, *serviceMemberPayload.EmailIsPreferred)
	suite.Equal(*resAddress.StreetAddress1, *serviceMemberPayload.ResidentialAddress.StreetAddress1)
	suite.Equal(*resAddress.County, *serviceMemberPayload.ResidentialAddress.County)
	suite.Equal(*backupAddress.StreetAddress1, *serviceMemberPayload.BackupMailingAddress.StreetAddress1)
	suite.Equal(*backupAddress.County, *serviceMemberPayload.BackupMailingAddress.County)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerSubmittedMove() {
	// Given: a logged in user
	user := factory.BuildDefaultUser(suite.DB())

	edipi := "2342342344"

	origAffiliation := models.AffiliationAIRFORCE
	newAffiliation := internalmessages.AffiliationARMY

	origDutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	newDutyLocation := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	newDutyLocationID := strfmt.UUID(newDutyLocation.ID.String())

	origFirstName := models.StringPointer("random string bla")
	newFirstName := models.StringPointer("John")

	origMiddleName := models.StringPointer("random string bla")
	newMiddleName := models.StringPointer("")

	origLastName := models.StringPointer("random string bla")
	newLastName := models.StringPointer("Doe")

	origSuffix := models.StringPointer("random string bla")
	newSuffix := models.StringPointer("Mr.")

	origTelephone := models.StringPointer("random string bla")
	newTelephone := models.StringPointer("555-555-5555")

	origSecondaryTelephone := models.StringPointer("random string bla")
	newSecondaryTelephone := models.StringPointer("555-555-5555")

	origPersonalEmail := models.StringPointer("wml@example.com")
	newPersonalEmail := models.StringPointer("example@email.com")

	origPhoneIsPreferred := models.BoolPointer(false)
	newPhoneIsPreferred := models.BoolPointer(true)

	origEmailIsPreferred := models.BoolPointer(true)
	newEmailIsPreferred := models.BoolPointer(false)

	newServiceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				UserID:             user.ID,
				Edipi:              &edipi,
				Affiliation:        &origAffiliation,
				FirstName:          origFirstName,
				MiddleName:         origMiddleName,
				LastName:           origLastName,
				Suffix:             origSuffix,
				Telephone:          origTelephone,
				SecondaryTelephone: origSecondaryTelephone,
				PersonalEmail:      origPersonalEmail,
				PhoneIsPreferred:   origPhoneIsPreferred,
				EmailIsPreferred:   origEmailIsPreferred,
			},
		},
		{
			Model:    origDutyLocation,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC:                 nil,
				DepartmentIndicator: nil,
				OrdersNumber:        nil,
				OrdersTypeDetail:    nil,
			},
		},
		{
			Model:    newServiceMember,
			LinkOnly: true,
		},
		{
			Model:    origDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    newDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	// The factory sets these values, fails if you try to blank them out via Customizations,
	// and gives defaults if you pass nil, so we have to set this after the creation.
	// This more closely resembles what orders would look like pre and post submission, before
	// a TOO gets to them.
	move.Orders.TAC = nil
	move.Orders.DepartmentIndicator = nil
	move.Orders.OrdersNumber = nil
	move.Orders.OrdersTypeDetail = nil

	suite.MustSave(&move.Orders)
	moveRouter := moverouter.NewMoveRouter()
	newSignedCertification := factory.BuildSignedCertification(nil, []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	err := moveRouter.Submit(suite.AppContextForTest(), &move, &newSignedCertification)
	suite.NoError(err)
	suite.MustSave(&move)

	resAddress := fakeAddressPayload()
	backupAddress := fakeAddressPayload()
	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi:                &edipi,
		BackupMailingAddress: backupAddress,
		ResidentialAddress:   resAddress,
		Affiliation:          &newAffiliation,
		EmailIsPreferred:     newEmailIsPreferred,
		FirstName:            newFirstName,
		LastName:             newLastName,
		MiddleName:           newMiddleName,
		PersonalEmail:        newPersonalEmail,
		PhoneIsPreferred:     newPhoneIsPreferred,
		SecondaryTelephone:   newSecondaryTelephone,
		Suffix:               newSuffix,
		Telephone:            newTelephone,
		CurrentLocationID:    &newDutyLocationID,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handlerConfig := suite.HandlerConfig()
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)
	handler := PatchServiceMemberHandler{handlerConfig}
	response := handler.Handle(params)

	suite.IsType(&servicememberop.PatchServiceMemberOK{}, response)
	okResponse := response.(*servicememberop.PatchServiceMemberOK)

	serviceMemberPayload := okResponse.Payload

	// These fields should not change (they should still be the original
	// values) after the move has been submitted.
	suite.Equal(origAffiliation, models.ServiceMemberAffiliation(*serviceMemberPayload.Affiliation))

	// These fields should change even if the move is submitted.
	suite.Equal(*newFirstName, *serviceMemberPayload.FirstName)
	suite.Equal(*newMiddleName, *serviceMemberPayload.MiddleName)
	suite.Equal(*newLastName, *serviceMemberPayload.LastName)
	suite.Equal(*newSuffix, *serviceMemberPayload.Suffix)
	suite.Equal(*newTelephone, *serviceMemberPayload.Telephone)
	suite.Equal(*newSecondaryTelephone, *serviceMemberPayload.SecondaryTelephone)
	suite.Equal(*newPersonalEmail, *serviceMemberPayload.PersonalEmail)
	suite.Equal(*newPhoneIsPreferred, *serviceMemberPayload.PhoneIsPreferred)
	suite.Equal(*newEmailIsPreferred, *serviceMemberPayload.EmailIsPreferred)

	suite.Equal(*resAddress.StreetAddress1, *serviceMemberPayload.ResidentialAddress.StreetAddress1)
	suite.Equal(*resAddress.County, *serviceMemberPayload.ResidentialAddress.County)
	suite.Equal(*backupAddress.StreetAddress1, *serviceMemberPayload.BackupMailingAddress.StreetAddress1)
	suite.Equal(*backupAddress.County, *serviceMemberPayload.BackupMailingAddress.County)

	// Then: we expect addresses to have been created
	addresses := []models.Address{}
	suite.DB().All(&addresses)
	suite.Equal(6, len(addresses))
	// Why 6?
	// Make duty locations +2 addresses each DL => 4
	// Patch service member +2 addresses added to service member => 2
	// Total => 6
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerWrongUser() {
	// Given: a logged in user
	user := factory.BuildDefaultUser(suite.DB())
	user2 := factory.BuildDefaultUser(suite.DB())

	var origEdipi = "2342342344"
	var newEdipi = "9999999999"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user2)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoServiceMember() {
	// Given: a logged in user
	user := factory.BuildDefaultUser(suite.DB())

	servicememberUUID := uuid.Must(uuid.NewV4())

	var newEdipi = "9999999999"

	patchPayload := internalmessages.PatchServiceMemberPayload{
		Edipi: &newEdipi,
	}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(servicememberUUID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)

	suite.Assertions.Equal(http.StatusNotFound, errResponse.Code)
}

func (suite *HandlerSuite) TestPatchServiceMemberHandlerNoChange() {
	// Given: a logged in user with a servicemember
	user := factory.BuildDefaultUser(suite.DB())

	var origEdipi = "4444444444"
	newServiceMember := models.ServiceMember{
		UserID: user.ID,
		Edipi:  &origEdipi,
	}
	suite.MustSave(&newServiceMember)

	patchPayload := internalmessages.PatchServiceMemberPayload{}

	req := httptest.NewRequest("PATCH", "/service_members/some_id", nil)
	req = suite.AuthenticateRequest(req, newServiceMember)

	params := servicememberop.PatchServiceMemberParams{
		HTTPRequest:               req,
		ServiceMemberID:           strfmt.UUID(newServiceMember.ID.String()),
		PatchServiceMemberPayload: &patchPayload,
	}

	handler := PatchServiceMemberHandler{suite.HandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(&servicememberop.PatchServiceMemberOK{}, response)
}

func (suite *HandlerSuite) TestShowServiceMemberOrders() {
	order1 := factory.BuildOrder(suite.DB(), nil, nil)
	order2 := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    order1.ServiceMember,
			LinkOnly: true,
		},
	}, nil)

	req := httptest.NewRequest("GET", "/service_members/some_id/current_orders", nil)
	req = suite.AuthenticateRequest(req, order1.ServiceMember)

	params := servicememberop.ShowServiceMemberOrdersParams{
		HTTPRequest:     req,
		ServiceMemberID: strfmt.UUID(order1.ServiceMemberID.String()),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	handler := ShowServiceMemberOrdersHandler{handlerConfig}

	response := handler.Handle(params)

	suite.IsType(&servicememberop.ShowServiceMemberOrdersOK{}, response)
	okResponse := response.(*servicememberop.ShowServiceMemberOrdersOK)
	responsePayload := okResponse.Payload

	// Should return the most recently created order
	suite.Equal(order2.ID.String(), responsePayload.ID.String())
}
