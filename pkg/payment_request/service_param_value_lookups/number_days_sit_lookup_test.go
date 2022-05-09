package serviceparamvaluelookups

import (
	"time"

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
	var serviceItemDDASITNine models.MTOServiceItem
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

		reServiceDOFSIT = testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
				Name: "Dom. Origin 1st Day SIT",
			},
		})

		reServiceDOASIT = testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
				Name: "Dom. Origin Add'l SIT",
			},
		})

		reServiceDDFSIT := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
				Name: "Dom. Destination 1st Day SIT",
			},
		})

		reServiceDDASIT = testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDDASIT,
				Name: "Dom. Destination Add'l SIT",
			},
		})

		moveTaskOrderOne = testdatagen.MakeDefaultMove(suite.DB())
		moveTaskOrderTwo = testdatagen.MakeDefaultMove(suite.DB())
		moveTaskOrderThree = testdatagen.MakeDefaultMove(suite.DB())
		moveTaskOrderFour = testdatagen.MakeDefaultMove(suite.DB())
		moveTaskOrderFive := testdatagen.MakeDefaultMove(suite.DB())

		mtoShipmentOne := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentTwo := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentThree := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentFour := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentFive := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentSix := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentSeven = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentEight := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentNine := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentTen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderOne,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentEleven := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderTwo,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentTwelve := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderThree,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentThirteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderThree,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentFourteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderFour,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentFifteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderFive,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

		mtoShipmentSixteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: moveTaskOrderFive,
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusSubmitted,
				SITDaysAllowance: &defaultSITDaysAllowance,
			},
		})

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

		serviceItemDOFSITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentOne,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentFour,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSix,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITFive = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentEight,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITSeven = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &originSITEntryDateOne,
				SITDepartureDate: &originSITDepartureDateTwo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentTen,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITEight = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &originSITEntryDateOne,
				SITDepartureDate: &originSITDepartureDateTwo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderTwo,
			MTOShipment: mtoShipmentEleven,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITNine := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderThree,
			MTOShipment: mtoShipmentThirteen,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOFSITTen := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentFifteen,
			ReService:   reServiceDOFSIT,
		})

		serviceItemDOASITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateTwo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentOne,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITTwo = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentOne,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITFour = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &originSITDepartureDateOne,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITFive = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateTwo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentFour,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITSix = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSix,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITEight := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &originSITDepartureDateOne,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentEight,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITNine = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentEight,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITTen := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderThree,
			MTOShipment: mtoShipmentThirteen,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITEleven = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &originSITDepartureDateThree,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderThree,
			MTOShipment: mtoShipmentThirteen,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITTwelve := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentFifteen,
			ReService:   reServiceDOASIT,
		})

		serviceItemDOASITThirteen := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &originSITDepartureDateOne,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentFifteen,
			ReService:   reServiceDOASIT,
		})

		serviceItemDDFSITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentTwo,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentFive,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSix,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITFive = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentNine,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITSeven = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &destinationSITEntryDateOne,
				SITDepartureDate: &destinationSITDepartureDateTwo,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentTen,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITEight := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFour,
			MTOShipment: mtoShipmentFourteen,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDFSITNine := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentSixteen,
			ReService:   reServiceDDFSIT,
		})

		serviceItemDDASITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateTwo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentTwo,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITTwo = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentTwo,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITFour = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentThree,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITFive = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateTwo,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentFive,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITSix = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSix,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITEight := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &destinationSITDepartureDateOne,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentNine,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITNine = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentNine,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITTen = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFour,
			MTOShipment: mtoShipmentFourteen,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITEleven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentSixteen,
			ReService:   reServiceDDASIT,
		})

		serviceItemDDASITTwelve := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITDepartureDate: &destinationSITDepartureDateThree,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderFive,
			MTOShipment: mtoShipmentSixteen,
			ReService:   reServiceDDASIT,
		})

		cost := unit.Cents(20000)

		paymentRequestOne := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestOne,
			MTOServiceItem: serviceItemDOFSITOne,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestOne,
			MTOServiceItem: serviceItemDOASITOne,
		})

		paymentRequestTwo := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestTwo,
			MTOServiceItem: serviceItemDDFSITOne,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestTwo,
			MTOServiceItem: serviceItemDDASITOne,
		})

		paymentRequestThree := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  3,
			},
			Move: moveTaskOrderOne,
		})

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDOFSITTwo, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDOASITThree, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDDFSITTwo, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestThree, serviceItemDDASITThree, "2021-11-11", "2021-11-20")

		paymentRequestFour := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  4,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestFour,
			MTOServiceItem: serviceItemDOFSITThree,
		})

		paymentRequestFive := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  5,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestFive,
			MTOServiceItem: serviceItemDDFSITThree,
		})

		paymentRequestSix := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  6,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestSix,
			MTOServiceItem: serviceItemDOFSITFour,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestSix,
			MTOServiceItem: serviceItemDDFSITFour,
		})

		paymentRequestSeven = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  7,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestSeven,
			MTOServiceItem: serviceItemDOFSITFive,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestSeven,
			MTOServiceItem: serviceItemDDFSITFive,
		})

		paymentRequestEight := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  8,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestEight,
			MTOServiceItem: serviceItemDOFSITSix,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestEight,
			MTOServiceItem: serviceItemDOASITEight,
		})

		paymentRequestNine := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  9,
			},
			Move: moveTaskOrderOne,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestNine,
			MTOServiceItem: serviceItemDDFSITSix,
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestNine,
			MTOServiceItem: serviceItemDDASITEight,
		})

		paymentRequestTen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
			Move: moveTaskOrderThree,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestTen,
			MTOServiceItem: serviceItemDOFSITNine,
		})

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestTen, serviceItemDOASITTen, "2021-11-11", "2021-11-20")

		paymentRequestEleven := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
			Move: moveTaskOrderFour,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestEleven,
			MTOServiceItem: serviceItemDDFSITEight,
		})

		paymentRequestTwelve := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
			Move: moveTaskOrderFive,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestTwelve,
			MTOServiceItem: serviceItemDOFSITTen,
		})

		paymentServiceItemParamTwo := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "29",
			},
		}
		testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDOASITTwelve.ReService.Code,
			paymentServiceItemParamTwo,
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
				PaymentRequest: paymentRequestTwelve,
				MTOServiceItem: serviceItemDOASITTwelve,
			})

		paymentServiceItemParamThree := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "32",
			},
		}
		testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDOASITThirteen.ReService.Code,
			paymentServiceItemParamThree,
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusPaid,
				},
				PaymentRequest: paymentRequestTwelve,
				MTOServiceItem: serviceItemDOASITThirteen,
			})

		paymentRequestThirteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusPaid,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: moveTaskOrderFive,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestThirteen,
			MTOServiceItem: serviceItemDDFSITNine,
		})

		paymentRequestFourteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  3,
			},
			Move: moveTaskOrderFive,
		})
		paymentServiceItemParamFour := []testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameNumberDaysSIT,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "27",
			},
		}
		testdatagen.MakePaymentServiceItemWithParams(
			suite.DB(),
			serviceItemDDASITEleven.ReService.Code,
			paymentServiceItemParamFour,
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: &cost,
					Status:     models.PaymentServiceItemStatusDenied,
				},
				PaymentRequest: paymentRequestFourteen,
				MTOServiceItem: serviceItemDDASITEleven,
			})

		paymentRequestFifteen = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  10,
			},
			Move: moveTaskOrderOne,
		})

		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITTwo, "2021-10-11", "2021-10-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITFour, "2021-10-21", "2021-10-30")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDDASITFour, "2021-11-01", "2021-11-10")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITSix, "2021-11-01", "2021-11-10")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDDASITSix, "2021-11-11", "2021-11-20")
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestFifteen, serviceItemDOASITNine, "2021-11-11", "2021-11-20")

		paymentRequestSixteen = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  1,
			},
			Move: moveTaskOrderTwo,
		})

		paymentRequestSeventeen = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: moveTaskOrderThree,
		})
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestSeventeen, serviceItemDOASITEleven, "2021-11-21", "2021-11-30")

		paymentRequestEighteen = testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: moveTaskOrderFour,
		})
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestEighteen, serviceItemDDASITTen, "2021-11-11", "2021-11-20")

		paymentRequestNineteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         true,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  4,
			},
			Move: moveTaskOrderFive,
		})
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
		serviceItemDOASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDOASIT,
		})
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

		serviceItemDDASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDDASIT,
		})
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestSeven, serviceItemDDASITSeven, "2021-12-01", "2021-12-10")

		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITSeven, paymentRequestSeven.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.NoError(err)
	})

	suite.Run("an MTO Shipment already has an Origin MTO Service Item with a SIT Departure Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASITNine, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "already has an Origin MTO Service Item")
		suite.Contains(err.Error(), "with a SIT Departure Date")
	})

	suite.Run("an MTO Shipment already has a Destination MTO Service Item with a SIT Departure Date", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDDASITNine, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "already has a Destination MTO Service Item")
		suite.Contains(err.Error(), "with a SIT Departure Date")
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
		paymentRequestOverlapping := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: move,
		})
		suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequestOverlapping, serviceItemDOASIT, "2020-07-25", "2020-08-10")
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, serviceItemDOASIT, paymentRequestOverlapping.ID, move.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
	suite.Run("it shouldn't matter if dates from rejected payment requests overlap with current payment request", func() {
		setupTestData()

		move, serviceItemDOASIT, paymentRequest := suite.setupMoveWithAddlDaysSITAndPaymentRequest(reServiceDOFSIT, originSITEntryDateOne, reServiceDOASIT, "2020-07-21", "2020-07-30")

		paymentRequestRejected := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusReviewedAllRejected,
				RejectionReason: nil,
				SequenceNumber:  1 + paymentRequest.SequenceNumber,
			},
			Move: move,
		})
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

		approvedDays := 400
		testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: serviceItemDOASIT.MTOShipment,
			SITExtension: models.SITExtension{
				ApprovedDays: &approvedDays,
			},
		})

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

		paymentRequestOverlapping := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
				SequenceNumber:  2,
			},
			Move: move,
		})
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
	paymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
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
	testdatagen.MakePaymentServiceItemWithParams(
		suite.DB(),
		serviceItem.ReService.Code,
		paymentServiceItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     status,
			},
			PaymentRequest: paymentRequest,
			MTOServiceItem: serviceItem,
		})
}

func (suite *ServiceParamValueLookupsSuite) makeAdditionalDaysSITPaymentServiceItem(paymentRequest models.PaymentRequest, serviceItem models.MTOServiceItem, startDate string, endDate string) {
	suite.makeAdditionalDaysSITPaymentServiceItemWithStatus(paymentRequest, serviceItem, startDate, endDate, models.PaymentServiceItemStatusPaid)
}

// setupMoveWithAddlDaysSITAndPaymentRequest creates a move with a single shipment, a Domestic Additional Days
// SIT service item, and a payment request for that service item.
func (suite *ServiceParamValueLookupsSuite) setupMoveWithAddlDaysSITAndPaymentRequest(sitFirstDayReService models.ReService, sitEntryDate time.Time, sitAdditionalDaysReService models.ReService, sitAdditionalDaysStartDate string, sitAdditionalDaysEndDate string) (models.Move, models.MTOServiceItem, models.PaymentRequest) {
	defaultSITDaysAllowance := 90
	move := testdatagen.MakeDefaultMove(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:           models.MTOShipmentStatusSubmitted,
			SITDaysAllowance: &defaultSITDaysAllowance,
		},
	})
	move.MTOShipments = models.MTOShipments{
		shipment,
	}
	serviceItemFirstDaySIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &sitEntryDate,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: shipment,
		ReService:   sitFirstDayReService,
	})
	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &sitEntryDate,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: shipment,
		ReService:   sitAdditionalDaysReService,
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPaid,
			RejectionReason: nil,
			SequenceNumber:  1,
		},
		Move: move,
	})
	cost := unit.Cents(20000)
	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemFirstDaySIT,
	})
	suite.makeAdditionalDaysSITPaymentServiceItem(paymentRequest, serviceItem, sitAdditionalDaysStartDate, sitAdditionalDaysEndDate)
	return move, serviceItem, paymentRequest
}
