package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func MakeTransportationAccountingCode(db *pop.Connection, assertions Assertions) models.TransportationAccountingCode {
	transportationAccountingCode := models.TransportationAccountingCode{
		ID:        uuid.UUID{000000},
		TAC:       "EO1",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now().Add(-72 * time.Hour),
	}
	mergeModels(&transportationAccountingCode, assertions.TransportationAccountingCode)
	mustCreate(db, &transportationAccountingCode, assertions.Stub)

	return transportationAccountingCode
}

func MakeDefaultTranportationAccountingCode(db *pop.Connection) models.TransportationAccountingCode {
	return MakeTransportationAccountingCode(db, Assertions{})
}
