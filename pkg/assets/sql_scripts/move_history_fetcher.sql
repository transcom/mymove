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
			AND NOT (audit_history.event_name = 'updateMTOStatusServiceCounselingCompleted' AND audit_history.changed_data = '{"status": "APPROVED"}')
				-- Not including status update to 'Approval' on mto_shipment layer above ppm_shipment when PPM is counseled.
				-- That is not needed for move history UI.
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
				'name', re_services.name,
				'shipment_type', move_shipments.shipment_type,
				'shipment_id_abbr', move_shipments.shipment_id_abbr
				)
			)::TEXT AS context
		FROM
			mto_service_item_customer_contacts
		JOIN service_items_customer_contacts on service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
		JOIN move_service_items on move_service_items.id = service_items_customer_contacts.mtoservice_item_id
		JOIN re_services ON move_service_items.re_service_id = re_services.id
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
			AND (audit_history.event_name <> 'deleteShipment' OR audit_history.event_name IS NULL)
				-- This event name is used to delete the parent shipment and child agent logs are unnecessary.
				-- NULLS are not counted in comparisons, so we include those as well.
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
	ppms (ppm_id, shipment_type, shipment_id, w2_address_id) AS (
		SELECT
			audit_history.object_id,
			move_shipments.shipment_type,
			move_shipments.id,
			ppm_shipments.w2_address_id
		FROM
			audit_history
		JOIN ppm_shipments ON audit_history.object_id = ppm_shipments.id
		JOIN move_shipments ON move_shipments.id = ppm_shipments.shipment_id
	),
	ppm_logs AS (
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'shipment_type', ppms.shipment_type,
						'shipment_id_abbr', (CASE WHEN ppms.shipment_id IS NOT NULL THEN LEFT(ppms.shipment_id::TEXT, 5) ELSE NULL END),
						'w2_address', (SELECT row_to_json(x) FROM (SELECT * FROM addresses WHERE addresses.id = CAST(ppms.w2_address_id AS UUID)) x)::TEXT
					)
				)
			)::TEXT AS context,
			COALESCE(ppms.shipment_id::TEXT, NULL)::TEXT AS context_id
		FROM
			audit_history
		JOIN ppms ON ppms.ppm_id = audit_history.object_id
		WHERE audit_history.table_name = 'ppm_shipments'
		GROUP BY
			ppms.shipment_id, audit_history.id
	),
	move_sits (sit_address_updates_id, address_id, service_name, shipment_type, shipment_id, original_address_id, office_remarks, contractor_remarks)  AS (
		SELECT
			audit_history.object_id,
			sit_address_updates.new_address_id,
			re_services.name,
			move_shipments.shipment_type,
			move_shipments.id,
			move_service_items.sit_destination_original_address_id,
			sit_address_updates.office_remarks,
			sit_address_updates.contractor_remarks
		FROM audit_history
			JOIN sit_address_updates ON audit_history.object_id = sit_address_updates.id AND audit_history.table_name = 'sit_address_updates'
			JOIN move_service_items ON move_service_items.id = sit_address_updates.mto_service_item_id
			JOIN move_shipments ON move_shipments.id = move_service_items.mto_shipment_id
			JOIN re_services ON move_service_items.re_service_id = re_services.id
	),
	sit_logs AS (
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'shipment_type', move_sits.shipment_type,
						'shipment_id_abbr', (CASE WHEN move_sits.shipment_id IS NOT NULL THEN LEFT(move_sits.shipment_id::TEXT, 5) ELSE NULL END),
						'sit_destination_address_final', (SELECT row_to_json(x) FROM (SELECT * FROM addresses WHERE addresses.id = CAST(move_sits.address_id AS UUID)) x)::TEXT,
						'sit_destination_address_initial', (SELECT row_to_json(x) FROM (SELECT * FROM addresses WHERE addresses.id = CAST(move_sits.original_address_id AS UUID)) x)::TEXT,
						'office_remarks', move_sits.office_remarks,
						'contractor_remarks', move_sits.contractor_remarks,
						'name', move_sits.service_name
					)
				)
			)::TEXT AS context,
			COALESCE(move_sits.shipment_id::TEXT, NULL)::TEXT AS context_id
		FROM
			audit_history
				JOIN move_sits ON move_sits.sit_address_updates_id = audit_history.object_id
		WHERE audit_history.table_name = 'sit_address_updates'
		GROUP BY
			move_sits.shipment_id, audit_history.id
	),
	file_uploads (user_upload_id, filename, upload_type, shipment_type, shipment_id_abbr, expense_type) AS (
		-- orders uploads have the document id the uploaded orders id column
		SELECT
			user_uploads.id,
			uploads.filename,
			'orders',
			NULL,
			NULL,
			NULL
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
			'amendedOrders',
			NULL,
			NULL,
			NULL
		FROM user_uploads
			JOIN documents ON user_uploads.document_id = documents.id
			JOIN move_orders ON move_orders.uploaded_amended_orders_id = documents.id
			LEFT JOIN uploads ON user_uploads.upload_id = uploads.id
		WHERE documents.service_member_id = move_orders.service_member_id

		UNION
		SELECT
			user_uploads.id,
			uploads.filename,
			CASE WHEN weight_tickets.empty_document_id = documents.id THEN 'emptyWeightTicket'
				 WHEN weight_tickets.full_document_id = documents.id THEN 'fullWeightTicket'
				 WHEN weight_tickets.proof_of_trailer_ownership_document_id = documents.id THEN 'trailerWeightTicket'
				 WHEN progear_weight_tickets.document_id = documents.id AND progear_weight_tickets.belongs_to_self = true THEN 'proGearWeightTicket'
				 WHEN progear_weight_tickets.document_id = documents.id AND (progear_weight_tickets.belongs_to_self IS NULL OR progear_weight_tickets.belongs_to_self = false) THEN 'spouseProGearWeightTicket'
				 WHEN moving_expenses.document_id = documents.id THEN 'expenseReceipt'
				 ELSE '' END,
			move_shipments.shipment_type::TEXT,
			move_shipments.shipment_id_abbr,
			CASE WHEN moving_expenses.document_id = documents.id THEN moving_expenses.moving_expense_type::text
				 ELSE NULL END
		FROM user_uploads
			JOIN documents ON user_uploads.document_id = documents.id
			LEFT JOIN weight_tickets ON weight_tickets.empty_document_id = documents.id OR weight_tickets.full_document_id = documents.id OR weight_tickets.proof_of_trailer_ownership_document_id = documents.id
			LEFT JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id
			LEFT JOIN moving_expenses ON moving_expenses.document_id = documents.id
			JOIN ppm_shipments ON ppm_shipments.id = COALESCE(weight_tickets.ppm_shipment_id, progear_weight_tickets.ppm_shipment_id, moving_expenses.ppm_shipment_id)
			JOIN move_shipments ON ppm_shipments.shipment_id = move_shipments.id
			JOIN uploads ON user_uploads.upload_id = uploads.id
	),
	file_uploads_logs as (
		SELECT
		    audit_history.*,
			json_agg(
				json_build_object(
					'filename', filename,
					'upload_type', upload_type,
					'shipment_type', shipment_type,
					'shipment_id_abbr', shipment_id_abbr,
					'moving_expense_type', expense_type
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
	-- document_review_items grabs historical data for weight ticket, pro-gear
	-- weight tickets, and moving expense tickets
	document_review_items (doc_id, shipment_type, shipment_id, moving_expense_type) AS (
		SELECT COALESCE(wt.id, pwt.id, me.id),
			ppms.shipment_type,
			ppms.shipment_id,
			me.moving_expense_type
		FROM audit_history ah
		LEFT JOIN weight_tickets wt ON ah.object_id = wt.id
		LEFT JOIN progear_weight_tickets pwt ON ah.object_id = pwt.id
		LEFT JOIN moving_expenses me on ah.object_id = me.id
		JOIN ppms ON ppms.ppm_id = COALESCE(wt.ppm_shipment_id, pwt.ppm_shipment_id, me.ppm_shipment_id)
	),
	document_review_logs AS (
	    SELECT
	        audit_history.*,
	            jsonb_agg(
	                jsonb_strip_nulls(
	                    jsonb_build_object(
	                        'shipment_type', document_review_items.shipment_type,
	                        'shipment_id_abbr', (CASE WHEN document_review_items.shipment_id IS NOT NULL THEN LEFT(document_review_items.shipment_id::TEXT, 5) ELSE NULL END),
	                        'moving_expense_type', document_review_items.moving_expense_type
	                    )
	                )
	            )::TEXT AS context,
	        COALESCE(document_review_items.shipment_id::TEXT, NULL)::TEXT AS context_id
	    FROM
	        audit_history
	    JOIN document_review_items ON document_review_items.doc_id = audit_history.object_id
	    WHERE audit_history.table_name = 'weight_tickets' OR audit_history.table_name = 'progear_weight_tickets' OR audit_history.table_name = 'moving_expenses'
	    GROUP BY
	        document_review_items.doc_id, document_review_items.shipment_id, audit_history.id
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
			sit_logs
		UNION
		SELECT
			*
		FROM
			ppm_logs
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
		UNION
		SELECT
			*
		FROM
			document_review_logs


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
