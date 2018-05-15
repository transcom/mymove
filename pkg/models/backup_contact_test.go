package models_test

import (
	"fmt"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_BackupContactCreate() {
	t := suite.T()

	serviceMember, err := testdatagen.MakeServiceMember(suite.db)
	if err != nil {
		t.Fatalf("could not create service member: %v", err)
	}

	newContact := models.BackupContact{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
		Name:            "name",
		Email:           "email@example.com",
		Permission:      internalmessages.BackupContactPermissionEDIT,
	}

	verrs, err := suite.db.ValidateAndCreate(&newContact)

	if err != nil {
		fmt.Println(err)
		t.Fatal("could not save BackupContact", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_BackupContactValidations() {
	contact := &models.BackupContact{}

	var expErrors = map[string][]string{
		"name":       {"Name can not be blank."},
		"email":      {"Email can not be blank."},
		"permission": {"Permission can not be blank."},
	}

	suite.verifyValidationErrors(contact, expErrors)
}

func (suite *ModelSuite) Test_FetchBackupContact() {
	t := suite.T()

	serviceMember1, _ := testdatagen.MakeServiceMember(suite.db)
	serviceMember2, _ := testdatagen.MakeServiceMember(suite.db)
	reqApp := app.MyApp

	backupContact := models.BackupContact{
		ServiceMemberID: serviceMember1.ID,
		Name:            "name",
		Email:           "email@example.com",
		Permission:      internalmessages.BackupContactPermissionEDIT,
	}
	suite.mustSave(&backupContact)

	shouldSucceed, err := models.FetchBackupContact(suite.db, serviceMember1.User, reqApp, backupContact.ID)
	if err != nil || !uuid.Equal(backupContact.ID, shouldSucceed.ID) {
		t.Errorf("failed retrieving own backup contact: %v", err)
	}

	_, err = models.FetchBackupContact(suite.db, serviceMember2.User, reqApp, backupContact.ID)
	if err == nil {
		t.Error("should have failed getting other user's contact")
	}
}
