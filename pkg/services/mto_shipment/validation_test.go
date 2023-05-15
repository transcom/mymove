package mtoshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MTOShipmentServiceSuite) TestvalidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := validatorFunc(func(_ appcontext.AppContext, _ *models.MTOShipment, _ *models.MTOShipment) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := validatorFunc(func(_ appcontext.AppContext, _ *models.MTOShipment, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite *MTOShipmentServiceSuite) TestvalidateShipment() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := validatorFunc(func(_ appcontext.AppContext, _ *models.MTOShipment, _ *models.MTOShipment) error {
			return nil
		})

		err := validateShipment(suite.AppContextForTest(), nil, nil, []validator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := validatorFunc(func(_ appcontext.AppContext, _ *models.MTOShipment, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateShipment(suite.AppContextForTest(), nil, nil, []validator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the weight ticket.")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := validatorFunc(func(_ appcontext.AppContext, _ *models.MTOShipment, _ *models.MTOShipment) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Weight ticket not found.")
		})

		err := validateShipment(suite.AppContextForTest(), nil, nil, []validator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Weight ticket not found.")
	})
}
