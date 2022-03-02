package move

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveFetcher struct {
}

// NewMoveFetcher creates a new moveFetcher service
func NewMoveFetcher() services.MoveFetcher {
	return &moveFetcher{}
}

//FetchMove retrieves a Move if it is visible for a given locator
func (f moveFetcher) FetchMove(appCtx appcontext.AppContext, locator string, searchParams *services.MoveFetcherParams) (*models.Move, error) {
	move := &models.Move{}
	query := appCtx.DB().Where("locator = $1", locator)

	if searchParams == nil || !searchParams.IncludeHidden {
		query.Where("show = TRUE")
	}

	err := query.First(move)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.Move{}, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.Move{}, apperror.NewQueryError("Move", err, "")
		}
	}

	return move, nil
}
