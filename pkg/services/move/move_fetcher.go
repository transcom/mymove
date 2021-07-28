package move

import (
	"database/sql"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveFetcher struct {
	db *pop.Connection
}

// NewMoveFetcher creates a new moveFetcher service
func NewMoveFetcher(db *pop.Connection) services.MoveFetcher {
	return &moveFetcher{db}
}

//FetchOrder retrieves a Move if it is visible for a given locator
func (f moveFetcher) FetchMove(locator string, searchParams *services.MoveFetcherParams) (*models.Move, error) {
	move := &models.Move{}
	query := f.db.Where("locator = $1", locator)

	if searchParams == nil || !searchParams.IncludeHidden {
		query.Where("show = TRUE")
	}

	err := query.First(move)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.Move{}, services.NotFoundError{}
		default:
			return &models.Move{}, err
		}
	}

	return move, nil
}
