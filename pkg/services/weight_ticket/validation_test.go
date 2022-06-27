package weightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite WeightTicketSuite) TestWeightTicketValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := weightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.WeightTicket, _ *models.WeightTicket) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := weightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.WeightTicket, _ *models.WeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite WeightTicketSuite) TestValidateWeightTicket() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := weightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.WeightTicket, _ *models.WeightTicket) error {
			return nil
		})

		err := validateWeightTicket(suite.AppContextForTest(), nil, nil, []weightTicketValidator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := weightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.WeightTicket, _ *models.WeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateWeightTicket(suite.AppContextForTest(), nil, nil, []weightTicketValidator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the weight ticket.")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := weightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.WeightTicket, _ *models.WeightTicket) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Weight ticket not found.")
		})

		err := validateWeightTicket(suite.AppContextForTest(), nil, nil, []weightTicketValidator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Weight ticket not found.")
	})
}
