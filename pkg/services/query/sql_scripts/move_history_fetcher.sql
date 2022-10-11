WITH moves AS (
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
	service_item_customer_contacts AS (
		SELECT
			mto_service_item_customer_contacts.*
		FROM
			mto_service_item_customer_contacts
		JOIN mto_service_items on mto_service_items.id = mto_service_item_customer_contacts.mto_service_item_id
		JOIN moves ON moves.id = mto_service_items.move_id
	),
	service_item_customer_contacts_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN service_item_customer_contacts ON service_item_customer_contacts.id = audit_history.object_id
			AND audit_history."table_name" = 'mto_service_item_customer_contacts'
	),
	service_item_dimensions AS (
		SELECT
			mto_service_item_dimensions.*,
			json_agg(json_build_object('name',
					re_services.name,
					'shipment_type',
					mto_shipments.shipment_type))::TEXT AS context
		FROM
			mto_service_item_dimensions
		JOIN mto_service_items on mto_service_items.id = mto_service_item_dimensions.mto_service_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON moves.id = mto_service_items.move_id
		GROUP BY
			mto_service_item_dimensions.id
	),
	service_item_dimensions_logs AS  (
		SELECT
			audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
		JOIN service_item_dimensions ON service_item_dimensions.id = audit_history.object_id
			AND audit_history."table_name" = 'mto_service_item_dimensions'
	),
	pickup_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'address_type', 'pickupAddress'::TEXT,
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
	secondary_pickup_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'address_type', 'secondaryPickupAddress'::TEXT,
					'shipment_type', shipments.shipment_type
				)
			)::TEXT AS context,
			shipments.id::text AS context_id
		FROM
			audit_history
				JOIN shipments ON shipments.secondary_pickup_address_id = audit_history.object_id
				AND audit_history. "table_name" = 'addresses'
		GROUP BY
			shipments.id, audit_history.id
	),
	destination_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'address_type', 'destinationAddress'::TEXT,
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
	secondary_destination_address_logs AS (
		SELECT
			audit_history.*,
			json_agg(
				json_build_object(
					'address_type', 'secondaryDestinationAddress'::TEXT,
					'shipment_type', shipments.shipment_type
				)
			)::TEXT AS context,
			shipments.id::text AS context_id
		FROM
			audit_history
				JOIN shipments ON shipments.secondary_delivery_address_id = audit_history.object_id
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
					payment_service_items.status,
					'shipment_id',
					mto_shipments.id::TEXT,
					'shipment_type', mto_shipments.shipment_type))::TEXT AS context,
			payment_requests.id AS id,
			payment_requests.move_id,
			payment_requests.payment_request_number
		FROM
			payment_requests
		JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
		JOIN mto_service_items ON mto_service_items.id = mto_service_item_id
		LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
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
	proof_of_service_docs AS (
		SELECT
			proof_of_service_docs.*,
			json_agg(json_build_object(
				'payment_request_number',
				payment_requests.payment_request_number::TEXT))::TEXT AS context
		FROM
			proof_of_service_docs
				JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
		GROUP BY proof_of_service_docs.id
	),
	proof_of_service_docs_logs AS (
		SELECT
			audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
			JOIN proof_of_service_docs ON proof_of_service_docs.id = audit_history.object_id
				AND audit_history. "table_name" = 'proof_of_service_docs'
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
			json_agg(json_build_object('shipment_type',
				shipments.shipment_type,
				'payment_request_number',
				payment_requests.payment_request_number))::TEXT AS context
		FROM
			reweighs
			JOIN shipments ON reweighs.shipment_id = shipments.id
			LEFT JOIN payment_requests ON shipments.move_id = payment_requests.move_id
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
	service_members AS (
		SELECT service_members.*
		FROM
			service_members
		WHERE
			service_members.id = (SELECT service_member_id FROM orders)
	),
	service_members_logs as (
		SELECT audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN service_members ON service_members.id = audit_history.object_id
			AND audit_history."table_name" = 'service_members'
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
			secondary_pickup_address_logs
		UNION ALL
		SELECT
			*
		FROM
			destination_address_logs
		UNION ALL
		SELECT
			*
		FROM
			secondary_destination_address_logs
		UNION ALL
		SELECT
			*
		FROM
			service_item_logs
		UNION ALL
		SELECT
			*
		FROM
			service_item_customer_contacts_logs
		UNION ALL
		SELECT
			*
		FROM
			service_item_dimensions_logs
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
			proof_of_service_docs_logs
		UNION ALL
		SELECT
			*
		FROM
			move_logs
		UNION ALL
		SELECT
		 	*
		FROM
			service_members_logs
	) SELECT DISTINCT
		combined_logs.*,
		COALESCE(office_users.first_name, prime_user_first_name) AS session_user_first_name,
		office_users.last_name AS session_user_last_name,
		office_users.email AS session_user_email,
		office_users.telephone AS session_user_telephone
	FROM
		combined_logs
		LEFT JOIN users_roles ON session_userid = users_roles.user_id
		LEFT JOIN roles ON users_roles.role_id = roles.id
		LEFT JOIN office_users ON office_users.user_id = session_userid
		LEFT JOIN (
			SELECT 'Prime' AS prime_user_first_name
			) prime_users ON roles.role_type = 'prime'
	ORDER BY
		action_tstamp_tx DESC
