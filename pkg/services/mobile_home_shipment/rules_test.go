package mobilehomeshipment

import (
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *MobileHomeShipmentSuite) TestValidationRules() {
	suite.Run("checkMTOShipmentID", func() {
		suite.Run("success", func() {
			newMobileHomeShipment := models.MobileHome{ShipmentID: uuid.Must(uuid.NewV4())}
			testCases := map[string]struct {
				newMobileHomeShipment models.MobileHome
				oldMobileHomeShipment *models.MobileHome
			}{
				"create": {
					newMobileHomeShipment: newMobileHomeShipment,
					oldMobileHomeShipment: nil,
				},
				"update": {
					newMobileHomeShipment: newMobileHomeShipment,
					oldMobileHomeShipment: &models.MobileHome{ShipmentID: newMobileHomeShipment.ShipmentID},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newMobileHomeShipment, testCase.oldMobileHomeShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newMobileHomeShipment models.MobileHome
				oldMobileHomeShipment *models.MobileHome
			}{
				"create": {
					newMobileHomeShipment: models.MobileHome{},
					oldMobileHomeShipment: nil,
				},
				"update": {
					newMobileHomeShipment: models.MobileHome{ShipmentID: id1},
					oldMobileHomeShipment: &models.MobileHome{ShipmentID: id2},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newMobileHomeShipment, testCase.oldMobileHomeShipment, nil)
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

	suite.Run("checkMobileHomeShipmentID", func() {
		suite.Run("success", func() {
			id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newMobileHomeShipment models.MobileHome
				oldMobileHomeShipment *models.MobileHome
			}{
				"create": {
					newMobileHomeShipment: models.MobileHome{},
					oldMobileHomeShipment: nil,
				},
				"update": {
					newMobileHomeShipment: models.MobileHome{ID: id},
					oldMobileHomeShipment: &models.MobileHome{ID: id},
				},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkMobileHomeShipmentID().Validate(suite.AppContextForTest(), testCase.newMobileHomeShipment, testCase.oldMobileHomeShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			id1 := uuid.Must(uuid.NewV4())
			id2 := uuid.Must(uuid.NewV4())

			testCases := map[string]struct {
				newMobileHomeShipment models.MobileHome
				oldMobileHomeShipment *models.MobileHome
				verr                  bool
			}{
				"create": {
					newMobileHomeShipment: models.MobileHome{ID: id1},
					oldMobileHomeShipment: nil,
					verr:                  true,
				},
				"update": {
					newMobileHomeShipment: models.MobileHome{ID: id1},
					oldMobileHomeShipment: &models.MobileHome{ID: id2},
					verr:                  true,
				}}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkMobileHomeShipmentID().Validate(suite.AppContextForTest(), testCase.newMobileHomeShipment, testCase.oldMobileHomeShipment, nil)
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
			newMobileHomeShipment := models.MobileHome{
				ShipmentID:     shipmentID,
				Year:           models.IntPointer(2000),
				Make:           models.StringPointer("Mobile Home Make"),
				Model:          models.StringPointer("Mobile Home Model"),
				LengthInInches: models.IntPointer(300),
				WidthInInches:  models.IntPointer(108),
				HeightInInches: models.IntPointer(72),
			}

			err := checkRequiredFields().Validate(suite.AppContextForTest(), newMobileHomeShipment, nil, nil)
			suite.NilOrNoVerrs(err)
		})

		suite.Run("Failure - New shipment", func() {
			testCases := []struct {
				desc     string
				shipment models.MobileHome
				errorKey string
				errorMsg string
			}{
				{
					"Missing Year Expected value",
					models.MobileHome{
						ShipmentID:     shipmentID,
						Make:           models.StringPointer("Mobile Home Make"),
						Model:          models.StringPointer("Mobile Home Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						Year:           models.IntPointer(-1),
					},
					"year",
					"cannot be a zero or a negative value"},
				{
					"Missing Make Expected value",
					models.MobileHome{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Model:          models.StringPointer("Mobile Home Model"),
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
					models.MobileHome{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Mobile Home Make"),
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
					models.MobileHome{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Mobile Home Make"),
						Model:          models.StringPointer("Mobile Home Model"),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: models.IntPointer(72),
						LengthInInches: nil,
					},
					"lengthInInches",
					"cannot be a zero or a negative value",
				},
				{
					"Missing WidthInInches Expected value",
					models.MobileHome{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Mobile Home Make"),
						Model:          models.StringPointer("Mobile Home Model"),
						LengthInInches: models.IntPointer(300),
						HeightInInches: models.IntPointer(72),
						WidthInInches:  nil,
					},
					"widthInInches",
					"cannot be a zero or a negative value"},
				{
					"Missing HeightInInches Expected value",
					models.MobileHome{
						ShipmentID:     shipmentID,
						Year:           models.IntPointer(2000),
						Make:           models.StringPointer("Mobile Home Make"),
						Model:          models.StringPointer("Mobile Home Model"),
						LengthInInches: models.IntPointer(300),
						WidthInInches:  models.IntPointer(108),
						HeightInInches: nil,
					},
					"heightInInches",
					"cannot be a zero or a negative value"},
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
