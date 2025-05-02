package movehistory

import (
	"database/sql"
	"strings"

	"github.com/gobuffalo/pop/v6"
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

// FetchMoveHistory retrieves a Move's history if it is visible for a given locator
func (f moveHistoryFetcher) FetchMoveHistory(appCtx appcontext.AppContext, params *services.FetchMoveHistoryParams, useDatabaseProcInstead bool) (*models.MoveHistory, int64, error) {
	var rawQuery string
	if useDatabaseProcInstead {
		// casting types to match function declared params
		rawQuery = `SELECT * FROM fetch_move_history(
					$1::text,
					$2::integer,
					$3::integer,
					$4::text,
					$5::text)`
	} else {
		var qerr error
		rawQuery, qerr = query.GetSQLQueryByName("move_history_fetcher")
		if qerr != nil {
			return &models.MoveHistory{}, 0, apperror.NewQueryError("AuditHistory", qerr, "")
		}
	}

	locator := params.Locator
	if params.Page == nil {
		params.Page = models.Int64Pointer(1)
	}
	if params.PerPage == nil {
		params.PerPage = models.Int64Pointer(20)
	}

	audits := models.AuditHistories{}
	var err error
	var query *pop.Query
	var totalCount int64
	if useDatabaseProcInstead {
		query = appCtx.DB().RawQuery(
			rawQuery,
			params.Locator,
			int(*params.Page),
			int(*params.PerPage),
			nil,
			nil,
		)
	} else {
		query = appCtx.DB().RawQuery(rawQuery, locator).Paginate(int(*params.Page), int(*params.PerPage))
	}

	err = query.All(&audits)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveHistory{}, 0, apperror.NewNotFoundError(uuid.Nil, "move locator "+params.Locator)
		default:
			// Catch the proc case
			if strings.Contains(err.Error(), "Move record not found for") {
				return &models.MoveHistory{}, 0, apperror.NewNotFoundError(uuid.Nil, "move locator "+params.Locator)
			}
			return &models.MoveHistory{}, 0, apperror.NewQueryError("AuditHistory", err, err.Error())
		}
	}

	// bypassing the paginator when using the db func does not give us the count back
	// we can use the temp table that is created from the db func to grab the count instead
	// else we will use the non-db paginator entry size
	if useDatabaseProcInstead {
		countQuery := "SELECT COUNT(*) FROM audit_hist_temp"
		err = appCtx.DB().RawQuery(countQuery).First(&totalCount)
		if err != nil {
			return &models.MoveHistory{}, 0, apperror.NewQueryError("AuditHistory Count", err, "")
		}
	} else {
		if query != nil && query.Paginator != nil {
			totalCount = int64(query.Paginator.TotalEntriesSize)
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
		AuditHistories: audits,
	}

	return &moveHistory, int64(totalCount), nil
}
