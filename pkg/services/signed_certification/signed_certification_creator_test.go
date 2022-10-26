package signedcertification

import (
	"fmt"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SignedCertificationSuite) TestCreateSignedCertification() {
	suite.Run("Returns an InvalidInputError if there's an issue with the input data", func() {
		creator := NewSignedCertificationCreator()

		newSignedCertification, err := creator.CreateSignedCertification(suite.AppContextForTest(), models.SignedCertification{})

		suite.Nil(newSignedCertification)

		if suite.Error(err) {
			suite.IsType(apperror.InvalidInputError{}, err)
			suite.Equal("Invalid input found while validating the signed certification.", err.Error())
		}
	})

	suite.Run("Returns a transaction error if one is raised when validating the data", func() {
		// It's hard to trigger this without skipping the validation that the default creator runs, so we'll use a
		// creator that doesn't have any validation so we can trigger the model-level validation.
		creator := &signedCertificationCreator{}

		newSignedCertification, createErr := creator.CreateSignedCertification(suite.AppContextForTest(), models.SignedCertification{})

		suite.Nil(newSignedCertification)

		if suite.Error(createErr) {
			suite.IsType(apperror.InvalidInputError{}, createErr)
			suite.Equal("Invalid input found while creating the signed certification.", createErr.Error())
		}
	})

	suite.Run("Returns a transaction error if one is raised when creating the record", func() {
		// It's hard to trigger this without skipping the validation that the default creator runs, so we'll use a
		// creator that doesn't have any validation so we can trigger the model-level validation.
		creator := &signedCertificationCreator{}

		signedCertification := testdatagen.MakeDefaultSignedCertification(suite.DB())

		newSignedCertification, createErr := creator.CreateSignedCertification(suite.AppContextForTest(), signedCertification)

		suite.Nil(newSignedCertification)

		if suite.Error(createErr) {
			suite.IsType(apperror.QueryError{}, createErr)
			suite.Equal("Unable to create signed certification", createErr.Error())
		}
	})

	suite.Run("Can successfully create a signed certification", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		serviceMember := move.Orders.ServiceMember

		shipmentCertType := models.SignedCertificationTypeSHIPMENT

		signedCertification := models.SignedCertification{
			SubmittingUserID:  serviceMember.User.ID,
			MoveID:            move.ID,
			CertificationType: &shipmentCertType,
			CertificationText: "I certify that the information I have provided is true and complete to the best of my knowledge.",
			Signature:         fmt.Sprintf("%s %s", *serviceMember.FirstName, *serviceMember.LastName),
			Date:              testdatagen.NextValidMoveDate,
		}

		creator := NewSignedCertificationCreator()

		newSignedCertification, err := creator.CreateSignedCertification(suite.AppContextForTest(), signedCertification)

		suite.Nil(err)

		if suite.NotNil(newSignedCertification) {
			suite.NotNil(newSignedCertification.ID)
		}
	})
}
