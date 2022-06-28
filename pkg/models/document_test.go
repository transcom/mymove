package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_DocumentCreate() {
	serviceMember := testdatagen.MakeStubbedServiceMember(suite.DB())

	document := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := document.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) Test_DocumentValidations() {
	document := &models.Document{}

	var expErrors = map[string][]string{
		"service_member_id": {"ServiceMemberID can not be blank."},
	}

	suite.verifyValidationErrors(document, expErrors)
}

func (suite *ModelSuite) TestFetchDocument() {
	t := suite.T()

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	session := auth.Session{
		UserID:          serviceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}
	document := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := suite.DB().ValidateAndSave(&document)
	if err != nil {
		t.Errorf("could not save UserUpload: %v", err)
		return
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	doc, _ := models.FetchDocument(suite.DB(), &session, document.ID, false)
	suite.Equal(doc.ID, document.ID)
}

func (suite *ModelSuite) TestFetchDeletedDocument() {
	t := suite.T()

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	session := auth.Session{
		UserID:          serviceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}

	deletedAt := time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC)
	document := models.Document{
		ServiceMemberID: serviceMember.ID,
		DeletedAt:       &deletedAt,
	}

	verrs, err := suite.DB().ValidateAndSave(&document)
	if err != nil {
		t.Errorf("could not save UserUpload: %v", err)
		return
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	doc, _ := models.FetchDocument(suite.DB(), &session, document.ID, false)

	// fetches a nil document
	suite.Equal(doc.ID, uuid.Nil)
	suite.Equal(doc.ServiceMemberID, uuid.Nil)

	doc2, _ := models.FetchDocument(suite.DB(), &session, document.ID, true)

	// fetches a nil document
	suite.Equal(doc2.ID, document.ID)
	suite.Equal(doc2.ServiceMemberID, serviceMember.ID)
}
