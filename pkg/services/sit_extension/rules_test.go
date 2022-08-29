package sitextension

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *SitExtensionServiceSuite) TestValidationRules() {
	suite.Run("checkShipmentID", func() {
		suite.Run("success", func() {
			sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}
			err := checkShipmentID().Validate(suite.AppContextForTest(), sit, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			var sit models.SITExtension
			err := checkShipmentID().Validate(suite.AppContextForTest(), sit, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "MTOShipmentID")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkRequiredFields", func() {
		//takes an app context& sit extension
		//returns a verification error
		suite.Run("success", func() {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
			})
			sitExtension := testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
				MTOShipment: shipment,
				SITExtension: models.SITExtension{
					MTOShipmentID: shipment.ID,
					RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
					Status:        models.SITExtensionStatusApproved,
					RequestedDays: 90,
				},
			})

			err := checkRequiredFields().Validate(suite.AppContextForTest(), sitExtension, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.NoVerrs(verr)
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})

		suite.Run("failure", func() {
			var sit models.SITExtension
			err := checkRequiredFields().Validate(suite.AppContextForTest(), sit, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "RequestedDays")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkSITExtensionPending - Success", func() {
		// Testing: There is no new sit extension
		sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})
		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Success after existing SIT is Approved", func() {
		// Testing: There is no new sit extension
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})

		// Approved Status SIT Extension
		// Changed Request Reason from the default
		testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipment,
			SITExtension: models.SITExtension{
				MTOShipmentID: shipment.ID,
				RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
				Status:        models.SITExtensionStatusApproved,
				RequestedDays: 90,
			},
		})
		sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Success after existing SIT is Denied", func() {
		// Testing: There is no new sit extension
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})

		// Denied SIT Extension
		testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipment,
			SITExtension: models.SITExtension{
				MTOShipmentID: shipment.ID,
				RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
				Status:        models.SITExtensionStatusDenied,
				RequestedDays: 90,
			},
		})
		sit := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4())}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), sit, &shipment)

		suite.NoError(err)
	})

	suite.Run("checkSITExtensionPending - Failure", func() {
		// Testing: There is a SIT extension and trying to be created
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})

		// Create SIT Extension #1 in DB
		// Change default status to Pending:
		testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipment,
			SITExtension: models.SITExtension{
				MTOShipmentID: shipment.ID,
				RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
				Status:        models.SITExtensionStatusPending,
				RequestedDays: 90,
			},
		})
		// Object we are trying to add to DB
		newSIT := models.SITExtension{MTOShipmentID: uuid.Must(uuid.NewV4()), Status: models.SITExtensionStatusPending, RequestedDays: 4}

		err := checkSITExtensionPending().Validate(suite.AppContextForTest(), newSIT, &shipment)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("checkPrimeAvailability - Failure", func() {
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.SITExtension{}, nil)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal("Not found while looking for Prime-available Shipment", err.Error())
	})

	suite.Run("checkPrimeAvailability - Success", func() {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()), // Move status is automatically set to APPROVED
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.SITExtension{}, &shipment)
		suite.NoError(err)
	})
}
