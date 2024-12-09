package paymentrequest

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestListShipmentPaymentSITBalance() {
	service := NewPaymentRequestShipmentsSITBalance()

	setUpShipmentWith120DaysOfAuthorizedSIT := func(db *pop.Connection) (models.Move, models.MTOShipment) {
		move := factory.BuildAvailableToPrimeMove(db, []factory.Customization{
			{
				Model: models.Move{
					Locator: "PARSIT",
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(db, []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		thirtyDaySITExtensionRequest := 30
		factory.BuildSITDurationUpdate(db, []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.SITDurationUpdate{
					Status:        models.SITExtensionStatusApproved,
					RequestedDays: thirtyDaySITExtensionRequest,
					ApprovedDays:  &thirtyDaySITExtensionRequest,
				},
			},
		}, nil)
		return move, shipment
	}

	suite.Run("returns only pending SIT status when there are no previous payments", func() {
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

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
		originEntryDate := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
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
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &originEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
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
		suite.Equal(89, pendingSITBalance.TotalSITDaysRemaining)
		// Authorized end date should be 120 days from the entry date.
		// Only add 119 to the entry date because that last day counts as one too
		// To understand this, use a date range calculator and set start date
		// to AUG 12 2024, end date to DEC 9 2024, and select to "Include end date in calculation"
		suite.Equal(doasit.SITEntryDate.AddDate(0, 0, 119).String(), pendingSITBalance.TotalSITEndDate.UTC().String())
	})

	suite.Run("calculates pending destination SIT balance when origin was invoiced previously", func() {
		// Set up a move with a shipment that has a 120 days of authorized SIT
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

		// Create a reviewed payment request
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

		// Attach origin SIT service items to the reviewed payment request.
		// Create the accompanying dates that uses up 30 of the 120 authorized days of SIT.
		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
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
		originDepartureDate := originEntryDate.AddDate(0, 0, 30)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &originEntryDate,
					SITDepartureDate: &originDepartureDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
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

		// Create a pending payment request
		pendingPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status:         models.PaymentRequestStatusPending,
					SequenceNumber: 2,
				},
			},
		}, nil)

		// Create the destination SIT service items and attach them to the pending payment request.
		// Set up the dates so that it went into storage 60 days ago and it is still in storage now.
		destinationEntryDate := time.Date(year, month, day-60, 0, 0, 0, 0, time.UTC)
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
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &destinationEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
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

		// 30 is what was previously billed on the reviewed payment request
		suite.Equal(30, *pendingSITBalance.PreviouslyBilledDays)
		suite.Equal(paymentEndDate.String(), pendingSITBalance.PreviouslyBilledEndDate.String())

		// 60 days is the pending amount on the pending payment request
		suite.Equal(60, pendingSITBalance.PendingSITDaysInvoiced)
		suite.Equal(destinationPaymentEndDate.String(), pendingSITBalance.PendingBilledEndDate.String())

		suite.Equal(120, pendingSITBalance.TotalSITDaysAuthorized)
		// 120 total authorized - 31 from origin SIT - 61 from destination SIT = 28 SIT days remaining
		suite.Equal(28, pendingSITBalance.TotalSITDaysRemaining)

		// 120 authorized - 31 already used - 1 to be inclusive of the last day = 88
		suite.Equal(ddasit.SITEntryDate.AddDate(0, 0, 88).String(), pendingSITBalance.TotalSITEndDate.UTC().String())
	})

	suite.Run("ignores including previously denied service items in SIT balance", func() {
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

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
		originDepartureDate := originEntryDate.AddDate(0, 0, 60)

		// Create the corresponding origin delivery service item that departed after 60 days of storage.
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &originEntryDate,
					SITDepartureDate: &originDepartureDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
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

		destinationEntryDate := time.Date(year, month, day-15, 0, 0, 0, 0, time.UTC)
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
		// Build the accompanying Destination Delivery service item
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &destinationEntryDate,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
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
		suite.Equal(43, pendingSITBalance.TotalSITDaysRemaining)
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

	suite.Run("returns SIT balance for payment request when there is only a future SIT no current SIT", func() {

		// Set up a move with a shipment that has a 120 days of authorized SIT
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

		year, month, day := time.Now().Date()
		// originEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		originEntryDate := time.Date(year, month, day+90, 0, 0, 0, 0, time.UTC)

		shipment.OriginSITAuthEndDate = &originEntryDate

		// Attach origin SIT service items to the reviewed payment request.
		// Create the accompanying dates that uses up 30 of the 120 authorized days of SIT.
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

		// Create a reviewed payment request
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
			{
				Model:    doasit,
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

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)

		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
	})

	suite.Run("returns SIT balance for payment request when there is only a past SIT no current SIT", func() {

		// Set up a move with a shipment that has a 120 days of authorized SIT
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

		year, month, day := time.Now().Date()
		sitDepartureDate := time.Date(year, month, day-50, 0, 0, 0, 0, time.UTC)
		originEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)

		shipment.OriginSITAuthEndDate = &sitDepartureDate

		// Attach origin SIT service items to the reviewed payment request.
		// Create the accompanying dates that uses up 30 of the 120 authorized days of SIT.
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &originEntryDate,
					SITDepartureDate: &sitDepartureDate,
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

		// Create a reviewed payment request
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
			{
				Model:    doasit,
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

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)

		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
	})
	suite.Run("returns SIT balance for payment request when there is only a past SIT with an unapproved payment service item and no current SIT", func() {

		// Set up a move with a shipment that has a 120 days of authorized SIT
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

		year, month, day := time.Now().Date()
		sitDepartureDate := time.Date(year, month, day-50, 0, 0, 0, 0, time.UTC)
		originEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)

		shipment.OriginSITAuthEndDate = &sitDepartureDate

		// Attach origin SIT service items to the reviewed payment request.
		// Create the accompanying dates that uses up 30 of the 120 authorized days of SIT.
		doasit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:           models.MTOServiceItemStatusApproved,
					SITEntryDate:     &originEntryDate,
					SITDepartureDate: &sitDepartureDate,
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

		// Create a reviewed payment request
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
			{
				Model:    doasit,
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

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)

		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
	})
	suite.Run("returns SIT balance for payment request when there is only a future SIT with a unapproved payment service item. No current SIT", func() {

		// Set up a move with a shipment that has a 120 days of authorized SIT
		move, shipment := setUpShipmentWith120DaysOfAuthorizedSIT(suite.DB())

		year, month, day := time.Now().Date()
		// originEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		originEntryDate := time.Date(year, month, day+90, 0, 0, 0, 0, time.UTC)

		shipment.OriginSITAuthEndDate = &originEntryDate

		// Attach origin SIT service items to the reviewed payment request.
		// Create the accompanying dates that uses up 30 of the 120 authorized days of SIT.
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

		// Create a reviewed payment request
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
			{
				Model:    doasit,
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

		sitBalances, err := service.ListShipmentPaymentSITBalance(suite.AppContextForTest(), reviewedPaymentRequest.ID)
		suite.NoError(err)

		suite.Len(sitBalances, 1)

		pendingSITBalance := sitBalances[0]
		suite.Equal(shipment.ID.String(), pendingSITBalance.ShipmentID.String())
	})
}
