package reweigh

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ReweighSuite) TestValidationRules() {
	suite.Run("checkShipmentID", func() {
		suite.Run("success", func() {
			newReweigh := models.Reweigh{ShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				newReweigh models.Reweigh
				oldReweigh *models.Reweigh
			}{
				"create": {
					newReweigh: newReweigh,
					oldReweigh: nil,
				},
				"update": {
					newReweigh: newReweigh,
					oldReweigh: &models.Reweigh{ShipmentID: newReweigh.ShipmentID},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newReweigh, testCase.oldReweigh, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newReweigh models.Reweigh
				oldReweigh *models.Reweigh
			}{
				"create": {
					newReweigh: models.Reweigh{},
					oldReweigh: nil,
				},
				"update": {
					newReweigh: models.Reweigh{ShipmentID: id},
					oldReweigh: &models.Reweigh{},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newReweigh, testCase.oldReweigh, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(verr.HasAny())
						suite.Contains(verr.Keys(), "ShipmentID")
					default:
						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
					}
				})
			}
		})
	})

	suite.Run("checkReweighID", func() {
		suite.Run("success", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newReweigh models.Reweigh
				oldReweigh *models.Reweigh
			}{
				"create": {
					newReweigh: models.Reweigh{},
					oldReweigh: nil,
				},
				"update": {
					newReweigh: models.Reweigh{ID: id},
					oldReweigh: &models.Reweigh{ID: id},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkReweighID().Validate(suite.AppContextForTest(), testCase.newReweigh, testCase.oldReweigh, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})
		//
		suite.Run("failure", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newReweigh models.Reweigh
				oldReweigh *models.Reweigh
				verr       bool
			}{
				"create": {
					newReweigh: models.Reweigh{ID: id},
					oldReweigh: nil,
					verr:       true,
				},
				"update": {
					newReweigh: models.Reweigh{ID: id},
					oldReweigh: &models.Reweigh{ID: uuid.Must(uuid.NewV4())},
					verr:       false,
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkReweighID().Validate(suite.AppContextForTest(), testCase.newReweigh, testCase.oldReweigh, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(testCase.verr, "expected something other than a *validate.Errors type")
						suite.Contains(verr.Keys(), "ID")
					default:
						suite.False(testCase.verr, "expected a *validate.Errors: %t - naid %s", err, testCase.newReweigh.ID)
					}
				})
			}

		})
	})

	suite.Run("checkRequiredFields", func() {
		requestedAt := time.Now()
		requestedBy := models.ReweighRequesterPrime

		oldReweigh := &models.Reweigh{
			RequestedAt: requestedAt,
			RequestedBy: requestedBy,
		}

		suite.Run("success", func() {
			requestedAt := time.Now()
			requestedBy := models.ReweighRequesterPrime

			reweigh := models.Reweigh{
				RequestedAt: requestedAt,
				RequestedBy: requestedBy,
			}

			err := checkRequiredFields().Validate(suite.AppContextForTest(), reweigh, oldReweigh, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.NoVerrs(verr)
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})

		// Test unsuccessful check for required info
		suite.Run("failure", func() {
			requestedAt := new(time.Time) // this is the zero time, what we need to nullify the field
			requestedBy := new(models.ReweighRequester)

			reweigh := models.Reweigh{
				RequestedAt: *requestedAt,
				RequestedBy: *requestedBy,
			}

			err := checkRequiredFields().Validate(suite.AppContextForTest(), reweigh, oldReweigh, nil)
			switch verr := err.(type) {
			case *validate.Errors:
				suite.False(verr.HasAny())
				suite.Empty(verr.Keys())
			default:
				suite.Failf("expected *validate.Errs", "%v", err)
			}
		})
	})

	suite.Run("checkPrimeAvailability - Failure because not available to prime", func() {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: nil,
			},
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.Reweigh{}, nil, &shipment)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("Not found while looking for Prime-available Shipment with id: %s", shipment.ID), err.Error())
	})

	suite.Run("checkPrimeAvailability - Failure because external vendor shipment", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
			},
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.Reweigh{}, nil, &shipment)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("Not found while looking for Prime-available Shipment with id: %s", shipment.ID), err.Error())
	})

	suite.Run("checkPrimeAvailability - Success", func() {
		currentTime := time.Now()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &currentTime,
			},
		})
		checker := movetaskorder.NewMoveTaskOrderChecker()
		err := checkPrimeAvailability(checker).Validate(suite.AppContextForTest(), models.Reweigh{}, nil, &shipment)
		suite.NoError(err)
	})
}
