package pptasapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	"github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	report "github.com/transcom/mymove/pkg/services/pptas_report"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
)

func NewPPTASAPI(handlerConfig handlers.HandlerConfig) *pptasops.MymoveAPI {
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

	pptasAPI.MovesPptasReportsHandler = PPTASReportsHandler{
		HandlerConfig:          handlerConfig,
		PPTASReportListFetcher: report.NewPPTASReportListFetcher(ppmEstimator, moveFetcher, tacFetcher, loaFetcher),
	}

	return pptasAPI
}
