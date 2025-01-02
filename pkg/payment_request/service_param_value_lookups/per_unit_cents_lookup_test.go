package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestPerUnitCentsLookup() {
	key := models.ServiceItemParamNamePerUnitCents
	var mtoServiceItem models.MTOServiceItem
	setupTestData := func(serviceCode models.ReServiceCode) {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: serviceCode,
				},
			},
		}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

	}

	suite.Run("success - returns perUnitCent value for IHPK", func() {
		setupTestData(models.ReServiceCodeIHPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "6997")
	})

	suite.Run("success - returns perUnitCent value for IHUPK", func() {
		setupTestData(models.ReServiceCodeIHUPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "752")
	})

	suite.Run("success - returns perUnitCent value for ISLH", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Des Moines",
					State:          "IA",
					PostalCode:     "50314",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(perUnitCents, "1605")
	})

	suite.Run("failure - unauthorized service code", func() {
		setupTestData(models.ReServiceCodeDUPK)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		perUnitCents, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Equal(perUnitCents, "")
	})
}
