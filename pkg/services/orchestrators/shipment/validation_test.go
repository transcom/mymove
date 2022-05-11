package shipment

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite ShipmentSuite) TestValidateShipment() {
	suite.Run("Returns InvalidInputError if a validator function returns an error", func() {
		err := validateShipment(suite.AppContextForTest(), models.MTOShipment{}, basicShipmentChecks()...)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the shipment.")
	})

	suite.Run("Returns nil if no validator function returns an error", func() {
		err := validateShipment(suite.AppContextForTest(), models.MTOShipment{ShipmentType: models.MTOShipmentTypeHHG}, basicShipmentChecks()...)

		suite.NoError(err)
	})
}
