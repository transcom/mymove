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

	moveTaskOrder := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	originSITEntryDate := time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	// originSITEntryDate2 := time.Date(2020, time.August, 20, 0, 0, 0, 0, time.UTC)
	originSITDepartureDate := time.Date(2020, time.September, 20, 0, 0, 0, 0, time.UTC)
	// originSITDepartureDate2 := time.Date(2020, time.October, 20, 0, 0, 0, 0, time.UTC)

	destinationSITEntryDate := time.Date(2020, time.October, 30, 0, 0, 0, 0, time.UTC)
	// destinationSITEntryDate2 := time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC)
	destinationSITDepartureDate := time.Date(2021, time.December, 30, 0, 0, 0, 0, time.UTC)
	// destinationSITDepartureDate2 := time.Date(2021, time.January, 30, 0, 0, 0, 0, time.UTC)

	serviceItemDOFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &originSITEntryDate,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDOFSIT,
	})

	serviceItemDOASIT1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASIT2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDOASIT,
	})

	serviceItemDOASIT3 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITDepartureDate: &originSITDepartureDate,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDOASIT,
	})

	serviceItemDDFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDate,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDDFSIT,
	})

	serviceItemDDASIT1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &destinationSITEntryDate,
			Status:       models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDDASIT,
	})

	serviceItemDDASIT2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITDepartureDate: &destinationSITDepartureDate,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        moveTaskOrder,
		MTOShipment: moveTaskOrder.MTOShipments[0],
		ReService:   reServiceDDASIT,
	})

	paymentRequest1 := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	cost := unit.Cents(20000)

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequest1,
		MTOServiceItem: serviceItemDOFSIT,
	})

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusPaid,
		},
		PaymentRequest: paymentRequest1,
		MTOServiceItem: serviceItemDOASIT1,
	})

	paymentRequest2 := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusDenied,
		},
		PaymentRequest: paymentRequest2,
		MTOServiceItem: serviceItemDOASIT2,
	})

	paymentRequest3 := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusRequested,
		},
		PaymentRequest: paymentRequest3,
		MTOServiceItem: serviceItemDOASIT3,
	})

	paymentRequest4 := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusRequested,
		},
		PaymentRequest: paymentRequest4,
		MTOServiceItem: serviceItemDDFSIT,
	})

	paymentRequest5 := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
			Status:     models.PaymentServiceItemStatusRequested,
		},
		PaymentRequest: paymentRequest5,
		MTOServiceItem: serviceItemDDASIT1,
	})

	suite.T().Run("lookup Number of Days SIT for an MTO Shipment with MTO Service Items with no departure date", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, serviceItemDDASIT2.ID, paymentRequest2.ID, moveTaskOrder.ID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("2", valueStr)
	})
}
