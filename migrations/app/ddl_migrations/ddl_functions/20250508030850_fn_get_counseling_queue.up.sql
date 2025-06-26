-- B-23520  Tevin Adams  initial function creation of get_counseling_queue

DROP FUNCTION IF EXISTS get_counseling_queue;
CREATE OR REPLACE FUNCTION get_counseling_queue(
	user_gbloc TEXT DEFAULT NULL,
	customer_name TEXT DEFAULT NULL,
	edipi TEXT DEFAULT NULL,
	emplid TEXT DEFAULT NULL,
	m_status TEXT[] DEFAULT NULL,
	move_code TEXT DEFAULT NULL,
	requested_move_date TIMESTAMP DEFAULT NULL,
	date_submitted TIMESTAMP DEFAULT NULL,
	branch TEXT DEFAULT NULL,
	origin_duty_location TEXT DEFAULT NULL,
	counseling_office TEXT DEFAULT NULL,
	sc_assigned_user TEXT DEFAULT NULL,
	has_safety_privilege BOOLEAN DEFAULT FALSE,
	page INTEGER DEFAULT 1,
	per_page INTEGER DEFAULT 20,
	sort TEXT DEFAULT NULL,
	sort_direction TEXT DEFAULT NULL
)
RETURNS TABLE (
	  id UUID,
    show BOOLEAN,
    locator TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE,
    orders_id UUID,
    status TEXT,
    locked_by UUID,
    sc_counseling_assigned_id UUID,
    counseling_transportation_office_id UUID,
    sc_assigned JSONB,
    counseling_transportation_office JSONB,
    orders JSONB,
    mto_shipments JSONB,
    total_count BIGINT,
    mtos_earliest_requested_pickup_date TIMESTAMP WITH TIME zone,
    mtos_earliest_requested_delivery_date TIMESTAMP WITH TIME zone,
    ppms_earliest_expected_departure_date TIMESTAMP WITH TIME zone,
    ppm_shipments JSONB,
    updated_at TIMESTAMP WITH TIME ZONE
) LANGUAGE plpgsql AS $$
DECLARE
  sql_query 	TEXT;
  offset_val 	INT;
  sort_ord    	TEXT := 'ASC';
  sort_col 		TEXT := 'm.submitted_at';
BEGIN
  page     := COALESCE(page,     1);
  per_page := COALESCE(per_page, 20);
  IF page     < 1 THEN page     := 1;  END IF;
  IF per_page < 1 THEN per_page := 20; END IF;
  offset_val := (page - 1) * per_page;

  sql_query := '
    SELECT
    m.id AS id,
    m.show AS show,
    m.locator::TEXT AS locator,
    m.submitted_at::TIMESTAMP WITH TIME ZONE AS submitted_at,
    m.orders_id AS orders_id,
    m.status::TEXT AS status,
    m.locked_by AS locked_by,
    m.sc_counseling_assigned_id AS sc_counseling_assigned_id,
    m.counseling_transportation_office_id AS counseling_transportation_office_id,
		json_build_object(
				''id'', sc_user.id,
				''first_name'', sc_user.first_name,
				''last_name'', sc_user.last_name
		)::JSONB AS sc_assigned,
		json_build_object(
				''id'', co.id,
				''name'', co.name
		)::JSONB AS counseling_transportation_office,
		json_build_object(
				''id'', o.id,
				''orders_type'', o.orders_type,
				''origin_duty_location_gbloc'', o.gbloc,
				''department_indicator'',       o.department_indicator,
				''service_member'', json_build_object(
					''id'',           sm.id,
					''first_name'',   sm.first_name,
          			''last_name'',    sm.last_name,
          			''edipi'',        sm.edipi,
          			''emplid'',       sm.emplid,
          			''affiliation'',  sm.affiliation
				),
				''origin_duty_location'', json_build_object(
					''id'',   origin_duty_locations.id,
					''name'', origin_duty_locations.name,
					''address'', json_build_object(
						''street_address_1'', addr.street_address_1,
                  		''street_address_2'', addr.street_address_2,
                  		''city'',           addr.city,
                  		''state'',          addr.state,
                  		''postal_code'',    addr.postal_code
					)
				)
		  )::JSONB AS orders,
      COALESCE(ms_agg.mto_shipments)::JSONB AS mto_shipments,
	  COUNT(*) OVER() AS total_count,
	  ms_agg.mtos_earliest_requested_pickup_date::timestamptz AS mtos_earliest_requested_pickup_date,
	  ms_agg.mtos_earliest_requested_delivery_date::timestamptz AS mtos_earliest_requested_delivery_date,
	  ppm_agg.ppms_earliest_expected_departure_date::timestamptz AS ppms_earliest_expected_departure_date,
    COALESCE(ppm_agg.ppm_shipments, ''[]'')::JSONB AS ppm_shipments,
    m.updated_at::TIMESTAMP WITH TIME ZONE AS updated_at
    From moves m
    JOIN orders o ON m.orders_id = o.id
    JOIN service_members sm ON o.service_member_id = sm.id
    JOIN duty_locations AS origin_duty_locations ON o.origin_duty_location_id = origin_duty_locations.id
    LEFT JOIN transportation_offices co
          ON m.counseling_transportation_office_id = co.id
    LEFT JOIN addresses addr
          ON origin_duty_locations.address_id = addr.id
    LEFT JOIN office_users AS sc_user ON m.sc_counseling_assigned_id = sc_user.id
 	LEFT JOIN LATERAL (
                SELECT json_agg(
                       	json_build_object(
							''id'', ms.id,
							''shipment_type'', ms.shipment_type,
                            ''status'', ms.status,
                            ''requested_pickup_date'', TO_CHAR(ms.requested_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''scheduled_pickup_date'', TO_CHAR(ms.scheduled_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''approved_date'', TO_CHAR(ms.approved_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''prime_estimated_weight'', ms.prime_estimated_weight
						)
                       )::JSONB AS mto_shipments,
						MIN(ms.requested_pickup_date) AS mtos_earliest_requested_pickup_date,
            MIN(ms.requested_delivery_date) AS mtos_earliest_requested_delivery_date
                FROM   mto_shipments ms
                WHERE  ms.move_id = m.id
                ) ms_agg ON TRUE
	LEFT JOIN LATERAL (
				SELECT json_agg(
					json_build_object(
						''id'', ms3.id,
						''shipment_type'', ms3.shipment_type,
            			''expected_departure_date'', TO_CHAR(ps2.expected_departure_date, ''YYYY-MM-DD"T00:00:00Z"''),
						''shipment_id'', ps2.shipment_id
					)
				) ::JSONB AS ppm_shipments,
        MIN(ps2.expected_departure_date) AS ppms_earliest_expected_departure_date
				FROM mto_shipments ms3
                JOIN ppm_shipments ps2 ON ps2.shipment_id = ms3.id
                WHERE  ms3.move_id = m.id
	) ppm_agg ON TRUE
    JOIN move_to_gbloc ON move_to_gbloc.move_id = m.id
    WHERE m.show = TRUE
  ';

  IF user_gbloc IS NOT NULL THEN
          sql_query := sql_query || ' AND move_to_gbloc.gbloc = $1 ';
  END IF;
  IF customer_name IS NOT NULL THEN
          sql_query := sql_query || ' AND (
              sm.first_name || '' '' || sm.last_name ILIKE ''%'' || $2 || ''%''
              OR sm.last_name || '' '' || sm.first_name ILIKE ''%'' || $2 || ''%''
          )';
  END IF;

  IF edipi IS NOT NULL THEN
    sql_query := sql_query || ' AND sm.edipi ILIKE ''%'' || $3 || ''%'' ';
  END IF;

  IF emplid IS NOT NULL THEN
      sql_query := sql_query || ' AND sm.emplid ILIKE ''%'' || $4 || ''%'' ';
  END IF;
--
  IF m_status IS NOT NULL THEN
    sql_query := sql_query || ' AND m.status = ANY($5) ';
  ELSE
	sql_query := sql_query || ' AND m.status = ''NEEDS SERVICE COUNSELING''';
  END IF;

  IF move_code IS NOT NULL THEN
    sql_query := sql_query || ' AND m.locator ILIKE ''%'' || $6 || ''%'' ';
  END IF;

  IF requested_move_date IS NOT NULL THEN
        sql_query := sql_query || '
        AND EXISTS (
            SELECT 1
            FROM mto_shipments ms2
            LEFT JOIN ppm_shipments ppm2 ON ppm2.shipment_id = ms2.id
            WHERE ms2.move_id = m.id
              AND (
                   ms2.requested_pickup_date::DATE = $7::DATE
                    OR ppm2.expected_departure_date::DATE = $7::DATE
                    OR (ms2.shipment_type = ''HHG_OUTOF_NTS''
                        AND ms2.requested_delivery_date::DATE = $7::DATE)
              )
        )';
    END IF;

  IF date_submitted IS NOT NULL THEN
      sql_query := sql_query || ' AND m.submitted_at::DATE = $8::DATE ';
  END IF;

  IF branch IS NOT NULL THEN
    sql_query := sql_query || ' AND sm.affiliation ILIKE ''%'' || $9 || ''%'' ';
  END IF;

  IF origin_duty_location IS NOT NULL THEN
    sql_query := sql_query || ' AND origin_duty_locations.name ILIKE ''%'' || $10 || ''%'' ';
  END IF;

  IF counseling_office IS NOT NULL THEN
    sql_query := sql_query || ' AND co.name ILIKE ''%'' || $11 || ''%'' ';
  END IF;

  IF sc_assigned_user IS NOT NULL THEN
    sql_query := sql_query || ' AND (sc_user.first_name || '' '' || sc_user.last_name) ILIKE ''%'' || $12 || ''%'' ';
  END IF;

  IF has_safety_privilege IS NOT TRUE THEN
      sql_query := sql_query || ' AND o.orders_type != ''SAFETY'' ';
  END IF;

  IF sort_direction IS NOT NULL AND lower(sort_direction) = 'desc' THEN
	sort_ord = 'desc';
  END IF;

  IF sort IS NULL THEN
    sql_query := sql_query || format(' ORDER BY submitted_at %s', sort_ord);
  ELSE
    CASE sort
      WHEN 'locator' THEN
        sql_query := sql_query || format(' ORDER BY m.locator %s', sort_ord);
      WHEN 'lastName', 'customerName' THEN
        sql_query := sql_query
          || format(' ORDER BY sm.last_name %1$s, sm.first_name %1$s', sort_ord);
      WHEN 'edipi' THEN
        sql_query := sql_query || format(' ORDER BY sm.edipi %s', sort_ord);
      WHEN 'emplid' THEN
        sql_query := sql_query || format(' ORDER BY sm.emplid %s', sort_ord);
      WHEN 'status' THEN
        sql_query := sql_query || format(' ORDER BY m.status %s', sort_ord);
      WHEN 'submittedAt' THEN
        sql_query := sql_query || format(' ORDER BY m.submitted_at %s', sort_ord);
      WHEN 'branch' THEN
        sql_query := sql_query || format(' ORDER BY sm.affiliation %s', sort_ord);
      WHEN 'originDutyLocation' THEN
        sql_query := sql_query || format(' ORDER BY origin_duty_locations.name %s', sort_ord);
      WHEN 'assignedTo' THEN
        sql_query := sql_query || format(' ORDER BY sc_user.last_name %s', sort_ord);
      WHEN 'counselingOffice' THEN
        sql_query := sql_query || format(' ORDER BY co.name %s', sort_ord);
      WHEN 'requestedMoveDates' THEN
        sql_query := sql_query || format(
          ' ORDER BY COALESCE('
          || 'ms_agg.mtos_earliest_requested_pickup_date,'
          || 'ppm_agg.ppms_earliest_expected_departure_date,'
          || 'ms_agg.mtos_earliest_requested_delivery_date'
          || ') %s',
          sort_ord
        );
          ELSE
        sql_query := sql_query || format(' ORDER BY %s %s', sort_col, sort_ord);
    END CASE;
  END IF;

 sql_query := sql_query || format(' LIMIT %s OFFSET %s', per_page, offset_val);

  RETURN QUERY EXECUTE sql_query
    USING
      user_gbloc,
      customer_name,
      edipi,
      emplid,
      m_status,
      move_code,
      requested_move_date,
      date_submitted,
      branch,
      origin_duty_location,
      counseling_office,
      sc_assigned_user;
END;
$$;