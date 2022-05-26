package ppmshipment

import (
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PPMShipmentSuite) TestValidationRules() {
	suite.Run("checkShipmentType", func() {
		suite.Run("success", func() {
			err := checkShipmentType().Validate(suite.AppContextForTest(), models.PPMShipment{}, nil, &models.MTOShipment{ShipmentType: models.MTOShipmentTypePPM})
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkShipmentType().Validate(suite.AppContextForTest(), models.PPMShipment{}, nil, &models.MTOShipment{ShipmentType: models.MTOShipmentTypeHHG})
			switch verr := err.(type) {
			case *validate.Errors:
				suite.True(verr.HasAny())
				suite.Contains(verr.Keys(), "ShipmentType")
			default:
				suite.Failf("expected *validate.Errors", "%t - %v", err, err)
			}
		})
	})

	suite.Run("checkMTOShipmentID", func() {
		suite.Run("success", func() {
			newPPMShipment := models.PPMShipment{ShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: newPPMShipment,
					oldPPMShipment: nil,
				},
				"update": {
					newPPMShipment: newPPMShipment,
					oldPPMShipment: &models.PPMShipment{ShipmentID: newPPMShipment.ShipmentID},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: models.PPMShipment{},
					oldPPMShipment: nil,
				},
				"update": {
					newPPMShipment: models.PPMShipment{ShipmentID: id1},
					oldPPMShipment: &models.PPMShipment{ShipmentID: id2},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
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

	suite.Run("checkPPMShipmentID", func() {
		suite.Run("success", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: models.PPMShipment{},
					oldPPMShipment: nil,
				},
				"update": {
					newPPMShipment: models.PPMShipment{ID: id},
					oldPPMShipment: &models.PPMShipment{ID: id},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkPPMShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
				verr           bool
			}{
				"create": {
					newPPMShipment: models.PPMShipment{ID: id1},
					oldPPMShipment: nil,
					verr:           true,
				},
				"update": {
					newPPMShipment: models.PPMShipment{ID: id1},
					oldPPMShipment: &models.PPMShipment{ID: id2},
					verr:           true,
				}}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkPPMShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(verr.HasAny())
						suite.Contains(verr.Keys(), "ID")
					default:
						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
					}
				})
			}

		})
	})
	suite.Run("CheckRequiredFields()", func() {
		expectedTime := time.Now()
		pickupPostal := "99999"
		destPostalcode := "99999"
		sitExpected := false
		shipmentID := uuid.Must(uuid.NewV4())

		suite.Run("success", func() {
			newPPMShipment := models.PPMShipment{
				ShipmentID:            shipmentID,
				ExpectedDepartureDate: expectedTime,
				PickupPostalCode:      pickupPostal,
				DestinationPostalCode: destPostalcode,
				SITExpected:           &sitExpected,
			}

			err := checkRequiredFields().Validate(suite.AppContextForTest(), newPPMShipment, nil, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("Failure - New shipment", func() {
			testCases := []struct {
				desc     string
				shipment models.PPMShipment
				errorKey string
				errorMsg string
			}{
				{
					"Missing expected departure date",
					models.PPMShipment{
						ShipmentID:            shipmentID,
						PickupPostalCode:      pickupPostal,
						DestinationPostalCode: destPostalcode,
						SITExpected:           &sitExpected,
					},
					"expectedDepartureDate",
					"cannot be a zero value"},
				{
					"Missing pickup postal code",
					models.PPMShipment{
						ShipmentID:            shipmentID,
						ExpectedDepartureDate: expectedTime,
						DestinationPostalCode: destPostalcode,
						SITExpected:           &sitExpected,
					},
					"pickupPostalCode",
					"cannot be nil or empty",
				},
				{
					"Missing destination postal code",
					models.PPMShipment{
						ShipmentID:            shipmentID,
						ExpectedDepartureDate: expectedTime,
						PickupPostalCode:      pickupPostal,
						SITExpected:           &sitExpected,
					},
					"destinationPostalCode",
					"cannot be nil or empty",
				},
				{
					"Missing SIT Expected value",
					models.PPMShipment{
						ShipmentID:            shipmentID,
						ExpectedDepartureDate: expectedTime,
						PickupPostalCode:      pickupPostal,
						DestinationPostalCode: destPostalcode,
						SITExpected:           nil,
					},
					"sitExpected",
					"cannot be nil",
				},
			}

			for _, tc := range testCases {
				tc := tc
				suite.Run(tc.desc, func() {
					err := checkRequiredFields().Validate(suite.AppContextForTest(), tc.shipment, nil, nil)

					switch verr := err.(type) {
					case *validate.Errors:
						suite.Equal(1, verr.Count())

						errorMsg, hasErrKey := verr.Errors[tc.errorKey]

						suite.True(hasErrKey)
						suite.Equal(tc.errorMsg, strings.Join(errorMsg, ""))
					default:
						suite.Failf("expected *validate.Errs", "%v", err)
					}
				})
			}
		})
	})

	suite.Run("CheckAdvance()", func() {
		suite.Run("Success", func() {
			suite.Run("advance set", func() {
				shipmentID := uuid.Must(uuid.NewV4())
				advanceRequested := false
				newAdvance := unit.Cents(10000)
				newAdvanceRequested := true
				estimatedIncentive := unit.Cents(17000)

				oldPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &advanceRequested,
					Advance:            nil,
				}

				newPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &newAdvanceRequested,
					Advance:            &newAdvance,
				}

				err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
				suite.NilOrNoVerrs(err)
			})
			suite.Run("advance set for first time", func() {
				shipmentID := uuid.Must(uuid.NewV4())
				newAdvance := unit.Cents(10000)
				newAdvanceRequested := true
				estimatedIncentive := unit.Cents(17000)

				oldPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   nil,
					Advance:            nil,
				}

				newPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &newAdvanceRequested,
					Advance:            &newAdvance,
				}

				err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
				suite.NilOrNoVerrs(err)
			})
			suite.Run("advanceRequested set from true to nil", func() {
				shipmentID := uuid.Must(uuid.NewV4())
				advanceRequested := true
				advance := unit.Cents(10000)
				newAdvanceRequested := false
				estimatedIncentive := unit.Cents(17000)

				oldPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &advanceRequested,
					Advance:            &advance,
				}

				newPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &newAdvanceRequested,
					Advance:            nil,
				}

				err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
				suite.NilOrNoVerrs(err)
			})
			suite.Run("advanceRequested set from nil to false", func() {
				shipmentID := uuid.Must(uuid.NewV4())
				newAdvanceRequested := false
				estimatedIncentive := unit.Cents(17000)

				oldPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   nil,
					Advance:            nil,
				}

				newPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   &newAdvanceRequested,
					Advance:            nil,
				}

				err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
				suite.NilOrNoVerrs(err)
			})
			suite.Run("advance stays nil during update", func() {
				shipmentID := uuid.Must(uuid.NewV4())
				estimatedIncentive := unit.Cents(17000)

				oldPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   nil,
					Advance:            nil,
				}

				newPPMShipment := models.PPMShipment{
					ShipmentID:         shipmentID,
					EstimatedIncentive: &estimatedIncentive,
					AdvanceRequested:   nil,
					Advance:            nil,
				}

				err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
				suite.NilOrNoVerrs(err)
			})
		})

		suite.Run("Failure", func() {
			id := uuid.Must(uuid.NewV4())
			estimatedIncentive := unit.Cents(17000)
			falsePointer := models.BoolPointer(false)
			truePointer := models.BoolPointer(true)
			zeroAdvance := unit.Cents(0)
			lessThanOneAdvance := unit.Cents(1) // amount less than $1
			normalAdvance := unit.Cents(10000)  // below 60%
			highAdvance := unit.Cents(12000)    // above 60%

			defaultOldShipmentValues := models.PPMShipment{
				ShipmentID:         id,
				EstimatedIncentive: &estimatedIncentive,
				AdvanceRequested:   nil,
				Advance:            nil,
			}

			testCases := map[string]struct {
				oldPPMShipment   models.PPMShipment
				newPPMShipment   models.PPMShipment
				expectedErrorMsg string
			}{
				"advance was requested but amount set to 0": {
					oldPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   falsePointer,
						Advance:            nil,
					},
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   truePointer,
						Advance:            &zeroAdvance,
					},
					expectedErrorMsg: "Advance cannot be a value less than $1",
				},
				"advance wasn't requested but amount isn't nil": {
					oldPPMShipment: defaultOldShipmentValues,
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   falsePointer,
						Advance:            &normalAdvance,
					},
					expectedErrorMsg: "Advance must be nil because of the advance requested value",
				},
				"advance set for greater than 60% of estimated incentive": {
					oldPPMShipment: defaultOldShipmentValues,
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   truePointer,
						Advance:            &highAdvance,
					},
					expectedErrorMsg: "Advance can not be greater than 60% of the estimated incentive",
				},
				"advance amount less than 1": {
					oldPPMShipment: defaultOldShipmentValues,
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   truePointer,
						Advance:            &lessThanOneAdvance,
					},
					expectedErrorMsg: "Advance cannot be a value less than $1",
				},
				"advance requested is nil but amount is not nil": {
					oldPPMShipment: defaultOldShipmentValues,
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   nil,
						Advance:            &normalAdvance,
					},
					expectedErrorMsg: "Advance must be nil because of the advance requested value",
				},
				"advance requested is true while advance is nil": {
					oldPPMShipment: defaultOldShipmentValues,
					newPPMShipment: models.PPMShipment{
						ShipmentID:         id,
						EstimatedIncentive: &estimatedIncentive,
						AdvanceRequested:   truePointer,
						Advance:            nil,
					},
					expectedErrorMsg: "An advance amount is required",
				},
			}

			for name, testCase := range testCases {
				name := name
				testCase := testCase

				suite.Run(name, func() {
					err := checkAdvance().Validate(suite.AppContextForTest(), testCase.newPPMShipment, &testCase.oldPPMShipment, nil)

					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(verr.HasAny())
						suite.Contains(verr.Keys(), "advance")
						suite.Equal(testCase.expectedErrorMsg, verr.Error())
					default:
						suite.Failf("expected *validate.Errors", "%t - %v", err, err)
					}
				})
			}
		})
	})

	suite.Run("CheckEstimatedWeight()", func() {
		suite.Run("success estimatedWeight set", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			estimatedWeight := unit.Pound(4000)
			estimatedIncentive := unit.Cents(17000)

			oldPPMShipment := models.PPMShipment{
				ShipmentID:      shipmentID,
				EstimatedWeight: &estimatedWeight,
			}

			newPPMShipment := models.PPMShipment{
				ShipmentID:         shipmentID,
				EstimatedWeight:    &estimatedWeight,
				EstimatedIncentive: &estimatedIncentive,
			}

			err := checkEstimatedWeight().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure estimatedWeight cannot be nil", func() {
			shipmentID := uuid.Must(uuid.NewV4())

			oldPPMShipment := models.PPMShipment{
				ShipmentID:      shipmentID,
				EstimatedWeight: nil,
			}

			newPPMShipment := models.PPMShipment{
				ShipmentID:      shipmentID,
				EstimatedWeight: nil,
			}

			err := checkEstimatedWeight().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
			suite.Equal("cannot be empty", err.Error())
		})
	})
}
