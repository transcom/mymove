package boatshipment

import (
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *BoatShipmentSuite) TestValidationRules() {
	suite.Run("checkShipmentType", func() {
		suite.Run("success", func() {
			err := checkShipmentType().Validate(suite.AppContextForTest(), models.BoatShipment{}, nil, &models.MTOShipment{ShipmentType: models.MTOShipmentTypeBoatHaulAway})
			suite.NilOrNoVerrs(err)
		})

		suite.Run("failure", func() {
			err := checkShipmentType().Validate(suite.AppContextForTest(), models.BoatShipment{}, nil, &models.MTOShipment{ShipmentType: models.MTOShipmentTypeHHG})
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
			newBoatShipment := models.BoatShipment{ShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				newBoatShipment models.BoatShipment
				oldBoatShipment *models.BoatShipment
			}{
				"create": {
					newBoatShipment: newBoatShipment,
					oldBoatShipment: nil,
				},
				"update": {
					newBoatShipment: newBoatShipment,
					oldBoatShipment: &models.BoatShipment{ShipmentID: newBoatShipment.ShipmentID},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newBoatShipment, testCase.oldBoatShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newBoatShipment models.BoatShipment
				oldBoatShipment *models.BoatShipment
			}{
				"create": {
					newBoatShipment: models.BoatShipment{},
					oldBoatShipment: nil,
				},
				"update": {
					newBoatShipment: models.BoatShipment{ShipmentID: id1},
					oldBoatShipment: &models.BoatShipment{ShipmentID: id2},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newBoatShipment, testCase.oldBoatShipment, nil)
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

	suite.Run("checkBoatShipmentID", func() {
		suite.Run("success", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newBoatShipment models.BoatShipment
				oldBoatShipment *models.BoatShipment
			}{
				"create": {
					newBoatShipment: models.BoatShipment{},
					oldBoatShipment: nil,
				},
				"update": {
					newBoatShipment: models.BoatShipment{ID: id},
					oldBoatShipment: &models.BoatShipment{ID: id},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkBoatShipmentID().Validate(suite.AppContextForTest(), testCase.newBoatShipment, testCase.oldBoatShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newBoatShipment models.BoatShipment
				oldBoatShipment *models.BoatShipment
				verr            bool
			}{
				"create": {
					newBoatShipment: models.BoatShipment{ID: id1},
					oldBoatShipment: nil,
					verr:            true,
				},
				"update": {
					newBoatShipment: models.BoatShipment{ID: id1},
					oldBoatShipment: &models.BoatShipment{ID: id2},
					verr:            true,
				}}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkBoatShipmentID().Validate(suite.AppContextForTest(), testCase.newBoatShipment, testCase.oldBoatShipment, nil)
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
		shipmentID := uuid.Must(uuid.NewV4())

		suite.Run("success", func() {
			newBoatShipment := models.BoatShipment{
				ShipmentID:     shipmentID,
				Type:           models.BoatShipmentTypeHaulAway,
				Year:           models.IntPointer(2000),
				Make:           models.StringPointer("Boat Make"),
				Model:          models.StringPointer("Boat Model"),
				LengthInInches: models.IntPointer(300),
				WidthInInches:  models.IntPointer(108),
				HeightInInches: models.IntPointer(72),
				HasTrailer:     models.BoolPointer(true),
				IsRoadworthy:   models.BoolPointer(false),
			}

			err := checkRequiredFields().Validate(suite.AppContextForTest(), newBoatShipment, nil, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("Failure - New shipment", func() {
			testCases := []struct {
				desc     string
				shipment models.BoatShipment
				errorKey string
				errorMsg string
			}{
				{
					"Missing Year Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						Year:           models.IntPointer(-1),
					},
					"year",
					"cannot be a zero or a negative value"},
				{
					"Missing Make Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						Make:           models.StringPointer(""),
					},
					"make",
					"cannot be empty",
				},
				{
					"Missing Model Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						Model:          models.StringPointer(""),
					},
					"model",
					"cannot be empty",
				},
				{
					"Missing LengthInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						LengthInInches: nil,
					},
					"lengthInInches",
					"cannot be a zero or a negative value",
				},
				{
					"Invalid LengthInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(10),
						WidthInInches:  models.IntPointer(10),
						HeightInInches: models.IntPointer(10),
					},
					"lengthInInches",
					"One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77."},
				{
					"Missing WidthInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(300),
						HeightInInches: models.IntPointer(72),
						WidthInInches:  nil,
					},
					"widthInInches",
					"cannot be a zero or a negative value"},
				{
					"Invalid WidthInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(10),
						WidthInInches:  models.IntPointer(10),
						HeightInInches: models.IntPointer(10),
					},
					"widthInInches",
					"One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77."},
				{
					"Missing HeightInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: nil,
					},
					"heightInInches",
					"cannot be a zero or a negative value"},
				{
					"Invalid HeightInInches Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(10),
						WidthInInches:  models.IntPointer(10),
						HeightInInches: models.IntPointer(10),
					},
					"heightInInches",
					"One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77."},
				{

					"Invalid isRoadworthy Expected value",
					models.BoatShipment{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Boat Make"),
						Model:          models.StringPointer("Boat Model"),
						LengthInInches: models.IntPointer(1000),
						WidthInInches:  models.IntPointer(10),
						HeightInInches: models.IntPointer(10),
						HasTrailer:     models.BoolPointer(true),
						IsRoadworthy:   nil,
					},
					"isRoadworthy",
					"isRoadworthy is required if hasTrailer is true"},
			}

			for _, tc := range testCases {
				tc := tc
				suite.Run(tc.desc, func() {
					err := checkRequiredFields().Validate(suite.AppContextForTest(), tc.shipment, nil, nil)

					switch verr := err.(type) {
					case *validate.Errors:

						errorMsg, hasErrKey := verr.Errors[tc.errorKey]

						suite.True(hasErrKey)
						suite.Equal(tc.errorMsg, strings.Join(errorMsg, ""))
						if tc.desc == "Invalid LengthInInches Expected value" || tc.desc == "Invalid WidthInInches Expected value" || tc.desc == "Invalid HeightInInches Expected value" {
							suite.Equal(3, verr.Count())
						} else {
							suite.Equal(1, verr.Count())
						}
					default:
						suite.Failf("expected *validate.Errs", "%v", err)
					}
				})
			}
		})
	})
}
