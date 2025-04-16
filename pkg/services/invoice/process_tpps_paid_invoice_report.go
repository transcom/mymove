package invoice

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
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
func NewTPPSPaidInvoiceReportProcessor() services.TPPSPaidInvoiceReportProcessor {
	return &tppsPaidInvoiceReportProcessor{}
}

// ProcessFile parses a TPPS paid invoice report response and updates the payment request status
func (t *tppsPaidInvoiceReportProcessor) ProcessFile(appCtx appcontext.AppContext, TPPSPaidInvoiceReportFilePath string, stringTPPSPaidInvoiceReport string) error {

	if TPPSPaidInvoiceReportFilePath == "" {
		appCtx.Logger().Info("No valid filepath found to process TPPS Paid Invoice Report", zap.String("TPPSPaidInvoiceReportFilePath", TPPSPaidInvoiceReportFilePath))
		return nil
	}
	tppsPaidInvoiceReport := tppsReponse.TPPSData{}

	appCtx.Logger().Info(fmt.Sprintf("Processing filepath: %s\n", TPPSPaidInvoiceReportFilePath))

	tppsData, err := tppsPaidInvoiceReport.Parse(appCtx, TPPSPaidInvoiceReportFilePath)
	if err != nil {
		appCtx.Logger().Error("unable to parse TPPS paid invoice report", zap.Error(err))
		return fmt.Errorf("unable to parse TPPS paid invoice report")
	}

	if tppsData != nil {
		appCtx.Logger().Info(fmt.Sprintf("Successfully parsed data from the TPPS paid invoice report: %s", TPPSPaidInvoiceReportFilePath))
		verrs, processedRowCount, errorProcessingRowCount, err := t.StoreTPPSPaidInvoiceReportInDatabase(appCtx, tppsData)
		if err != nil {
			return err
		} else if verrs.HasAny() {
			return verrs
		} else {
			appCtx.Logger().Info("Stored TPPS Paid Invoice Report information in the database")
			appCtx.Logger().Info(fmt.Sprintf("Rows successfully stored in DB: %d", processedRowCount))
			appCtx.Logger().Info(fmt.Sprintf("Rows not stored in DB due to foreign key constraint or other error: %d", errorProcessingRowCount))
		}

		var paymentRequestWithStatusUpdatedToPaid = map[string]string{}

		// For the data in the TPPS Paid Invoice Report, find the payment requests that match the
		// invoice numbers of the rows in the report and update the payment request status to PAID
		updatedPaymentRequestStatusCount := 0
		for _, tppsDataForOnePaymentRequest := range tppsData {
			var paymentRequest models.PaymentRequest

			err = appCtx.DB().Q().
				Where("payment_requests.payment_request_number = ?", tppsDataForOnePaymentRequest.InvoiceNumber).
				First(&paymentRequest)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					appCtx.Logger().Warn(fmt.Sprintf("No matching existing payment request found for invoice number %s, can't update status to PAID", tppsDataForOnePaymentRequest.InvoiceNumber))
					continue
				} else {
					appCtx.Logger().Error(fmt.Sprintf("Database error while looking up payment request for invoice number %s", tppsDataForOnePaymentRequest.InvoiceNumber), zap.Error(err))
					continue
				}
			}

			if paymentRequest.ID == uuid.Nil {
				appCtx.Logger().Error(fmt.Sprintf("Invalid payment request ID for invoice number %s", tppsDataForOnePaymentRequest.InvoiceNumber))
				continue
			}

			_, paymentRequestExistsInUpdatedStatusMap := paymentRequestWithStatusUpdatedToPaid[paymentRequest.ID.String()]
			if !paymentRequestExistsInUpdatedStatusMap {
				paymentRequest.Status = models.PaymentRequestStatusPaid
				err = appCtx.DB().Update(&paymentRequest)
				if err != nil {
					appCtx.Logger().Info(fmt.Sprintf("Failure updating payment request %s to PAID status", paymentRequest.PaymentRequestNumber))
					continue
				} else {
					if tppsDataForOnePaymentRequest.InvoiceNumber != uuid.Nil.String() && paymentRequest.ID != uuid.Nil {
						t.logTPPSInvoiceReportWithPaymentRequest(appCtx, tppsDataForOnePaymentRequest, paymentRequest)
					}
					updatedPaymentRequestStatusCount += 1
					paymentRequestWithStatusUpdatedToPaid[paymentRequest.ID.String()] = paymentRequest.PaymentRequestNumber
				}
			}
		}
		appCtx.Logger().Info(fmt.Sprintf("Payment requests that had status updated to PAID in DB: %d", updatedPaymentRequestStatusCount))

		return nil
	} else {
		appCtx.Logger().Info("No TPPS Paid Invoice Report data was parsed, so no data was stored in the database")
	}

	return nil
}

func (t *tppsPaidInvoiceReportProcessor) EDIType() models.EDIType {
	return models.TPPSPaidInvoiceReport
}

func (t *tppsPaidInvoiceReportProcessor) logTPPSInvoiceReportWithPaymentRequest(appCtx appcontext.AppContext, tppsResponse tppsReponse.TPPSData, paymentRequest models.PaymentRequest) {
	appCtx.Logger().Info("Updated payment request status to PAID",
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

func (t *tppsPaidInvoiceReportProcessor) StoreTPPSPaidInvoiceReportInDatabase(appCtx appcontext.AppContext, tppsData []tppsReponse.TPPSData) (*validate.Errors, int, int, error) {
	var verrs *validate.Errors
	var failedEntries []error
	DateParamFormat := "2006-01-02"
	processedRowCount := 0
	errorProcessingRowCount := 0

	for _, tppsEntry := range tppsData {
		timeOfTPPSCreatedDocumentDate, err := time.Parse(DateParamFormat, tppsEntry.TPPSCreatedDocumentDate)
		if err != nil {
			appCtx.Logger().Warn("unable to parse TPPSCreatedDocumentDate", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		timeOfSellerPaidDate, err := time.Parse(DateParamFormat, tppsEntry.SellerPaidDate)
		if err != nil {
			appCtx.Logger().Warn("unable to parse SellerPaidDate", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		invoiceTotalChargesInMillicents, err := priceToMillicents(tppsEntry.InvoiceTotalCharges)
		if err != nil {
			appCtx.Logger().Warn("unable to parse InvoiceTotalCharges", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		intLineBillingUnits, err := strconv.Atoi(tppsEntry.LineBillingUnits)
		if err != nil {
			appCtx.Logger().Warn("unable to parse LineBillingUnits", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		lineUnitPriceInMillicents, err := priceToMillicents(tppsEntry.LineUnitPrice)
		if err != nil {
			appCtx.Logger().Warn("unable to parse LineUnitPrice", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		lineNetChargeInMillicents, err := priceToMillicents(tppsEntry.LineNetCharge)
		if err != nil {
			appCtx.Logger().Warn("unable to parse LineNetCharge", zap.String("invoiceNumber", tppsEntry.InvoiceNumber), zap.Error(err))
			failedEntries = append(failedEntries, fmt.Errorf("invoiceNumber %s: %v", tppsEntry.InvoiceNumber, err))
			continue
		}

		txnErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
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
				if isForeignKeyConstraintViolation(err) {
					appCtx.Logger().Warn(fmt.Sprintf("skipping entry due to missing foreign key reference for invoice number %s", tppsEntry.InvoiceNumber))
					failedEntries = append(failedEntries, fmt.Errorf("invoice number %s: foreign key constraint violation", tppsEntry.InvoiceNumber))
					errorProcessingRowCount += 1
					return fmt.Errorf("rolling back transaction to prevent blocking")
				}

				appCtx.Logger().Error(fmt.Sprintf("failed to save entry for invoice number %s", tppsEntry.InvoiceNumber), zap.Error(err))
				failedEntries = append(failedEntries, fmt.Errorf("invoice number %s: %v", tppsEntry.InvoiceNumber, err))
				errorProcessingRowCount += 1
				return fmt.Errorf("rolling back transaction to prevent blocking")
			}

			appCtx.Logger().Info(fmt.Sprintf("successfully saved entry in DB for invoice number: %s", tppsEntry.InvoiceNumber))
			processedRowCount += 1
			return nil
		})

		if txnErr != nil {
			appCtx.Logger().Error(fmt.Sprintf("transaction error for invoice number %s", tppsEntry.InvoiceNumber), zap.Error(txnErr))
			errorProcessingRowCount += 1
		}
	}

	if len(failedEntries) > 0 {
		for _, err := range failedEntries {
			appCtx.Logger().Error("failed entry", zap.Error(err))
		}
	}

	// Return verrs but not a hard failure so we can process the rest of the entries
	return verrs, processedRowCount, errorProcessingRowCount, nil
}

func isForeignKeyConstraintViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23503"
	}
	return false
}
