package ghcapi

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	linesofaccountingop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/lines_of_accounting"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// Handler to handle requests for a line of accounting based on orders issue date,
// order's department indicator, and Transportation Accounting Code (TAC)
type LinesOfAccountingRequestLineOfAccountingHandler struct {
	handlers.HandlerConfig
	services.LineOfAccountingFetcher
}

// Handle requesting (Fetching) a line of accounting from a request payload
// It takes in the parameters of:
// - TAC
// - Order's department indicator
// - EffectiveDate
// And uses these parameters to filter the correct Line of Accounting for the provided TAC. It does this by filtering
// through both TAC and LOAs based on the provided code and effective date. The 'Effective Date' is the date
// that can be either the orders issued date (For HHG shipments), MTO approval date (For NTS shipments),
// or even the current date for NTS shipments with no approval yet (Just providing a preview to the office users per customer request)
// Effective date is used to find "Active" TGET data by searching for the TACs and LOAs with begin and end dates containing this date
func (h LinesOfAccountingRequestLineOfAccountingHandler) Handle(params linesofaccountingop.RequestLineOfAccountingParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body

			if payload == nil {
				err := apperror.NewBadDataError("Invalid request for lines of accounting: params body is nil")
				appCtx.Logger().Error(err.Error())
				return linesofaccountingop.NewRequestLineOfAccountingBadRequest(), err
			}

			if payload.DepartmentIndicator == nil {
				err := apperror.NewBadDataError("Invalid request for lines of accounting: department indicator is nil")
				appCtx.Logger().Error(err.Error())
				return linesofaccountingop.NewRequestLineOfAccountingBadRequest(), err
			}

			loas, err := h.LineOfAccountingFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicator(*payload.DepartmentIndicator), time.Time(payload.EffectiveDate), payload.TacCode, appCtx)
			if err != nil {
				if err == sql.ErrNoRows {
					// Either TAC or LOA service objects triggered a sql err for now rows
					// This error check will currently never be triggered, but in the case
					// of the service object being updated in the future, this will catch it and keep the API giving good errors
					// instead of defaulting to an internal server error

					errMsg := fmt.Sprintf("Unable to find any lines of accounting based on the provided parameters: departmentIndicator=%s, ordersIssueDate=%s, tacCode=%s", *payload.DepartmentIndicator, time.Time(payload.EffectiveDate), payload.TacCode)
					appCtx.Logger().Info(errMsg)
					// Do not return any payload here as no LOA was found
					return linesofaccountingop.NewRequestLineOfAccountingOK(), nil
				}
				return linesofaccountingop.NewRequestLineOfAccountingInternalServerError(), err
			}
			if len(loas) == 0 {
				// No LOAs were identified with the provided parameters
				// Return an empty 200 and log the error
				errMsg := fmt.Sprintf("Unable to find any lines of accounting based on the provided parameters: departmentIndicator=%s, ordersIssueDate=%s, tacCode=%s", *payload.DepartmentIndicator, time.Time(payload.EffectiveDate), payload.TacCode)
				appCtx.Logger().Info(errMsg)
				return linesofaccountingop.NewRequestLineOfAccountingOK(), nil
			}

			// pick first one (sorted by FBMC, loa_bgn_dt, tac_fy_txt) inside the service object
			loa := loas[0]

			returnPayload := payloads.LineOfAccounting(&loa)

			return linesofaccountingop.NewRequestLineOfAccountingOK().WithPayload(returnPayload), nil
		})
}
