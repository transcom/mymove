package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestPortZipLookup() {
	key := models.ServiceItemParamNamePortZip
	var mtoServiceItem models.MTOServiceItem
	setupTestData := func(serviceCode models.ReServiceCode, portID uuid.UUID) {
		if serviceCode == models.ReServiceCodePOEFSC {
			mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: serviceCode,
					},
				},
				{
					Model: models.MTOServiceItem{
						POELocationID: &portID,
					},
				},
			}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
		} else {
			mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: serviceCode,
					},
				},
				{
					Model: models.MTOServiceItem{
						PODLocationID: &portID,
					},
				},
			}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
		}
	}

	suite.Run("success - returns PortZip value for POEFSC", func() {
		port := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "SEA",
				},
			},
		}, nil)
		setupTestData(models.ReServiceCodePOEFSC, port.ID)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		portZip, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(portZip, port.UsPostRegionCity.UsprZipID)
	})

	suite.Run("success - returns PortZip value for PODFSC", func() {
		port := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		setupTestData(models.ReServiceCodePODFSC, port.ID)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		portZip, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal(portZip, port.UsPostRegionCity.UsprZipID)
	})

	suite.Run("success - returns PortZip value for Port Code 4E1 for PPMs", func() {
		port := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "4E1",
				},
			},
		}, nil)

		contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualPickupDate: models.TimePointer(time.Now()),
					MarketCode:       models.MarketCodeInternational,
				},
			},
			{
				Model: models.Address{
					StreetAddress1: "Tester Address",
					City:           "Tulsa",
					State:          "OK",
					PostalCode:     "74133",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "JBER",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		portZipLookup := PortZipLookup{
			ServiceItem: mtoServiceItem,
		}

		appContext := suite.AppContextForTest()
		portZip, err := portZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppm.ShipmentID,
			ContractID:    contractYear.ContractID,
		})
		suite.NoError(err)
		suite.Equal(portZip, port.UsPostRegionCity.UsprZipID)
	})

	suite.Run("returns nothing if shipment is HHG and service item does not have port info", func() {
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					POELocationID: nil,
				},
			},
		}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		portZip, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Equal(portZip, "")
	})
}
