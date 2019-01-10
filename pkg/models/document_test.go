package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_DocumentCreate() {
	t := suite.T()

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	document := models.Document{
		ServiceMemberID: serviceMember.ID,
	}

	verrs, err := suite.DB().ValidateAndSave(&document)

	if err != nil {
		t.Fatalf("could not save Document: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_DocumentValidations() {
	document := &models.Document{}

	var expErrors = map[string][]string{
		"service_member_id": {"ServiceMemberID can not be blank."},
	}

	suite.verifyValidationErrors(document, expErrors)
}
