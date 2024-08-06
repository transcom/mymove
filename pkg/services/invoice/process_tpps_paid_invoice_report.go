package invoice

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	tppsReponse "github.com/transcom/mymove/pkg/edi/tpps_paid_invoice_report"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type tppsPaidInvoiceReportProcessor struct {
}

// NewTPPSPaidInvoiceReportProcessor returns a new TPPS paid invoice report processor
func NewTPPSPaidInvoiceReportProcessor() services.SyncadaFileProcessor {

	return &tppsPaidInvoiceReportProcessor{}
}

// ProcessFile parses a TPPS paid invoice report response and updates the payment request status
func (e *tppsPaidInvoiceReportProcessor) ProcessFile(appCtx appcontext.AppContext, _ string, stringTPPSPaidInvoiceReport string) error {
	tppsPaidInvoiceReport := tppsReponse.EDI{}

	tppsData, err := tppsPaidInvoiceReport.Parse(stringTPPSPaidInvoiceReport)
	if err != nil {
		appCtx.Logger().Error("unable to parse TPPS paid invoice report", zap.Error(err))
		return fmt.Errorf("unable to parse TPPS paid invoice report")
	} else {
		appCtx.Logger().Info("Successfully parsed TPPS Paid Invoice Report and found ", zap.Int("entries", len(tppsData)))
	}

	return nil
}

func (e *tppsPaidInvoiceReportProcessor) EDIType() models.EDIType {
	return models.TPPSPaidInvoiceReport
}
