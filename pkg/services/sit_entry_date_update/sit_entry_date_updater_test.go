package sitentrydateupdate

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *UpdateSitEntryDateServiceSuite) TestUpdateSitEntryDate() {

	updater := NewSitEntryDateUpdater()

	// setting up a shipment model with multiple service items since sister items are checked and updated
	setupModels := func() (models.MTOServiceItem, models.MTOServiceItem) {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		ddfServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)
		ddaServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
		}, nil)
		return ddfServiceItem, ddaServiceItem
	}

	setupInternationalModels := func() (models.MTOServiceItem, models.MTOServiceItem) {
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					MarketCode: models.MarketCodeInternational,
				},
			},
		}, nil)
		idfServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDFSIT,
				},
			},
		}, nil)
		idaServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDASIT,
				},
			},
		}, nil)
		return idfServiceItem, idaServiceItem
	}

	// Test not found error
	suite.Run("Not Found Error", func() {
		ddfServiceItem, _ := setupModels()
		notFoundServiceItem := models.SITEntryDateUpdate{
			ID:           ddfServiceItem.ID,
			SITEntryDate: ddfServiceItem.SITEntryDate,
		}
		notFoundUUID, err := uuid.NewV4()
		suite.NoError(err)
		notFoundServiceItem.ID = notFoundUUID

		updatedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &notFoundServiceItem)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Not Found Error - international", func() {
		ddfServiceItem, _ := setupInternationalModels()
		notFoundServiceItem := models.SITEntryDateUpdate{
			ID:           ddfServiceItem.ID,
			SITEntryDate: ddfServiceItem.SITEntryDate,
		}
		notFoundUUID, err := uuid.NewV4()
		suite.NoError(err)
		notFoundServiceItem.ID = notFoundUUID

		updatedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &notFoundServiceItem)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	// Test successful update of both service items
	suite.Run("Successful update of service items", func() {
		ddfServiceItem, ddaServiceItem := setupModels()
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           ddfServiceItem.ID,
			SITEntryDate: ddfServiceItem.SITEntryDate,
		}
		newSitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)
		newSitEntryDateNextDay := newSitEntryDate.Add(24 * time.Hour)

		updatedServiceItem.SITEntryDate = &newSitEntryDate
		ddaServiceItem.SITEntryDate = &newSitEntryDateNextDay

		changedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(ddfServiceItem.ID, updatedServiceItem.ID)
		suite.Equal(updatedServiceItem.SITEntryDate.Local(), changedServiceItem.SITEntryDate.Local())
		suite.Equal(ddaServiceItem.SITEntryDate.Local(), newSitEntryDateNextDay.Local())
	})

	suite.Run("Successful update of service items - international", func() {
		ddfServiceItem, ddaServiceItem := setupInternationalModels()
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           ddfServiceItem.ID,
			SITEntryDate: ddfServiceItem.SITEntryDate,
		}
		newSitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)
		newSitEntryDateNextDay := newSitEntryDate.Add(24 * time.Hour)

		updatedServiceItem.SITEntryDate = &newSitEntryDate
		ddaServiceItem.SITEntryDate = &newSitEntryDateNextDay

		changedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(ddfServiceItem.ID, updatedServiceItem.ID)
		suite.Equal(updatedServiceItem.SITEntryDate.Local(), changedServiceItem.SITEntryDate.Local())
		suite.Equal(ddaServiceItem.SITEntryDate.Local(), newSitEntryDateNextDay.Local())
	})

}
