package movehistory

import (
	"database/sql"

	"github.com/go-openapi/swag"
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
func (f moveHistoryFetcher) FetchMoveHistory(appCtx appcontext.AppContext, params *services.FetchMoveHistoryParams) (*models.MoveHistory, int64, error) {
	rawQuery := `WITH moves AS (
		SELECT
			moves.*
		FROM
			moves
		WHERE locator = $1
	),
	shipments AS (
		SELECT
			mto_shipments.*
		FROM
			mto_shipments
		WHERE move_id = (SELECT moves.id FROM moves)
	),
	shipment_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN shipments ON shipments.id = audit_history.object_id
				AND audit_history."table_name" = 'mto_shipments'
	),
	move_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN moves ON audit_history.table_name = 'moves'
			AND audit_history.object_id = moves.id
	),
	orders AS (
		SELECT
			orders.*
		FROM
			orders
		JOIN moves ON moves.orders_id = orders.id
		WHERE orders.id = (SELECT moves.orders_id FROM moves)
	),
	orders_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN orders ON orders.id = audit_history.object_id
			AND audit_history."table_name" = 'orders'
	),
	service_items AS (
		SELECT
			mto_service_items.*, re_services.name, mto_shipments.shipment_type
		FROM
			mto_service_items
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipment_id
		WHERE mto_shipments.move_id = (SELECT moves.id FROM moves)
	),
	service_item_logs AS (
		SELECT
			audit_history.*,
			json_build_object('name', service_items.name, 'shipment_type', service_items.shipment_type)::TEXT AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN service_items ON service_items.id = audit_history.object_id
			AND audit_history."table_name" = 'mto_service_items'
	),
	pickup_address_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			shipments.id::text AS context_id
		FROM
			audit_history
		JOIN shipments ON shipments.pickup_address_id = audit_history.object_id
			AND audit_history. "table_name" = 'addresses'
	),
	destination_address_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			shipments.id::text AS context_id
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
			service_item_logs
		UNION ALL
		SELECT
			*
		FROM
			shipment_logs
		UNION ALL
		SELECT
			*
		FROM
			orders_logs
		UNION ALL
		SELECT
			*
		FROM
			move_logs
	) SELECT DISTINCT
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
		LEFT JOIN office_users ON office_users.user_id = session_userid
	ORDER BY
		action_tstamp_tx DESC`

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
