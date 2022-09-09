package movehistory

import (
	"database/sql"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type moveHistoryFetcher struct {
}

// NewMoveHistoryFetcher creates a new MoveHistoryFetcher service
func NewMoveHistoryFetcher() services.MoveHistoryFetcher {
	return &moveHistoryFetcher{}
}

//FetchMoveHistory retrieves a Move's history if it is visible for a given locator
func (f moveHistoryFetcher) FetchMoveHistory(appCtx appcontext.AppContext, params *services.FetchMoveHistoryParams) (*models.MoveHistory, int64, error) {
	// dot, _ := dotsql.LoadFromFile("pkg/sql/move_history_fetcher.sql")
	// queries := dot.QueryMap()
	// rawQuery := queries["move_history_fetcher"]

	rawQuery, queryErr := query.GetQueryString("move_history_fetcher")
	if queryErr != nil {
		return &models.MoveHistory{}, 0, queryErr
	}

	audits := &models.AuditHistories{}
	locator := params.Locator
	if params.Page == nil {
		params.Page = swag.Int64(1)
	}
	if params.PerPage == nil {
		params.PerPage = swag.Int64(20)
	}

	query := appCtx.DB().RawQuery(rawQuery, locator).Paginate(int(*params.Page), int(*params.PerPage))
	err := query.All(audits)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.MoveHistory{}, 0, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.MoveHistory{}, 0, apperror.NewQueryError("AuditHistory", err, "")
		}
	}

	var move models.Move
	err = appCtx.DB().Q().Where("locator = $1", locator).First(&move)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.MoveHistory{}, 0, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.MoveHistory{}, 0, apperror.NewQueryError("Move", err, "")
		}
	}

	moveHistory := models.MoveHistory{
		ID:             move.ID,
		Locator:        move.Locator,
		ReferenceID:    move.ReferenceID,
		AuditHistories: *audits,
	}

	return &moveHistory, int64(query.Paginator.TotalEntriesSize), nil
}
