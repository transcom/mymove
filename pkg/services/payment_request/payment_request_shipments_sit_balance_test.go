package paymentrequest

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestListShipmentPaymentSITBalance() {
	service := NewPaymentRequestShipmentsSITBalance()

	suite.Run("returns only pending SIT status when there are no previous payments", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Locator:            "PARSIT",
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		oneHundredAndTwentyDays := 120
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
			Move: move,
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: paymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     paymentRequest,
			MTOServiceItem:     doasit,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     paymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), paymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)
		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
		suite.Nil(pendingSITBalance.PreviouslyBilledDays)
		suite.Equal(30, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(paymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())
		suite.Equal(120, pendingSITBalance.TotalSITDaysAuthorized)
		suite.Equal(90, pendingSITBalance.TotalSITDaysRemaining)
		suite.Equal(paymentEndDate.AddDate(0, 0, 91).String(), pendingSITBalance.TotalSITEndDate.String())
	})

	suite.Run("calculates pending destination SIT balance when origin was invoiced previously", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		oneHundredAndTwentyDays := 120
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
			Move: move,
		})

		reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: reviewedPaymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status:         models.PaymentRequestStatusReviewed,
				SequenceNumber: 2,
			},
		})

		destinationEntryDate := time.Date(year, month, day-89, 0, 0, 0, 0, time.UTC)
		ddasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOServiceItem: ddasit,
			Move:           move,
		})

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationPaymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "60",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), pendingPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)
		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
		suite.Equal(30, *pendingSITBalance.PreviouslyBilledDays)
		suite.Equal(paymentEndDate.String(), pendingSITBalance.PreviouslyBilledEndDate.String())
		suite.Equal(60, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(destinationPaymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())
		suite.Equal(120, pendingSITBalance.TotalSITDaysAuthorized)
		suite.Equal(30, pendingSITBalance.TotalSITDaysRemaining)
		suite.Equal(destinationPaymentEndDate.AddDate(0, 0, 31).String(), pendingSITBalance.TotalSITEndDate.String())
	})

	suite.Run("ignores including previously denied service items in SIT balance", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		oneHundredAndTwentyDays := 120
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
			Move: move,
		})

		reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusDenied,
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: reviewedPaymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status:         models.PaymentRequestStatusReviewed,
				SequenceNumber: 2,
			},
		})

		destinationEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		ddasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOServiceItem: ddasit,
			Move:           move,
		})

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationPaymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "60",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), pendingPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)
		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
		suite.Equal(120, pendingSITBalance.TotalSITDaysAuthorized)
		suite.Equal(60, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(60, pendingSITBalance.TotalSITDaysRemaining)
		suite.Equal(destinationPaymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())
		suite.Nil(pendingSITBalance.PreviouslyBilledDays)
	})

	suite.Run("returns nil for reviewed payment request without SIT service items", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		oneHundredAndTwentyDays := 120
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
			Move: move,
		})

		reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
			PaymentRequest: reviewedPaymentRequest,
			MTOShipment:    shipment,
			Move:           move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)
		suite.Nil(sitBalances)
	})

	suite.Run("returns nil for pending payment request without SIT service items", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		oneHundredAndTwentyDays := 120
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &oneHundredAndTwentyDays,
			},
			Move: move,
		})

		pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOShipment:    shipment,
			Move:           move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), pendingPaymentRequest.ID)
		suite.NoError(err)
		suite.Nil(sitBalances)
	})

	suite.Run("returns zero authorized days for pending payment request shipment without a set SITDaysAllowance", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:       models.MTOShipmentStatusApproved,
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
			},
			Move: move,
		})

		pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		})

		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOShipment:    shipment,
			Move:           move,
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     doasit,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), pendingPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)
		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
		suite.Equal(0, pendingSITBalance.TotalSITDaysAuthorized)
		suite.Equal(30, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(0, pendingSITBalance.TotalSITDaysRemaining)
		suite.Equal(paymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())
		suite.Nil(pendingSITBalance.PreviouslyBilledDays)
	})

	suite.Run("returns zero authorized days for reviewed payment request shipment without a set SITDaysAllowance", func() {
		availableToPrimeAt := time.Now()
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &availableToPrimeAt,
			},
		})

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:       models.MTOShipmentStatusApproved,
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
			},
			Move: move,
		})

		reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
			PaymentRequest: reviewedPaymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)
		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
		suite.Equal(0, pendingSITBalance.TotalSITDaysAuthorized)
		suite.Equal(30, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(0, pendingSITBalance.TotalSITDaysRemaining)
		suite.Equal(paymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())
		suite.Equal(30, *pendingSITBalance.PreviouslyBilledDays)
	})
}
