package signedcertification

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *SignedCertificationSuite) TestSignedCertificationValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := signedCertificationValidatorFunc(func(_ appcontext.AppContext, _ models.SignedCertification, _ *models.SignedCertification) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), models.SignedCertification{}, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with an error", func() {
		verrs := validate.NewErrors()

		verrs.Add("ID", "fake error")

		validator := signedCertificationValidatorFunc(func(_ appcontext.AppContext, _ models.SignedCertification, _ *models.SignedCertification) error {
			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), models.SignedCertification{}, nil)

		suite.Error(err)
		suite.Equal(verrs, err)
	})
}

func (suite *SignedCertificationSuite) TestValidateSignedCertification() {
	suite.Run("Runs validation and returns nil if no errors", func() {
		checkAlwaysReturnNil := signedCertificationValidatorFunc(func(_ appcontext.AppContext, _ models.SignedCertification, _ *models.SignedCertification) error {
			return nil
		})

		err := validateSignedCertification(suite.AppContextForTest(), models.SignedCertification{}, nil, checkAlwaysReturnNil)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns error if there is an input error", func() {
		checkAlwaysReturnValidationErr := signedCertificationValidatorFunc(func(_ appcontext.AppContext, _ models.SignedCertification, _ *models.SignedCertification) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateSignedCertification(suite.AppContextForTest(), models.SignedCertification{}, nil, checkAlwaysReturnValidationErr)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the signed certification.", err.Error())
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := signedCertificationValidatorFunc(func(_ appcontext.AppContext, _ models.SignedCertification, _ *models.SignedCertification) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "for SignedCertification.")
		})

		err := validateSignedCertification(suite.AppContextForTest(), models.SignedCertification{}, nil, checkAlwaysReturnOtherError)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found for SignedCertification.")
	})
}
