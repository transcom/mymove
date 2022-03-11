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
	query := `WITH moves AS (
		SELECT
			moves.*
		FROM
			moves
		WHERE locator = $1
	),
	changed_addresses AS (
		SELECT DISTINCT
			unnest(ARRAY [old_data->>'pickup_address_id',
										   old_data->>'destination_address_id']) AS ID
		FROM
			audit_history
			JOIN moves ON audit_history.table_name = 'mto_shipments'
				AND audit_history.old_data ->> 'move_id' = moves.id::text
	),
	shipments AS (
		SELECT
			mto_shipments.*
		FROM
			mto_shipments
	),
	shipment_logs AS (
		SELECT
			audit_history.*,
			'shipments' AS context
		FROM
			audit_history
		JOIN shipments ON shipments.id = audit_history.object_id
			AND audit_history. "table_name" = 'mto_shipments'
	),
	move_logs AS (
		SELECT
			audit_history.*,
			'moves' AS context
		FROM
			audit_history
		JOIN moves ON audit_history.table_name = 'moves'
			AND audit_history.object_id = moves.id
	),
	pickup_address_logs AS (
		SELECT
			audit_history.*,
			'pickup_address' AS context
		FROM
			audit_history
		JOIN shipments ON shipments.pickup_address_id = audit_history.object_id
			AND audit_history. "table_name" = 'addresses'
	),
	destination_address_logs AS (
		SELECT
			audit_history.*,
			'destination_address' AS context
		FROM
			audit_history
		JOIN shipments ON shipments.destination_address_id = audit_history.object_id
			AND audit_history. "table_name" = 'addresses'
	),
	combined_logs AS (
		SELECT
			*
		FROM
			pickup_address_logs
		UNION ALL
		SELECT
			*
		FROM
			destination_address_logs
		UNION ALL
		SELECT
			*
		FROM
			shipment_logs
		UNION ALL
		SELECT
			*
		FROM
			move_logs
	)
	SELECT DISTINCT
		combined_logs.*,
		office_users.first_name AS session_user_first_name,
		office_users.last_name AS session_user_last_name,
		office_users.email AS session_user_email,
		office_users.telephone AS session_user_telephone
	FROM
		combined_logs
		LEFT JOIN users_roles ON session_userid = users_roles.user_id
		LEFT JOIN roles ON users_roles.role_id = roles.id
			AND(roles.role_type = 'transportation_ordering_officer'
				OR roles.role_type = 'transportation_invoicing_officer'
				OR roles.role_type = 'ppm_office_users'
				OR role_type = 'services_counselor'
				OR role_type = 'contracting_officer')
		LEFT JOIN office_users ON office_users.user_id = session_userid;`
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
