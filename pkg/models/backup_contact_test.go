package models_test

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_BackupContactCreate() {
	serviceMember := factory.BuildServiceMember(nil, nil, []factory.Trait{factory.GetTraitServiceMemberSetIDs})

	newContact := models.BackupContact{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
		FirstName:       "firstName",
		LastName:        "lastName",
		Email:           "email@example.com",
		Permission:      models.BackupContactPermissionEDIT,
	}

	verrs, err := newContact.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) Test_BackupContactValidations() {
	contact := &models.BackupContact{}

	var expErrors = map[string][]string{
		"first_name": {"FirstName can not be blank."},
		"last_name":  {"LastName can not be blank."},
		"email":      {"Email can not be blank."},
		"permission": {"Permission can not be blank."},
	}

	suite.verifyValidationErrors(contact, expErrors, nil)
}

func (suite *ModelSuite) Test_FetchBackupContact() {
	t := suite.T()

	serviceMember1 := factory.BuildServiceMember(suite.DB(), nil, nil)
	serviceMember2 := factory.BuildServiceMember(suite.DB(), nil, nil)

	backupContact := models.BackupContact{
		ServiceMemberID: serviceMember1.ID,
		FirstName:       "firstName",
		LastName:        "lastName",
		Email:           "email@example.com",
		Permission:      models.BackupContactPermissionEDIT,
	}
	suite.MustSave(&backupContact)

	session := &auth.Session{
		UserID:          serviceMember1.UserID,
		ServiceMemberID: serviceMember1.ID,
		ApplicationName: auth.MilApp,
	}
	shouldSucceed, err := models.FetchBackupContact(suite.DB(), session, backupContact.ID)
	if err != nil || backupContact.ID != shouldSucceed.ID {
		t.Errorf("failed retrieving own backup contact: %v", err)
	}

	session.UserID = serviceMember2.UserID
	session.ServiceMemberID = serviceMember2.ID
	_, err = models.FetchBackupContact(suite.DB(), session, backupContact.ID)
	if err == nil {
		t.Error("should have failed getting other user's contact")
	}
}
