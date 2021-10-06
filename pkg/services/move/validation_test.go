package move

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MoveServiceSuite) TestMoveValidation() {
	appCtx := suite.TestAppContext()

	suite.Run("checkMoveVisibility", func() {
		show := true
		hide := false
		testCases := map[string]struct {
			move   models.Move
			delta  *models.Move
			result error
		}{
			"Success - Move is visible": {
				move: models.Move{
					Show: &show,
				},
				delta:  nil,
				result: nil,
			},
			"Success - Move is being reactivated": {
				move: models.Move{
					Show: nil,
				},
				delta: &models.Move{
					Show: &show,
				},
				result: nil,
			},
			"Fail - Move is not visible": {
				move: models.Move{
					Show: &hide,
				},
				delta:  nil,
				result: services.NotFoundError{},
			},
			"Fail - Move is being deactivated": {
				move: models.Move{
					Show: &show,
				},
				delta: &models.Move{
					Show: &hide,
				},
				result: services.NotFoundError{},
			},
		}
		for name, test := range testCases {
			suite.T().Run(name, func(t *testing.T) {
				err := checkMoveVisibility().Validate(appCtx, test.move, test.delta)
				suite.IsType(test.result, err)
			})
		}
	})

	suite.Run("checkPrimeAvailability", func() {
		now := time.Now()
		testCases := map[string]struct {
			move   models.Move
			delta  *models.Move
			result error
		}{
			"Success - Move is available to Prime": {
				move: models.Move{
					AvailableToPrimeAt: &now,
				},
				delta:  nil,
				result: nil,
			},
			"Success - Move is being made available to Prime": {
				move: models.Move{
					AvailableToPrimeAt: nil,
				},
				delta: &models.Move{
					AvailableToPrimeAt: &now,
				},
				result: nil,
			},
			"Fail - Move is not available to the Prime": {
				move: models.Move{
					AvailableToPrimeAt: &time.Time{},
				},
				delta:  nil,
				result: services.NotFoundError{},
			},
			"Fail - Move is being made unavailable to the Prime": {
				move: models.Move{
					AvailableToPrimeAt: &now,
				},
				delta: &models.Move{
					AvailableToPrimeAt: &time.Time{},
				},
				result: services.NotFoundError{},
			},
		}
		for name, test := range testCases {
			suite.T().Run(name, func(t *testing.T) {
				err := checkPrimeAvailability().Validate(appCtx, test.move, test.delta)
				suite.IsType(test.result, err)
			})
		}
	})
}
