package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeLineOfAccounting(db *pop.Connection, assertions Assertions) models.LineOfAccounting {
	loa := models.LineOfAccounting{
		ID:        uuid.UUID{000000},
		UpdatedAt: time.Now(),
		CreatedAt: time.Now().Add(-72 * time.Hour),
	}

	mergeModels(&loa, assertions.LineOfAccounting)
	mustCreate(db, &loa, assertions.Stub)

	return loa
}

func MakeDefualtLineOfAccounting(db *pop.Connection) models.LineOfAccounting {
	return MakeLineOfAccounting(db, Assertions{})
}
