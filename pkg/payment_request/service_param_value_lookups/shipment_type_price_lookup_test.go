package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestShipmentTypePriceLookup() {
	key := models.ServiceItemParamNameNTSPackingFactor

	setupShipmentTypeLookupData := func(code models.ReServiceCode) (models.Move, models.MTOServiceItem) {
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		startDate := time.Now().Add(-24 * time.Hour)
		endDate := time.Now().Add(24 * time.Hour)
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		reService := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: code,
				},
			},
		}, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		return move, mtoServiceItem
	}

	suite.Run("lookup success", func() {
		// Set up INPK data and then look for it
		move, mtoServiceItem := setupShipmentTypeLookupData(models.ReServiceCodeINPK)
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		factor, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.NotEmpty(factor)
		suite.Equal("1.45", factor) // Hard coded dummy factor not truncated within the db
	})
	suite.Run("lookup fail", func() {
		// Set up FSC data and then look for it
		// This will fail because FSC does not have a market factor
		move, mtoServiceItem := setupShipmentTypeLookupData(models.ReServiceCodeFSC)
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)
		factor, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Empty(factor)
	})
	suite.Run("lookup optional", func() {
		// IHPK doesn't need the INPK ShipmentTypePrice if
		// IHPK is the one being priced and not INPK
		move, mtoServiceItem := setupShipmentTypeLookupData(models.ReServiceCodeIHPK)
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)

		factor, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Empty(factor)
	})
}
