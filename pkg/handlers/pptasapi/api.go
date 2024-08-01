package pptasapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	"github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	report "github.com/transcom/mymove/pkg/services/report"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
)

func NewPPTASApiHandler(handlerConfig handlers.HandlerConfig) http.Handler {

	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasAPI := pptasops.NewMymoveAPI(pptasSpec)
	pptasAPI.ServeError = handlers.ServeCustomError

	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	moveFetcher := move.NewMoveFetcher()
	tacFetcher := transportationaccountingcode.NewTransportationAccountingCodeFetcher()
	loaFetcher := lineofaccounting.NewLinesOfAccountingFetcher(tacFetcher)

	pptasAPI.MovesReportsHandler = ReportsHandler{
		HandlerConfig:     handlerConfig,
		ReportListFetcher: report.NewReportListFetcher(ppmEstimator, moveFetcher, tacFetcher, loaFetcher),
	}

	return pptasAPI.Serve(nil)
}
