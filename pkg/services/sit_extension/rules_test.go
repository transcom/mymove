package sitextension

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SitExtensionServiceSuite) TestValidationRules() {
	suite.Run("checkShipmentID", func() {
		suite.Run("success", func() {
			sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				sit models.SITExtension
			}{
				"create": {
					sit: sit,
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(appcontext.NewAppContext(suite.DB(), suite.logger), testCase.sit, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			//id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				sit models.SITExtension
			}{
				"create": {
					sit: models.SITExtension{},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(appcontext.NewAppContext(suite.DB(), suite.logger), testCase.sit, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(verr.HasAny())
						suite.Contains(verr.Keys(), "MTOShipmentID")
					default:
						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
					}
				})
			}
		})
	})

	suite.Run("checkPrimeAvailability - Failure", func() {
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(appcontext.NewAppContext(suite.DB(), suite.logger), models.SITExtension{}, nil)
		suite.NotNil(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal("not found while looking for Prime-available Shipment", err.Error())
	})

	suite.Run("checkPrimeAvailability - Success", func() {
		currentTime := time.Now()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &currentTime,
				Status:             models.MoveStatusAPPROVED,
			},
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(appcontext.NewAppContext(suite.DB(), suite.logger), models.SITExtension{}, &shipment)
		suite.NoError(err)
	})
}
