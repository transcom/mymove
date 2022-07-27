package movingexpense

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type movingExpenseUpdater struct {
	checks []movingExpenseValidator
}

func NewMovingExpenseUpdater() services.MovingExpenseUpdater {
	return &movingExpenseUpdater{
		checks: updateChecks(),
	}
}

func (f *movingExpenseUpdater) UpdateMovingExpense(appCtx appcontext.AppContext, movingExpense models.MovingExpense, eTag string) (*models.MovingExpense, error) {
	return nil, nil
}
