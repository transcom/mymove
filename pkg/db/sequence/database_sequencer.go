package sequence

import (
	"github.com/transcom/mymove/pkg/appcontext"
)

// Using PostgreSQL Sequences
// https://www.postgresql.org/docs/10/sql-createsequence.html

// databaseSequencer represents a database-backed sequence with the given connection and sequence name
type databaseSequencer struct {
	sequenceName string
}

// NextVal returns the next value of the given sequence
func (ds databaseSequencer) NextVal(appCtx appcontext.AppContext) (int64, error) {
	var nextVal int64
	err := appCtx.DB().RawQuery("SELECT nextval($1);", ds.sequenceName).First(&nextVal)
	return nextVal, err
}

// SetVal sets the current value of a sequence
func (ds databaseSequencer) SetVal(appCtx appcontext.AppContext, val int64) error {
	err := appCtx.DB().RawQuery("SELECT setval($1, $2)", ds.sequenceName, val).Exec()
	return err
}

// NewDatabaseSequencer is a factory for creating a new database-backed sequencer
func NewDatabaseSequencer(sequenceName string) Sequencer {
	return &databaseSequencer{sequenceName}
}
