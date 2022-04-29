package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTransportationAccountingCode creates a single transportation access code
func MakeTransportationAccountingCode(db *pop.Connection, assertions Assertions) models.TransportationAccountingCode {
	transportationAccountingCode := models.TransportationAccountingCode{
		TAC: "E01A",
	}

	mergeModels(&transportationAccountingCode, assertions.TransportationAccountingCode)

	mustCreate(db, &transportationAccountingCode, assertions.Stub)

	return transportationAccountingCode
}

// MakeDefaultTransportationAccountingCode makes a TransportationAccountingCode with default values
func MakeDefaultTransportationAccountingCode(db *pop.Connection) models.TransportationAccountingCode {
	return MakeTransportationAccountingCode(db, Assertions{})
}
