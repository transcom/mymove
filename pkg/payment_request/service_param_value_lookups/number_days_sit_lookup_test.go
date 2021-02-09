package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestNumberDaysSITLookup() {
	key := models.ServiceItemParamNameNumberDaysSIT

	reServiceDOFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DOFSIT",
			Name: "Dom. Origin 1st Day SIT",
		},
	})

	reServiceDOASIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DOASIT",
			Name: "Dom. Origin Add'l SIT",
		},
	})

	reServiceDDFSIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DDFSIT",
			Name: "Dom. Destination 1st Day SIT",
		},
	})

	reServiceDDASIT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DDASIT",
			Name: "Dom. Destination Add'l SIT",
		},
	})

	moveTaskOrderOne := testdatagen.MakeDefaultMove(suite.DB())
	moveTaskOrderTwo := testdatagen.MakeDefaultMove(suite.DB())
	moveTaskOrderThree := testdatagen.MakeDefaultMove(suite.DB())
	moveTaskOrderFour := testdatagen.MakeDefaultMove(suite.DB())
	moveTaskOrderFive := testdatagen.MakeDefaultMove(suite.DB())

	mtoShipmentOne := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentTwo := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentThree := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFour := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFive := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentSix := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentSeven := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentEight := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentNine := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentTen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderOne,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentEleven := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderTwo,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentTwelve := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderThree,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentThirteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderThree,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFourteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderFour,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFifteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderFive,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentSixteen := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrderFive,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
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

	serviceItemDOFSITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDOFSITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate:     &originSITEntryDateOne,
			SITDepartureDate: &originSITDepartureDateTwo,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentTen,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOFSITEight := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDOASITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDOASITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITDepartureDate: &originSITDepartureDateOne,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentFour,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDOASITNine := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDOASITEleven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDDFSITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDDFSITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDDASITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDDASITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentFive,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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

	serviceItemDDASITNine := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrderOne,
		MTOShipment: mtoShipmentNine,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITTen := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
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
	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestThree,
		MTOServiceItem: serviceItemDOFSITTwo,
	})

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestThree,
		MTOServiceItem: serviceItemDOASITThree,
	})

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestThree,
		MTOServiceItem: serviceItemDDFSITTwo,
	})

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestThree,
		MTOServiceItem: serviceItemDDASITThree,
	})

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

	paymentRequestSeven := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
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

	paymentServiceItemParamOne := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameNumberDaysSIT,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "29",
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		suite.DB(),
		serviceItemDOASITTen.ReService.Code,
		paymentServiceItemParamOne,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &cost,
				Status:     models.PaymentServiceItemStatusPaid,
			},
			PaymentRequest: paymentRequestTen,
			MTOServiceItem: serviceItemDOASITTen,
		})

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

	paymentRequestFifteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         true,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			RejectionReason: nil,
			SequenceNumber:  10,
		},
		Move: moveTaskOrderOne,
	})

	paymentRequestSixteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         true,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			RejectionReason: nil,
			SequenceNumber:  1,
		},
		Move: moveTaskOrderTwo,
	})

	paymentRequestSeventeen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         true,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			RejectionReason: nil,
			SequenceNumber:  2,
		},
		Move: moveTaskOrderThree,
	})

	paymentRequestEighteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         true,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			RejectionReason: nil,
			SequenceNumber:  2,
		},
		Move: moveTaskOrderFour,
	})

	paymentRequestNineteen := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			IsFinal:         true,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			RejectionReason: nil,
			SequenceNumber:  4,
		},
		Move: moveTaskOrderFive,
	})

	suite.T().Run("an MTO Shipment has multiple Origin MTO Service Items with different SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITTwo.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment has multiple Destination MTO Service Items with different SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITTwo.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment has multiple Origin MTO Service Items with identical SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITFour.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has multiple Destination MTO Service Items with identical SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITFour.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment already has an Origin MTO Service Item with a different SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITFive.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment already has a Destination MTO Service Item with a different SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITFive.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment already has an Origin MTO Service Item with an identical SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITSix.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment already has a Destination MTO Service Item with an identical SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITSix.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has Origin MTO Service Items but non with a SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOFSITFive.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)

		serviceItemDOASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDOASIT,
		})

		paramLookup, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITSeven.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has Destination MTO Service Items but non with a SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDFSITFive.ID, paymentRequestSeven.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)

		serviceItemDDASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrderOne,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDDASIT,
		})

		paramLookup, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITSeven.ID, paymentRequestSeven.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment already has an Origin MTO Service Item with a SIT Departure Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITNine.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment already has a Destination MTO Service Item with a SIT Departure Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITNine.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment has an SIT Entry Date and SIT Departure Date on the same MTO Service Item", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOFSITSeven.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)

		paramLookup, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDFSITSeven.ID, paymentRequestFifteen.ID, moveTaskOrderOne.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO with one MTO Shipment with one DOFSIT payment service item", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOFSITEight.ID, paymentRequestSixteen.ID, moveTaskOrderTwo.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(key)
		suite.NoError(err)
		suite.Equal("0", value)
	})

	suite.T().Run("an MTO with more than one MTO Shipment", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITEleven.ID, paymentRequestSeventeen.ID, moveTaskOrderThree.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(key)
		suite.NoError(err)
		suite.Equal("10", value)
	})

	suite.T().Run("an MTO with an MTO Shipment with no SIT Departure Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITTen.ID, paymentRequestEighteen.ID, moveTaskOrderFour.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(key)
		suite.NoError(err)
		suite.Equal("29", value)
	})

	suite.T().Run("an MTO with an MTO Shipment that has more SIT days than the MTO has remaining", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITTwelve.ID, paymentRequestNineteen.ID, moveTaskOrderFive.ID, nil)
		suite.FatalNoError(err)

		value, err := paramLookup.ServiceParamValue(key)
		suite.NoError(err)
		suite.Equal("27", value)
	})
}
