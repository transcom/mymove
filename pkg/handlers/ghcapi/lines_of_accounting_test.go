package ghcapi

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	linesofaccountingop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/lines_of_accounting"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	"github.com/transcom/mymove/pkg/services/mocks"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
)

func (suite *HandlerSuite) TestLinesOfAccountingRequestLineOfAccountingHandler() {
	// LOA service fetcher for the API
	createLoaFetcher := func() services.LineOfAccountingFetcher {
		tacFetcher := transportationaccountingcode.NewTransportationAccountingCodeFetcher()
		return lineofaccounting.NewLinesOfAccountingFetcher(tacFetcher)
	}

	// Good TAC and LOA linked by LoaSysId
	buildGoodTacAndLoa := func() {
		now := time.Now()
		startDate := now.AddDate(-1, 0, 0)
		endDate := now.AddDate(1, 0, 0)
		tacCode := "GOOD"

		loa := factory.BuildLineOfAccounting(suite.DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaBgnDt:   &startDate,
					LoaEndDt:   &endDate,
					LoaSysID:   models.StringPointer("1234567890"),
					LoaHsGdsCd: models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
				},
			},
		}, nil)
		factory.BuildTransportationAccountingCodeWithoutAttachedLoa(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               tacCode,
					TrnsprtnAcntBgnDt: &startDate,
					TrnsprtnAcntEndDt: &endDate,
					TacFnBlModCd:      models.StringPointer("1"),
					LoaSysID:          loa.LoaSysID,
				},
			},
		}, nil)
	}

	// Create test cases without mocks
	testCases := []struct {
		name                         string
		shouldGenerateGoodTacAndLoa  bool
		departmentIndicator          *string
		effectiveDate                *strfmt.Date
		tacCode                      *string
		expectedResponse             middleware.Responder
		shouldTheBodyBeCompletelyNil bool
		expectedStatus               int
	}{
		{
			name:                         "Returns a 200 on fetching a good TAC and LOA",
			shouldGenerateGoodTacAndLoa:  true,
			departmentIndicator:          models.StringPointer(models.DepartmentIndicatorARMY.String()),
			effectiveDate:                (*strfmt.Date)(models.TimePointer(time.Now())),
			tacCode:                      models.StringPointer("GOOD"),
			shouldTheBodyBeCompletelyNil: false,
			expectedResponse:             &linesofaccountingop.RequestLineOfAccountingOK{},
		},
		{
			name:                         "Returns a 400 on nil body",
			shouldGenerateGoodTacAndLoa:  false,
			departmentIndicator:          nil,
			effectiveDate:                nil,
			tacCode:                      nil,
			shouldTheBodyBeCompletelyNil: true,
			expectedResponse:             &linesofaccountingop.RequestLineOfAccountingBadRequest{},
		},
		{
			name:                         "Returns a 400 on nil department indicator",
			shouldGenerateGoodTacAndLoa:  false,
			departmentIndicator:          nil,
			effectiveDate:                nil,
			tacCode:                      nil,
			shouldTheBodyBeCompletelyNil: false,
			expectedResponse:             &linesofaccountingop.RequestLineOfAccountingBadRequest{},
		},
		{
			name:                         "Return 200 if TAC cannot be found",
			shouldGenerateGoodTacAndLoa:  false,
			departmentIndicator:          models.StringPointer(models.DepartmentIndicatorARMY.String()),
			effectiveDate:                (*strfmt.Date)(models.TimePointer(time.Now())),
			tacCode:                      models.StringPointer("BAD"), // This may break in the future if TAC codes are enforced to be a minimum of 4 characters
			shouldTheBodyBeCompletelyNil: false,
			expectedResponse:             &linesofaccountingop.RequestLineOfAccountingOK{},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.shouldGenerateGoodTacAndLoa {
				buildGoodTacAndLoa()
			}

			loaFetcher := createLoaFetcher()
			handler := LinesOfAccountingRequestLineOfAccountingHandler{
				HandlerConfig:           suite.HandlerConfig(),
				LineOfAccountingFetcher: loaFetcher,
			}

			// Handle nil body
			var departmentIndicator *ghcmessages.DepartmentIndicator
			if tc.departmentIndicator != nil {
				deptIndicator := ghcmessages.DepartmentIndicator(*tc.departmentIndicator)
				departmentIndicator = &deptIndicator
			}

			var effectiveDate strfmt.Date
			if tc.effectiveDate != nil {
				effectiveDate = *tc.effectiveDate
			}

			var tacCode string
			if tc.tacCode != nil {
				tacCode = *tc.tacCode
			}

			var body *ghcmessages.FetchLineOfAccountingPayload
			if !tc.shouldTheBodyBeCompletelyNil {
				body = &ghcmessages.FetchLineOfAccountingPayload{
					DepartmentIndicator: departmentIndicator,
					EffectiveDate:       effectiveDate,
					TacCode:             tacCode,
				}
			}

			req := httptest.NewRequest("POST", "/lines-of-accounting", nil)
			params := linesofaccountingop.RequestLineOfAccountingParams{
				HTTPRequest: req,
				Body:        body,
			}

			response := handler.Handle(params)
			suite.IsType(tc.expectedResponse, response)
		})
	}

	// Run mock tests
	suite.Run("Returns 200 on LOA fetcher giving sql err rows not found", func() {
		mockLoaFetcher := &mocks.LineOfAccountingFetcher{}
		handler := LinesOfAccountingRequestLineOfAccountingHandler{
			HandlerConfig:           suite.HandlerConfig(),
			LineOfAccountingFetcher: mockLoaFetcher,
		}
		req := httptest.NewRequest("POST", "/lines-of-accounting", nil)
		departmentIndicator := ghcmessages.DepartmentIndicator(models.DepartmentIndicatorARMY)
		params := linesofaccountingop.RequestLineOfAccountingParams{
			HTTPRequest: req,
			Body: &ghcmessages.FetchLineOfAccountingPayload{
				DepartmentIndicator: &departmentIndicator,
				EffectiveDate:       strfmt.Date(time.Now()),
				TacCode:             "MOCK",
			},
		}

		mockLoaFetcher.On("FetchLongLinesOfAccounting", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]models.LineOfAccounting{}, sql.ErrNoRows)

		response := handler.Handle(params)
		suite.IsType(&linesofaccountingop.RequestLineOfAccountingOK{}, response)
	})

	suite.Run("Returns 500 on LOA fetcher erroring without sql no rows", func() {
		mockLoaFetcher := &mocks.LineOfAccountingFetcher{}
		handler := LinesOfAccountingRequestLineOfAccountingHandler{
			HandlerConfig:           suite.HandlerConfig(),
			LineOfAccountingFetcher: mockLoaFetcher,
		}
		req := httptest.NewRequest("POST", "/lines-of-accounting", nil)
		departmentIndicator := ghcmessages.DepartmentIndicator(models.DepartmentIndicatorARMY)
		params := linesofaccountingop.RequestLineOfAccountingParams{
			HTTPRequest: req,
			Body: &ghcmessages.FetchLineOfAccountingPayload{
				DepartmentIndicator: &departmentIndicator,
				EffectiveDate:       strfmt.Date(time.Now()),
				TacCode:             "MOCK",
			},
		}

		mockLoaFetcher.On("FetchLongLinesOfAccounting", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]models.LineOfAccounting{}, errors.New("Mock error"))

		response := handler.Handle(params)
		suite.IsType(&linesofaccountingop.RequestLineOfAccountingInternalServerError{}, response)
	})
}
