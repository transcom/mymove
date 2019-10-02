package ghcrateengine

import "github.com/gobuffalo/pop"

// RateEngine encapsulates the TSP rate engine process
type GHCRateEngine struct {
	db     *pop.Connection
	logger Logger
}

func NewGHCRateEngine(db *pop.Connection, logger Logger) GHCRateEngine {
	return GHCRateEngine{
		db:     db,
		logger: logger,
	}
}
