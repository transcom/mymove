package gunsafeweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GunSafeWeightTicketSuite) TestGunSafeWeightTicketValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite *GunSafeWeightTicketSuite) TestValidateGunSafeWeightTicket() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
			return nil
		})

		err := validateGunSafeWeightTicket(suite.AppContextForTest(), nil, nil, []gunSafeWeightTicketValidator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateGunSafeWeightTicket(suite.AppContextForTest(), nil, nil, []gunSafeWeightTicketValidator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the gunSafe weight ticket.")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "GunSafe weight ticket not found.")
		})

		err := validateGunSafeWeightTicket(suite.AppContextForTest(), nil, nil, []gunSafeWeightTicketValidator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "GunSafe weight ticket not found.")
	})
}
