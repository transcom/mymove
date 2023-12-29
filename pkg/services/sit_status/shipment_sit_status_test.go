package sitstatus

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/unit"
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

	type localSubtestData struct {
		shipment             models.MTOShipment
		sitCustomerContacted time.Time
		sitRequestedDelivery time.Time
		eTag                 string
		planner              *mocks.Planner
	}

	makeSubtestData := func(addService bool, serviceCode models.ReServiceCode, estimatedWeight unit.Pound) (subtestData *localSubtestData) {
		subtestData = &localSubtestData{}

		shipmentSITAllowance := int(90)
		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		subtestData.shipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
		}, nil)

		subtestData.sitCustomerContacted = time.Now()
		year, month, day = time.Now().Add(time.Hour * 24 * 7).Date()
		subtestData.sitRequestedDelivery = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		subtestData.eTag = etag.GenerateEtag(subtestData.shipment.UpdatedAt)
		subtestData.planner = &mocks.Planner{}
		subtestData.planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		if addService {
			year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
			aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			customerContactDatePlusFive := subtestData.sitCustomerContacted.AddDate(0, 0, GracePeriodDays)

			factory := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model:    subtestData.shipment,
					LinkOnly: true,
				},
				{
					Model: models.MTOServiceItem{
						SITEntryDate:     &aMonthAgo,
						Status:           models.MTOServiceItemStatusApproved,
						SITDepartureDate: &customerContactDatePlusFive,
						UpdatedAt:        aMonthAgo,
					},
				},
				{
					Model: models.ReService{
						Code: serviceCode,
					},
				},
			}, nil)

			subtestData.shipment.MTOServiceItems = models.MTOServiceItems{factory}
		}

		return subtestData
	}

	suite.Run("calculates allowance end date for a shipment currently in Destination SIT", func() {
		subtestData := makeSubtestData(true, models.ReServiceCodeDDFSIT, unit.Pound(1400))
		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(&subtestData.sitCustomerContacted, sitStatus.CurrentSIT.SITCustomerContacted)
		suite.Equal(&subtestData.sitRequestedDelivery, sitStatus.CurrentSIT.SITRequestedDelivery)
		suite.NotEqual(&subtestData.shipment.MTOServiceItems[0].UpdatedAt, aMonthAgo)
	})

	suite.Run("calculates allowance end date and requested delivery date for a shipment currently in Origin SIT", func() {
		subtestData := makeSubtestData(true, models.ReServiceCodeDOFSIT, unit.Pound(1400))
		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(&subtestData.sitCustomerContacted, sitStatus.CurrentSIT.SITCustomerContacted)
		suite.Equal(&subtestData.sitRequestedDelivery, sitStatus.CurrentSIT.SITRequestedDelivery)
		suite.NotEqual(&subtestData.shipment.UpdatedAt, aMonthAgo)
		suite.NotEqual(&subtestData.shipment.MTOServiceItems[0].UpdatedAt, aMonthAgo)
	})

	suite.Run("calculate requested delivery date with sitDepartureDate before customer contact date plus grade period", func() {
		subtestData := makeSubtestData(false, models.ReServiceCodeDOFSIT, unit.Pound(1400))
		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		customerContactDatePlusThree := subtestData.sitCustomerContacted.AddDate(0, 0, GracePeriodDays-2)

		factory := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.shipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &aMonthAgo,
					Status:           models.MTOServiceItemStatusApproved,
					SITDepartureDate: &customerContactDatePlusThree,
					UpdatedAt:        aMonthAgo,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		subtestData.shipment.MTOServiceItems = models.MTOServiceItems{factory}
		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		suite.Equal(&subtestData.sitCustomerContacted, sitStatus.CurrentSIT.SITCustomerContacted)
		suite.Equal(&subtestData.sitRequestedDelivery, sitStatus.CurrentSIT.SITRequestedDelivery)
		suite.NotEqual(&subtestData.shipment.UpdatedAt, aMonthAgo)
		suite.NotEqual(&subtestData.shipment.MTOServiceItems[0].UpdatedAt, aMonthAgo)
	})

	suite.Run("failure test for calculate allowance with stale etag", func() {
		subtestData := makeSubtestData(false, models.ReServiceCodeDOFSIT, unit.Pound(1400))
		year, month, day := time.Now().Add(time.Hour * 24 * -15).Date()
		oldDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		subtestData.eTag = etag.GenerateEtag(oldDate)

		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)

		suite.Error(err)
		suite.Nil(sitStatus)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("failure test for calculate allowance with no service items", func() {
		subtestData := makeSubtestData(false, models.ReServiceCodeDOFSIT, unit.Pound(1400))
		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)

		suite.Error(err)
		suite.Nil(sitStatus)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("failure test for calculate allowance with no current SIT", func() {
		subtestData := makeSubtestData(false, models.ReServiceCodeCS, unit.Pound(1400))
		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)

		suite.Error(err)
		suite.Nil(sitStatus)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("failure test for ghc transit time query", func() {
		subtestData := makeSubtestData(true, models.ReServiceCodeDOFSIT, unit.Pound(20000))

		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)
		suite.Error(err)
		suite.Nil(sitStatus)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("failure test for ZipTransitDistance", func() {
		subtestData := makeSubtestData(true, models.ReServiceCodeDOFSIT, unit.Pound(1400))
		subtestData.planner = &mocks.Planner{}
		subtestData.planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, apperror.UnprocessableEntityError{})

		sitStatus, err := sitStatusService.CalculateSITAllowanceRequestedDates(suite.AppContextForTest(), subtestData.shipment, subtestData.planner,
			&subtestData.sitCustomerContacted, &subtestData.sitRequestedDelivery, subtestData.eTag)
		suite.Error(err)
		suite.Nil(sitStatus)
		suite.IsType(apperror.UnprocessableEntityError{}, err)
	})

}
