package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	reportViolationop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/report_violations"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/services/mocks"
	reportviolationservice "github.com/transcom/mymove/pkg/services/report_violation"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetReportViolationByIDHandler() {
	// 200 response
	suite.Run("Successful fetch (integration) test", func() {

		handlerConfig := suite.HandlerConfig()
		fetcher := reportviolationservice.NewReportViolationFetcher()
		handler := GetReportViolationsHandler{
			HandlerConfig:          handlerConfig,
			ReportViolationFetcher: fetcher,
		}

		reportViolation := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})

		request := httptest.NewRequest("GET", fmt.Sprintf("/report-violations/%s",
			reportViolation.ReportID.String()), nil)
		params := reportViolationop.GetReportViolationsByReportIDParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(reportViolation.ReportID.String()),
		}
		response := handler.Handle(params)

		suite.IsType(&reportViolationop.GetReportViolationsByReportIDOK{}, response)
	})

	// 404 response
	suite.Run("404 response when service returns not found", func() {
		badID, _ := uuid.NewV4()

		handlerConfig := suite.HandlerConfig()
		mockFetcher := mocks.ReportViolationFetcher{}
		handler := GetReportViolationsHandler{
			HandlerConfig:          handlerConfig,
			ReportViolationFetcher: &mockFetcher,
		}

		mockFetcher.On("FetchReportViolationsByReportID",
			mock.AnythingOfType("*appcontext.appContext"),
			badID,
		).Return(nil, apperror.QueryError{})

		request := httptest.NewRequest("GET", fmt.Sprintf("/report-violations/%s",
			badID.String()), nil)
		params := reportViolationop.GetReportViolationsByReportIDParams{
			HTTPRequest: request,
			ReportID:    strfmt.UUID(badID.String()),
		}
		response := handler.Handle(params)

		suite.IsType(&reportViolationop.GetReportViolationsByReportIDInternalServerError{}, response)
	})
}
func (suite *HandlerSuite) TestAssociateReportViolationsHandler() {
	suite.Run("Successful POST", func() {

		handlerConfig := suite.HandlerConfig()
		creator := &mocks.ReportViolationsCreator{}
		handler := AssociateReportViolationsHandler{
			HandlerConfig:           handlerConfig,
			ReportViolationsCreator: creator,
		}

		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})

		body := ghcmessages.AssociateReportViolations{Violations: []strfmt.UUID{strfmt.UUID(violation.ID.String())}}

		request := httptest.NewRequest("POST", fmt.Sprintf("/moves/report-violations/%s", report.ID.String()), nil)

		params := reportViolationop.AssociateReportViolationsParams{
			HTTPRequest: request,
			Body:        &body,
		}

		creator.On("AssociateReportViolations",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.ReportViolations"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&reportViolationop.AssociateReportViolationsNoContent{}, response)
	})

	suite.Run("Unsuccessful POST", func() {

		handlerConfig := suite.HandlerConfig()
		creator := &mocks.ReportViolationsCreator{}
		handler := AssociateReportViolationsHandler{
			HandlerConfig:           handlerConfig,
			ReportViolationsCreator: creator,
		}

		report := testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		violation := testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})

		body := ghcmessages.AssociateReportViolations{Violations: []strfmt.UUID{strfmt.UUID(violation.ID.String())}}

		request := httptest.NewRequest("POST", fmt.Sprintf("/moves/report-violations/%s", report.ID.String()), nil)

		params := reportViolationop.AssociateReportViolationsParams{
			HTTPRequest: request,
			Body:        &body,
		}

		creator.On("AssociateReportViolations",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.ReportViolations"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(fmt.Errorf("error")).Once()

		response := handler.Handle(params)

		suite.Assertions.IsType(&reportViolationop.AssociateReportViolationsInternalServerError{}, response)
	})
}
