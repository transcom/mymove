-- B-22911 Beth introduced a move history sql refactor for us to swap
-- out with the pop query to be more efficient

set client_min_messages = debug;
set session statement_timeout = '10000s';

CREATE OR REPLACE FUNCTION fetch_move_history (move_code text, page integer DEFAULT 1, per_page integer DEFAULT 20, sort text DEFAULT NULL::text, sort_direction text DEFAULT NULL::text)
 RETURNS TABLE (id uuid,
				schema_name text,
				table_name text,
				relid oid,
				object_id uuid,
				session_userid uuid,
				event_name text,
				action_tstamp_tx timestamptz,
				action_tstamp_stm timestamptz,
				action_tstamp_clk timestamptz,
				transaction_id int8,
				client_query text,
				"action" text,
				old_data jsonb,
				changed_data jsonb,
				statement_only bool,
				context text,
				context_id text,
				move_id uuid,
 				shipment_id uuid,
 				session_user_first_name text,
 				session_user_last_name text,
 				session_user_email text,
 				session_user_telephone text,
 				seq_num int)
AS
$BODY$

DECLARE
	v_count INT;
	v_move_id UUID;
	sql_query TEXT;
	offset_value INTEGER;
	sort_column TEXT;
	sort_order TEXT;
v_rowcount int;

BEGIN

	IF page < 1 THEN
		page := 1;
	END IF;

	IF per_page < 1 THEN
		per_page := 20;
	END IF;

	offset_value := (page - 1) * per_page;

	SELECT moves.id into v_move_id
	  FROM moves
	 WHERE moves.locator = move_code;

	IF v_move_id is null THEN
		RAISE EXCEPTION 'Move record not found for %', move_id;
		RETURN;
	END IF;

	DROP TABLE IF EXISTS audit_hist_temp;

	CREATE TEMP TABLE audit_hist_temp
		(id uuid,
		schema_name text,
		table_name text,
		relid oid,
		object_id uuid,
		session_userid uuid,
		event_name text,
		action_tstamp_tx timestamptz,
		action_tstamp_stm timestamptz,
		action_tstamp_clk timestamptz,
		transaction_id int8,
		client_query text,
		"action" text,
		old_data jsonb,
		changed_data jsonb,
		statement_only bool,
		seq_num int,
		context text,
		context_id text,
		move_id uuid,
		shipment_id uuid);

	CREATE INDEX audit_hist_temp_session_userid ON audit_hist_temp (session_userid);

	--moves
	INSERT INTO audit_hist_temp
	SELECT
		audit_history.*,
		jsonb_agg(jsonb_strip_nulls(
			jsonb_build_object(
				'closeout_office_name',
				(SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.closeout_office_id)),
				'counseling_office_name',
				(SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.counseling_transportation_office_id)),
				'assigned_office_user_first_name',
				(SELECT office_users.first_name FROM office_users WHERE office_users.id IN (uuid(c.sc_assigned_id), uuid(c.too_assigned_id), uuid(c.tio_assigned_id))),
				'assigned_office_user_last_name',
				(SELECT office_users.last_name FROM office_users WHERE office_users.id IN (uuid(c.sc_assigned_id), uuid(c.too_assigned_id), uuid(c.tio_assigned_id)))
			))
		)::TEXT AS context,
		NULL AS context_id,
		audit_history.object_id::uuid AS move_id,
		NULL as shipment_id
	FROM
		audit_history
	JOIN jsonb_to_record(audit_history.changed_data) as c(closeout_office_id TEXT, counseling_transportation_office_id TEXT, sc_assigned_id TEXT, too_assigned_id TEXT, tio_assigned_id TEXT) ON TRUE
	WHERE audit_history.table_name = 'moves'
		-- Remove log for when shipment_seq_num updates
	AND NOT (audit_history.event_name = NULL AND audit_history.changed_data::TEXT LIKE '%shipment_seq_num%' AND LENGTH(audit_history.changed_data::TEXT) < 25)
    AND audit_history.object_id = v_move_id
	group by audit_history.id;

	--mto_shipments
	select count(*) into v_count
	  from mto_shipments a
	  join moves b on a.move_id = b.id
	 where b.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
					audit_history.*,
					NULLIF(
						jsonb_agg(jsonb_strip_nulls(
							jsonb_build_object(
								'shipment_type', mto_shipments.shipment_type,
								'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
								'shipment_locator', mto_shipments.shipment_locator
							)
						))::TEXT, '[{}]'::TEXT
					) AS context,
					NULL AS context_id,
					mto_shipments.move_id,
					mto_shipments.id
				FROM
					audit_history
					JOIN mto_shipments ON mto_shipments.id = audit_history.object_id
					JOIN moves on mto_shipments.move_id = v_move_id
				WHERE audit_history.table_name = 'mto_shipments'
					AND NOT (audit_history.event_name = 'updateMTOStatusServiceCounselingCompleted' AND audit_history.changed_data = '{"status": "APPROVED"}')
						-- Not including status update to 'Approval' on mto_shipment layer above ppm_shipment when PPM is counseled.
						-- That is not needed for move history UI.
					AND NOT (audit_history.event_name = 'submitMoveForApproval' AND audit_history.changed_data = '{"status": "SUBMITTED"}')
						-- Not including update on mto_shipment for ppm_shipment when submitted
						-- handled on seperate event
					AND NOT (audit_history.event_name = NULL AND audit_history.changed_data::TEXT LIKE '%shipment_locator%' AND LENGTH(audit_history.changed_data::TEXT) < 35)
				GROUP BY audit_history.id, mto_shipments.move_id, mto_shipments.id;

	END IF;

	--orders
	select count(*) into v_count
	  from orders a
	  join moves b on a.id = b.orders_id
 	 where b.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
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
 			NULL AS context_id,
			v_move_id as move_id,
			null as shipment_id
		FROM
			audit_history
		JOIN orders on orders.id = audit_history.object_id
		JOIN moves on orders.id = moves.orders_id
		JOIN jsonb_to_record(audit_history.changed_data) as c(origin_duty_location_id TEXT, new_duty_location_id TEXT) on TRUE
		WHERE audit_history.table_name = 'orders'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id;

	END IF;

	--mto_service_items
	select count(*) into v_count
	  from mto_service_items a
	  join moves b on a.move_id = b.id
 	 where b.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
			NULL AS context_id,
			moves.id as move_id,
			mto_shipments.id as shipment_id
		FROM
			audit_history
			JOIN mto_service_items ON mto_service_items.id = audit_history.object_id
			JOIN re_services ON mto_service_items.re_service_id = re_services.id
			LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
			JOIN moves ON moves.id = mto_service_items.move_id
		WHERE audit_history.table_name = 'mto_service_items'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id, mto_service_items.id, moves.id, mto_shipments.id;

	END IF;

	--service_item_customer_contacts
	select count(*) into v_count
	  from mto_service_item_customer_contacts
	  JOIN service_items_customer_contacts on service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
	  JOIN mto_service_items on mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
      JOIN moves ON moves.id = mto_service_items.move_id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
			NULL AS context_id,
			moves.id as move_id,
			mto_shipments.id as shipment_id
		FROM audit_history
		JOIN mto_service_item_customer_contacts ON mto_service_item_customer_contacts.id = audit_history.object_id
		JOIN service_items_customer_contacts on service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
		JOIN mto_service_items on mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
			LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON moves.id = mto_service_items.move_id
		WHERE audit_history.table_name = 'mto_service_item_customer_contacts'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id, mto_service_item_customer_contacts.id, moves.id, mto_shipments.id;

	END IF;

	--service_item_dimensions
	select count(*) into v_count
	  from mto_service_item_dimensions
	  JOIN mto_service_items on mto_service_items.id = mto_service_item_dimensions.mto_service_item_id
	  LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
	  LEFT JOIN moves on moves.id = mto_shipments.move_id
	 WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_locator', mto_shipments.shipment_locator
			))::TEXT AS context,
			NULL AS context_id,
			moves.id as move_id,
			mto_shipments.id as shipment_id
		FROM
			audit_history
		JOIN mto_service_item_dimensions ON mto_service_item_dimensions.id = audit_history.object_id
		JOIN mto_service_items ON mto_service_items.id = mto_service_item_dimensions.mto_service_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON mto_shipments.move_id = moves.id
		WHERE audit_history.table_name = 'mto_service_item_dimensions'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id, mto_service_item_dimensions.id, moves.id, mto_shipments.id;

	END IF;

	--entitlements
	select count(*) into v_count
	  from entitlements
	  JOIN orders on entitlements.id = orders.entitlement_id
	  JOIN moves on orders.id = moves.orders_id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id,
			moves.id as move_id,
			NULL as shipment_id
		FROM
			audit_history
			JOIN entitlements ON entitlements.id = audit_history.object_id
			JOIN orders on entitlements.id = orders.entitlement_id
	  		JOIN moves on orders.id = moves.orders_id
		WHERE audit_history.table_name = 'entitlements'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id, entitlements.id, moves.id;

	END IF;

	--payment_requests
	select count(*) into v_count
	  from payment_requests
	 where payment_requests.move_id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'price', payment_service_items.price_cents::TEXT,
				'status', payment_service_items.status,
				'shipment_id', mto_shipments.id::TEXT,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_type', mto_shipments.shipment_type,
				'shipment_locator', mto_shipments.shipment_locator,
				'rejection_reason', payment_service_items.rejection_reason
				)
			)::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id as move_id,
			mto_shipments.id as shipment_id
		FROM
			audit_history
			JOIN payment_requests ON payment_requests.id = audit_history.object_id
			JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
			JOIN mto_service_items ON mto_service_items.id = payment_service_items.mto_service_item_id
			LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
			JOIN re_services ON mto_service_items.re_service_id = re_services.id
		WHERE audit_history.table_name = 'payment_requests'
		  AND payment_requests.move_id = v_move_id
		GROUP BY
			audit_history.id, payment_requests.id, payment_requests.move_id, mto_shipments.id;

	END IF;

	--payment_service_items
	select count(*) into v_count
	  from payment_requests
	  JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
	 where payment_requests.move_id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'price', payment_service_items.price_cents::TEXT,
				'status', payment_service_items.status,
				'rejection_reason', payment_service_items.rejection_reason,
				'paid_at', payment_service_items.paid_at,
				'shipment_id', mto_shipments.id::TEXT,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_type', mto_shipments.shipment_type,
				'shipment_locator', mto_shipments.shipment_locator,
				'rejection_reason', payment_service_items.rejection_reason
				)
			)::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id as move_id,
			mto_shipments.id as shipment_id
		FROM
			audit_history
			JOIN payment_service_items ON payment_service_items.id = audit_history.object_id
			JOIN payment_requests ON payment_service_items.payment_request_id = payment_requests.id
			JOIN mto_service_items ON mto_service_items.id = payment_service_items.mto_service_item_id
			LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
			JOIN re_services ON mto_service_items.re_service_id = re_services.id
		WHERE audit_history.table_name = 'payment_service_items'
		  AND payment_requests.move_id = v_move_id
		GROUP BY
			audit_history.id, payment_requests.id, payment_requests.move_id, mto_shipments.id;

	END IF;

	--proof_of_service_docs
	select count(*) into v_count
	  from proof_of_service_docs
	  join payment_requests on proof_of_service_docs.payment_request_id = payment_requests.id
	 where payment_requests.move_id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'payment_request_number',
				payment_requests.payment_request_number::TEXT))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id as move_id,
			NULL as shipment_id
		FROM
			audit_history
			JOIN proof_of_service_docs ON proof_of_service_docs.id = audit_history.object_id
			JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
		WHERE audit_history.table_name = 'proof_of_service_docs'
		  AND payment_requests.move_id = v_move_id
		GROUP BY
			audit_history.id, proof_of_service_docs.id, payment_requests.move_id;

	END IF;

	--mto_agents
	select count(*) into v_count
	  from mto_agents
	  join mto_shipments on mto_agents.mto_shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where mto_shipments.move_id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
			NULL AS context_id,
			mto_shipments.move_id as move_id,
			mto_shipments.id as shipment_id
		FROM audit_history
		JOIN mto_agents ON mto_agents.id = audit_history.object_id
		join mto_shipments on mto_agents.mto_shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE audit_history.table_name = 'mto_agents'
		 AND mto_shipments.move_id = v_move_id
		 AND (audit_history.event_name <> 'deleteShipment' OR audit_history.event_name IS NULL)
				-- This event name is used to delete the parent shipment and child agent logs are unnecessary.
				-- NULLS are not counted in comparisons, so we include those as well.
		GROUP BY audit_history.id, mto_agents.id, mto_shipments.move_id, mto_shipments.id;

	END IF;

	--reweighs
	select count(*) into v_count
	  from reweighs
	  join mto_shipments on reweighs.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where mto_shipments.move_id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'payment_request_number', payment_requests.payment_request_number,
				'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
			NULL AS context_id,
			mto_shipments.move_id as move_id,
			mto_shipments.id as shipment_id
		FROM audit_history
		JOIN reweighs ON reweighs.id = audit_history.object_id
		JOIN mto_shipments ON reweighs.shipment_id = mto_shipments.id
		JOIN moves on mto_shipments.move_id = moves.id
		LEFT JOIN payment_requests ON mto_shipments.move_id = payment_requests.move_id
	   WHERE audit_history.table_name = 'reweighs'
		 AND mto_shipments.move_id = v_move_id
	   GROUP BY audit_history.id, reweighs.id, mto_shipments.move_id, mto_shipments.id;

	END IF;

	--service_members
	select count(*) into v_count
	  from service_members
	  join orders on service_members.id = orders.service_member_id
	  join moves on orders.id = moves.orders_id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
				NULLIF(
				jsonb_agg(jsonb_strip_nulls(
					jsonb_build_object(
						'current_duty_location_name',
						(SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.duty_location_id))
					)
				))::TEXT, '[{}]'::TEXT
			) AS context,
			NULL AS context_id,
			moves.id as move_id,
			NULL as shipment_id
		FROM
			audit_history
			JOIN service_members ON service_members.id = audit_history.object_id
			JOIN jsonb_to_record(audit_history.changed_data) as c(duty_location_id TEXT) on TRUE
			join orders on service_members.id = orders.service_member_id
	  		join moves on orders.id = moves.orders_id
		WHERE audit_history.table_name = 'service_members'
		  AND moves.id = v_move_id
		GROUP BY audit_history.id, service_members.id, moves.id;

	END IF;

	--ppm_shipments
	select count(*) into v_count
	  from ppm_shipments
	  join mto_shipments on ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(ppm_shipments.shipment_id::TEXT, 5),
						'w2_address', (SELECT row_to_json(x) FROM (SELECT * FROM addresses WHERE addresses.id = CAST(ppm_shipments.w2_address_id AS UUID)) x)::TEXT,
						'shipment_locator', mto_shipments.shipment_locator,
						'pickup_postal_address_id', ppm_shipments.pickup_postal_address_id,
						'secondary_pickup_postal_address_id', ppm_shipments.secondary_pickup_postal_address_id
					)
				)
			)::TEXT AS context,
			COALESCE(ppm_shipments.shipment_id::TEXT, NULL)::TEXT AS context_id,
			moves.id as move_id,
			mto_shipments.id as shipment_id
		FROM audit_history
		JOIN ppm_shipments ON audit_history.object_id = ppm_shipments.id
		join mto_shipments on ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE audit_history.table_name = 'ppm_shipments'
	     AND moves.id = v_move_id
	   GROUP BY ppm_shipments.id, audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'destination_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'destinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'destination_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - secondary destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'secondaryDestinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - tertiary destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_delivery_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'tertiaryDestinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_delivery_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'pickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - secondary pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'secondaryPickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - tertiary pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join mto_shipments ON a2.object_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'tertiaryPickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join mto_shipments ON a2.object_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'pickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM secondary pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'secondaryPickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM tertiary pickup
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'tertiaryPickupAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'destinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM secondary destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'secondary_destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'secondaryDestinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'secondary_destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - PPM tertiary destination
	select count(*) into v_count
	  from audit_history a1 --addresses
	  join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	  join ppm_shipments ON a2.object_id = ppm_shipments.id
	  join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	  join moves on mto_shipments.move_id = moves.id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			a1.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'tertiaryDestinationAddress',
						'shipment_type', mto_shipments.shipment_type,
						'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
						'shipment_locator', mto_shipments.shipment_locator
					)
				)
			)::TEXT AS context,
			mto_shipments.id::TEXT AS context_id,	--shipment address
			moves.id as move_id,
			mto_shipments.id as shipment_id
		from audit_history a1 --addresses
	    join audit_hist_temp a2 ON (a2.changed_data->>'tertiary_destination_postal_address_id')::uuid = a1.object_id AND a1."table_name" = 'addresses' --mto_shipments changes
	    join ppm_shipments ON a2.object_id = ppm_shipments.id
	    join mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
	    join moves on mto_shipments.move_id = moves.id
	   WHERE a1.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY a1.id, moves.id, mto_shipments.id;

	END IF;

	--addresses - service member residential
	select count(*) into v_count
	  from addresses
	  join service_members ON service_members.residential_address_id = addresses.id
	  join orders ON service_members.id = orders.service_member_id
	  join moves on orders.id = moves.orders_id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'residentialAddress'
					)
				)
			)::TEXT AS context,
			service_members.id::TEXT AS context_id,	--service members address
			moves.id as move_id,
			NULL as shipment_id
		from audit_history
	    join service_members ON service_members.residential_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
	    join orders ON service_members.id = orders.service_member_id
	    join moves on orders.id = moves.orders_id
	   WHERE audit_history.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY audit_history.id, moves.id, service_members.id;

	END IF;

	--addresses - service member backup mailing
	select count(*) into v_count
	  from addresses
	  join service_members ON service_members.backup_mailing_address_id = addresses.id
	  join orders ON service_members.id = orders.service_member_id
	  join moves on orders.id = moves.orders_id
	 where moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_strip_nulls(
					jsonb_build_object(
						'address_type', 'backupMailingAddress'
					)
				)
			)::TEXT AS context,
			service_members.id::TEXT AS context_id,	--service members address
			moves.id as move_id,
			NULL as shipment_id
		from audit_history
	    join service_members ON service_members.backup_mailing_address_id = audit_history.object_id AND audit_history."table_name" = 'addresses'
	    join orders ON service_members.id = orders.service_member_id
	    join moves on orders.id = moves.orders_id
	   WHERE audit_history.table_name = 'addresses'
	     AND moves.id = v_move_id
		GROUP BY audit_history.id, moves.id, service_members.id;

	END IF;

	--file_uploads - orders
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN orders ON orders.uploaded_orders_id = documents.id
		AND documents.service_member_id = orders.service_member_id
	  join moves on orders.id = moves.orders_id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'orders'
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		NULL as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN orders ON orders.uploaded_orders_id = documents.id
		  AND documents.service_member_id = orders.service_member_id
	    join moves on orders.id = moves.orders_id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id;

	END IF;

	--file_uploads - amended orders
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN orders ON orders.uploaded_amended_orders_id = documents.id
		AND documents.service_member_id = orders.service_member_id
	  join moves on orders.id = moves.orders_id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'amendedOrders'
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		NULL as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN orders ON orders.uploaded_amended_orders_id = documents.id
		  AND documents.service_member_id = orders.service_member_id
	    join moves on orders.id = moves.orders_id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id;

	END IF;

	--file_uploads - empty weight ticket
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN weight_tickets ON weight_tickets.empty_document_id = documents.id
	  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'emptyWeightTicket',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN weight_tickets ON weight_tickets.empty_document_id = documents.id
	    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--file_uploads - full weight ticket
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
	  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'fullWeightTicket',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
	    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;


	--file_uploads - trailer weight ticket
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN weight_tickets ON weight_tickets.proof_of_trailer_ownership_document_id = documents.id
	  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'trailerWeightTicket',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN weight_tickets ON weight_tickets.proof_of_trailer_ownership_document_id = documents.id
	    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;


	--file_uploads - pro gear weight ticket
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND progear_weight_tickets.belongs_to_self = true
	  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'proGearWeightTicket',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND progear_weight_tickets.belongs_to_self = true
	    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;


	--file_uploads - spouse pro gear weight ticket
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND coalesce(progear_weight_tickets.belongs_to_self,false) = false
	  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'spouseProGearWeightTicket',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND coalesce(progear_weight_tickets.belongs_to_self,false) = false
	    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--file_uploads - expense receipt
	select count(*) into v_count
	  from user_uploads
	  JOIN documents ON user_uploads.document_id = documents.id
	  JOIN moving_expenses ON moving_expenses.document_id = documents.id
	  JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
	  JOIN uploads ON user_uploads.upload_id = uploads.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			json_agg(
				json_build_object(
					'filename', uploads.filename,
					'upload_type', 'expenseReceipt',
					'shipment_type', mto_shipments.shipment_type,
					'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
					'moving_expense_type', moving_expenses.moving_expense_type::TEXT,
					'shipment_locator', mto_shipments.shipment_locator
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN user_uploads ON user_uploads.id = audit_history.object_id
	    JOIN documents ON user_uploads.document_id = documents.id
	    JOIN moving_expenses ON moving_expenses.document_id = documents.id
	    JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
	    JOIN uploads ON user_uploads.upload_id = uploads.id
	   WHERE audit_history.table_name = 'user_uploads'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--backup contacts
	select count(*) into v_count
	  from backup_contacts
	  JOIN service_members ON service_members.id = backup_contacts.service_member_id
	  JOIN orders on orders.service_member_id = service_members.id
	  join moves on moves.orders_id = orders.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT audit_history.*,
			NULL AS context,
		NULL AS context_id,
		moves.id as move_id,
		NULL as shipment_id
		FROM audit_history
		JOIN backup_contacts ON backup_contacts.id = audit_history.object_id
	    JOIN service_members ON service_members.id = backup_contacts.service_member_id
	    JOIN orders on orders.service_member_id = service_members.id
	    join moves on moves.orders_id = orders.id
	   WHERE audit_history.table_name = 'backup_contacts'
		 AND moves.id = v_move_id
	   GROUP BY audit_history.id, moves.id;

	END IF;

	--document review items - weight tickets
	select count(*) into v_count
	  from audit_history
	  JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
	  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
     WHERE moves.id = v_move_id
	   AND audit_history.table_name = 'weight_tickets';

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
	        audit_history.*,
	            jsonb_agg(
	                jsonb_strip_nulls(
	                    jsonb_build_object(
	                        'shipment_type', mto_shipments.shipment_type,
	                        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
							'shipment_locator', mto_shipments.shipment_locator
	                    )
	                )
	            )::TEXT AS context,
	    COALESCE(mto_shipments.id::TEXT, NULL)::TEXT AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
	    FROM audit_history
	    JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
	    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
       WHERE moves.id = v_move_id
	     AND audit_history.table_name = 'weight_tickets'
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--document review items - pro gear weight tickets
	select count(*) into v_count
	  from audit_history
	  JOIN progear_weight_tickets ON progear_weight_tickets.id = audit_history.object_id
	  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
     WHERE moves.id = v_move_id
	   AND audit_history.table_name = 'progear_weight_tickets';

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
	        audit_history.*,
	            jsonb_agg(
	                jsonb_strip_nulls(
	                    jsonb_build_object(
	                        'shipment_type', mto_shipments.shipment_type,
	                        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
							'shipment_locator', mto_shipments.shipment_locator
	                    )
	                )
	            )::TEXT AS context,
	    COALESCE(mto_shipments.id::TEXT, NULL)::TEXT AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
	    FROM audit_history
	    JOIN progear_weight_tickets ON progear_weight_tickets.id = audit_history.object_id
	    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
       WHERE moves.id = v_move_id
	     AND audit_history.table_name = 'progear_weight_tickets'
	   GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--document review items - moving expenses
	select count(*) into v_count
	  from audit_history
	  JOIN moving_expenses ON moving_expenses.id = audit_history.object_id
	  JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
	  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	  join moves on mto_shipments.move_id = moves.id
     WHERE moves.id = v_move_id
	   AND audit_history.table_name = 'moving_expenses';

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
	        audit_history.*,
	            jsonb_agg(
	                jsonb_strip_nulls(
	                    jsonb_build_object(
	                        'shipment_type', mto_shipments.shipment_type,
	                        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
	                        'moving_expense_type', moving_expenses.moving_expense_type,
							'shipment_locator', mto_shipments.shipment_locator
	                    )
	                )
	            )::TEXT AS context,
	    COALESCE(mto_shipments.id::TEXT, NULL)::TEXT AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
	    FROM audit_history
	    JOIN moving_expenses ON moving_expenses.id = audit_history.object_id
	    JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
	    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
	    join moves on mto_shipments.move_id = moves.id
       WHERE moves.id = v_move_id
	     AND audit_history.table_name = 'moving_expenses'
	  GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

	--gsr appeals
	select count(*) into v_count
	  from gsr_appeals
	  JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
	  LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
	  LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
	  JOIN moves ON evaluation_reports.move_id = moves.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(
				jsonb_build_object(
					'evaluation_report_type', evaluation_reports.type,
					'violation_paragraph_number', pws_violations.paragraph_number,
					'violation_title', pws_violations.title,
					'violation_summary', pws_violations.requirement_summary
				)
			)::TEXT AS context,
		NULL AS context_id,
		moves.id as move_id,
		NULL as shipment_id
		FROM audit_history
		JOIN gsr_appeals ON gsr_appeals.id = audit_history.object_id
		JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
	    LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
	    LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
	    JOIN moves ON evaluation_reports.move_id = moves.id
		WHERE audit_history.table_name = 'gsr_appeals'
 		  AND moves.id = v_move_id
		GROUP BY audit_history.id, moves.id;

	END IF;

	--shipment address updates
	select count(*) into v_count
	  from shipment_address_updates
	  JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
	  JOIN moves ON mto_shipments.move_id = moves.id
     WHERE moves.id = v_move_id;

	IF v_count > 0 THEN

		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'status', shipment_address_updates.status
				)
			)::TEXT AS context,
			NULL AS context_id,
		moves.id as move_id,
		mto_shipments.id as shipment_id
		FROM audit_history
		JOIN shipment_address_updates ON shipment_address_updates.id = audit_history.object_id
	    JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
	    JOIN moves ON mto_shipments.move_id = moves.id
		WHERE audit_history.table_name = 'shipment_address_updates'
 		  AND moves.id = v_move_id
		GROUP BY audit_history.id, moves.id, mto_shipments.id;

	END IF;

raise debug 'Start return query %', clock_timestamp();

	RETURN QUERY
	SELECT
		x.id,
		x.schema_name,
		x.table_name,
		x.relid,
		x.object_id,
		x.session_userid,
		x.event_name,
		x.action_tstamp_tx,
		x.action_tstamp_stm,
		x.action_tstamp_clk,
		x.transaction_id,
		x.client_query,
		x."action",
		x.old_data,
		x.changed_data,
		x.statement_only,
		x.context,
		x.context_id,
		x.move_id,
		x.shipment_id,
		COALESCE(office_users.first_name, prime_user_first_name, service_members.first_name) AS session_user_first_name,
		COALESCE(office_users.last_name, service_members.last_name) AS session_user_last_name,
		COALESCE(office_users.email, service_members.personal_email) AS session_user_email,
		COALESCE(office_users.telephone, service_members.telephone) AS session_user_telephone,
		x.seq_num
	FROM audit_hist_temp x
	LEFT JOIN users_roles ON x.session_userid = users_roles.user_id
	LEFT JOIN roles ON users_roles.role_id = roles.id
	LEFT JOIN office_users ON office_users.user_id = x.session_userid
	LEFT JOIN service_members ON service_members.user_id = x.session_userid
	LEFT JOIN (
		SELECT 'Prime' AS prime_user_first_name
	) prime_users ON roles.role_type = 'prime'
	ORDER BY x.action_tstamp_tx DESC
	LIMIT per_page OFFSET offset_value;

select count(*) into v_count from audit_hist_temp;
raise debug 'Total recs %', v_count;

end

$BODY$
language plpgsql;
