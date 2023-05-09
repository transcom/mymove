package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&reportViolationop.GetReportViolationsByReportIDOK{}, response)
		payload := response.(*reportViolationop.GetReportViolationsByReportIDOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&reportViolationop.GetReportViolationsByReportIDInternalServerError{}, response)
		payload := response.(*reportViolationop.GetReportViolationsByReportIDInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
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

		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&reportViolationop.AssociateReportViolationsNoContent{}, response)

		// Validate outgoing payload: no payload
	})

	suite.Run("Unsuccessful POST", func() {

		handlerConfig := suite.HandlerConfig()
		creator := &mocks.ReportViolationsCreator{}
		handler := AssociateReportViolationsHandler{
			HandlerConfig:           handlerConfig,
			ReportViolationsCreator: creator,
		}

		report := factory.BuildEvaluationReport(suite.DB(), nil, nil)
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

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&reportViolationop.AssociateReportViolationsInternalServerError{}, response)
		payload := response.(*reportViolationop.AssociateReportViolationsInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
