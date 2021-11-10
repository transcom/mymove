package sequence

import "github.com/transcom/mymove/pkg/appcontext"

// Sequencer provides an interface for generating sequence numbers.
type Sequencer interface {
	NextVal(appCtx appcontext.AppContext) (int64, error)
	SetVal(appCtx appcontext.AppContext, val int64) error
}
