package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeAccessCode creates a single AccessCode
func MakeAccessCode(db *pop.Connection, assertions Assertions) models.AccessCode {
	accessCode := models.AccessCode{
		Code:      models.GenerateLocator(),
		MoveType:  models.SelectedMoveTypePPM,
		CreatedAt: time.Now(),
	}

	mergeModels(&accessCode, assertions.AccessCode)

	mustCreate(db, &accessCode, assertions.Stub)

	return accessCode
}
