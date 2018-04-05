package rateengine

import (
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
	"time"
)

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db     *pop.Connection
	logger *zap.Logger
	date   time.Time
}

func (re *RateEngine) determineCWT(weight int) (cwt int) {
	return weight / 100
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger, date time.Time) *RateEngine {
	return &RateEngine{db: db, logger: logger, date: date}
}
