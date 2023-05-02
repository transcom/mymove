package paymentrequest

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestListShipmentPaymentSITBalance() {
	service := NewPaymentRequestShipmentsSITBalance()

	suite.Run("returns only pending SIT status when there are no previous payments", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "PARSIT",
				},
			},
		}, nil)

		oneHundredAndTwentyDays := 120
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &oneHundredAndTwentyDays,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		oneHundredAndTwentyDays := 120
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &oneHundredAndTwentyDays,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reviewedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status:         models.PaymentRequestStatusReviewed,
					SequenceNumber: 2,
				},
			},
		}, nil)

		destinationEntryDate := time.Date(year, month, day-89, 0, 0, 0, 0, time.UTC)
		ddasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &destinationEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationPaymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "60",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		oneHundredAndTwentyDays := 120
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &oneHundredAndTwentyDays,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reviewedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusDenied,
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status:         models.PaymentRequestStatusReviewed,
					SequenceNumber: 2,
				},
			},
		}, nil)

		destinationEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		ddasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &destinationEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: destinationPaymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "60",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    ddasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    ddasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		oneHundredAndTwentyDays := 120
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &oneHundredAndTwentyDays,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reviewedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)
		suite.Nil(sitBalances)
	})

	suite.Run("returns nil for pending payment request without SIT service items", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		oneHundredAndTwentyDays := 120
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &oneHundredAndTwentyDays,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			}, {
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			}, {
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			}, {
				Model:    shipment,
				LinkOnly: true,
			}, {
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), pendingPaymentRequest.ID)
		suite.NoError(err)
		suite.Nil(sitBalances)
	})

	suite.Run("returns zero authorized days for pending payment request shipment without a set SITDaysAllowance", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusApproved,
					ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
		}, nil)

		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			}, {
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			}, {
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			}, {
				Model:    shipment,
				LinkOnly: true,
			}, {
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    pendingPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:       models.MTOShipmentStatusApproved,
					ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		reviewedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: originEntryDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestStart,
				},
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: paymentEndDate.Format("2006-01-02"),
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameSITPaymentRequestEnd,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
		}, nil)

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		factory.BuildPaymentServiceItemParam(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItemParam{
					Value: "30",
				},
			},
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameNumberDaysSIT,
				},
			},
			{
				Model:    doasitParam.PaymentServiceItem,
				LinkOnly: true,
			},
			{
				Model:    reviewedPaymentRequest,
				LinkOnly: true,
			},
			{
				Model:    doasit,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

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
