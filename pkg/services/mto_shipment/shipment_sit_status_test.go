package mtoshipment

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestShipmentSITStatus() {
	sitStatusService := NewShipmentSITStatus()

	suite.Run("returns nil when the shipment has no service items", func() {
		submittedShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{})

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), submittedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("returns nil when the shipment has no SIT service items", func() {
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusApproved,
				MTOServiceItems: testdatagen.MakeMTOServiceItems(suite.DB()),
			},
		})

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("returns nil when the shipment has a SIT service item with entry date in the future", func() {
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})

		nextWeek := time.Now().Add(time.Hour * 24 * 7)
		futureSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &nextWeek,
				Status:       models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{futureSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})

	suite.Run("includes SIT service item that has departed storage", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		dopsit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &aMonthAgo,
				SITDepartureDate: &fifteenDaysAgo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dopsit}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)
		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(dopsit.ID.String(), sitStatus.PastSITs[0].ID.String())

		suite.Equal(15, sitStatus.TotalSITDaysUsed)
		suite.Equal(75, sitStatus.TotalDaysRemaining)
		suite.Equal("", sitStatus.Location) // No current SIT so it will receive the zero value of empty
	})

	suite.Run("calculates status for a shipment currently in SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dopsit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &aMonthAgo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dopsit}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.Location)
		suite.Equal(30, sitStatus.TotalSITDaysUsed)
		suite.Equal(60, sitStatus.TotalDaysRemaining)
		suite.Equal(30, sitStatus.DaysInSIT)
		suite.Equal(aMonthAgo.String(), sitStatus.SITEntryDate.String())
		suite.Nil(sitStatus.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())
		suite.Len(sitStatus.PastSITs, 0)
	})

	suite.Run("combines SIT days sum for shipment with past and current SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		pastDOPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &aMonthAgo,
				SITDepartureDate: &fifteenDaysAgo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDOPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &aWeekAgo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOPSIT, currentDOPSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(OriginSITLocation, sitStatus.Location)
		suite.Equal(22, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(68, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.SITEntryDate.String())
		suite.Nil(sitStatus.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOPSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("combines SIT days sum for shipment with past origin and current destination SIT", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		fifteenDaysAgo := aMonthAgo.Add(time.Hour * 24 * 15)
		pastDOPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &aMonthAgo,
				SITDepartureDate: &fifteenDaysAgo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &aWeekAgo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOPSIT, currentDDPSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(DestinationSITLocation, sitStatus.Location)
		suite.Equal(22, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(68, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.SITEntryDate.String())
		suite.Nil(sitStatus.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOPSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("returns negative days remaining when days in SIT exceeds shipment allowance", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * 30 * -6).Date()
		sixMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		threeMonthsAgo := sixMonthsAgo.Add(time.Hour * 24 * 30 * 3)
		pastDOPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &sixMonthsAgo,
				SITDepartureDate: &threeMonthsAgo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &aWeekAgo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOPSIT, currentDDPSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)

		suite.NotNil(sitStatus)

		suite.Equal(DestinationSITLocation, sitStatus.Location)
		suite.Equal(97, sitStatus.TotalSITDaysUsed) // 15 days from previous SIT, 7 days from the current
		suite.Equal(-7, sitStatus.TotalDaysRemaining)
		suite.Equal(7, sitStatus.DaysInSIT)
		suite.Equal(aWeekAgo.String(), sitStatus.SITEntryDate.String())
		suite.Nil(sitStatus.SITDepartureDate)
		suite.Equal(approvedShipment.ID.String(), sitStatus.ShipmentID.String())

		suite.Len(sitStatus.PastSITs, 1)
		suite.Equal(pastDOPSIT.ID.String(), sitStatus.PastSITs[0].ID.String())
	})

	suite.Run("excludes SIT service items that have not been approved by the TOO", func() {
		shipmentSITAllowance := int(90)
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &shipmentSITAllowance,
			},
		})

		year, month, day := time.Now().Add(time.Hour * 24 * 30 * -6).Date()
		sixMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		threeMonthsAgo := sixMonthsAgo.Add(time.Hour * 24 * 30 * 3)
		pastDOPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &sixMonthsAgo,
				SITDepartureDate: &threeMonthsAgo,
				Status:           models.MTOServiceItemStatusRejected,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})

		year, month, day = time.Now().Add(time.Hour * 24 * -7).Date()
		aWeekAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		currentDDPSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: approvedShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &aWeekAgo,
				Status:       models.MTOServiceItemStatusRejected,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		})

		approvedShipment.MTOServiceItems = models.MTOServiceItems{pastDOPSIT, currentDDPSIT}

		sitStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), approvedShipment)
		suite.NoError(err)
		suite.Nil(sitStatus)
	})
}
