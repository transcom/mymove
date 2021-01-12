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

//FetchMoveOrder retrieves a Move for a given locator
func (f moveFetcher) FetchMove(locator string) (*models.Move, error) {
	move := &models.Move{}
	if err := f.db.Where("locator = $1", locator).First(move); err != nil {
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