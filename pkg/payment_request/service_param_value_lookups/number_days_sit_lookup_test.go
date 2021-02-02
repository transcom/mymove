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

	moveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())

	mtoShipmentOne := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentTwo := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentThree := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFour := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentFive := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentSix := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	mtoShipmentSeven := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: moveTaskOrder,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	moveTaskOrder.MTOShipments = models.MTOShipments{
		mtoShipmentOne,
		mtoShipmentTwo,
		mtoShipmentThree,
		mtoShipmentFour,
		mtoShipmentFive,
		mtoShipmentSix,
		mtoShipmentSeven,
	}

	originSITEntryDateOne := time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	originSITEntryDateTwo := time.Date(2020, time.August, 20, 0, 0, 0, 0, time.UTC)
	originSITDepartureDateOne := time.Date(2020, time.September, 20, 0, 0, 0, 0, time.UTC)
	// originSITDepartureDateTwo := time.Date(2020, time.October, 20, 0, 0, 0, 0, time.UTC)

	destinationSITEntryDateOne := time.Date(2020, time.October, 30, 0, 0, 0, 0, time.UTC)
	destinationSITEntryDateTwo := time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC)
	// destinationSITDepartureDateOne := time.Date(2020, time.December, 30, 0, 0, 0, 0, time.UTC)
	// // destinationSITDepartureDateTwo := time.Date(2021, time.January, 30, 0, 0, 0, 0, time.UTC)

	serviceItemDOFSITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentOne,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOFSITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOFSITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentFour,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOFSITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSix,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOFSITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSeven,
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOASITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentOne,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentOne,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITDepartureDate: &originSITDepartureDateOne,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentFour,
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSix,
		ReService:   reServiceDOASIT,
	})

	serviceItemDDFSITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentTwo,
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDFSITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDFSITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentFive,
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDFSITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSix,
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDFSITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSeven,
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDASITOne := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentTwo,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITTwo := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentTwo,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITThree := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITFour := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentThree,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITFive := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateTwo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentFive,
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASITSix := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDateOne,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: mtoShipmentSix,
		ReService:   reServiceDDASIT,
	})

	cost := unit.Cents(20000)

	paymentRequestOne := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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

	paymentRequestTwo := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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

	paymentRequestThree := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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

	paymentRequestFour := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestFour,
		MTOServiceItem: serviceItemDOFSITThree,
	})

	paymentRequestFive := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequestFive,
		MTOServiceItem: serviceItemDDFSITThree,
	})

	paymentRequestSix := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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

	paymentRequestSeven := testdatagen.MakeDefaultPaymentRequest(suite.DB())
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

	suite.T().Run("an MTO Shipment has multiple Origin MTO Service Items with different SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITTwo.ID, paymentRequestOne.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment has multiple Destination MTO Service Items with different SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITTwo.ID, paymentRequestTwo.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment has multiple Origin MTO Service Items with identical SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITFour.ID, paymentRequestThree.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has multiple Destination MTO Service Items with identical SIT Entry Dates", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITFour.ID, paymentRequestThree.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment already has an Origin MTO Service Item with a different SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITFive.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment already has a Destination MTO Service Item with a different SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITFive.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})

	suite.T().Run("an MTO Shipment already has an Origin MTO Service Item with an identical SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITSix.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment already has a Destination MTO Service Item with an identical SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITSix.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has Origin MTO Service Items but non with a SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOFSITFive.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)

		serviceItemDOASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &originSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrder,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDOASIT,
		})

		paramLookup, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDOASITSeven.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})

	suite.T().Run("an MTO Shipment has Destination MTO Service Items but non with a SIT Entry Date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDFSITFive.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)

		serviceItemDDASITSeven := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate: &destinationSITEntryDateOne,
				Status:       models.MTOServiceItemStatusApproved,
			},
			Move:        moveTaskOrder,
			MTOShipment: mtoShipmentSeven,
			ReService:   reServiceDDASIT,
		})

		paramLookup, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASITSeven.ID, paymentRequestFour.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.NoError(err)
	})
}
