package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// TransportationAccountingCode model struct that represents transportation accounting codes
// TODO: Update this model and internal use to reflect incoming TransportationAccountingCode model updates.
// Don't forget to update the MakeDefaultTransportationAccountingCode function inside of the testdatagen package.
type TransportationAccountingCode struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TAC       string    `json:"tac" db:"tac"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (t TransportationAccountingCode) TableName() string {
	return "transportation_accounting_codes"
}

func MapTransportationAccountingCodeFileRecordToInternalStruct(tacFileRecord TransportationAccountingCodeTextFileRecord) TransportationAccountingCode {
	return TransportationAccountingCode{
		TAC: tacFileRecord.TRNSPRTN_ACNT_CD,
	}
}
