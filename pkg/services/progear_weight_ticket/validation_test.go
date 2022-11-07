package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite ProgearWeightTicketSuite) TestProgearWeightTicketValidatorFuncValidate() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

func (suite ProgearWeightTicketSuite) TestValidateProgearWeightTicket() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
			return nil
		})

		err := validateProgearWeightTicket(suite.AppContextForTest(), nil, nil, []progearWeightTicketValidator{checkAlwaysReturnNil}...)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validateProgearWeightTicket(suite.AppContextForTest(), nil, nil, []progearWeightTicketValidator{checkAlwaysReturnValidationErr}...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the weight ticket.")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, _ *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Weight ticket not found.")
		})

		err := validateProgearWeightTicket(suite.AppContextForTest(), nil, nil, []progearWeightTicketValidator{checkAlwaysReturnOtherError}...)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Weight ticket not found.")
	})
}
