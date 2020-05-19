package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// EstimateCalculator is the exported interface for calculating the PPM estimate
//go:generate mockery -name EstimateCalculator
type EstimateCalculator interface {
	CalculateEstimate(ppm *models.PersonallyProcuredMove, moveID uuid.UUID) error
}
