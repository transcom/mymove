package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
)

// SetInvoiceSequenceNumber sets the invoice sequence number for a given SCAC/year.
func SetInvoiceSequenceNumber(db *pop.Connection, scac string, year int, sequenceNumber int) error {
	if len(scac) == 0 {
		return errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	sql := `INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
			VALUES ($1, $2, $3)
		ON CONFLICT (standard_carrier_alpha_code, year)
		DO
			UPDATE
				SET sequence_number = $3
				WHERE trackers.standard_carrier_alpha_code = $1 AND trackers.year = $2
	`

	return db.RawQuery(sql, scac, year, sequenceNumber).Exec()
}
