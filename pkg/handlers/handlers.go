package handlers

import (
	"github.com/markbates/pop"
	"go.uber.org/zap"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewHandlerContext returns a new HandlerContext with its private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}
