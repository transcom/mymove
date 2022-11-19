WITH move AS (
		SELECT
			moves.*
		FROM
			moves
		WHERE
			moves.locator = $1
	),
	move_shipments AS (
		SELECT
			mto_shipments.*, LEFT(mto_shipments.id::TEXT, 5) AS shipment_id_abbr
		FROM
			mto_shipments
		WHERE
			mto_shipments.move_id = (
				SELECT
					move.id
				FROM
					move)
	),
	shipment_logs AS (
		SELECT
			audit_history.*,
			NULLIF(
				jsonb_agg(jsonb_strip_nulls(
					jsonb_build_object(
						'shipment_type', move_shipments.shipment_type,
						'shipment_id_abbr', move_shipments.shipment_id_abbr
					)
				))::TEXT, '[{}]'::TEXT
			) AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_shipments ON move_shipments.id = audit_history.object_id
		WHERE audit_history.table_name = 'mto_shipments'
		GROUP BY audit_history.id
	),
	move_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN move ON audit_history.object_id = move.id
		WHERE audit_history.table_name = 'moves'
	),
	move_orders AS (
		SELECT
			orders.*
		FROM
			orders
		WHERE orders.id = (SELECT move.orders_id FROM move)
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
					jsonb_build_object(
						'origin_duty_location_name',
						(SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.origin_duty_location_id)),
						'new_duty_location_name',
						(SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.new_duty_location_id))
					)
				))::TEXT, '[{}]'::TEXT
			) AS context,
 			NULL AS context_id
		FROM
			audit_history
		JOIN move_orders ON move_orders.id = audit_history.object_id
		JOIN jsonb_to_record(audit_history.changed_data) as c(origin_duty_location_id TEXT, new_duty_location_id TEXT) on TRUE
		WHERE audit_history.table_name = 'orders'
		GROUP BY audit_history.id
	),
	move_service_items AS (
		SELECT
			mto_service_items.*
		FROM
			mto_service_items
		WHERE
			mto_service_items.move_id = (SELECT move.id FROM move)
	),
	service_item_logs AS (
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr
				)
			)::TEXT AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_service_items ON move_service_items.id = audit_history.object_id
			JOIN re_services ON move_service_items.re_service_id = re_services.id
			LEFT JOIN move_shipments ON move_service_items.mto_shipment_id = move_shipments.id
		WHERE audit_history.table_name = 'mto_service_items'
		GROUP BY
				audit_history.id, move_service_items.id
	),
	service_item_customer_contacts AS (
		SELECT
			mto_service_item_customer_contacts.*,
			jsonb_agg(jsonb_build_object(
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr
				)
			)::TEXT AS context
		FROM
			mto_service_item_customer_contacts
		JOIN move_service_items on move_service_items.id = mto_service_item_customer_contacts.mto_service_item_id
			LEFT JOIN move_shipments ON move_service_items.mto_shipment_id = move_shipments.id
		JOIN move ON move.id = move_service_items.move_id
		GROUP BY mto_service_item_customer_contacts.id
	),
	service_item_customer_contacts_logs AS (
		SELECT
			audit_history.*,
			service_item_customer_contacts.context AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN service_item_customer_contacts ON service_item_customer_contacts.id = audit_history.object_id
			WHERE audit_history.table_name = 'mto_service_item_customer_contacts'
	),
	service_item_dimensions AS (
		SELECT
			mto_service_item_dimensions.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr
				)
			)::TEXT AS context
		FROM
			mto_service_item_dimensions
		JOIN move_service_items on move_service_items.id = mto_service_item_dimensions.mto_service_item_id
		JOIN re_services ON move_service_items.re_service_id = re_services.id
			LEFT JOIN move_shipments ON move_service_items.mto_shipment_id = move_shipments.id
		GROUP BY mto_service_item_dimensions.id
	),
	service_item_dimensions_logs AS  (
		SELECT
			audit_history.*,
			service_item_dimensions.context AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN service_item_dimensions ON service_item_dimensions.id = audit_history.object_id
		WHERE audit_history.table_name = 'mto_service_item_dimensions'

	),
	move_entitlements AS (
		SELECT
			entitlements.*
		FROM
			entitlements
	WHERE
		entitlements.id = (
			SELECT
				entitlement_id
			FROM
				move_orders)
	),
	entitlements_logs AS (
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_entitlements ON move_entitlements.id = audit_history.object_id
		WHERE audit_history.table_name = 'entitlements'
	),
	move_payment_requests AS (
		SELECT
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'price', payment_service_items.price_cents::TEXT,
				'status', payment_service_items.status,
				'shipment_id', move_shipments.id::TEXT,
				'shipment_id_abbr', move_shipments.shipment_id_abbr,
				'shipment_type', move_shipments.shipment_type
				)
			)::TEXT AS context,
			payment_requests.id AS id,
			payment_requests.move_id,
			payment_requests.payment_request_number
		FROM
			payment_requests
			JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
			JOIN move_service_items ON move_service_items.id = payment_service_items.mto_service_item_id
			LEFT JOIN move_shipments ON move_shipments.id = move_service_items.mto_shipment_id
			JOIN re_services ON move_service_items.re_service_id = re_services.id
		WHERE
			payment_requests.move_id = (SELECT move.id FROM move)
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
			JOIN move_payment_requests ON move_payment_requests.id = audit_history.object_id
		WHERE audit_history.table_name = 'payment_requests'
	),
	move_proof_of_service_docs AS (
		SELECT
			proof_of_service_docs.*,
			jsonb_agg(jsonb_build_object(
				'payment_request_number',
				move_payment_requests.payment_request_number::TEXT))::TEXT AS context
		FROM
			proof_of_service_docs
				JOIN move_payment_requests ON proof_of_service_docs.payment_request_id = move_payment_requests.id
		GROUP BY proof_of_service_docs.id
	),
	proof_of_service_docs_logs AS (
		SELECT
			audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_proof_of_service_docs ON move_proof_of_service_docs.id = audit_history.object_id
		WHERE audit_history.table_name = 'proof_of_service_docs'
	),
	agents AS (
		SELECT
			mto_agents.id,
			jsonb_agg(jsonb_build_object(
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr
				)
			)::TEXT AS context
		FROM
			mto_agents
			JOIN move_shipments ON mto_agents.mto_shipment_id = move_shipments.id
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
		WHERE audit_history.table_name = 'mto_agents'
	),
	move_reweighs AS (
		SELECT
			reweighs.id,
			jsonb_agg(jsonb_build_object(
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr,
				'payment_request_number', move_payment_requests.payment_request_number
				)
			)::TEXT AS context
		FROM
			reweighs
			JOIN move_shipments ON reweighs.shipment_id = move_shipments.id
			LEFT JOIN move_payment_requests ON move_shipments.move_id = move_payment_requests.move_id
		GROUP BY
			reweighs.id
	),
	reweigh_logs as (
		SELECT audit_history.*,
			context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_reweighs ON move_reweighs.id = audit_history.object_id
		WHERE audit_history.table_name = 'reweighs'
	),
	move_service_members AS (
		SELECT service_members.*
		FROM
			service_members
		WHERE
			service_members.id = (SELECT move_orders.service_member_id FROM move_orders)
	),
	service_members_logs as (
		SELECT audit_history.*,
				NULLIF(
				jsonb_agg(jsonb_strip_nulls(
					jsonb_build_object(
						'current_duty_location_name',
						(SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.duty_location_id))
					)
				))::TEXT, '[{}]'::TEXT
			) AS context,
			NULL AS context_id
		FROM
			audit_history
			JOIN move_service_members ON move_service_members.id = audit_history.object_id
			JOIN jsonb_to_record(audit_history.changed_data) as c(duty_location_id TEXT) on TRUE
		WHERE audit_history.table_name = 'service_members'
		GROUP BY audit_history.id
	),
	move_addresses (address_id, address_type, shipment_type, shipment_id, service_member_id)  AS (
		SELECT
			audit_history.object_id,
			'destinationAddress',
			move_shipments.shipment_type,
			move_shipments.id::TEXT,
			NULL
		FROM audit_history
			JOIN move_shipments ON move_shipments.destination_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
		UNION
		SELECT
			audit_history.object_id,
			'secondaryDestinationAddress',
			move_shipments.shipment_type,
			move_shipments.id::TEXT,
			NULL
		FROM audit_history
			JOIN move_shipments ON move_shipments.secondary_delivery_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
		UNION
		SELECT
			audit_history.object_id,
			'pickupAddress',
			move_shipments.shipment_type,
			move_shipments.id::TEXT,
			NULL
		FROM audit_history
			JOIN move_shipments ON move_shipments.pickup_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
		UNION
		SELECT
			audit_history.object_id,
			'secondaryPickupAddress',
			move_shipments.shipment_type,
			move_shipments.id::TEXT,
			NULL
		FROM audit_history
			JOIN move_shipments ON move_shipments.secondary_pickup_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
		UNION
		SELECT
			audit_history.object_id,
			'residentialAddress',
			NULL,
			NULL,
			move_service_members.id::TEXT
		FROM audit_history
			JOIN move_service_members ON move_service_members.residential_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
		UNION
		SELECT
			audit_history.object_id,
			'backupMailingAddress',
			NULL,
			NULL,
			move_service_members.id::TEXT
		FROM audit_history
			JOIN move_service_members ON move_service_members.backup_mailing_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
	),
	address_logs AS (
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', move_addresses.address_type,
						'shipment_type', move_addresses.shipment_type,
						'shipment_id_abbr', (CASE WHEN move_addresses.shipment_id IS NOT NULL THEN LEFT(move_addresses.shipment_id::TEXT, 5) ELSE NULL END)
					)
				)
			)::TEXT AS context,
			COALESCE(move_addresses.shipment_id::TEXT, move_addresses.service_member_id::TEXT, NULL)::TEXT AS context_id
		FROM
			audit_history
				JOIN move_addresses ON move_addresses.address_id = audit_history.object_id
		WHERE audit_history.table_name = 'addresses'
		GROUP BY
			move_addresses.shipment_id, move_addresses.service_member_id, audit_history.id
	),
	file_uploads (user_upload_id, filename, upload_type) AS (
		-- orders uploads have the document id the uploaded orders id column
		SELECT
			user_uploads.id,
			uploads.filename,
			'orders'
		FROM user_uploads
			JOIN documents ON user_uploads.document_id = documents.id
			JOIN move_orders ON move_orders.uploaded_orders_id = documents.id
			LEFT JOIN uploads ON user_uploads.upload_id = uploads.id
		WHERE documents.service_member_id = move_orders.service_member_id

		-- amended orders have the document id in the uploaded amended orders id column
		UNION
		SELECT
			user_uploads.id,
			uploads.filename,
			'amendedOrders'
		FROM user_uploads
			JOIN documents ON user_uploads.document_id = documents.id
			JOIN move_orders ON move_orders.uploaded_amended_orders_id = documents.id
			LEFT JOIN uploads ON user_uploads.upload_id = uploads.id
		WHERE documents.service_member_id = move_orders.service_member_id
	),
	file_uploads_logs as (
		SELECT
		    audit_history.*,
			json_agg(
				json_build_object(
					'filename', filename,
					'upload_type', upload_type
				)
			)::TEXT AS context,
		NULL AS context_id
		FROM
			audit_history
				JOIN file_uploads ON user_upload_id = audit_history.object_id
		WHERE audit_history.table_name = 'user_uploads'
		GROUP BY audit_history.id
	),
	move_backup_contacts AS (
		SELECT backup_contacts.*
		FROM
			backup_contacts
		WHERE
			backup_contacts.service_member_id = (SELECT id FROM move_service_members)
	),
	backup_contacts_logs as (
		SELECT audit_history.*,
			NULL AS context,
			NULL AS context_id
		FROM
			audit_history
		JOIN move_backup_contacts ON move_backup_contacts.id = audit_history.object_id
	),
	combined_logs AS (
		SELECT
			*
		FROM
			address_logs
		UNION
		SELECT
			*
		FROM
			service_item_logs
		UNION
		SELECT
			*
		FROM
			service_item_customer_contacts_logs
		UNION
		SELECT
			*
		FROM
			service_item_dimensions_logs
		UNION
		SELECT
			*
		FROM
			shipment_logs
		UNION
		SELECT
			*
		FROM
			entitlements_logs
		UNION
		SELECT
			*
		FROM
			reweigh_logs
		UNION
		SELECT
			*
		FROM
			orders_logs
		UNION
		SELECT
			*
		FROM
			agents_logs
		UNION
		SELECT
			*
		FROM
			payment_requests_logs
		UNION
		SELECT
			*
		FROM
			proof_of_service_docs_logs
		UNION
		SELECT
			*
		FROM
			move_logs
		UNION
		SELECT
		 	*
		FROM
			service_members_logs
		UNION
		SELECT
		    *
		FROM
		    file_uploads_logs
		UNION
		SELECT
		 	*
		FROM
			backup_contacts_logs


	) SELECT DISTINCT
		combined_logs.*,
		COALESCE(office_users.first_name, prime_user_first_name, service_members.first_name) AS session_user_first_name,
		COALESCE(office_users.last_name, service_members.last_name) AS session_user_last_name,
		COALESCE(office_users.email, service_members.personal_email) AS session_user_email,
		COALESCE(office_users.telephone, service_members.telephone) AS session_user_telephone
FROM
	combined_logs
		LEFT JOIN users_roles ON session_userid = users_roles.user_id
		LEFT JOIN roles ON users_roles.role_id = roles.id
		LEFT JOIN office_users ON office_users.user_id = session_userid
		LEFT JOIN service_members ON service_members.user_id = session_userid
		LEFT JOIN (
			SELECT 'Prime' AS prime_user_first_name
			) prime_users ON roles.role_type = 'prime'
	ORDER BY
		action_tstamp_tx DESC