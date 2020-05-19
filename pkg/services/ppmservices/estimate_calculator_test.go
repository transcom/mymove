package ppmservices

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMServiceSuite) TestCalculateEstimate() {
	// Subtests:
	// estimate calculation success
	// bad moveID fails
	// bad origin zip fails (90210, 90210)
	// bad origin duty station zip fails (90210, 90210)
	// bad ppm weight estimate (0?) fails
	// bad ppm value (sit charge) fails

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID: uuid.FromStringOrNil("5b9645ea-8dae-40ab-9d25-46c0b00e6f98"),
		},
	})

	suite.T().Run("calculates ppm estimate success", func(t *testing.T) {
		pickupZip := "90210"
		ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				Move:             move,
				MoveID:           move.ID,
				PickupPostalCode: &pickupZip,
			},
		})
		// new duty station zip: 30813
		planner := route.NewTestingPlanner(300)
		calculator := NewEstimateCalculator(suite.DB(), suite.logger, planner)
		err := calculator.CalculateEstimate(&ppm, move.ID)
		suite.NoError(err)
		//suite.Equal(300, ppm.PlannedSITMax)
		//suite.Equal(300, ppm.SITMax)
		//suite.Equal(2000, ppm.IncentiveEstimateMin)
		//suite.Equal(3000, ppm.IncentiveEstimateMax)
	})

	suite.T().Run("receives a bad moveID fails", func(t *testing.T) {
		ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				Move:   move,
				MoveID: move.ID,
			},
		})
		planner := route.NewTestingPlanner(300)
		calculator := NewEstimateCalculator(suite.DB(), suite.logger, planner)

		nonExistentMoveID, err := uuid.FromString("2ef27bd2-97ae-4808-96cb-0cadd7f48972")
		if err != nil {
			suite.logger.Fatal("failure to get uuid from string")
		}
		err = calculator.CalculateEstimate(&ppm, nonExistentMoveID)
		suite.Error(err)
	})

	suite.T().Run("", func(t *testing.T) {

	})

	suite.T().Run("", func(t *testing.T) {

	})

	suite.T().Run("", func(t *testing.T) {

	})
}
