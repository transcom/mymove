package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestNumberDaysSITLookup() {
	key := models.ServiceItemParamNameNumberDaysSIT
	defaultSITDaysAllowance := 90

	var serviceItemDOASITTwo models.MTOServiceItem
	var serviceItemDOASITSix models.MTOServiceItem
	var serviceItemDOASITFour models.MTOServiceItem
	var serviceItemDOASITFive models.MTOServiceItem
	var serviceItemDOASITNine models.MTOServiceItem
	var serviceItemDOASITEleven models.MTOServiceItem

	var serviceItemDOFSITSeven models.MTOServiceItem
	var serviceItemDOFSITEight models.MTOServiceItem
	var serviceItemDOFSITFive models.MTOServiceItem

	var serviceItemDDASITTwo models.MTOServiceItem
	var serviceItemDDASITFour models.MTOServiceItem
	var serviceItemDDASITFive models.MTOServiceItem
	var serviceItemDDASITSix models.MTOServiceItem
	var serviceItemDDASITTen models.MTOServiceItem

	var serviceItemDDFSITFive models.MTOServiceItem
	var serviceItemDDFSITSeven models.MTOServiceItem

	var paymentRequestSeven models.PaymentRequest
	var paymentRequestFifteen models.PaymentRequest
	var paymentRequestSixteen models.PaymentRequest
	var paymentRequestSeventeen models.PaymentRequest
	var paymentRequestEighteen models.PaymentRequest

	var reServiceDOASIT models.ReService
	var reServiceDDASIT models.ReService
	var reServiceDOFSIT models.ReService

	var mtoShipmentSeven models.MTOShipment

	var moveTaskOrderOne models.Move
	var moveTaskOrderTwo models.Move
	var moveTaskOrderThree models.Move
	var moveTaskOrderFour models.Move

	originSITEntryDateOne := time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	originSITEntryDateTwo := time.Date(2020, time.August, 20, 0, 0, 0, 0, time.UTC)
	originSITDepartureDateOne := time.Date(2020, time.September, 20, 0, 0, 0, 0, time.UTC)
	originSITDepartureDateTwo := time.Date(2020, time.July, 21, 0, 0, 0, 0, time.UTC)
	originSITDepartureDateThree := time.Date(2020, time.August, 29, 0, 0, 0, 0, time.UTC)

	destinationSITEntryDateOne := time.Date(2020, time.October, 30, 0, 0, 0, 0, time.UTC)
	destinationSITEntryDateTwo := time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC)
	destinationSITDepartureDateOne := time.Date(2020, time.December, 30, 0, 0, 0, 0, time.UTC)
	destinationSITDepartureDateTwo := time.Date(2020, time.October, 31, 0, 0, 0, 0, time.UTC)
	destinationSITDepartureDateThree := time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC)

	setupTestData := func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		reServiceDOFSIT = factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
					Name: "Dom. Origin 1st Day SIT",
				},
			},
		}, nil)

		reServiceDOASIT = factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
					Name: "Dom. Origin Add'l SIT",
				},
			},
		}, nil)

		reServiceDDFSIT := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
					Name: "Dom. Destination 1st Day SIT",
				},
			},
		}, nil)

		reServiceDDASIT = factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
					Name: "Dom. Destination Add'l SIT",
				},
			},
		}, nil)

		moveTaskOrderOne = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		moveTaskOrderTwo = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		moveTaskOrderThree = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		moveTaskOrderFour = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		moveTaskOrderFive := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentOne := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentTwo := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentThree := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentFour := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentFive := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentSix := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentSeven = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentEight := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentNine := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentTen := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentEleven := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderTwo,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentTwelve := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentThirteen := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentFourteen := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderFour,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentFifteen := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		mtoShipmentSixteen := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusSubmitted,
					SITDaysAllowance: &defaultSITDaysAllowance,
				},
			},
		}, nil)

		moveTaskOrderOne.MTOShipments = models.MTOShipments{
			mtoShipmentOne,
			mtoShipmentTwo,
			mtoShipmentThree,
			mtoShipmentFour,
			mtoShipmentFive,
			mtoShipmentSix,
			mtoShipmentSeven,
			mtoShipmentEight,
			mtoShipmentNine,
			mtoShipmentTen,
		}

		moveTaskOrderTwo.MTOShipments = models.MTOShipments{
			mtoShipmentEleven,
		}

		moveTaskOrderThree.MTOShipments = models.MTOShipments{
			mtoShipmentTwelve,
			mtoShipmentThirteen,
		}

		moveTaskOrderFour.MTOShipments = models.MTOShipments{
			mtoShipmentFourteen,
		}

		moveTaskOrderFive.MTOShipments = models.MTOShipments{
			mtoShipmentFifteen,
			mtoShipmentSixteen,
		}

		serviceItemDOFSITOne := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentOne,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITTwo := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITThree := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFour,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITFour := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSix,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITFive = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSeven,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITSix := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentEight,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITSeven = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &originSITEntryDateOne,
					SITDepartureDate: &originSITDepartureDateTwo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentTen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITEight = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &originSITEntryDateOne,
					SITDepartureDate: &originSITDepartureDateTwo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderTwo,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentEleven,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITNine := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThirteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOFSITTen := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFifteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITOne := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateTwo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentOne,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITTwo = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentOne,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITThree := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITFour = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &originSITDepartureDateOne,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITFive = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateTwo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFour,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITSix = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSix,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITEight := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &originSITDepartureDateOne,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentEight,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITNine = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentEight,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITTen := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThirteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITEleven = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &originSITDepartureDateThree,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThirteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITTwelve := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFifteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDOASITThirteen := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &originSITDepartureDateOne,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFifteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITOne := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentTwo,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITTwo := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITThree := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFive,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITFour := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSix,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITFive = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSeven,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITSix := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentNine,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITSeven = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &destinationSITEntryDateOne,
					SITDepartureDate: &destinationSITDepartureDateTwo,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentTen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITEight := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFour,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFourteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDFSITNine := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSixteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDFSIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITOne := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateTwo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentTwo,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITTwo = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentTwo,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITThree := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITFour = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentThree,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITFive = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateTwo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFive,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITSix = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSix,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITEight := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &destinationSITDepartureDateOne,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentNine,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITTen = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFour,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentFourteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITEleven := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSixteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		serviceItemDDASITTwelve := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &destinationSITDepartureDateThree,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSixteen,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)

		cost := unit.Cents(20000)

		paymentRequestOne := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  1,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestOne,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITOne,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestOne,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOASITOne,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestTwo := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestTwo,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITOne,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestTwo,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDASITOne,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestThree := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  3,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDOFSITTwo, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDOASITThree, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDDFSITTwo, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDDASITThree, "2021-11-11", "2021-11-20")

		paymentRequestFour := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  4,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestFour,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITThree,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestFive := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  5,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestFive,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITThree,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestSix := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  6,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestSix,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITFour,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestSix,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITFour,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestSeven = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  7,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestSeven,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITFive,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestSeven,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITFive,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestEight := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  8,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestEight,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITSix,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestEight,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOASITEight,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestNine := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  9,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestNine,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITSix,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestNine,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDASITEight,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestTen := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  1,
				},
			},
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestTen,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITNine,
				LinkOnly: true,
			},
		}, nil)

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestTen, serviceItemDOASITTen, "2021-11-11", "2021-11-20")

		paymentRequestEleven := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  1,
				},
			},
			{
				Model:    moveTaskOrderFour,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestEleven,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITEight,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestTwelve := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  1,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestTwelve,
				LinkOnly: true,
			}, {
				Model:    serviceItemDOFSITTen,
				LinkOnly: true,
			},
		}, nil)

		paymentServiceItemParamTwo := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "29",
			},
		}
		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDOASITTwelve.ReService.Code,
			paymentServiceItemParamTwo,
			[]factory.Customization{
				{
					Model: models.PaymentServiceItem{
						PriceCents: &cost,
						Status:     models.PaymentServiceItemStatusPaid,
					},
				},
				{
					Model:    paymentRequestTwelve,
					LinkOnly: true,
				},
				{
					Model:    serviceItemDOASITTwelve,
					LinkOnly: true,
				},
			}, nil)

		paymentServiceItemParamThree := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "32",
			},
		}
		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDOASITThirteen.ReService.Code,
			paymentServiceItemParamThree,
			[]factory.Customization{
				{
					Model: models.PaymentServiceItem{
						PriceCents: &cost,
						Status:     models.PaymentServiceItemStatusPaid,
					},
				},
				{
					Model:    paymentRequestTwelve,
					LinkOnly: true,
				},
				{
					Model:    serviceItemDOASITThirteen,
					LinkOnly: true,
				},
			}, nil)

		paymentRequestThirteen := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusPaid,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
			}, {
				Model:    paymentRequestThirteen,
				LinkOnly: true,
			}, {
				Model:    serviceItemDDFSITNine,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestFourteen := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  3,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
		}, nil)
		paymentServiceItemParamFour := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "27",
			},
		}
		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDDASITEleven.ReService.Code,
			paymentServiceItemParamFour,
			[]factory.Customization{
				{
					Model: models.PaymentServiceItem{
						PriceCents: &cost,
						Status:     models.PaymentServiceItemStatusDenied,
					},
				},
				{
					Model:    paymentRequestFourteen,
					LinkOnly: true,
				},
				{
					Model:    serviceItemDDASITEleven,
					LinkOnly: true,
				},
			}, nil)

		paymentRequestFifteen = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  10,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
		}, nil)

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITTwo, "2021-10-11", "2021-10-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITFour, "2021-10-21", "2021-10-30")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDDASITFour, "2021-11-01", "2021-11-10")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITSix, "2021-11-01", "2021-11-10")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDDASITSix, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITNine, "2021-11-11", "2021-11-20")

		paymentRequestSixteen = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  1,
				},
			},
			{
				Model:    moveTaskOrderTwo,
				LinkOnly: true,
			},
		}, nil)

		paymentRequestSeventeen = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    moveTaskOrderThree,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestSeventeen, serviceItemDOASITEleven, "2021-11-21", "2021-11-30")

		paymentRequestEighteen = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    moveTaskOrderFour,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestEighteen, serviceItemDDASITTen, "2021-11-11", "2021-11-20")

		paymentRequestNineteen := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         true,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  4,
				},
			},
			{
				Model:    moveTaskOrderFive,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestNineteen, serviceItemDDASITTwelve, "2021-11-11", "2021-11-20")
	}

	suite.Run("an MTO Shipment has multiple Origin MTO Service Items with different SIT Entry Dates", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITTwo, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "multiple Origin MTO Service Items with different SIT Entry Dates")
	})

	suite.Run("an MTO Shipment has multiple Destination MTO Service Items with different SIT Entry Dates", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITTwo, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "multiple Destination MTO Service Items with different SIT Entry Dates")
	})

	// TODO can we support this case? the test data has 2 DOASIT service items, does that even make sense?
	suite.Run("an MTO Shipment has multiple Origin MTO Service Items with identical SIT Entry Dates", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITFour, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	// TODO can we support this case? the test data has 2 DDASIT service items on the shipment, does that even make sense?
	suite.Run("an MTO Shipment has multiple Destination MTO Service Items with identical SIT Entry Dates", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITFour, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment already has an Origin MTO Service Item with a different SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITFive, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "already has an Origin MTO Service Item with a different SIT Entry Date")
	})

	suite.Run("an MTO Shipment already has a Destination MTO Service Item with a different SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITFive, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "already has a Destination MTO Service Item with a different SIT Entry Date")
	})

	suite.Run("an MTO Shipment already has an Origin MTO Service Item with an identical SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITSix, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment already has a Destination MTO Service Item with an identical SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITSix, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment has Origin MTO Service Items but none with a SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOFSITFive, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		// Test that it fails
		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "does not have an Origin MTO Service Item with a SIT Entry Date")

		// Now test that it succeeds after we add a service item with entry date
		serviceItemDOASITSeven := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &originSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSeven,
				LinkOnly: true,
			},
			{
				Model:    reServiceDOASIT,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITSeven, "2021-11-21", "2021-11-30")

		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITSeven, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment has Destination MTO Service Items but none with a SIT Entry Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDFSITFive, paymentRequestSeven.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "does not have a Destination MTO Service Item with a SIT Entry Date")

		serviceItemDDASITSeven := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &destinationSITEntryDateOne,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    moveTaskOrderOne,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentSeven,
				LinkOnly: true,
			},
			{
				Model:    reServiceDDASIT,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestSeven, serviceItemDDASITSeven, "2021-12-01", "2021-12-10")

		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITSeven, paymentRequestSeven.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment only has a First Day SIT MTO Service Item", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOFSITSeven, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "failed to find a PaymentServiceItem for MTOServiceItem")

		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDFSITSeven, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "failed to find a PaymentServiceItem for MTOServiceItem")
	})

	suite.Run("an MTO with one MTO Shipment with one DOFSIT payment service item", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOFSITEight, paymentRequestSixteen.ID, moveTaskOrderTwo.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "failed to find a PaymentServiceItem for MTOServiceItem")
		suite.Equal("", value)
	})

	suite.Run("an MTO with more than one MTO Shipment", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITEleven, paymentRequestSeventeen.ID, moveTaskOrderThree.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Equal("10", value)
	})

	suite.Run("an MTO with an MTO Shipment with no SIT Departure Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITTen, paymentRequestEighteen.ID, moveTaskOrderFour.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
		suite.Equal("10", value)
	})

	suite.Run("simple date calculation", func() {
		setupTestData()

		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2020-07-30")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		days, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)

		suite.Equal("10", days)
	})
	suite.Run("invalid start date", func() {
		setupTestData()

		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "not a date", "2020-07-30")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "failed to parse SITPaymentRequestStart as a date")
	})
	suite.Run("invalid end date", func() {
		setupTestData()

		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-01", "not a date")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "failed to parse SITPaymentRequestEnd as a date")
	})
	suite.Run("overlapping dates should error", func() {
		setupTestData()

		move, serviceItemDOASIT, _ := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2020-07-30")
		paymentRequestOverlapping := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestOverlapping, serviceItemDOASIT, "2020-07-25", "2020-08-10")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequestOverlapping.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
	suite.Run("it shouldn't matter if dates from rejected payment requests overlap with current payment request", func() {
		setupTestData()

		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2020-07-30")

		paymentRequestRejected := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusReviewedAllRejected,
					RejectionReason: nil,
					SequenceNumber:  1 + paymentRequest.SequenceNumber,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		suite.makeAdditionalDaysSITPaymentServiceItemWithStatus(paymentRequestRejected, serviceItemDOASIT, "2020-07-21", "2020-07-30", models.PaymentServiceItemStatusDenied)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
	})

	suite.Run("Requests for SIT additional days past the allowance for the shipment should be rejected", func() {
		setupTestData()

		// End date is a year in the future in order to exceed allowance
		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2021-07-30")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})

	suite.Run("Requests for SIT additional days past the original allowance should be accepted if they are covered by extensions", func() {
		setupTestData()

		// End date is a year in the future in order to make sure we exceed the allowance.
		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2021-07-30")

		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    serviceItemDOASIT.MTOShipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					ApprovedDays: models.IntPointer(400),
				},
			},
		}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
	})

	suite.Run("SIT days remaining calculation should account for first day in SIT", func() {
		setupTestData()

		// The Additional Days SIT service item should be for exactly the allowed amount,
		// So the first day in SIT will put it over the limit.
		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2020-10-18")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequest.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})

	suite.Run("SIT Additional Days cannot start on the same day as the end of a previously billed date range", func() {
		setupTestData()

		move, serviceItemDOASIT, _ := suite.setupMoveWithAddlDaysSITAndPaymentRequest(
			reServiceDOFSIT,
			originSITEntryDateOne,
			reServiceDOASIT,
			"2020-07-21", "2020-07-30")

		paymentRequestOverlapping := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
					SequenceNumber:  2,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// Previously billed DOASIT ends on 2020-07-30. This one starts on that same date, so the lookup should fail.
		suite.makeAdditionalDaysSITPaymentServiceItem(
			paymentRequestOverlapping,
			serviceItemDOASIT,
			"2020-07-30", "2020-08-15")

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequestOverlapping.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
}

func (suite *ServiceParamValueLookupsSuite) makeAdditionalDaysSITPaymentServiceItemWithStatus(paymentRequest models.PaymentRequest, serviceItem models.MTOServiceItem, startDate string, endDate string, status models.PaymentServiceItemStatus) {
	cost := unit.Cents(20000)
	paymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameSITPaymentRequestStart,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   startDate,
		},
		{
			Key:     models.ServiceItemParamNameSITPaymentRequestEnd,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   endDate,
		},
	}
	factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		serviceItem.ReService.Code,
		paymentServiceItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     status,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)
}

func (suite *ServiceParamValueLookupsSuite) makeAdditionalDaysSITPaymentServiceItem(paymentRequest models.PaymentRequest, serviceItem models.MTOServiceItem, startDate string, endDate string) {
	suite.makeAdditionalDaysSITPaymentServiceItemWithStatus(paymentRequest, serviceItem, startDate, endDate, models.PaymentServiceItemStatusPaid)
}

// setupMoveWithAddlDaysSITAndPaymentRequest creates a move with a single shipment, a Domestic Additional Days
// SIT service item, and a payment request for that service item.
func (suite *ServiceParamValueLookupsSuite) setupMoveWithAddlDaysSITAndPaymentRequest(sitFirstDayReService models.ReService, sitEntryDate time.Time, sitAdditionalDaysReService models.ReService, sitAdditionalDaysStartDate string, sitAdditionalDaysEndDate string) (models.Move, models.MTOServiceItem, models.PaymentRequest) {
	defaultSITDaysAllowance := 90
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		},
	}, nil)
	move.MTOShipments = models.MTOShipments{
		shipment,
	}
	serviceItemFirstDaySIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				SITEntryDate: &sitEntryDate,
				Status:       models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    sitFirstDayReService,
			LinkOnly: true,
		},
	}, nil)
	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				SITEntryDate: &sitEntryDate,
				Status:       models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model:    sitAdditionalDaysReService,
			LinkOnly: true,
		},
	}, nil)

	paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
	cost := unit.Cents(20000)
	factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFirstDaySIT,
			LinkOnly: true,
		},
	}, nil)
	suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequest, serviceItem, sitAdditionalDaysStartDate, sitAdditionalDaysEndDate)
	return move, serviceItem, paymentRequest
}
