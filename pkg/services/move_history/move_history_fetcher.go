package movehistory

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveHistoryFetcher struct {
}

// NewMoveHistoryFetcher creates a new MoveHistoryFetcher service
func NewMoveHistoryFetcher() services.MoveHistoryFetcher {
	return &moveHistoryFetcher{}
}

//FetchMoveHistory retrieves a Move's history if it is visible for a given locator
func (f moveHistoryFetcher) FetchMoveHistory(appCtx appcontext.AppContext, locator string) (*models.MoveHistory, error) {
	moveHistory := models.MoveHistory{}
	audits := &models.AuditHistories{}
	query := `WITH moves AS
  (SELECT moves.*
    FROM moves
   WHERE id = '$1'),
changed_addresses AS
  (SELECT DISTINCT unnest(array[changed_data->>'pickup_address_id',
								changed_data->>'destination_address_id']) AS ID
     FROM audit_history
     JOIN moves ON audit_history.table_name = 'mto_shipments'
      AND audit_history.changed_data->>'move_id' = moves.id::text)
-- the move history
SELECT audit_history.*
  FROM audit_history
  JOIN moves ON audit_history.table_name = 'moves'
   AND audit_history.object_id = moves.id
UNION ALL
-- history of the associated order
SELECT audit_history.*
  FROM audit_history
  JOIN moves ON audit_history.table_name = 'orders'
   AND audit_history.object_id = moves.orders_id
UNION ALL
-- history of everything that has a FK to the move
SELECT audit_history.*
  FROM audit_history
  JOIN moves ON changed_data->>'move_id' = moves.id::text
UNION ALL
-- history of related addresses
SELECT audit_history.*
  FROM audit_history
  JOIN changed_addresses ON
       audit_history.table_name = 'addresses'
   AND (audit_history.object_id = changed_addresses.id::uuid);
`
	err := appCtx.DB().RawQuery(query, locator).All(audits)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.MoveHistory{}, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.MoveHistory{}, apperror.NewQueryError("AuditHistory", err, "")
		}
	}

	// TODO build up MoveHistory and return to caller

	return &moveHistory, nil
}
