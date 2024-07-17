package pptasapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/pptasapi"
	pptasops "github.com/transcom/mymove/pkg/gen/pptasapi/pptasoperations"
	"github.com/transcom/mymove/pkg/handlers"
	report "github.com/transcom/mymove/pkg/services/report"
)

func NewPPTASAPI(handlerConfig handlers.HandlerConfig) *pptasops.MymoveAPI {
	pptasSpec, err := loads.Analyzed(pptasapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	pptasAPI := pptasops.NewMymoveAPI(pptasSpec)
	pptasAPI.ServeError = handlers.ServeCustomError

	pptasAPI.MovesListReportsHandler = ListReportsHandler{
		HandlerConfig:     handlerConfig,
		ReportListFetcher: report.NewReportListFetcher(),
	}

	return pptasAPI
}
