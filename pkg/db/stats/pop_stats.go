package stats

import (
	"database/sql"
	"errors"

	"github.com/transcom/mymove/pkg/appcontext"
)

type statser interface {
	Stats() sql.DBStats
}

// DBStats returns the sql.DBStats for the configured pop connection
func DBStats(appCtx appcontext.AppContext) (sql.DBStats, error) {
	if dbWithStats, ok := appCtx.DB().Store.(statser); ok {
		return dbWithStats.Stats(), nil
	}
	return sql.DBStats{}, errors.New("Cannot get stats from pop.Connection")
}
