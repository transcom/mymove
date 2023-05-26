package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primev2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ListMove payload
func ListMove(move *models.Move) *primev2messages.ListMove {
	if move == nil {
		return nil
	}
	payload := &primev2messages.ListMove{
		ID:                 "blahahahaha",
		MoveCode:           move.Locator,
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		OrderID:            strfmt.UUID(move.OrdersID.String()),
		ReferenceID:        *move.ReferenceID,
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
	}

	if move.PPMEstimatedWeight != nil {
		payload.PpmEstimatedWeight = int64(*move.PPMEstimatedWeight)
	}

	if move.PPMType != nil {
		payload.PpmType = *move.PPMType
	}

	return payload
}

// ListMoves payload
func ListMoves(moves *models.Moves) []*primev2messages.ListMove {
	payload := make(primev2messages.ListMoves, len(*moves))

	for i, m := range *moves {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListMove(&copyOfM)
	}
	return payload
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *primev2messages.Error {
	payload := primev2messages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}
