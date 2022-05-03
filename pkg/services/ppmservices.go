package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/models"
)

// DEPRECATED
// EstimateCalculator is the exported interface for calculating the legacy Personally Procured Move estimate
//go:generate mockery --name EstimateCalculator --disable-version-string
type EstimateCalculator interface {
	CalculateEstimates(appCtx appcontext.AppContext, ppm *models.PersonallyProcuredMove, moveID uuid.UUID) (int64, rateengine.CostComputation, error)
}
