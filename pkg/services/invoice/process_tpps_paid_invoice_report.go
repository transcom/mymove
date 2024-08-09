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
func (t *tppsPaidInvoiceReportProcessor) ProcessFile(appCtx appcontext.AppContext, _ string, stringTPPSPaidInvoiceReport string) error {
	tppsPaidInvoiceReport := tppsReponse.EDI{}

	tppsData, err := tppsPaidInvoiceReport.Parse(stringTPPSPaidInvoiceReport)
	if err != nil {
		appCtx.Logger().Error("unable to parse TPPS paid invoice report", zap.Error(err))
		return fmt.Errorf("unable to parse TPPS paid invoice report")
	} else {
		appCtx.Logger().Info("Successfully parsed TPPS Paid Invoice Report")
	}

	appCtx.Logger().Info("RECEIVED: TPPS Paid Invoice Report Processor received a TPPS Paid Invoice Report")

	if tppsData != nil {
		verrs, errs := t.StoreTPPSPaidInvoiceReportInDatabase(appCtx, tppsData, "", stringTPPSPaidInvoiceReport)
		if err != nil {
			return errs
		}
		if verrs.HasAny() {
			return verrs
		}
	}

	return nil
}

func (t *tppsPaidInvoiceReportProcessor) EDIType() models.EDIType {
	return models.TPPSPaidInvoiceReport
}

func getPriceParts(rawPrice string) (int, int, int, error) {
	// Get rid of a dollar sign if there is one.
	basePrice := strings.Replace(rawPrice, "$", "", -1)

	// Split the string on the decimal point.
	priceParts := strings.Split(basePrice, ".")
	if len(priceParts) != 2 {
		return 0, 0, 0, fmt.Errorf("expected 2 price parts but found %d for price [%s]", len(priceParts), rawPrice)
	}

	integerPart, err := strconv.Atoi(priceParts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("could not convert integer part of price [%s]", rawPrice)
	}

	expectedDecimalPlaces := 0
	if len(priceParts[1]) > 0 {
		expectedDecimalPlaces = len(priceParts[1])
	}

	fractionalPart, err := strconv.Atoi(priceParts[1])
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
	if expectedDecimalPlaces == 0 {
		// . cents = 0 millicents
		millicents = (integerPart * 100000)
	} else if expectedDecimalPlaces == 1 {
		// .5 cents = 50000 millicents
		millicents = (integerPart * 100000) + (fractionalPart * 10000)
	} else if expectedDecimalPlaces == 2 {
		// .25 cents = 25000 millicents
		millicents = (integerPart * 100000) + (fractionalPart * 1000)
	} else if expectedDecimalPlaces == 3 {
		// .025 cents = 25000 millicents
		millicents = (integerPart * 100000) + (fractionalPart * 100)
	} else if expectedDecimalPlaces == 4 {
		// .0025 cents = 250 millicents
		millicents = (integerPart * 100000) + (fractionalPart * 10)
	}
	return millicents, nil
}

func (t *tppsPaidInvoiceReportProcessor) StoreTPPSPaidInvoiceReportInDatabase(appCtx appcontext.AppContext, tppsData []tppsReponse.TPPSData, _ string, stringTPPSPaidInvoiceReport string) (*validate.Errors, error) {
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
			}
			invoiceTotalChargesInMillicents, err := priceToMillicents(tppsEntry.InvoiceTotalCharges)
			if err != nil {
				appCtx.Logger().Info("unable to parse InvoiceTotalCharges from TPPS paid invoice report", zap.Error(err))
			}
			intLineBillingUnits, err := strconv.Atoi(tppsEntry.LineBillingUnits)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineBillingUnits from TPPS paid invoice report", zap.Error(err))
			}
			lineUnitPriceInMillicents, err := priceToMillicents(tppsEntry.LineUnitPrice)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineUnitPrice from TPPS paid invoice report", zap.Error(err))
			}
			lineNetChargeInMillicents, err := priceToMillicents(tppsEntry.LineNetCharge)
			if err != nil {
				appCtx.Logger().Info("unable to parse LineNetCharge from TPPS paid invoice report", zap.Error(err))
			}

			tppsEntryModel := models.TPPSPaidInvoiceReportEntry{
				InvoiceNumber:                   tppsEntry.InvoiceNumber,
				TPPSCreatedDocumentDate:         timeOfTPPSCreatedDocumentDate,
				SellerPaidDate:                  timeOfSellerPaidDate,
				InvoiceTotalChargesInMillicents: unit.Millicents(invoiceTotalChargesInMillicents),
				LineDescription:                 tppsEntry.LineDescription,
				ProductDescription:              tppsEntry.ProductDescription,
				LineBillingUnits:                intLineBillingUnits,
				LineUnitPrice:                   unit.Millicents(lineUnitPriceInMillicents),
				LineNetCharge:                   unit.Millicents(lineNetChargeInMillicents),
				POTCN:                           tppsEntry.POTCN,
				LineNumber:                      tppsEntry.LineNumber,
				FirstNoteCode:                   tppsEntry.FirstNoteCode,
				FirstNoteDescription:            tppsEntry.FirstNoteCodeDescription,
				FirstNoteCodeTo:                 tppsEntry.FirstNoteTo,
				FirstNoteCodeMessage:            tppsEntry.FirstNoteMessage,
				SecondNoteCode:                  tppsEntry.SecondNoteCode,
				SecondNoteDescription:           tppsEntry.SecondNoteCodeDescription,
				SecondNoteCodeTo:                tppsEntry.SecondNoteTo,
				SecondNoteCodeMessage:           tppsEntry.SecondNoteMessage,
				ThirdNoteCode:                   tppsEntry.ThirdNoteCode,
				ThirdNoteDescription:            tppsEntry.ThirdNoteCodeDescription,
				ThirdNoteCodeTo:                 tppsEntry.ThirdNoteTo,
				ThirdNoteCodeMessage:            tppsEntry.ThirdNoteMessage,
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
