package ppmshipment

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PPMShipmentSuite) TestValidationRules() {
	suite.Run("checkMTOShipmentID", func() {
		suite.Run("success", func() {
			//newPPMShipment := models.PPMShipment{ShipmentID: uuid.Must(uuid.NewV4())}
			ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
				Stub: true,
			})
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: ppmShipment,
					oldPPMShipment: nil,
				},
				//"update": {
				//	newPPMShipment: newPPMShipment,
				//	oldPPMShipment: &models.PPMShipment{ShipmentID: newPPMShipment.ShipmentID},
				//},
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})

		suite.Run("failure", func() {
			ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
				Stub: true,
			})
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: ppmShipment,
					oldPPMShipment: nil,
				},
				//"update": {
				//	newPPMShipment: models.PPMShipment{ShipmentID: id},
				//	oldPPMShipment: &models.PPMShipment{},
				//},
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
			ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
				Stub: true,
			})
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
			}{
				"create": {
					newPPMShipment: ppmShipment,
					oldPPMShipment: nil,
				},
				// Add Update Test case here
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkPPMShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					suite.NilOrNoVerrs(err)
				})
			}
		})
		//
		suite.Run("failure", func() {
			ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
				Stub: true,
			})
			//id := uuid.Must(uuid.NewV4())
			testCases := map[string]struct {
				newPPMShipment models.PPMShipment
				oldPPMShipment *models.PPMShipment
				verr           bool
			}{
				"create": {
					newPPMShipment: ppmShipment,
					oldPPMShipment: nil,
					verr:           true,
				},
				//"update": Add Update Test Case here
			}
			for name, testCase := range testCases {
				suite.Run(name, func() {
					err := checkPPMShipmentID().Validate(suite.AppContextForTest(), testCase.newPPMShipment, testCase.oldPPMShipment, nil)
					switch verr := err.(type) {
					case *validate.Errors:
						suite.True(testCase.verr, "expected something other than a *validate.Errors type")
						suite.Contains(verr.Keys(), "ID")
					default:
						suite.False(testCase.verr, "expected a *validate.Errors: %t - naid %s", err, testCase.newPPMShipment.ID)
					}
				})
			}

		})
	})
	suite.Run("CheckRequiredFields()", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		moveRouter := moverouter.NewMoveRouter()
		mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		//subtestData := suite.createSubtestData(testdatagen.Assertions{
		//	PPMshipment: models.PPMShipment{
		//		ExpectedDepartureDate: &expectedTime,
		//		PickupPostalCode:      &pickupPostal,
		//		DestinationPostalCode: &destPostalcode,
		//		SitExpected:           &sitExpected,
		//	}
		//})

		suite.Run("Success", func() {
			expectedTime := time.Now()
			pickupPostal := "99999"
			destPostalcode := "99999"
			sitExpected := false

			//newPPMShipment := models.PPMShipment{
			//	ExpectedDepartureDate: &expectedTime,
			//	PickupPostalCode:      &pickupPostal,
			//	DestinationPostalCode: &destPostalcode,
			//	SitExpected:           &sitExpected,
			//}

			ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
				PPMShipment: models.PPMShipment{
					ExpectedDepartureDate: &expectedTime,
					PickupPostalCode:      &pickupPostal,
					DestinationPostalCode: &destPostalcode,
					SitExpected:           &sitExpected,
				},
				Stub: true,
			})

			ppmShipmentCreator := NewPPMShipmentCreator(mtoShipmentCreator)
			createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentCheck(suite.AppContextForTest(), &ppmShipment)

			suite.Nil(err)
			suite.NotNil(createdPPMShipment)
		})
	})
}
