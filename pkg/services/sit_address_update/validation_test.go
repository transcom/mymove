package sitaddressupdate

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite SITAddressUpdateServiceSuite) TestSITAddressUpdateValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, _ *models.SITAddressUpdate) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), &models.SITAddressUpdate{})

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, _ *models.SITAddressUpdate) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), &models.SITAddressUpdate{})

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite SITAddressUpdateServiceSuite) TestValidateSITAddressUpdate() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, _ *models.SITAddressUpdate) error {
			return nil
		})

		err := validateSITAddressUpdate(suite.AppContextForTest(), &models.SITAddressUpdate{}, []sitAddressUpdateValidator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, _ *models.SITAddressUpdate) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateSITAddressUpdate(suite.AppContextForTest(), &models.SITAddressUpdate{}, []sitAddressUpdateValidator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the SIT address update.")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, _ *models.SITAddressUpdate) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "SITAddressUpdate not found.")
		})

		err := validateSITAddressUpdate(suite.AppContextForTest(), &models.SITAddressUpdate{}, []sitAddressUpdateValidator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "SITAddressUpdate not found.")
	})
}
