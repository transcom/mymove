package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_DocumentCreate() {
	serviceMember := factory.BuildServiceMember(nil, nil, []factory.Trait{factory.GetTraitServiceMemberSetIDs})

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

	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
	session := auth.Session{
		UserID:          serviceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}
	userUpload := factory.BuildUserUpload(suite.DB(), nil, nil)
	err := models.DeleteUserUpload(suite.DB(), &userUpload)
	suite.Nil(err)
	userUploads := models.UserUploads{userUpload}

	document := models.Document{
		ID:              *userUpload.DocumentID,
		ServiceMemberID: serviceMember.ID,
		UserUploads:     userUploads,
	}
	userUpload.DocumentID = &document.ID

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
	suite.Equal(0, len(doc.UserUploads))
}

func (suite *ModelSuite) TestFetchDeletedDocument() {
	t := suite.T()

	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
	session := auth.Session{
		UserID:          serviceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	}

	userUpload := factory.BuildUserUpload(suite.DB(), nil, nil)
	err := models.DeleteUserUpload(suite.DB(), &userUpload)
	suite.Nil(err)
	userUploads := models.UserUploads{userUpload}

	deletedAt := time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC)
	document := models.Document{
		ID:              *userUpload.DocumentID,
		ServiceMemberID: serviceMember.ID,
		DeletedAt:       &deletedAt,
		UserUploads:     userUploads,
	}
	userUpload.DocumentID = &document.ID

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
	suite.Equal(1, len(doc2.UserUploads))
}
