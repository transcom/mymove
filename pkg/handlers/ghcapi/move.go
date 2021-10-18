package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveHandler gets a move by locator
type GetMoveHandler struct {
	handlers.HandlerContext
	services.MoveFetcher
}

// Handle handles the getMove by locator request
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	locator := params.Locator
	if locator == "" {
		return moveop.NewGetMoveBadRequest()
	}

	move, err := h.FetchMove(appCtx, locator, nil)

	if err != nil {
		logger.Error("Error retrieving move by locator", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return moveop.NewGetMoveNotFound()
		default:
			return moveop.NewGetMoveInternalServerError()
		}
	}

	payload := payloads.Move(move)
	return moveop.NewGetMoveOK().WithPayload(payload)
}

type CreateFinancialReviewFlagHandler struct {
	handlers.HandlerContext
	financialReviewFlagCreator services.MoveFinancialReviewFlagCreator
}

// Handle flags a move for financial review
func (h CreateFinancialReviewFlagHandler) Handle(params moveop.FlagMoveForFinancialReviewParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	moveID := uuid.FromStringOrNil(params.MoveID.String())
	if moveID == uuid.Nil {
		errorMessage := fmt.Sprintf("unable to parse moveID as UUID: %s", params.MoveID.String())
		payload := payloadForValidationError("Unable to flag move for financial review", errorMessage, h.GetTraceID(), validate.NewErrors())
		return moveop.NewFlagMoveForFinancialReviewUnprocessableEntity().WithPayload(payload)
	}

	remarks := params.Body.Remarks
	if remarks == nil || *remarks == "" {
		payload := payloadForValidationError("Unable to flag move for financial review", "missing or empty remarks field", h.GetTraceID(), validate.NewErrors())
		return moveop.NewFlagMoveForFinancialReviewUnprocessableEntity().WithPayload(payload)
	}
	move, err := h.financialReviewFlagCreator.CreateFinancialReviewFlag(appCtx, moveID, *remarks)

	if err != nil {
		logger.Error("Error flagging move for financial review", zap.Error(err))
		switch err.(type) {
		case apperror.NotFoundError:
			return moveop.NewFlagMoveForFinancialReviewNotFound()
		default:
			return moveop.NewFlagMoveForFinancialReviewInternalServerError()
		}
	}

	payload := payloads.Move(move)
	return moveop.NewFlagMoveForFinancialReviewOK().WithPayload(payload)
}
