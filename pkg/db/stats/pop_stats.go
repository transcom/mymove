package stats

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/gobuffalo/pop/v5"
	"github.com/jmoiron/sqlx"
)

func DB(c *pop.Connection) (*sqlx.DB, error) {
	// *sigh* pop does not expose DBStats, so use reflection to get
	// access

	// the store has *sqlx.DB as the first field
	dbi := reflect.ValueOf(c.Store).Elem().Field(0).Interface()
	if db, ok := dbi.(*sqlx.DB); ok {
		return db, nil
	}
	return nil, errors.New("Cannot get db field from pop.Connection")
}

// DBStats returns the sql.DBStats for the configured pop connection
func DBStats(c *pop.Connection) (sql.DBStats, error) {
	db, err := DB(c)
	if err != nil {
		return sql.DBStats{}, err
	}
	return db.DB.Stats(), nil
}
