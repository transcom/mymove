package services

import "github.com/transcom/mymove/pkg/models"

// AccessCodeValidator is the service object interface for ValidateAccessCode
type AccessCodeValidator interface {
	ValidateAccessCode(code string, moveType models.SelectedMoveType) (*models.AccessCode, bool, error)
}
