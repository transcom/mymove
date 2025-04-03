package sitentrydateupdate

import (
	"fmt"
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
		idfServiceItem, _ := setupInternationalModels()
		notFoundServiceItem := models.SITEntryDateUpdate{
			ID:           idfServiceItem.ID,
			SITEntryDate: idfServiceItem.SITEntryDate,
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

	suite.Run("Fails to update when DOFSIT entry date is after DOFSIT departure date", func() {
		today := models.TimePointer(time.Now())
		tomorrow := models.TimePointer(time.Now())
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		dofsitServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: tomorrow,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           dofsitServiceItem.ID,
			SITEntryDate: models.TimePointer(tomorrow.AddDate(0, 0, 1)),
		}
		_, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Entry Date (%s) must be before the SIT Departure Date (%s)",
			updatedServiceItem.SITEntryDate.Format("2006-01-02"),
			dofsitServiceItem.SITDepartureDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Fails to update when DOFSIT entry date is the same as DOFSIT departure date", func() {
		today := models.TimePointer(time.Now())
		tomorrow := models.TimePointer(time.Now())
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		dofsitServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: tomorrow,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           dofsitServiceItem.ID,
			SITEntryDate: tomorrow,
		}
		_, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Entry Date (%s) must be before the SIT Departure Date (%s)",
			updatedServiceItem.SITEntryDate.Format("2006-01-02"),
			dofsitServiceItem.SITDepartureDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Fails to update when DDFSIT entry date is after DDFSIT departure date", func() {
		today := models.TimePointer(time.Now())
		tomorrow := models.TimePointer(time.Now())
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		ddfsitServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: tomorrow,
				},
			},
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
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           ddfsitServiceItem.ID,
			SITEntryDate: models.TimePointer(tomorrow.AddDate(0, 0, 1)),
		}
		_, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Entry Date (%s) must be before the SIT Departure Date (%s)",
			updatedServiceItem.SITEntryDate.Format("2006-01-02"),
			ddfsitServiceItem.SITDepartureDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("Fails to update when DDFSIT entry date is the same as DDFSIT departure date", func() {
		today := models.TimePointer(time.Now())
		tomorrow := models.TimePointer(time.Now())
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		ddfsitServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     today,
					SITDepartureDate: tomorrow,
				},
			},
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
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           ddfsitServiceItem.ID,
			SITEntryDate: tomorrow,
		}
		_, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)
		suite.Error(err)
		expectedError := fmt.Sprintf(
			"the SIT Entry Date (%s) must be before the SIT Departure Date (%s)",
			updatedServiceItem.SITEntryDate.Format("2006-01-02"),
			ddfsitServiceItem.SITDepartureDate.Format("2006-01-02"),
		)
		suite.Contains(err.Error(), expectedError)
	})
	suite.Run("Successful update of service items - international", func() {
		idfServiceItem, idaServiceItem := setupInternationalModels()
		updatedServiceItem := models.SITEntryDateUpdate{
			ID:           idfServiceItem.ID,
			SITEntryDate: idfServiceItem.SITEntryDate,
		}
		newSitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)
		newSitEntryDateNextDay := newSitEntryDate.Add(24 * time.Hour)

		updatedServiceItem.SITEntryDate = &newSitEntryDate
		idaServiceItem.SITEntryDate = &newSitEntryDateNextDay

		changedServiceItem, err := updater.UpdateSitEntryDate(suite.AppContextForTest(), &updatedServiceItem)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(idfServiceItem.ID, updatedServiceItem.ID)
		suite.Equal(updatedServiceItem.SITEntryDate.Local(), changedServiceItem.SITEntryDate.Local())
		suite.Equal(idaServiceItem.SITEntryDate.Local(), newSitEntryDateNextDay.Local())
	})

}
