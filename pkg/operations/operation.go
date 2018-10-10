package operations

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

// Operation is a struct on the service object layer to handle move cancelations
type Operation struct {
	DB      *pop.Connection
	Logger  *zap.Logger
	Session *auth.Session
	Verrs   *validate.Errors
	Err     error
}

func (op *Operation) hadErrors(verrs *validate.Errors, saveErr error) bool {
	if saveErr != nil {
		op.Err = errors.Wrap(saveErr, "error saving model")
		return true
	}
	if verrs.HasAny() {
		op.Verrs = verrs
		op.Err = errors.New("Model validation failure")
		return true
	}
	return false
}
