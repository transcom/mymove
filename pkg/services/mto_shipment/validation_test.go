package mtoshipment

import (
	"context"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *MTOShipmentServiceSuite) TestUpdateValidations() {
	suite.Run("checkStatus", func() {
		testCases := map[models.MTOShipmentStatus]bool{
			"":                                            true,
			models.MTOShipmentStatusDraft:                 true,
			models.MTOShipmentStatusSubmitted:             true,
			"random_junk_status":                          false,
			models.MTOShipmentStatusApproved:              false,
			models.MTOShipmentStatusRejected:              false,
			models.MTOShipmentStatusCancellationRequested: false,
			models.MTOShipmentStatusCanceled:              false,
			models.MTOShipmentStatusDiversionRequested:    false,
		}
		for status, allowed := range testCases {
			suite.Run("status "+string(status), func() {
				err := checkStatus().Validate(
					context.Background(),
					&models.MTOShipment{Status: status},
					nil,
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
