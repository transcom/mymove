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
				SitExpected:           sitExpected,
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
						SitExpected:           sitExpected,
					},
					"expectedDepartureDate",
					"cannot be a zero value"},
				{
					"Missing pickup postal code",
					models.PPMShipment{
						ShipmentID:            shipmentID,
						ExpectedDepartureDate: expectedTime,
						DestinationPostalCode: destPostalcode,
						SitExpected:           sitExpected,
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
						SitExpected:           sitExpected,
					},
					"destinationPostalCode",
					"cannot be nil or empty",
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
		suite.Run("success advance set", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			advanceRequested := false
			newAdvance := unit.Cents(10000)
			newAdvanceRequested := true
			estimatedIncentive := int32(17000)

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
		suite.Run("success advance set for first time", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			newAdvance := unit.Cents(10000)
			newAdvanceRequested := true
			estimatedIncentive := int32(17000)

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
		suite.Run("success advanceRequested set from true to nil", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			advanceRequested := true
			advance := unit.Cents(10000)
			newAdvanceRequested := false
			estimatedIncentive := int32(17000)

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
		suite.Run("success advanceRequested set from nil to false", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			newAdvanceRequested := false
			estimatedIncentive := int32(17000)

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
		suite.Run("success advance stays nil during update", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			estimatedIncentive := int32(17000)

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

	suite.Run("CheckAdvance()", func() {
		suite.Run("failure - advance not set", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			advanceRequested := false
			newAdvance := unit.Cents(0)
			newAdvanceRequested := true
			estimatedIncentive := int32(17000)

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
			suite.Error(err)
		})

		suite.Run("failure - advance not nil", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			newAdvanceRequested := false
			newAdvance := unit.Cents(10000)
			estimatedIncentive := int32(17000)

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
			suite.Error(err)
		})

		suite.Run("failure - advance less than 1", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			newAdvanceRequested := true
			newAdvance := unit.Cents(1)
			estimatedIncentive := int32(17000)

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
			suite.Error(err)
		})
		suite.Run("failure - advance is not nil", func() {
			shipmentID := uuid.Must(uuid.NewV4())
			newAdvance := unit.Cents(10000)
			estimatedIncentive := int32(17000)

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
				Advance:            &newAdvance,
			}

			err := checkAdvance().Validate(suite.AppContextForTest(), newPPMShipment, &oldPPMShipment, nil)
			suite.Error(err)
		})
	})
}
