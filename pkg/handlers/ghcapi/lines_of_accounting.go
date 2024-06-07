package ghcapi

import (
	"database/sql"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	linesofaccountingop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/lines_of_accounting"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// Handler to handle requests for a line of accounting based on orders issue date,
// service member affiliation, and Transportation Accounting Code (TAC)
type LinesOfAccountingRequestLineOfAccountingHandler struct {
	handlers.HandlerConfig
	services.LineOfAccountingFetcher
}

// Handle requesting (Fetching) a line of accounting from a request payload
func (h LinesOfAccountingRequestLineOfAccountingHandler) Handle(params linesofaccountingop.RequestLineOfAccountingParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body

			if payload == nil {
				err := apperror.NewBadDataError("Invalid request for lines of accounting: params body is nil")
				appCtx.Logger().Error(err.Error())
				return linesofaccountingop.NewRequestLineOfAccountingBadRequest(), err
			}

			if payload.ServiceMemberAffiliation == nil {
				err := apperror.NewBadDataError("Invalid request for lines of accounting: service member affiliation is nil")
				return linesofaccountingop.NewRequestLineOfAccountingBadRequest(), err
			}

			loas, err := h.LineOfAccountingFetcher.FetchLongLinesOfAccounting(models.ServiceMemberAffiliation(*payload.ServiceMemberAffiliation), time.Time(payload.OrdersIssueDate), payload.TacCode, appCtx)
			if err != nil {
				if err == sql.ErrNoRows {
					// Either TAC or LOA service objects triggered a sql err for now rows
					// This error check will currently never be triggered, but in the case
					// of the service object being updated in the future, this will catch it and keep the API giving good errors
					// instead of defaulting to an internal server error
					errMsg := "Unable to find any lines of accounting based on the provided parameters"
					err := apperror.NewNotFoundError(uuid.Nil, errMsg)
					errPayload := &ghcmessages.Error{Message: &errMsg}

					return linesofaccountingop.NewRequestLineOfAccountingNotFound().WithPayload(errPayload), err

				}
				return linesofaccountingop.NewRequestLineOfAccountingInternalServerError(), err
			}
			if len(loas) == 0 {
				// No LOAs were identified with the provided paramters
				errMsg := "Unable to find any lines of accounting based on the provided parameters"
				err := apperror.NewNotFoundError(uuid.Nil, errMsg)
				errPayload := &ghcmessages.Error{Message: &errMsg}

				return linesofaccountingop.NewRequestLineOfAccountingNotFound().WithPayload(errPayload), err
			}

			// pick first one (sorted by FBMC, loa_bgn_dt, tac_fy_txt) inside the service object
			loa := loas[0]

			returnPayload := payloads.LineOfAccounting(&loa)

			return linesofaccountingop.NewRequestLineOfAccountingOK().WithPayload(returnPayload), nil
		})
}
