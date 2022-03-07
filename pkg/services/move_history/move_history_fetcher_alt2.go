package movehistory

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

//FetchMoveHistoryAlt2 retrieves a Move's history if it is visible for a given locator
func (f moveHistoryFetcher) FetchMoveHistoryAlt2(appCtx appcontext.AppContext, locator string) (*models.MoveHistory, error) {
	/*
		This query adds a `shipment_id` column to all rows and `context`.
		For rows where there are addresses, the shipment_id will be filled in.
		The column for `shipment_id` would be created during the query, not during the save to the table.


		This file is experimental and trying to see how to pull in more descriptive data about the audit history.
		There is a question out in Slack about how to get some more results in:
		* https://ustcdp3.slack.com/archives/CSGDM3NUW/p1644865434835809
		* https://trussworks.slack.com/archives/C015D41GPM2/p1644865640412049

		THIS FILE IS NOT BEING RUN IN THE POC, you can run it locally if you are curious about it, but there is no
		Go datatype/structure to handle the return for this, really have just been playing with it in the console.
	*/

	query := `WITH moves AS
         (SELECT moves.*
          FROM moves
          WHERE locator = $1),
     shipments AS
         (SELECT mto_shipments.*
          FROM mto_shipments
          JOIN moves ON mto_shipments.move_id = moves.id),
    shipment_addresses AS
    (
        (SELECT shipments.*, 'secondary_pickup_address_id' AS address_type
         FROM shipments
                  JOIN addresses on shipments.secondary_pickup_address_id = addresses.id)
        UNION
        (SELECT shipments.*, 'secondary_delivery_address_id' AS address_type
         FROM shipments
                  JOIN addresses on shipments.secondary_delivery_address_id = addresses.id)
        UNION
        (SELECT shipments.*, 'pickup_address_id' AS address_type
         FROM shipments
                  JOIN addresses on shipments.pickup_address_id = addresses.id)
        UNION
        (SELECT shipments.*, 'destination_address_id' AS address_type
         FROM shipments
                  JOIN addresses on shipments.destination_address_id = addresses.id)
    )
-- the move history
SELECT audit_history.*, null::uuid AS shipment_id, null::text AS context
FROM audit_history
    JOIN moves ON audit_history.table_name = 'moves'
    AND audit_history.object_id = moves.id
UNION ALL
-- history of the associated order
SELECT audit_history.*, null::uuid AS shipment_id, null::text AS context
FROM audit_history
    JOIN moves ON audit_history.table_name = 'orders'
    AND audit_history.object_id = moves.orders_id
UNION ALL
-- history of everything that has a FK to the move
SELECT audit_history.*, null::uuid AS shipment_id, null::text AS context
FROM audit_history
         JOIN moves ON changed_data->>'move_id' = moves.id::text
UNION ALL
-- history of related addresses
SELECT audit_history.*, shipment_addresses.id::uuid AS shipment_id, shipment_addresses.address_type as context
FROM audit_history
     JOIN shipment_addresses ON audit_history.table_name = 'addresses'
        AND audit_history.object_id IN (
                shipment_addresses.destination_address_id,
                shipment_addresses.pickup_address_id,
                shipment_addresses.secondary_delivery_address_id,
                shipment_addresses.secondary_pickup_address_id);
`
	audits := &models.AuditHistories{}
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

	var move models.Move
	err = appCtx.DB().Q().Where("locator = $1", locator).First(&move)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// Not found error expects an id but we're querying by locator
			return &models.MoveHistory{}, apperror.NewNotFoundError(uuid.Nil, "move locator "+locator)
		default:
			return &models.MoveHistory{}, apperror.NewQueryError("Move", err, "")
		}
	}

	moveHistory := models.MoveHistory{
		ID:             move.ID,
		Locator:        move.Locator,
		ReferenceID:    move.ReferenceID,
		AuditHistories: *audits,
	}

	return &moveHistory, nil
}
