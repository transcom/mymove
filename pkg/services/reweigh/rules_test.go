package reweigh

//
//import (
//	"context"
//	"fmt"
//
//	"github.com/gobuffalo/validate/v3"
//	"github.com/gofrs/uuid"
//	"github.com/transcom/mymove/pkg/models"
//	"github.com/transcom/mymove/pkg/services"
//	"time"
//)
//
//func (suite *ReweighSuite) TestValidationRules() {
//	suite.Run("checkShipmentID", func() {
//		suite.Run("success", func() {
//			newReweigh := models.Reweigh{ShipmentID: uuid.Must(uuid.NewV4())}
//			testCases := map[string]struct {
//				newReweigh models.Reweigh
//				oldReweigh *models.Reweigh
//			}{
//				"create": {
//					newReweigh: newReweigh,
//					oldReweigh: nil,
//				},
//				"update": {
//					newReweigh: newReweigh,
//					oldReweigh: &models.Reweigh{ShipmentID: newReweigh.ShipmentID},
//				},
//			}
//			for name, testCase := range testCases {
//				suite.Run(name, func() {
//					err := checkShipmentID().Validate(context.Background(), testCase.newReweigh, testCase.oldReweigh, nil)
//					suite.NilOrNoVerrs(err)
//				})
//			}
//		})
//
//		suite.Run("failure", func() {
//			id := uuid.Must(uuid.NewV4())
//			testCases := map[string]struct {
//				newReweigh models.Reweigh
//				oldReweigh *models.Reweigh
//			}{
//				"create": {
//					newReweigh: models.Reweigh{},
//					oldReweigh: nil,
//				},
//				"update": {
//					newReweigh: models.Reweigh{ShipmentID: id},
//					oldReweigh: &models.Reweigh{},
//				},
//			}
//			for name, testCase := range testCases {
//				suite.Run(name, func() {
//					err := checkShipmentID().Validate(context.Background(), testCase.newReweigh, testCase.oldReweigh, nil)
//					switch verr := err.(type) {
//					case *validate.Errors:
//						suite.True(verr.HasAny())
//						suite.Contains(verr.Keys(), "ShipmentID")
//					default:
//						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
//					}
//				})
//			}
//		})
//	})
//
//	suite.Run("checkReweighID", func() {
//		suite.Run("success", func() {
//			id := uuid.Must(uuid.NewV4())
//			testCases := map[string]struct {
//				newReweigh models.Reweigh
//				oldReweigh *models.Reweigh
//			}{
//				"create": {
//					newReweigh: models.Reweigh{},
//					oldReweigh: nil,
//				},
//				"update": {
//					newReweigh: models.Reweigh{ID: id},
//					oldReweigh: &models.Reweigh{ID: id},
//				},
//			}
//			for name, testCase := range testCases {
//				suite.Run(name, func() {
//					err := checkReweighID().Validate(context.Background(), testCase.newReweigh, testCase.oldReweigh, nil)
//					suite.NilOrNoVerrs(err)
//				})
//			}
//		})
////
//		suite.Run("failure", func() {
//			id := uuid.Must(uuid.NewV4())
//			testCases := map[string]struct {
//				newReweigh models.Reweigh
//				oldReweigh *models.Reweigh
//				verr bool
//			}{
//				"create": {
//					newReweigh: models.Reweigh{ID: id},
//					oldReweigh: nil,
//					verr: true,
//				},
//				"update": {
//					newReweigh: models.Reweigh{ID: id},
//					oldReweigh: &models.Reweigh{ID: uuid.Must(uuid.NewV4())},
//					verr: false,
//				},
//			}
//			for name, testCase := range testCases {
//				suite.Run(name, func() {
//					err := checkReweighID().Validate(context.Background(), testCase.newReweigh, testCase.oldReweigh, nil)
//					switch verr := err.(type) {
//					case *validate.Errors:
//						suite.True(testCase.verr, "expected something other than a *validate.Errors type")
//						suite.Contains(verr.Keys(), "ID")
//					default:
//						suite.False(testCase.verr, "expected a *validate.Errors: %t - naid %s", err, testCase.newReweigh.ID)
//					}
//				})
//			}
//
//		})
//	})
//
//	suite.Run("checkRequiredFields", func() {
//		requestedAt := time.Now()
//		requestedBy := models.ReweighRequesterPrime
//
//		oldReweigh := &models.Reweigh{
//			RequestedAt: requestedAt,
//			RequestedBy:  requestedBy,
//		}
//
//		suite.Run("success", func() {
//			requestedAt := time.Now()
//			requestedBy := models.ReweighRequesterPrime
//
//			reweigh := models.Reweigh{
//				RequestedAt: requestedAt,
//				RequestedBy: requestedBy,
//			}
//
//			err := checkRequiredFields().Validate(context.Background(), reweigh, oldReweigh, nil)
//			switch verr := err.(type) {
//			case *validate.Errors:
//				suite.NoVerrs(verr)
//			default:
//				suite.Failf("expected *validate.Errs", "%v", err)
//			}
//		})
//
//		// Test unsuccessful check for required info
//		suite.Run("failure", func() {
//			//requestedAt := models.ReweighRequester.IsZero
//			//time := time.Time{}
//			requestedAt := new(time.Time) // this is the zero time, what we need to nullify the field
//			requestedBy := new(models.ReweighRequester)
//
//			reweigh := models.Reweigh{
//				RequestedAt: *requestedAt,
//				RequestedBy: *requestedBy,
//			}
//
//			err := checkRequiredFields().Validate(context.Background(), reweigh, oldReweigh, nil)
//			switch verr := err.(type) {
//			case *validate.Errors:
//				suite.False(verr.HasAny())
//				suite.Empty(verr.Keys())
//			default:
//				suite.Failf("expected *validate.Errs", "%v", err)
//			}
//		})
//	})
//
//	suite.Run("checkPrimeAvailability", func() {
//		shipment := models.MTOShipment{
//			ID: uuid.Must(uuid.NewV4()),
//		}
//		checkHappy := primeFunc(func(uuid.UUID) (bool, error) {
//			return true, nil
//		})
//		checkError := primeFunc(func(uuid.UUID) (bool, error) {
//			return false, fmt.Errorf("forced")
//		})
//		testCases := map[string]struct {
//			check services.MoveTaskOrderChecker
//			ship  *models.MTOShipment
//			err   bool
//		}{
//			"happy": {
//				check: checkHappy,
//				ship:  &shipment,
//				err:   false,
//			},
//			"error": {
//				check: checkError,
//				ship:  &shipment,
//				err:   true,
//			},
//			"misused": {
//				check: checkHappy,
//				ship:  nil,
//				err:   true,
//			},
//		}
//
//		for name, testCase := range testCases {
//			suite.Run(name, func() {
//				err := checkPrimeAvailability(testCase.check).Validate(context.Background(), models.Reweigh{}, nil, testCase.ship)
//				if err == nil {
//					if testCase.err {
//						suite.Fail("expected error")
//					}
//					return
//				}
//				suite.IsType(services.NotFoundError{}, err)
//			})
//		}
//
//	})
//}
//
//type primeFunc func(uuid.UUID) (bool, error)
//
//func (fn primeFunc) MTOAvailableToPrime(id uuid.UUID) (bool, error) {
//	return fn(id)
//}
