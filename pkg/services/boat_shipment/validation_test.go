package boatshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// TestBoatShipmentValidatorFunc tests the boat shipment validator function
func (suite *BoatShipmentSuite) TestBoatShipmentValidatorFunc() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), models.BoatShipment{}, nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), models.BoatShipment{}, nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

// TestValidateBoatShipment tests the validateBoatShipment function
func (suite *BoatShipmentSuite) TestValidateBoatShipment() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
			return nil
		})

		err := validateBoatShipment(suite.AppContextForTest(), models.BoatShipment{}, nil, nil, checkAlwaysReturnNil)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		boatShipment := models.BoatShipment{}

		err := validateBoatShipment(suite.AppContextForTest(), boatShipment, &boatShipment, nil, checkAlwaysReturnValidationErr)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Boat shipment not found.")
		})

		err := validateBoatShipment(suite.AppContextForTest(), models.BoatShipment{}, nil, nil, checkAlwaysReturnOtherError)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Boat shipment not found.")
	})
}
