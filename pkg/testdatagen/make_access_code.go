package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAccessCode creates a single AccessCode
func MakeAccessCode(db *pop.Connection, assertions Assertions) models.AccessCode {
	defaultMoveType := models.SelectedMoveTypePPM
	selectedMoveType := assertions.AccessCode.MoveType
	if selectedMoveType == nil {
		selectedMoveType = &defaultMoveType
	}

	accessCode := models.AccessCode{
		Code:      models.GenerateLocator(),
		MoveType:  selectedMoveType,
		CreatedAt: time.Now(),
	}

	mergeModels(&accessCode, assertions.AccessCode)

	mustCreate(db, &accessCode)

	return accessCode
}
