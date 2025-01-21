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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
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

	suite.Run("failure - no port zip on service item", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
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

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
}
