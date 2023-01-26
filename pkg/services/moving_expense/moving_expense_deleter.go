package movingexpense

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseDeleter struct {
}

func NewMovingExpenseDeleter() services.MovingExpenseDeleter {
	return &movingExpenseDeleter{}
}

func (d *movingExpenseDeleter) DeleteMovingExpense(appCtx appcontext.AppContext, movingExpenseID uuid.UUID) error {
	return nil
}
