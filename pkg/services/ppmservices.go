package services

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/rateengine"

	"github.com/transcom/mymove/pkg/models"
)

// EstimateCalculator is the exported interface for calculating the PPM estimate
//go:generate mockery --name EstimateCalculator --disable-version-string
type EstimateCalculator interface {
	CalculateEstimates(ppm *models.PersonallyProcuredMove, moveID uuid.UUID, logger *zap.Logger) (int64, rateengine.CostComputation, error)
}
