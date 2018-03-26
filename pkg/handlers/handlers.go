package handlers

import (
	"io"

	"github.com/markbates/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db     *pop.Connection
	logger *zap.Logger
}

type fileStorer interface {
	Store(string, io.ReadSeeker, string) (*storage.StoreResult, error)
	Key(...string) string
	PresignedURL(string) (string, error)
}

// NewHandlerContext returns a new HandlerContext with its private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

type FileHandlerContext struct {
	*HandlerContext
	storage fileStorer
}

func NewFileHandlerContext(db *pop.Connection, logger *zap.Logger, storer fileStorer) FileHandlerContext {
	hc := NewHandlerContext(db, logger)
	return FileHandlerContext{
		HandlerContext: &hc,
		storage:        storer,
	}
}
