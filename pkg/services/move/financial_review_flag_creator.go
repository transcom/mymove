package move

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type financialReviewFlagCreator struct {
}

func NewFinancialReviewFlagCreator() services.MoveFinancialReviewFlagCreator {
	return &financialReviewFlagCreator{}
}

func (f financialReviewFlagCreator) CreateFinancialReviewFlag(appCtx appcontext.AppContext, moveID uuid.UUID, remarks string) (*models.Move, error) {
	return nil, fmt.Errorf("not implemented")
}
