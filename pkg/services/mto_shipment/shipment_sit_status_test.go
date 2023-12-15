package mtoshipment

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MTOShipmentServiceSuite) TestShipmentSITStatus() {
	sitStatusService := NewShipmentSITStatus()

	suite.Run("returns nil when the shipment has no service items", func() {
		submittedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), submittedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("returns nil when the shipment has no SIT service items", func() {
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
					// TODO: Come back and add these service items to customizations
					//MTOServiceItems: testdatagen.MakeMTOServiceItems(suite.DB()),
				},
			},
		}, nil)

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("returns SIT Status when the shipment has a SIT service item with entry date in the future", func() {
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		nextWeek := time.Now().Add(time.Hour * 24 * 7)
		futureSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &nextWeek,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{futureSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
	})

	suite.Run("includes SIT service item that has departed storage", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		dofsit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &aMonthAgo,
					SITDepartureDate: &fifteenDaysAgo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dofsit}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(dofsit.ID.String(), sitStatus.PastSITs[0].ID.String())

		suite.Equal(15, sitStatus.TotalSITDaysUsed)
		suite.Equal(15, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(75, sitStatus.TotalDaysRemaining)
		suite.Nil(sitStatus.CurrentSIT) // No current SIT since all SIT items have departed status
	})

	suite.Run("calculates status for a shipment currently in SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dofsit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aMonthAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dofsit}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(30, sitStatus.TotalSITDaysUsed)
		suite.Equal(30, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(60, sitStatus.TotalDaysRemaining)
		suite.Equal(30, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aMonthAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())
		suite.Len(sitStatus.PastSITs, 0)
	})

	suite.Run("combines SIT days sum for shipment with past and current SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &aMonthAgo,
					SITDepartureDate: &fifteenDaysAgo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aWeekAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOFSIT, currentDOFSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(22, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(22, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(68, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOFSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("combines SIT days sum for shipment with past origin and current destination SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &aMonthAgo,
					SITDepartureDate: &fifteenDaysAgo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aWeekAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOFSIT, currentDDFSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(DestinationSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(22, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(22, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(68, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOFSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("returns negative days remaining when days in SIT exceeds shipment allowance", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * 30 * -6).Date()
		sixMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		threeMonthsAgo := sixMonthsAgo.Add(time.Hour * 24 * 30 * 3)
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &sixMonthsAgo,
					SITDepartureDate: &threeMonthsAgo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aWeekAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOFSIT, currentDDFSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(DestinationSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(97, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(97, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(-7, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOFSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("excludes SIT service items that have not been approved by the TOO", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * 30 * -6).Date()
		sixMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		threeMonthsAgo := sixMonthsAgo.Add(time.Hour * 24 * 30 * 3)
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &sixMonthsAgo,
					SITDepartureDate: &threeMonthsAgo,
					Status:           models.MTOServiceItemStatusRejected,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aWeekAgo,
					Status:       models.MTOServiceItemStatusRejected,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOFSIT, currentDDFSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})
}
