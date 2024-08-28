package invoice

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	tppsReponse "github.com/transcom/mymove/pkg/edi/tpps_paid_invoice_report"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type tppsPaidInvoiceReportProcessor struct {
}

// TPPSData represents TPPS paid invoice report data
type TPPSData struct {
	InvoiceNumber             string
	TPPSCreatedDocumentDate   string
	SellerPaidDate            string
	InvoiceTotalCharges       string
	LineDescription           string
	ProductDescription        string
	LineBillingUnits          string
	LineUnitPrice             string
	LineNetCharge             string
	POTCN                     string
	LineNumber                string
	FirstNoteCode             string
	FirstNoteCodeDescription  string
	FirstNoteTo               string
	FirstNoteMessage          string
	SecondNoteCode            string
	SecondNoteCodeDescription string
	SecondNoteTo              string
	SecondNoteMessage         string
	ThirdNoteCode             string
	ThirdNoteCodeDescription  string
	ThirdNoteTo               string
	ThirdNoteMessage          string
}

// NewTPPSPaidInvoiceReportProcessor returns a new TPPS paid invoice report processor
func NewTPPSPaidInvoiceReportProcessor() services.SyncadaFileProcessor {

	return &tppsPaidInvoiceReportProcessor{}
}

// ProcessFile parses a TPPS paid invoice report response and updates the payment request status
func (t *tppsPaidInvoiceReportProcessor) ProcessFile(appCtx appcontext.AppContext, TPPSPaidInvoiceReportFilePath string, stringTPPSPaidInvoiceReport string) error {
	tppsPaidInvoiceReport := tppsReponse.TPPSData{}

	tppsData, err := tppsPaidInvoiceReport.Parse(TPPSPaidInvoiceReportFilePath, "")
	if err != nil {
		appCtx.Logger().Error("unable to parse TPPS paid invoice report", zap.Error(err))
		return fmt.Errorf("unable to parse TPPS paid invoice report")
	} else {
		appCtx.Logger().Info("Successfully parsed TPPS Paid Invoice Report")
	}

	appCtx.Logger().Info("RECEIVED: TPPS Paid Invoice Report Processor received a TPPS Paid Invoice Report")

	if tppsData != nil {
		verrs, errs := t.StoreTPPSPaidInvoiceReportInDatabase(appCtx, tppsData)
		if err != nil {
			return errs
		} else if verrs.HasAny() {
			return verrs
		} else {
			appCtx.Logger().Info("Successfully stored TPPS Paid Invoice Report information in the database")
		}

		transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			var paymentRequestWithStatusUpdatedToPaid = map[string]string{}

			// For the data in the TPPS Paid Invoice Report, find the payment requests that match the
			// invoice numbers of the rows in the report and update the payment request status to PAID
			for _, tppsDataForOnePaymentRequest := range tppsData {
				var paymentRequest models.PaymentRequest

				err = txnAppCtx.DB().Q().
					Where("payment_requests.payment_request_number = ?", tppsDataForOnePaymentRequest.InvoiceNumber).
					First(&paymentRequest)

				if err != nil {
					return err
				}

				// Since there can be many rows in a TPPS report that reference the same payment request, we want
				// to keep track of which payment requests we've already updated the status to PAID for and
				// only update it's status once, using a map to keep track of already updated payment requests
				_, paymentRequestExistsInUpdatedStatusMap := paymentRequestWithStatusUpdatedToPaid[paymentRequest.ID.String()]
				if !paymentRequestExistsInUpdatedStatusMap {
					paymentRequest.Status = models.PaymentRequestStatusPaid
					err = txnAppCtx.DB().Update(&paymentRequest)
					if err != nil {
						txnAppCtx.Logger().Error("failure updating payment request to PAID", zap.Error(err))
						return fmt.Errorf("failure updating payment request status to PAID: %w", err)
					}

					txnAppCtx.Logger().Info("SUCCESS: TPPS Paid Invoice Report Processor updated Payment Request to PAID status")
					t.logTPPSInvoiceReportWithPaymentRequest(txnAppCtx, tppsDataForOnePaymentRequest, paymentRequest)

					paymentRequestWithStatusUpdatedToPaid[paymentRequest.ID.String()] = paymentRequest.PaymentRequestNumber
				}
			}
			return nil
		})

		if transactionError != nil {
			appCtx.Logger().Error(transactionError.Error())
			return transactionError
		}
		return nil
	}

	return nil
}

func (t *tppsPaidInvoiceReportProcessor) EDIType() models.EDIType {
	return models.TPPSPaidInvoiceReport
}

func (t *tppsPaidInvoiceReportProcessor) logTPPSInvoiceReportWithPaymentRequest(appCtx appcontext.AppContext, tppsResponse tppsReponse.TPPSData, paymentRequest models.PaymentRequest) {
	appCtx.Logger().Info("TPPS Paid Invoice Report log",
		zap.String("TPPSPaidInvoiceReportEntry.InvoiceNumber", tppsResponse.InvoiceNumber),
		zap.String("PaymentRequestNumber", paymentRequest.PaymentRequestNumber),
		zap.String("PaymentRequest.Status", string(paymentRequest.Status)),
		zap.String("PaymentRequest.ID", paymentRequest.ID.String()),
	)
}

func getPriceParts(rawPrice string) (int, int, int, error) {
	// Get rid of a dollar sign if there is one.
	basePrice := strings.Replace(rawPrice, "$", "", -1)

	// Split the string on the decimal point.
	priceParts := strings.Split(basePrice, ".")

	integerPart, err := strconv.Atoi(priceParts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not convert integer part of price [%s]", rawPrice)
	}

	decimalsExistOnDollarAmount := len(priceParts) == 2
	expectedDecimalPlaces := 0
	if decimalsExistOnDollarAmount && len(priceParts[1]) > 0 {
		expectedDecimalPlaces = len(priceParts[1])
	}

	fractionalPart := 0
	if decimalsExistOnDollarAmount {
		fractionalPart, err = strconv.Atoi(priceParts[1])
	}
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not convert fractional part of price [%s]", rawPrice)
	}

	return integerPart, fractionalPart, expectedDecimalPlaces, nil
}

func priceToMillicents(rawPrice string) (int, error) {
	integerPart, fractionalPart, expectedDecimalPlaces, err := getPriceParts(rawPrice)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}
	var millicents int
	millicents = (integerPart * 100000)
	if expectedDecimalPlaces == 1 {
		millicents += (fractionalPart * 10000)
	} else if expectedDecimalPlaces == 2 {
		millicents += (fractionalPart * 1000)
	} else if expectedDecimalPlaces == 3 {
		millicents += (fractionalPart * 100)
	} else if expectedDecimalPlaces == 4 {
		millicents += (fractionalPart * 10)
	}
	return millicents, nil
}

func (t *tppsPaidInvoiceReportProcessor) StoreTPPSPaidInvoiceReportInDatabase(appCtx appcontext.AppContext, tppsData []tppsReponse.TPPSData) (*validate.Errors, error) {
	var verrs *validate.Errors
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		DateParamFormat := "2006-01-02"

		for _, tppsEntry := range tppsData {
			timeOfTPPSCreatedDocumentDate, err := time.Parse(DateParamFormat, tppsEntry.TPPSCreatedDocumentDate)
			if err != nil {
				appCtx.Logger().Info("unable to parse TPPSCreatedDocumentDate from TPPS paid invoice report", zap.Error(err))
			}
			timeOfSellerPaidDate, err := time.Parse(DateParamFormat, tppsEntry.SellerPaidDate)
			if err != nil {
				appCtx.Logger().Info("unable to parse SellerPaidDate from TPPS paid invoice report", zap.Error(err))
				return verrs
			}
			invoiceTotalChargesInMillicents, err := priceToMillicents(tppsEntry.InvoiceTotalCharges)
			if err != nil {
				appCtx.Logger().Info("unable to parse InvoiceTotalCharges from TPPS paid invoice report", zap.Error(err))
				return verrs
			}
			intLineBillingUnits, err := strconv.Atoi(tppsEntry.LineBillingUnits)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineBillingUnits from TPPS paid invoice report", zap.Error(err))
				return verrs
			}
			lineUnitPriceInMillicents, err := priceToMillicents(tppsEntry.LineUnitPrice)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineUnitPrice from TPPS paid invoice report", zap.Error(err))
				return verrs
			}
			lineNetChargeInMillicents, err := priceToMillicents(tppsEntry.LineNetCharge)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineNetCharge from TPPS paid invoice report", zap.Error(err))
				return verrs
			}

			tppsEntryModel := models.TPPSPaidInvoiceReportEntry{
				InvoiceNumber:                   tppsEntry.InvoiceNumber,
				TPPSCreatedDocumentDate:         &timeOfTPPSCreatedDocumentDate,
				SellerPaidDate:                  timeOfSellerPaidDate,
				InvoiceTotalChargesInMillicents: unit.Millicents(invoiceTotalChargesInMillicents),
				LineDescription:                 tppsEntry.LineDescription,
				ProductDescription:              tppsEntry.ProductDescription,
				LineBillingUnits:                intLineBillingUnits,
				LineUnitPrice:                   unit.Millicents(lineUnitPriceInMillicents),
				LineNetCharge:                   unit.Millicents(lineNetChargeInMillicents),
				POTCN:                           tppsEntry.POTCN,
				LineNumber:                      tppsEntry.LineNumber,
				FirstNoteCode:                   &tppsEntry.FirstNoteCode,             // #nosec G601
				FirstNoteDescription:            &tppsEntry.FirstNoteCodeDescription,  // #nosec G601
				FirstNoteCodeTo:                 &tppsEntry.FirstNoteTo,               // #nosec G601
				FirstNoteCodeMessage:            &tppsEntry.FirstNoteMessage,          // #nosec G601
				SecondNoteCode:                  &tppsEntry.SecondNoteCode,            // #nosec G601
				SecondNoteDescription:           &tppsEntry.SecondNoteCodeDescription, // #nosec G601
				SecondNoteCodeTo:                &tppsEntry.SecondNoteTo,              // #nosec G601
				SecondNoteCodeMessage:           &tppsEntry.SecondNoteMessage,         // #nosec G601
				ThirdNoteCode:                   &tppsEntry.ThirdNoteCode,             // #nosec G601
				ThirdNoteDescription:            &tppsEntry.ThirdNoteCodeDescription,  // #nosec G601
				ThirdNoteCodeTo:                 &tppsEntry.ThirdNoteTo,               // #nosec G601
				ThirdNoteCodeMessage:            &tppsEntry.ThirdNoteMessage,          // #nosec G601
			}

			verrs, err = txnAppCtx.DB().ValidateAndSave(&tppsEntryModel)
			if err != nil {
				appCtx.Logger().Error("failure saving entry from TPPS paid invoice report", zap.Error(err))
				return err
			}
		}

		return nil
	})

	if transactionError != nil {
		appCtx.Logger().Error(transactionError.Error())
		return verrs, transactionError
	}
	if verrs.HasAny() {
		appCtx.Logger().Error("unable to process TPPS paid invoice report", zap.Error(verrs))
		return verrs, nil
	}

	return nil, nil
}
