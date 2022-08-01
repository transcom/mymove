package movingexpense

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite MovingExpenseSuite) TestWeightTicketValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := movingExpenseValidatorFunc(func(_ appcontext.AppContext, _ *models.MovingExpense, _ *models.MovingExpense) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := movingExpenseValidatorFunc(func(_ appcontext.AppContext, _ *models.MovingExpense, _ *models.MovingExpense) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite MovingExpenseSuite) TestValidateMovingExpense() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := movingExpenseValidatorFunc(func(_ appcontext.AppContext, _ *models.MovingExpense, _ *models.MovingExpense) error {
			return nil
		})

		err := validateMovingExpense(suite.AppContextForTest(), nil, nil, []movingExpenseValidator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := movingExpenseValidatorFunc(func(_ appcontext.AppContext, _ *models.MovingExpense, _ *models.MovingExpense) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateMovingExpense(suite.AppContextForTest(), nil, nil, []movingExpenseValidator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input received. fake error")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := movingExpenseValidatorFunc(func(_ appcontext.AppContext, _ *models.MovingExpense, _ *models.MovingExpense) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Moving expense not found.")
		})

		err := validateMovingExpense(suite.AppContextForTest(), nil, nil, []movingExpenseValidator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Moving expense not found.")
	})
}
