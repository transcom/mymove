package shipmentorchestrator

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ShipmentSuite) TestCheckShipmentType() {
	suite.Run("Return an error for an empty shipment type", func() {
		shipment := models.MTOShipment{}

		err := checkShipmentType().Validate(suite.AppContextForTest(), shipment)

		suite.Error(err)
		suite.IsType(&validate.Errors{}, err)
		suite.Contains(err.Error(), "ShipmentType must be a valid type.")
	})

	validShipmentTypes := []models.MTOShipmentType{
		models.MTOShipmentTypeHHG,
		models.MTOShipmentTypeHHGLongHaulDom,
		models.MTOShipmentTypeHHGShortHaulDom,
		models.MTOShipmentTypeHHGIntoNTSDom,
		models.MTOShipmentTypeHHGOutOfNTSDom,
		models.MTOShipmentTypePPM,
	}

	for _, shipmentType := range validShipmentTypes {
		shipmentType := shipmentType

		suite.Run(fmt.Sprintf("Doesn't return an error if the shipment type is %s", shipmentType), func() {
			shipment := models.MTOShipment{ShipmentType: shipmentType}

			err := checkShipmentType().Validate(suite.AppContextForTest(), shipment)

			suite.NilOrNoVerrs(err)
		})
	}
}
