package shipment

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
		models.MTOShipmentTypeHHGIntoNTS,
		models.MTOShipmentTypeHHGOutOfNTS,
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

func (suite *ShipmentSuite) TestCheckShipmentStatus() {
	suite.Run("checkStatus", func() {
		testCases := map[models.MTOShipmentStatus]bool{
			"":                                true,
			models.MTOShipmentStatusDraft:     true,
			models.MTOShipmentStatusSubmitted: true,
			models.MTOShipmentStatusApproved:  true,
			models.MTOShipmentStatusCancellationRequested: true,
			models.MTOShipmentStatusDiversionRequested:    true,
			models.MTOShipmentStatusRejected:              false,
			models.MTOShipmentStatusCanceled:              false,
			models.MTOShipmentStatusTerminatedForCause:    false,
		}
		for status, allowed := range testCases {
			suite.Run("status "+string(status), func() {
				err := checkStatus().Validate(
					suite.AppContextForTest(),
					models.MTOShipment{Status: status},
				)
				if allowed {
					suite.Empty(err.Error())
				} else {
					suite.NotEmpty(err.Error())
				}
			})
		}
	})
}
