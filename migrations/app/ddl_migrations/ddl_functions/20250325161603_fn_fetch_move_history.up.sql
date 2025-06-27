-- B-22911 Beth introduced a move history sql refactor for us to swapnout with the pop query to be more efficient
-- B-22924  Daniel Jordan  adding sit_extension table to history and updating main func
-- B-23602  Beth Grohmann  fixed join in fn_populate_move_history_mto_shipments
-- B-23581 Paul Stonebraker updated assigned office user counselor columns
-- B 22696 Jon Spight added too destination assignments to history / Audit log
-- B-23623  Beth Grohmann  fetch_move_history - update final query to pull from all user tables

set client_min_messages = debug;
set session statement_timeout = '10000s';

-- ============================================
-- Sub-function: check and resolve move ID
-- ============================================
CREATE OR REPLACE FUNCTION fn_get_move_id(move_code TEXT)
RETURNS UUID AS $$
DECLARE
    v_move_id UUID;
BEGIN
    SELECT moves.id INTO v_move_id
    FROM moves
    WHERE moves.locator = move_code;

    IF v_move_id IS NULL THEN
        RAISE EXCEPTION 'Move record not found for %', move_code;
    END IF;

    RETURN v_move_id;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- Sub-function: create the temp table
-- ============================================
CREATE OR REPLACE FUNCTION fn_create_audit_temp_table()
RETURNS VOID AS $$
BEGIN
    DROP TABLE IF EXISTS audit_hist_temp;

    CREATE TEMP TABLE audit_hist_temp (
        id uuid PRIMARY KEY, -- Prevent grouping duplicates
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
        shipment_id uuid
    );

    CREATE INDEX audit_hist_temp_session_userid ON audit_hist_temp (session_userid);
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- Sub-function: populate move data
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_moves(v_move_id UUID)
RETURNS VOID AS '
DECLARE
    v_count INT;
BEGIN
    INSERT INTO audit_hist_temp
    SELECT
        audit_history.*,
        jsonb_agg(jsonb_strip_nulls(
            jsonb_build_object(
                ''closeout_office_name'',
                (SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.closeout_office_id)),
                ''counseling_office_name'',
                (SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.counseling_transportation_office_id)),
                ''assigned_office_user_first_name'',
                (SELECT office_users.first_name FROM office_users WHERE office_users.id IN (uuid(c.sc_counseling_assigned_id), uuid(c.sc_closeout_assigned_id), uuid(c.too_task_order_assigned_id), uuid(c.tio_assigned_id), uuid(c.too_destination_assigned_id))),
                ''assigned_office_user_last_name'',
                (SELECT office_users.last_name FROM office_users WHERE office_users.id IN (uuid(c.sc_counseling_assigned_id), uuid(c.sc_closeout_assigned_id), uuid(c.too_task_order_assigned_id), uuid(c.tio_assigned_id), uuid(c.too_destination_assigned_id)))
            ))
        )::TEXT AS context,
        NULL AS context_id,
        audit_history.object_id::uuid AS move_id,
        NULL AS shipment_id
    FROM
        audit_history
    JOIN jsonb_to_record(audit_history.changed_data) AS c(
        closeout_office_id TEXT,
        counseling_transportation_office_id TEXT,
        sc_counseling_assigned_id TEXT,
        sc_closeout_assigned_id TEXT,
        too_task_order_assigned_id TEXT,
        too_destination_assigned_id TEXT,
        tio_payment_request_assigned_id TEXT
    ) ON TRUE
    WHERE audit_history.table_name = ''moves''
        AND NOT (audit_history.event_name IS NULL AND audit_history.changed_data::TEXT LIKE ''%shipment_seq_num%'' AND LENGTH(audit_history.changed_data::TEXT) < 25)
        AND audit_history.object_id = v_move_id
    GROUP BY audit_history.id;
END;
'
LANGUAGE plpgsql;

-- ============================================
-- Sub-function: populate mto_shipments
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_mto_shipments(v_move_id UUID)
RETURNS VOID AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM mto_shipments a
    JOIN moves b ON a.move_id = b.id
    WHERE b.id = v_move_id;

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
            mto_shipments.move_id
        FROM
            audit_history
        JOIN mto_shipments ON mto_shipments.id = audit_history.object_id
        JOIN moves ON mto_shipments.move_id = moves.id
        WHERE audit_history.table_name = 'mto_shipments'
            AND moves.id = v_move_id
            AND NOT (audit_history.event_name = 'updateMTOStatusServiceCounselingCompleted' AND audit_history.changed_data = '{"status": "APPROVED"}')
            AND NOT (audit_history.event_name = 'submitMoveForApproval' AND audit_history.changed_data = '{"status": "SUBMITTED"}')
            AND NOT (audit_history.event_name IS NULL AND audit_history.changed_data::TEXT LIKE '%shipment_locator%' AND LENGTH(audit_history.changed_data::TEXT) < 35)
        GROUP BY audit_history.id, mto_shipments.move_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- Sub-function: populate orders
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_orders(v_move_id UUID)
RETURNS VOID AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM orders a
    JOIN moves b ON a.id = b.orders_id
    WHERE b.id = v_move_id;

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
            v_move_id AS move_id,
            NULL AS shipment_id
        FROM
            audit_history
        JOIN orders ON orders.id = audit_history.object_id
        JOIN moves ON orders.id = moves.orders_id
        JOIN jsonb_to_record(audit_history.changed_data) AS c(
            origin_duty_location_id TEXT,
            new_duty_location_id TEXT
        ) ON TRUE
        WHERE audit_history.table_name = 'orders'
            AND moves.id = v_move_id
        GROUP BY audit_history.id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- Sub-function: populate service items
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_service_items(v_move_id UUID)
RETURNS VOID AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM mto_service_items a
    JOIN moves b ON a.move_id = b.id
    WHERE b.id = v_move_id;

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
            moves.id AS move_id
        FROM
            audit_history
        JOIN mto_service_items ON mto_service_items.id = audit_history.object_id
        JOIN re_services ON mto_service_items.re_service_id = re_services.id
        LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
        JOIN moves ON moves.id = mto_service_items.move_id
        WHERE audit_history.table_name = 'mto_service_items'
            AND moves.id = v_move_id
        GROUP BY audit_history.id, mto_service_items.id, moves.id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate service item customer contacts
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_service_item_customer_contacts(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM mto_service_item_customer_contacts
	JOIN service_items_customer_contacts ON service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
	JOIN mto_service_items ON mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
	JOIN moves ON moves.id = mto_service_items.move_id
	WHERE moves.id = p_move_id;

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
			moves.id AS move_id
		FROM audit_history
		JOIN mto_service_item_customer_contacts ON mto_service_item_customer_contacts.id = audit_history.object_id
		JOIN service_items_customer_contacts ON service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
		JOIN mto_service_items ON mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON moves.id = mto_service_items.move_id
		WHERE audit_history.table_name = 'mto_service_item_customer_contacts'
		  AND moves.id = p_move_id
		GROUP BY audit_history.id, mto_service_item_customer_contacts.id, moves.id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate service item dimensions
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_service_item_dimensions(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM mto_service_item_dimensions
	JOIN mto_service_items ON mto_service_items.id = mto_service_item_dimensions.mto_service_item_id
	LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
	LEFT JOIN moves ON moves.id = mto_shipments.move_id
	WHERE moves.id = p_move_id;

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
			moves.id AS move_id
		FROM audit_history
		JOIN mto_service_item_dimensions ON mto_service_item_dimensions.id = audit_history.object_id
		JOIN mto_service_items ON mto_service_items.id = mto_service_item_dimensions.mto_service_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON mto_shipments.move_id = moves.id
		WHERE audit_history.table_name = 'mto_service_item_dimensions'
		  AND moves.id = p_move_id
		GROUP BY audit_history.id, mto_service_item_dimensions.id, moves.id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate entitlements
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_entitlements(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM entitlements
	JOIN orders ON entitlements.id = orders.entitlement_id
	JOIN moves ON orders.id = moves.orders_id
	WHERE moves.id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id,
			moves.id AS move_id,
			NULL AS shipment_id
		FROM audit_history
		JOIN entitlements ON entitlements.id = audit_history.object_id
		JOIN orders ON entitlements.id = orders.entitlement_id
		JOIN moves ON orders.id = moves.orders_id
		WHERE audit_history.table_name = 'entitlements'
		  AND moves.id = p_move_id
		GROUP BY audit_history.id, entitlements.id, moves.id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate payment requests
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_payment_requests(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM payment_requests
	WHERE payment_requests.move_id = p_move_id;

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
			))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id AS move_id
		FROM audit_history
		JOIN payment_requests ON payment_requests.id = audit_history.object_id
		JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
		JOIN mto_service_items ON mto_service_items.id = payment_service_items.mto_service_item_id
		LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		WHERE audit_history.table_name = 'payment_requests'
		  AND payment_requests.move_id = p_move_id
		GROUP BY audit_history.id, payment_requests.id, payment_requests.move_id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate payment service items
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_payment_service_items(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM payment_requests
	JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
	WHERE payment_requests.move_id = p_move_id;

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
				'shipment_locator', mto_shipments.shipment_locator
			))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id AS move_id
		FROM audit_history
		JOIN payment_service_items ON payment_service_items.id = audit_history.object_id
		JOIN payment_requests ON payment_service_items.payment_request_id = payment_requests.id
		JOIN mto_service_items ON mto_service_items.id = payment_service_items.mto_service_item_id
		LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		WHERE audit_history.table_name = 'payment_service_items'
		  AND payment_requests.move_id = p_move_id
		GROUP BY audit_history.id, payment_requests.id, payment_requests.move_id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate proof of service docs
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_proof_of_service_docs(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM proof_of_service_docs
	JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
	WHERE payment_requests.move_id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'payment_request_number', payment_requests.payment_request_number::TEXT
			))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id AS move_id,
			NULL AS shipment_id
		FROM audit_history
		JOIN proof_of_service_docs ON proof_of_service_docs.id = audit_history.object_id
		JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
		WHERE audit_history.table_name = 'proof_of_service_docs'
		  AND payment_requests.move_id = p_move_id
		GROUP BY audit_history.id, proof_of_service_docs.id, payment_requests.move_id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate mto agents
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_mto_agents(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM mto_agents
    JOIN mto_shipments ON mto_agents.mto_shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE mto_shipments.move_id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(jsonb_build_object(
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'shipment_locator', mto_shipments.shipment_locator
      ))::TEXT AS context,
      NULL AS context_id,
      mto_shipments.move_id AS move_id
    FROM audit_history
    JOIN mto_agents ON mto_agents.id = audit_history.object_id
    JOIN mto_shipments ON mto_agents.mto_shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = 'mto_agents'
      AND mto_shipments.move_id = p_move_id
      AND (audit_history.event_name <> 'deleteShipment' OR audit_history.event_name IS NULL)
    GROUP BY audit_history.id, mto_agents.id, mto_shipments.move_id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate reweighs
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_reweighs(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM reweighs
    JOIN mto_shipments ON reweighs.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE mto_shipments.move_id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(jsonb_build_object(
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'payment_request_number', payment_requests.payment_request_number,
        'shipment_locator', mto_shipments.shipment_locator
      ))::TEXT AS context,
      NULL AS context_id,
      mto_shipments.move_id AS move_id
    FROM audit_history
    JOIN reweighs ON reweighs.id = audit_history.object_id
    JOIN mto_shipments ON reweighs.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    LEFT JOIN payment_requests ON mto_shipments.move_id = payment_requests.move_id
    WHERE audit_history.table_name = 'reweighs'
      AND mto_shipments.move_id = p_move_id
    GROUP BY audit_history.id, reweighs.id, mto_shipments.move_id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate service members
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_service_members(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM service_members
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      NULLIF(
        jsonb_agg(jsonb_strip_nulls(
          jsonb_build_object(
            'current_duty_location_name',
            (SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.duty_location_id))
          )
        ))::TEXT,
        '[{}]'::TEXT
      ) AS context,
      NULL AS context_id,
      moves.id AS move_id,
      NULL AS shipment_id
    FROM audit_history
    JOIN service_members ON service_members.id = audit_history.object_id
    JOIN jsonb_to_record(audit_history.changed_data) AS c(duty_location_id TEXT) ON TRUE
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    WHERE audit_history.table_name = 'service_members'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, service_members.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate ppm shipments
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_ppm_shipments(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM ppm_shipments
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(
        jsonb_strip_nulls(
          jsonb_build_object(
            'shipment_type', mto_shipments.shipment_type,
            'shipment_id_abbr', LEFT(ppm_shipments.shipment_id::TEXT, 5),
            'w2_address', (
              SELECT row_to_json(x)
              FROM (SELECT * FROM addresses WHERE addresses.id = CAST(ppm_shipments.w2_address_id AS UUID)) x
            )::TEXT,
            'shipment_locator', mto_shipments.shipment_locator,
            'pickup_postal_address_id', ppm_shipments.pickup_postal_address_id,
            'secondary_pickup_postal_address_id', ppm_shipments.secondary_pickup_postal_address_id
          )
        )
      )::TEXT AS context,
      COALESCE(ppm_shipments.shipment_id::TEXT, NULL)::TEXT AS context_id,
      moves.id AS move_id
    FROM audit_history
    JOIN ppm_shipments ON audit_history.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = 'ppm_shipments'
      AND moves.id = p_move_id
    GROUP BY ppm_shipments.id, audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'destination_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'destination_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - secondary dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_secondary_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - tertiary dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_tertiary_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - secondary pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_secondary_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - tertiary pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_tertiary_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm sec pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_secondary_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm tert pickup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_tertiary_pickup(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

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
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_pickup_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'destination_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      a1.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'destinationAddress',
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'shipment_locator', mto_shipments.shipment_locator
      )))::TEXT AS context,
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'destination_postal_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm secondary dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_secondary_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_destination_postal_address_id')::uuid = a1.object_id
                             AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      a1.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'secondaryDestinationAddress',
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'shipment_locator', mto_shipments.shipment_locator
      )))::TEXT AS context,
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_destination_postal_address_id')::uuid = a1.object_id
                             AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - ppm tertiary dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_ppm_tertiary_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_destination_postal_address_id')::uuid = a1.object_id
                             AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      a1.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'tertiaryDestinationAddress',
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'shipment_locator', mto_shipments.shipment_locator
      )))::TEXT AS context,
      moves.id AS move_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'tertiary_destination_postal_address_id')::uuid = a1.object_id
                             AND a1.table_name = 'addresses'
    JOIN ppm_shipments ON a2.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
    GROUP BY a1.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - service member res
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_service_member_residential(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM addresses
    JOIN service_members ON service_members.residential_address_id = addresses.id
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'residentialAddress'
      )))::TEXT AS context,
      service_members.id::TEXT AS context_id,
      moves.id AS move_id,
      NULL AS shipment_id
    FROM audit_history
    JOIN service_members ON service_members.residential_address_id = audit_history.object_id AND audit_history.table_name = 'addresses'
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    WHERE moves.id = p_move_id
    GROUP BY audit_history.id, moves.id, service_members.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate addresses - service member backup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_service_member_backup_mailing(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM addresses
    JOIN service_members ON service_members.backup_mailing_address_id = addresses.id
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'backupMailingAddress'
      )))::TEXT AS context,
      service_members.id::TEXT AS context_id,
      moves.id AS move_id,
      NULL AS shipment_id
    FROM audit_history
    JOIN service_members ON service_members.backup_mailing_address_id = audit_history.object_id AND audit_history.table_name = 'addresses'
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    WHERE moves.id = p_move_id
    GROUP BY audit_history.id, moves.id, service_members.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - orders
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_orders(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN orders ON orders.uploaded_orders_id = documents.id
              AND documents.service_member_id = orders.service_member_id
  JOIN moves ON orders.id = moves.orders_id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'orders'
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN orders ON orders.uploaded_orders_id = documents.id
               AND documents.service_member_id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - amended orders
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_amended_orders(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN orders ON orders.uploaded_amended_orders_id = documents.id
              AND documents.service_member_id = orders.service_member_id
  JOIN moves ON orders.id = moves.orders_id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'amendedOrders'
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN orders ON orders.uploaded_amended_orders_id = documents.id
               AND documents.service_member_id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - empty weight
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_empty_weight(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN weight_tickets ON weight_tickets.empty_document_id = documents.id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'emptyWeightTicket',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN weight_tickets ON weight_tickets.empty_document_id = documents.id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - full weight
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_full_weight(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'fullWeightTicket',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - trailer weight
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_trailer_weight(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN weight_tickets ON weight_tickets.proof_of_trailer_ownership_document_id = documents.id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'trailerWeightTicket',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN weight_tickets ON weight_tickets.proof_of_trailer_ownership_document_id = documents.id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - pro gear
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_pro_gear(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND progear_weight_tickets.belongs_to_self = true
  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'proGearWeightTicket',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND progear_weight_tickets.belongs_to_self = true
    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - spouse pro gear
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_spouse_pro_gear(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND coalesce(progear_weight_tickets.belongs_to_self, false) = false
  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'spouseProGearWeightTicket',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN progear_weight_tickets ON progear_weight_tickets.document_id = documents.id AND coalesce(progear_weight_tickets.belongs_to_self, false) = false
    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate uploads - expense receipt
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_expense_receipt(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN moving_expenses ON moving_expenses.document_id = documents.id
  JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'expenseReceipt',
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'moving_expense_type', moving_expenses.moving_expense_type::TEXT,
             'shipment_locator', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN moving_expenses ON moving_expenses.document_id = documents.id
    JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate backup contacts
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_backup_contacts(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM backup_contacts
  JOIN service_members ON service_members.id = backup_contacts.service_member_id
  JOIN orders ON orders.service_member_id = service_members.id
  JOIN moves ON moves.orders_id = orders.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           NULL AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN backup_contacts ON backup_contacts.id = audit_history.object_id
    JOIN service_members ON service_members.id = backup_contacts.service_member_id
    JOIN orders ON orders.service_member_id = service_members.id
    JOIN moves ON moves.orders_id = orders.id
    WHERE audit_history.table_name = 'backup_contacts'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate doc review - weights
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_doc_review_weight(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM audit_history
  JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id
    AND audit_history.table_name = 'weight_tickets';

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           )))::TEXT AS context,
           moves.id AS move_id
    FROM audit_history
    JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
      AND audit_history.table_name = 'weight_tickets'
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate doc review - pro gear
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_doc_review_progear(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM audit_history
  JOIN progear_weight_tickets ON progear_weight_tickets.id = audit_history.object_id
  JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id
    AND audit_history.table_name = 'progear_weight_tickets';

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           )))::TEXT AS context,
           moves.id AS move_id
    FROM audit_history
    JOIN progear_weight_tickets ON progear_weight_tickets.id = audit_history.object_id
    JOIN ppm_shipments ON ppm_shipments.id = progear_weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
      AND audit_history.table_name = 'progear_weight_tickets'
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate doc review - expenses
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_doc_review_expenses(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM audit_history
  JOIN moving_expenses ON moving_expenses.id = audit_history.object_id
  JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id
    AND audit_history.table_name = 'moving_expenses';

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'moving_expense_type', moving_expenses.moving_expense_type,
             'shipment_locator', mto_shipments.shipment_locator
           )))::TEXT AS context,
           moves.id AS move_id
    FROM audit_history
    JOIN moving_expenses ON moving_expenses.id = audit_history.object_id
    JOIN ppm_shipments ON ppm_shipments.id = moving_expenses.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
      AND audit_history.table_name = 'moving_expenses'
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate gsr appeals
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_gsr_appeals(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM gsr_appeals
  JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
  LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
  LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
  JOIN moves ON evaluation_reports.move_id = moves.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_build_object(
             'evaluation_report_type', evaluation_reports.type,
             'violation_paragraph_number', pws_violations.paragraph_number,
             'violation_title', pws_violations.title,
             'violation_summary', pws_violations.requirement_summary
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN gsr_appeals ON gsr_appeals.id = audit_history.object_id
    JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
    LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
    LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
    JOIN moves ON evaluation_reports.move_id = moves.id
    WHERE audit_history.table_name = 'gsr_appeals'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ======================================================
-- Sub-function: populate shipment address updates
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_shipment_address_updates(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM shipment_address_updates
  JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_build_object(
             'status', shipment_address_updates.status
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id
    FROM audit_history
    JOIN shipment_address_updates ON shipment_address_updates.id = audit_history.object_id
    JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = 'shipment_address_updates'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate sit extension updates
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_sit_extensions(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM audit_history
  JOIN sit_extensions ON sit_extensions.id = audit_history.object_id
  JOIN mto_shipments ON mto_shipments.id = sit_extensions.mto_shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id
    AND audit_history.table_name = 'sit_extensions';

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
             'shipment_type', mto_shipments.shipment_type,
             'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
             'shipment_locator', mto_shipments.shipment_locator
           )))::TEXT AS context,
           moves.id AS move_id
    FROM audit_history
    JOIN sit_extensions ON sit_extensions.id = audit_history.object_id
    JOIN mto_shipments ON mto_shipments.id = sit_extensions.mto_shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
      AND audit_history.table_name = 'sit_extensions'
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;


-- ============================================
-- ============================================
-- Main Function: fetch_move_history
-- ============================================
-- ============================================
CREATE OR REPLACE FUNCTION public.fetch_move_history (
    move_code text,
    page integer DEFAULT 1,
    per_page integer DEFAULT 20,
    sort text DEFAULT NULL::text,
    sort_direction text DEFAULT NULL::text
)
RETURNS TABLE (
    id uuid,
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
    seq_num int
)
AS $$
DECLARE
    v_move_id UUID;
    offset_value INT := (GREATEST(page, 1) - 1) * GREATEST(per_page, 1);
    v_count INT;
BEGIN
    -- Resolve move ID
    v_move_id := fn_get_move_id(move_code);

    -- Create temp table
    PERFORM fn_create_audit_temp_table();

    -- Populate each section
    PERFORM fn_populate_move_history_moves(v_move_id);
    PERFORM fn_populate_move_history_mto_shipments(v_move_id);
    PERFORM fn_populate_move_history_orders(v_move_id);
    PERFORM fn_populate_move_history_service_items(v_move_id);
    PERFORM fn_populate_mto_agents(v_move_id);
    PERFORM fn_populate_reweighs(v_move_id);
    PERFORM fn_populate_service_members(v_move_id);
    PERFORM fn_populate_ppm_shipments(v_move_id);
    PERFORM fn_populate_addresses_destination(v_move_id);
    PERFORM fn_populate_addresses_secondary_destination(v_move_id);
    PERFORM fn_populate_addresses_tertiary_destination(v_move_id);
    PERFORM fn_populate_addresses_pickup(v_move_id);
    PERFORM fn_populate_addresses_secondary_pickup(v_move_id);
    PERFORM fn_populate_addresses_tertiary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_secondary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_tertiary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_destination(v_move_id);
    PERFORM fn_populate_addresses_ppm_secondary_destination(v_move_id);
    PERFORM fn_populate_addresses_ppm_tertiary_destination(v_move_id);
    PERFORM fn_populate_addresses_service_member_residential(v_move_id);
    PERFORM fn_populate_addresses_service_member_backup_mailing(v_move_id);
    PERFORM fn_populate_uploads_orders(v_move_id);
    PERFORM fn_populate_uploads_amended_orders(v_move_id);
    PERFORM fn_populate_uploads_empty_weight(v_move_id);
    PERFORM fn_populate_uploads_full_weight(v_move_id);
    PERFORM fn_populate_uploads_trailer_weight(v_move_id);
    PERFORM fn_populate_uploads_pro_gear(v_move_id);
    PERFORM fn_populate_uploads_spouse_pro_gear(v_move_id);
    PERFORM fn_populate_uploads_expense_receipt(v_move_id);
    PERFORM fn_populate_backup_contacts(v_move_id);
    PERFORM fn_populate_doc_review_weight(v_move_id);
    PERFORM fn_populate_doc_review_progear(v_move_id);
    PERFORM fn_populate_doc_review_expenses(v_move_id);
    PERFORM fn_populate_gsr_appeals(v_move_id);
    PERFORM fn_populate_shipment_address_updates(v_move_id);
    PERFORM fn_populate_move_history_entitlements(v_move_id);
    PERFORM fn_populate_move_history_proof_of_service_docs(v_move_id);
    PERFORM fn_populate_move_history_payment_service_items(v_move_id);
    PERFORM fn_populate_move_history_payment_requests(v_move_id);
    PERFORM fn_populate_move_history_service_item_dimensions(v_move_id);
    PERFORM fn_populate_move_history_service_item_customer_contacts(v_move_id);
    PERFORM fn_populate_sit_extensions(v_move_id);

    -- adding a CTE here to stop duplicate entries because of duplicate user_id values
    -- with this CTE we get one consolidated row of user details
    RETURN QUERY WITH user_info AS (
      SELECT
        ur.user_id,
        MAX(COALESCE(ou.first_name, prime_user_first_name)) AS first_name,
		MAX(ou.last_name) AS last_name,
		MAX(ou.email) AS email,
		MAX(ou.telephone) AS telephone
      FROM users_roles ur
	  LEFT JOIN roles r ON ur.role_id = r.id
      LEFT JOIN office_users ou ON ou.user_id = ur.user_id
	  LEFT JOIN (
			SELECT 'Prime' AS prime_user_first_name
			) prime_users ON r.role_type = 'prime'
      GROUP BY ur.user_id
    )
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
      COALESCE(ui.first_name, sm.first_name) AS session_user_first_name,
      COALESCE(ui.last_name, sm.last_name) AS session_user_last_name,
      COALESCE(ui.email, sm.personal_email) AS session_user_email,
      COALESCE(ui.telephone, sm.telephone) AS session_user_telephone,
      x.seq_num
    FROM audit_hist_temp x
    LEFT JOIN user_info ui ON ui.user_id = x.session_userid
	LEFT JOIN service_members sm ON sm.user_id = x.session_userid
    ORDER BY x.action_tstamp_tx DESC
    LIMIT per_page OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;
