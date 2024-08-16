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
		suite.Equal(15, sitStatus.TotalSITDaysUsed)
		suite.Equal(15, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(75, sitStatus.TotalDaysRemaining)
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

	// TODO:
	// The point of this test is to satisfy the following requirement by the PO
	// - If Origin/Destination SIT is no longer "current", as in it has departed in a day
	//	 prior to today, then as long as a second SIT has not been created then it will be returned
	//   as "current". This is because "CurrentSIT" populates a SIT dashboard modal on the frontend
	//	 and we currently do not handle more than 1 SIT on the UI. This is a shoehorn to
	//	 remedy SIT date retention per E-05849 https://www13.v1host.com/USTRANSCOM38/Epic.mvc/Summary?oidToken=Epic%3A994028
	// suite.Run("calculates status for a shipment with origin SIT in the past but no newer SIT is present", func() {
	// 	shipmentSITAllowance := int(90)
	// 	year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
	// 	aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	// 	testCases := []struct {
	// 		reServiceCode models.ReServiceCode
	// 	}{
	// 		{
	// 			reServiceCode: models.ReServiceCodeDOFSIT,
	// 		},
	// 		{
	// 			reServiceCode: models.ReServiceCodeDDFSIT,
	// 		},
	// 	}
	// 	for _, tc := range testCases {
	// 		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
	// 			{
	// 				Model: models.MTOShipment{
	// 					Status:           models.MTOShipmentStatusApproved,
	// 					SITDaysAllowance: &shipmentSITAllowance,
	// 				},
	// 			},
	// 		}, nil)
	// 		// Set sit to a month ago
	// 		aMonthAndADayAgo := aMonthAgo.AddDate(0, 0, -1).UTC()
	// 		aDayAgo := time.Now().AddDate(0, 0, -1).UTC()
	// 		sitServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
	// 			{
	// 				Model:    approvedShipment,
	// 				LinkOnly: true,
	// 			},
	// 			{
	// 				Model: models.MTOServiceItem{
	// 					SITEntryDate:     &aMonthAndADayAgo,
	// 					SITDepartureDate: &aDayAgo,
	// 					Status:           models.MTOServiceItemStatusApproved,
	// 				},
	// 			},
	// 			{
	// 				Model: models.ReService{
	// 					Code: tc.reServiceCode,
	// 				},
	// 			},
	// 		}, nil)

	// 		approvedShipment.MTOServiceItems = models.MTOServiceItems{sitServiceItem}

	// 		sitStatus, shipment, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
	// 		suite.NoError(err)
	// 		suite.NotNil(sitStatus)
	// 		suite.NotNil(sitStatus.CurrentSIT)
	// 		suite.Equal(30, sitStatus.TotalSITDaysUsed)
	// 		suite.Equal(30, sitStatus.CalculatedTotalDaysInSIT)
	// 		suite.Equal(60, sitStatus.TotalDaysRemaining)
	// 		suite.Equal(30, sitStatus.CurrentSIT.DaysInSIT)
	// 		suite.Equal(aMonthAgo.String(), sitStatus.CurrentSIT.SITEntryDate.String())
	// 		suite.Nil(sitStatus.CurrentSIT.SITDepartureDate)
	// 		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())
	// 		suite.Len(sitStatus.PastSITs, 0)
	// 		suite.NotNil(sitStatus.CurrentSIT.SITAuthorizedEndDate)
	// 		// check that shipment values impacted by current SIT get updated
	// 		if tc.reServiceCode == models.ReServiceCodeDOFSIT {
	// 			suite.Equal(OriginSITLocation, sitStatus.CurrentSIT.Location)
	// 			suite.Equal(&sitStatus.CurrentSIT.SITAuthorizedEndDate, shipment.OriginSITAuthEndDate)
	// 			suite.Nil(shipment.DestinationSITAuthEndDate)
	// 		}
	// 		if tc.reServiceCode == models.ReServiceCodeDDFSIT {
	// 			suite.Equal(DestinationSITLocation, sitStatus.CurrentSIT.Location)
	// 			suite.Equal(&sitStatus.CurrentSIT.SITAuthorizedEndDate, shipment.DestinationSITAuthEndDate)
	// 			suite.Nil(shipment.OriginSITAuthEndDate)
	// 		}
	// 	}

	// })

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
		suite.Equal(23, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(23, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(67, sitStatus.TotalDaysRemaining)
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
		suite.Equal(23, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(23, sitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(67, sitStatus.TotalDaysRemaining)
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
