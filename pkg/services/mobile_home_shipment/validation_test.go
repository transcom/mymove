package mobilehomeshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// TestMobileHomeShipmentValidatorFunc tests the mobile home shipment validator function
func (suite *MobileHomeShipmentSuite) TestMobileHomeShipmentValidatorFunc() {
	suite.Run("Calling Validate runs validation function with no errors", func() {
		validator := mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
			return nil
		})

		err := validator.Validate(suite.AppContextForTest(), models.MobileHome{}, nil, nil)

		suite.NoError(err)
	})

	suite.Run("Calling Validate runs validation function with errors", func() {
		validator := mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		err := validator.Validate(suite.AppContextForTest(), models.MobileHome{}, nil, nil)

		suite.Error(err)
		suite.Contains(err.Error(), "fake error")
	})
}

// TestValidateMobileHomeShipment tests the validateMobileHomeShipment function
func (suite *MobileHomeShipmentSuite) TestValidateMobileHomeShipment() {
	suite.Run("Runs validation and returns nil when there are no errors", func() {
		checkAlwaysReturnNil := mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
			return nil
		})

		err := validateMobileHomeShipment(suite.AppContextForTest(), models.MobileHome{}, nil, nil, checkAlwaysReturnNil)

		suite.NoError(err)
	})

	suite.Run("Runs validation and returns input errors", func() {
		checkAlwaysReturnValidationErr := mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
			verrs := validate.NewErrors()

			verrs.Add("ID", "fake error")

			return verrs
		})

		mobileHome := models.MobileHome{}

		err := validateMobileHomeShipment(suite.AppContextForTest(), mobileHome, &mobileHome, nil, checkAlwaysReturnValidationErr)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found")
	})

	suite.Run("Runs validation and returns other errors", func() {
		checkAlwaysReturnOtherError := mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
			return apperror.NewNotFoundError(uuid.Must(uuid.NewV4()), "Mobile Home shipment not found.")
		})

		err := validateMobileHomeShipment(suite.AppContextForTest(), models.MobileHome{}, nil, nil, checkAlwaysReturnOtherError)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Mobile Home shipment not found.")
	})
}
