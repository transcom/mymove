package sitstatus

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *SITStatusServiceSuite) TestShipmentSITStatus() {
	sitStatusService := NewShipmentSITStatus()

	suite.Run("returns the clamped values", func() {
		lowNum := 3
		highNum := 99
		clampedNumber, err := Clamp(1, lowNum, highNum)

		suite.NoError(err)
		suite.NotNil(clampedNumber)
	})

	suite.Run("returns nil when the shipment has no service items", func() {
		submittedShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), submittedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("returns nil when the shipment has no SIT service items", func() {
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
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

		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
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

		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(dofsit.ID.String(), sitStatus.PastSITs[0].ServiceItems[0].ID.String())
		suite.Equal(16, sitStatus.TotalSITDaysUsed)
		suite.Equal(16, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(74, sitStatus.TotalDaysRemaining)
		suite.Nil(sitStatus.CurrentSIT) // No current SIT since all SIT items have departed status
		// check that shipment values impacted by current SIT do not get updated since current SIT is nil
		suite.Nil(shipment.DestinationSITAuthEndDate)
		suite.Nil(shipment.OriginSITAuthEndDate)
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

		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(31, sitStatus.TotalSITDaysUsed)
		suite.Equal(31, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(59, sitStatus.TotalDaysRemaining)
		suite.Equal(31, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aMonthAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())
		suite.Len(sitStatus.PastSITs, 0)
		suite.NotNil(sitStatus.CurrentSIT.SITAuthorizedEndDate)
		// check that shipment values impacted by current SIT get updated
		suite.Equal(&sitStatus.CurrentSIT.SITAuthorizedEndDate, shipment.OriginSITAuthEndDate)
		suite.Nil(shipment.DestinationSITAuthEndDate)
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

		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(24, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(24, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(66, sitStatus.TotalDaysRemaining)
		suite.Equal(8, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOFSIT.ID.String(), sitStatus.PastSITs[0].ServiceItems[0].ID.String())

		// check that shipment values impacted by current SIT get updated
		suite.Equal(&sitStatus.CurrentSIT.SITAuthorizedEndDate, shipment.OriginSITAuthEndDate)
		suite.Nil(shipment.DestinationSITAuthEndDate)
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

		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(DestinationSITLocation, sitStatus.CurrentSIT.Location)
		suite.Equal(24, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(24, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(66, sitStatus.TotalDaysRemaining)
		suite.Equal(8, sitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOFSIT.ID.String(), sitStatus.PastSITs[0].ServiceItems[0].ID.String())
		// check that shipment values impacted by current SIT get updated
		suite.Equal(&sitStatus.CurrentSIT.SITAuthorizedEndDate, shipment.DestinationSITAuthEndDate)
		suite.Nil(shipment.OriginSITAuthEndDate)
	})

	suite.Run("When a SIT is created and there are no more remaining authorized days, the new SIT created in violation defaults to the past SITs authorized end date", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		ninetyOneDaysAgo := time.Now().AddDate(0, 0, -91).UTC().Truncate(24 * time.Hour)
		yesterday := time.Now().AddDate(0, 0, -1).UTC().Truncate(24 * time.Hour)

		// Create origin SIT that used all the allowance
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &ninetyOneDaysAgo,
					SITDepartureDate: &yesterday,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		today := time.Now().UTC().Truncate(24 * time.Hour)
		currentDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &today,
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

		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
		suite.Equal(ninetyOneDaysAgo.AddDate(0, 0, shipmentSITAllowance-1).UTC().Truncate(24*time.Hour), sitStatus.CurrentSIT.SITAuthorizedEndDate)
	})

	suite.Run("current sit calculates all used allowance properly", func() {
		// ensure that in the scenario that the SIT hasn't departed yet
		// and it has used all its allowance, that it still presents an accurate
		// authorized end date
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		eightyNineDaysAgo := time.Now().AddDate(0, 0, -89).UTC().Truncate(24 * time.Hour)

		// Create origin SIT that used all the allowance but has not yet departed
		currentDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &eightyNineDaysAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		today := time.Now().UTC().Truncate(24 * time.Hour)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{currentDOFSIT}

		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
		// Enforce that the authorized end date is today because it is the 90th day
		suite.Equal(today, sitStatus.CurrentSIT.SITAuthorizedEndDate)
	})

	suite.Run("if the authorized end date is calculated with a remaining allowance of zero it will default to the entry date", func() {
		// In the scenario that an authorized end date is calculated
		// with zero remaining authorized allowance, when we go to subtract
		// the last day to be inclusive, it will attempt to set
		// the authorized end date to before the entry date.
		// This ensures that this calculation will default to the entry date
		// rather than going backwards.
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		// Set past SIT date to fully use up the allowance
		ninetyDaysAgo := time.Now().AddDate(0, 0, -90).UTC().Truncate(24 * time.Hour)
		yesterday := time.Now().AddDate(0, 0, -1).UTC().Truncate(24 * time.Hour)

		// Create origin SIT that used all the allowance
		pastDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &ninetyDaysAgo,
					SITDepartureDate: &yesterday, // Count yesterday to equal 90 days used
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		// Create a new SIT with no remaining allowance
		today := time.Now().UTC().Truncate(24 * time.Hour)
		currentDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &today,
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

		// Calculate SIT status
		sitStatus, _, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(sitStatus.CurrentSIT.SITEntryDate, sitStatus.CurrentSIT.SITAuthorizedEndDate)
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

		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
		// check that shipment values impacted by current SIT do not get updated since current SIT is nil
		suite.Nil(shipment.DestinationSITAuthEndDate)
		suite.Nil(shipment.OriginSITAuthEndDate)
	})
}
