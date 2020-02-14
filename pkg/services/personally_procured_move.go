package services

import (
"github.com/gofrs/uuid"

"github.com/transcom/mymove/pkg/models"
)

// PersonallyProcuredMoveFetcher is the service object interface for FetchPersonallyProcuredMove
//go:generate mockery -name MoveTaskOrderFetcher
type PersonallyProcuredMoveFetcher interface {
	FetchPersonallyProcuredMove(personallyProcuredMoveID uuid.UUID) (*models.PersonallyProcuredMove, error)
}
