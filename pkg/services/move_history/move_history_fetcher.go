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
		WHERE
			locator = $1
	),
	shipments AS (
		SELECT
			mto_shipments.*
		FROM
			mto_shipments
		WHERE
			move_id = (
				SELECT
					moves.id
				FROM
					moves)
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
	WHERE
		orders.id = (
			SELECT
				moves.orders_id
			FROM
				moves)
	),
	-- Context is null if empty, {}, object
    -- Joining the jsonb changed_data for every record to surface duty location ids.
    -- Left join duty_locations since we don't expect origin/new duty locations to change every time.
    -- Convert changed_data.origin_duty_location_id and changed_data.new_duty_location_id to UUID type to take advantage of indexing.
	orders_logs AS (
		SELECT
			audit_history.*,
			NULLIF(
				jsonb_agg(jsonb_strip_nulls(
					jsonb_build_object('origin_duty_location_name', old_duty.name, 'new_duty_location_name', new_duty.name)
				))::TEXT, '[{}]'::TEXT
			) AS context,
 			NULL AS context_id
		FROM
			audit_history
		JOIN orders ON orders.id = audit_history.object_id
			AND audit_history."table_name" = 'orders'
		JOIN jsonb_to_record(audit_history.changed_data) as c(origin_duty_location_id TEXT, new_duty_location_id TEXT) on TRUE
		LEFT JOIN duty_locations AS old_duty on uuid(c.origin_duty_location_id) = old_duty.id
		LEFT JOIN duty_locations AS new_duty on uuid(c.new_duty_location_id) = new_duty.id
		GROUP BY audit_history.id
	),
	service_items AS (
		SELECT
			mto_service_items.id,
			json_agg(json_build_object('name',
					re_services.name,
					'shipment_type',
					mto_shipments.shipment_type))::TEXT AS context
		FROM
			mto_service_items
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		JOIN moves ON moves.id = mto_service_items.move_id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
    GROUP BY
			mto_service_items.id
	),
	service_item_logs AS (
		SELECT
			audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
			JOIN service_items ON service_items.id = audit_history.object_id
				AND audit_history. "table_name" = 'mto_service_items'
	),
	pickup_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'addressType', 'pickupAddress'::TEXT,
					'shipment_type', shipments.shipment_type
				)
				)::TEXT AS context,
			shipments.id::text AS context_id
		FROM
			audit_history
		JOIN shipments ON shipments.pickup_address_id = audit_history.object_id
			AND audit_history. "table_name" = 'addresses'
		GROUP BY
			shipments.id, audit_history.id
	),
	destination_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'addressType', 'destinationAddress'::TEXT,
					'shipment_type', shipments.shipment_type
				)
			)::TEXT AS context,
			shipments.id::text AS context_id
		FROM
			audit_history
		JOIN shipments ON shipments.destination_address_id = audit_history.object_id
			AND audit_history. "table_name" = 'addresses'
		GROUP BY
			shipments.id, audit_history.id
	),
	entitlements AS (
		SELECT
			entitlements.*
		FROM
			entitlements
	WHERE
		entitlements.id = (
			SELECT
				entitlement_id
			FROM
				orders)
	),
	entitlements_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN entitlements ON entitlements.id = audit_history.object_id
				AND audit_history. "table_name" = 'entitlements'
	),
	payment_requests AS (
		SELECT
			json_agg(json_build_object('name',
					re_services.name,
					'price',
					payment_service_items.price_cents::TEXT,
					'status',
					payment_service_items.status))::TEXT AS context,
			payment_requests.id AS id
		FROM
			payment_requests
		JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
		JOIN mto_service_items ON mto_service_items.id = mto_service_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
	WHERE
		payment_requests.move_id = (
			SELECT
				moves.id
			FROM
				moves)
		GROUP BY
			payment_requests.id
	),
	payment_requests_logs AS (
		SELECT DISTINCT
			audit_history.*,
			context AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN payment_requests ON payment_requests.id = audit_history.object_id
				AND audit_history. "table_name" = 'payment_requests'
	),
	agents AS (
		SELECT
			mto_agents.id,
			json_agg(json_build_object(
				'shipment_type',
				shipments.shipment_type))::TEXT AS context
		FROM
			mto_agents
			JOIN shipments ON mto_agents.mto_shipment_id = shipments.id
		GROUP BY
			mto_agents.id
	),
	agents_logs AS (
		SELECT
			audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
			JOIN agents ON agents.id = audit_history.object_id
				AND audit_history."table_name" = 'mto_agents'
	),
	reweighs AS (
		SELECT
			reweighs.id,
			json_agg(json_build_object(
					'shipment_type',
					shipments.shipment_type))::TEXT AS context
		FROM
			reweighs
			JOIN shipments ON reweighs.shipment_id = shipments.id
		GROUP BY
			reweighs.id
	),
	reweigh_logs as (
		SELECT audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
		JOIN reweighs ON reweighs.id = audit_history.object_id
			AND audit_history."table_name" = 'reweighs'
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
			entitlements_logs
		UNION ALL
		SELECT
			*
		FROM
			reweigh_logs
		UNION ALL
		SELECT
			*
		FROM
			orders_logs
		UNION ALL
		SELECT
			*
		FROM
			agents_logs
		UNION ALL
		SELECT
			*
		FROM
			payment_requests_logs
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
