package sitextension

import (
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
			err := checkShipmentID().Validate(appcontext.NewAppContext(suite.DB(), suite.logger), sit, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			var sit models.SITExtension
			err := checkShipmentID().Validate(appcontext.NewAppContext(suite.DB(), suite.logger), sit, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "MTOShipmentID")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkSITExtensionPending - Success", func() {
		// Testing: There is no new sit extension
		//sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}
		//err := checkSITExtensionPending().Validate
	})

	suite.Run("checkSITExtensionPending - Failure", func() {

	})

	suite.Run("checkPrimeAvailability - Failure", func() {
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(appcontext.NewAppContext(suite.DB(), suite.logger), models.SITExtension{}, nil)
		suite.NotNil(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Equal("not found while looking for Prime-available Shipment", err.Error())
	})

	suite.Run("checkPrimeAvailability - Success", func() {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(appcontext.NewAppContext(suite.DB(), suite.logger), models.SITExtension{}, &shipment)
		suite.NoError(err)
	})
}
