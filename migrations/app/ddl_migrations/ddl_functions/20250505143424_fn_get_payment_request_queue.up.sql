-- B-22543  Daniel Jordan  initial function creation
DROP FUNCTION IF EXISTS get_payment_request_queue(
  TEXT, TEXT, TEXT, TEXT, TEXT, TEXT, TEXT,
  payment_request_status[], TIMESTAMP, TEXT, TEXT,
  BOOLEAN, INTEGER, INTEGER, TEXT, TEXT
);

CREATE OR REPLACE FUNCTION get_payment_request_queue(
  user_gbloc                TEXT                     DEFAULT NULL,
  branch                    TEXT                     DEFAULT NULL,
  locator                   TEXT                     DEFAULT NULL,
  edipi                     TEXT                     DEFAULT NULL,
  emplid                    TEXT                     DEFAULT NULL,
  customer_name             TEXT                     DEFAULT NULL,
  origin_duty_location      TEXT                     DEFAULT NULL,
  status                    payment_request_status[] DEFAULT NULL,
  submitted_at              TIMESTAMP                DEFAULT NULL,
  tio_assigned_user         TEXT                     DEFAULT NULL,
  p_counseling_office       TEXT                     DEFAULT NULL,
  has_safety_privilege      BOOLEAN                  DEFAULT FALSE,
  page                      INTEGER                  DEFAULT 1,
  per_page                  INTEGER                  DEFAULT 20,
  sort                      TEXT                     DEFAULT NULL,
  sort_direction            TEXT                     DEFAULT NULL
)
RETURNS TABLE (
  payment_request   JSONB,
  move              JSONB,
  orders            JSONB,
  origin_to_office  JSONB,
  tio_user          JSONB,
  counseling_office JSONB,
  total_count       BIGINT
) LANGUAGE plpgsql AS $$
DECLARE
  sql_query   TEXT;
  offset_val  INT;
  sort_col    TEXT := 'pr.id';
  sort_ord    TEXT := 'ASC';
BEGIN
  page     := COALESCE(page,     1);
  per_page := COALESCE(per_page, 20);
  IF page     < 1 THEN page     := 1;  END IF;
  IF per_page < 1 THEN per_page := 20; END IF;
  offset_val := (page - 1) * per_page;

  IF sort IS NOT NULL THEN
    CASE sort
      WHEN 'locator'                 THEN sort_col := 'm.locator';
      WHEN 'submittedAt'             THEN sort_col := 'pr.created_at';
      WHEN 'branch'                  THEN sort_col := 'sm.affiliation';
      WHEN 'customerName'            THEN sort_col := 'sm.last_name';
      WHEN 'edipi'                   THEN sort_col := 'sm.edipi';
      WHEN 'emplid'                  THEN sort_col := 'sm.emplid';
      WHEN 'originDutyLocation'      THEN sort_col := 'origin_dl.name';
      WHEN 'assignedTo'              THEN sort_col := 'tio.first_name';
      WHEN 'counselingOffice'        THEN sort_col := 'co.name';
      WHEN 'status'                  THEN sort_col := 'pr.status';
      WHEN 'age'                     THEN sort_col := 'pr.created_at';
      ELSE NULL;
    END CASE;
  END IF;

  IF sort_direction IS NOT NULL AND lower(sort_direction) = 'desc' THEN
    sort_ord := 'DESC';
  END IF;

  sql_query := '
    SELECT
      jsonb_build_object(
        ''id'',                                   pr.id,
        ''is_final'',                             pr.is_final,
        ''rejection_reason'',                     pr.rejection_reason,
        ''created_at'',                           to_char(pr.created_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''updated_at'',                           to_char(pr.updated_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''move_id'',                              pr.move_id,
        ''status'',                               pr.status,
        ''requested_at'',                         to_char(pr.requested_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''reviewed_at'',                          to_char(pr.reviewed_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''sent_to_gex_at'',                       to_char(pr.sent_to_gex_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''received_by_gex_at'',                   to_char(pr.received_by_gex_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''paid_at'',                              to_char(pr.paid_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''payment_request_number'',               pr.payment_request_number,
        ''sequence_number'',                      pr.sequence_number,
        ''recalculation_of_payment_request_id'',  pr.recalculation_of_payment_request_id
      )::JSONB AS payment_request,
      jsonb_build_object(
        ''id'',             m.id,
        ''locator'',        m.locator,
        ''submitted_at'',   to_char(m.submitted_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''locked_by'',      m.locked_by,
        ''lock_expires_at'',to_char(m.lock_expires_at, ''YYYY-MM-DD"T"HH24:MI:SS.US"Z"''),
        ''ShipmentGBLOC'',
          COALESCE(
            jsonb_build_array(
              jsonb_build_object(''gbloc'', mtg.gbloc)
            ),
            ''[]''::JSONB
          )
      )::jsonb AS move,
      json_build_object(
        ''id'',                         o.id,
        ''orders_type'',                o.orders_type,
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
            ''id'',   origin_dl.id,
            ''name'', origin_dl.name,
            ''address'', json_build_object(
                ''street_address_1'', addr.street_address_1,
                ''street_address_2'', addr.street_address_2,
                ''city'',           addr.city,
                ''state'',          addr.state,
                ''postal_code'',     addr.postal_code
            )
        )
      )::JSONB AS orders,
      json_build_object(
        ''id'',   origin_to.id,
        ''name'', origin_to.name
      )::JSONB AS origin_to_office,
      CASE
        WHEN tio.id IS NULL THEN
            ''null''::JSONB
        ELSE
            json_build_object(
            ''id'',         tio.id,
            ''first_name'', tio.first_name,
            ''last_name'',  tio.last_name
            )::JSONB
      END AS tio_user,
      json_build_object(
        ''id'',   co.id,
        ''name'', co.name
      )::JSONB AS counseling_office,
      COUNT(*) OVER() AS total_count
    FROM payment_requests pr
    JOIN moves m     ON pr.move_id = m.id
    JOIN orders o    ON m.orders_id = o.id
    JOIN service_members sm ON o.service_member_id = sm.id
    JOIN duty_locations origin_dl
         ON o.origin_duty_location_id = origin_dl.id
    LEFT JOIN transportation_offices origin_to
         ON origin_dl.transportation_office_id = origin_to.id
    LEFT JOIN move_to_gbloc mtg
         ON mtg.move_id = m.id
    LEFT JOIN office_users tio
         ON m.tio_payment_request_assigned_id = tio.id
    LEFT JOIN transportation_offices co
         ON m.counseling_transportation_office_id = co.id
    LEFT JOIN addresses addr
        ON origin_dl.address_id = addr.id
    WHERE m.show = TRUE
  ';

  IF user_gbloc IS NOT NULL THEN
    sql_query := sql_query || ' AND mtg.gbloc = $1';
  END IF;

  IF branch IS NOT NULL THEN
    sql_query := sql_query || ' AND sm.affiliation ILIKE ''%''||$2||''%''';
    ELSE IF user_gbloc <> 'USMC' THEN
        sql_query := sql_query || ' AND sm.affiliation <> ''MARINES''';
    END IF;
  END IF;

  IF locator IS NOT NULL THEN
    sql_query := sql_query || ' AND m.locator ILIKE ''%'' || $3 || ''%''';
  END IF;

  IF edipi IS NOT NULL THEN
    sql_query := sql_query || ' AND sm.edipi ILIKE ''%'' || $4 || ''%''';
  END IF;

  IF emplid IS NOT NULL THEN
    sql_query := sql_query || ' AND sm.emplid ILIKE ''%'' || $5 || ''%''';
  END IF;

  IF customer_name IS NOT NULL THEN
    -- strip commas and search both “last first” and “first last”
    sql_query := sql_query || '
      AND (
        (sm.last_name || '' '' || sm.first_name) ILIKE ''%'' ||
          replace($6, '','', '''') || ''%''
     OR (sm.first_name || '' '' || sm.last_name) ILIKE ''%'' ||
          replace($6, '','', '''') || ''%''
      )
    ';
  END IF;

  IF origin_duty_location IS NOT NULL THEN
    sql_query := sql_query || ' AND origin_dl.name ILIKE ''%'' || $7 || ''%''';
  END IF;

  IF status IS NOT NULL THEN
    sql_query := sql_query || ' AND pr.status = ANY($8)';
  END IF;

  IF submitted_at IS NOT NULL THEN
    sql_query := sql_query || ' AND pr.created_at::DATE = $9::DATE';
  END IF;

  IF tio_assigned_user IS NOT NULL THEN
    sql_query := sql_query || ' AND (tio.first_name || '' '' || tio.last_name) ILIKE ''%'' || $10 || ''%''';
  END IF;

  IF p_counseling_office IS NOT NULL THEN
    sql_query := sql_query || ' AND co.name ILIKE ''%'' || $11 || ''%''';
  END IF;

  IF NOT has_safety_privilege THEN
    sql_query := sql_query || ' AND o.orders_type != ''SAFETY''';
  END IF;

  IF sort IS NULL THEN
    sql_query := sql_query || ' ORDER BY pr.created_at ASC';
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
        sql_query := sql_query || format(' ORDER BY pr.status %s', sort_ord);
      WHEN 'submittedAt' THEN
        sql_query := sql_query || format(' ORDER BY pr.created_at %s', sort_ord);
      WHEN 'branch' THEN
        sql_query := sql_query || format(' ORDER BY sm.affiliation %s', sort_ord);
      WHEN 'originDutyLocation' THEN
        sql_query := sql_query || format(' ORDER BY origin_dl.name %s', sort_ord);
      WHEN 'assignedTo' THEN
        sql_query := sql_query || format(' ORDER BY tio.first_name %s', sort_ord);
      WHEN 'counselingOffice' THEN
        sql_query := sql_query || format(' ORDER BY co.name %s', sort_ord);
      WHEN 'age' THEN
        sql_query := sql_query
          || format(
            ' ORDER BY pr.created_at %s',
            CASE WHEN lower(sort_ord) = 'asc' THEN 'DESC' ELSE 'ASC' END
          );
      ELSE
        sql_query := sql_query || format(' ORDER BY %s %s', sort_col, sort_ord);
    END CASE;
  END IF;

  IF sort IS NULL OR sort <> 'locator' THEN
    sql_query := sql_query || ', m.locator ASC';
  END IF;

  sql_query := sql_query || format(' LIMIT %s OFFSET %s', per_page, offset_val);

  RETURN QUERY EXECUTE sql_query
  USING
    user_gbloc,                -- $1
    branch,                    -- $2
    locator,                   -- $3
    edipi,                     -- $4
    emplid,                    -- $5
    customer_name,             -- $6
    origin_duty_location,      -- $7
    status,                    -- $8
    submitted_at,              -- $9
    tio_assigned_user,         -- $10
    p_counseling_office;       -- $11
END;
$$;
