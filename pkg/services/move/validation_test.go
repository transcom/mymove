package move

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MoveServiceSuite) TestMoveValidation() {
	appCtx := suite.AppContextForTest()

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
				result: apperror.NotFoundError{},
			},
			"Fail - Move is being deactivated": {
				move: models.Move{
					Show: &show,
				},
				delta: &models.Move{
					Show: &hide,
				},
				result: apperror.NotFoundError{},
			},
		}
		for name, test := range testCases {
			suite.Run(name, func() {
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
				result: apperror.NotFoundError{},
			},
			"Fail - Move is being made unavailable to the Prime": {
				move: models.Move{
					AvailableToPrimeAt: &now,
				},
				delta: &models.Move{
					AvailableToPrimeAt: &time.Time{},
				},
				result: apperror.NotFoundError{},
			},
		}
		for name, test := range testCases {
			suite.Run(name, func() {
				err := checkPrimeAvailability().Validate(appCtx, test.move, test.delta)
				suite.IsType(test.result, err)
			})
		}
	})
}
