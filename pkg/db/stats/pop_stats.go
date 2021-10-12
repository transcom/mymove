package stats

import (
	"database/sql"
	"errors"

	"github.com/gobuffalo/pop/v5"
)

type statser interface {
	Stats() sql.DBStats
}

// DBStats returns the sql.DBStats for the configured pop connection
func DBStats(c *pop.Connection) (sql.DBStats, error) {
	if dbWithStats, ok := c.Store.(statser); ok {
		return dbWithStats.Stats(), nil
	}
	return sql.DBStats{}, errors.New("Cannot get stats from pop.Connection")
}
