package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *pptasmessages.ClientError {
	instanceToUse := strfmt.UUID(traceID.String())
	payload := pptasmessages.ClientError{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: &instanceToUse,
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ListMove payload
func ListMove(move *models.Move) *pptasmessages.ListMove {
	if move == nil {
		return nil
	}
	payload := &pptasmessages.ListMove{
		ID:                 strfmt.UUID(move.ID.String()),
		MoveCode:           move.Locator,
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		OrderID:            strfmt.UUID(move.OrdersID.String()),
		ReferenceID:        *move.ReferenceID,
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
	}

	if move.PPMType != nil {
		payload.PpmType = *move.PPMType
	}

	return payload
}

// ListMoves payload
func ListMoves(moves *models.Moves) []*pptasmessages.ListMove {
	payload := make(pptasmessages.ListMoves, len(*moves))

	for i, m := range *moves {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListMove(&copyOfM)
	}
	return payload
}
