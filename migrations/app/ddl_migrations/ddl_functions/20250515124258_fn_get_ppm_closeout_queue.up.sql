-- B-23547 - Cameron Jewell - Initial PPM Closeout queue refactor from pop -> proc
--           Cameron Jewell - Remember that Navy, Marines, and Coast Guard have special filters!
--           Cameron Jewell - Their filter is applied in the return
DROP FUNCTION IF EXISTS get_ppm_closeout_queue;

CREATE OR REPLACE FUNCTION get_ppm_closeout_queue(
    user_gbloc                                      TEXT    DEFAULT NULL,
    customer_name_in                                TEXT    DEFAULT NULL,
    edipi                                           TEXT    DEFAULT NULL,
    emplid                                          TEXT    DEFAULT NULL,
    m_status                                        TEXT[]  DEFAULT NULL,
    move_code                                       TEXT    DEFAULT NULL,
    ppm_submitted_at_in                             DATE    DEFAULT NULL,
    branch                                          TEXT    DEFAULT NULL,
    moves_ppm_type_full_or_partial_ppm_filter       TEXT    DEFAULT NULL, -- Moves declaration of FULL or PARTIAL. Not to be confused with ppm_shipments.ppm_type
    origin_duty_location                            TEXT    DEFAULT NULL,
    counseling_office                               TEXT    DEFAULT NULL,
    destination_duty_location                       TEXT    DEFAULT NULL,
    ppm_closeout_location_filter                    TEXT    DEFAULT NULL,
    sc_assigned_user                                TEXT    DEFAULT NULL,
    has_safety_privilege                            BOOLEAN DEFAULT FALSE,
    page                                            INT     DEFAULT 1,
    per_page                                        INT     DEFAULT 20,
    sort                                            TEXT    DEFAULT NULL,
    sort_direction                                  TEXT    DEFAULT NULL
)
RETURNS TABLE(
    id                                  UUID,
    show                                BOOLEAN,
    locator                             TEXT,
    full_or_partial_ppm                 TEXT,
    orders_id                           UUID,
    locked_by                           UUID,
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    lock_expires_at TIMESTAMP WITH TIME ZONE,
    sc_closeout_assigned_id             UUID,
    counseling_transportation_office_id UUID,
    orders                              JSONB,
    ppm_shipments                       JSONB,
    counseling_transportation_office    JSONB,
    ppm_closeout_location               JSONB,
    sc_assigned                         JSONB,
    mto_shipments                       JSONB,
    status                              TEXT,
    total_count                         BIGINT
) AS '
DECLARE
    offset_value INT;
    sort_column  TEXT;
    sort_order   TEXT;
BEGIN
    offset_value := (GREATEST(page,1) - 1) * per_page;

    CASE sort
        WHEN ''locator''                  THEN sort_column := ''locator'';
        WHEN ''status''                   THEN sort_column := ''status'';
        WHEN ''customerName''             THEN sort_column := ''customer_name_out'';
        WHEN ''edipi''                    THEN sort_column := ''edipi_out'';
        WHEN ''emplid''                   THEN sort_column := ''emplid_out'';
        WHEN ''originDutyLocation''       THEN sort_column := ''origin_name'';
        WHEN ''destinationDutyLocation''  THEN sort_column := ''destination_name'';
        WHEN ''counselingOffice''         THEN sort_column := ''counseling_name'';
        WHEN ''ppmType''                  THEN sort_column := ''full_or_partial_ppm'';
        WHEN ''ppmStatus''                THEN sort_column := ''ppm_status'';
        WHEN ''assignedTo''               THEN sort_column := ''counselor_name'';
        WHEN ''closeoutInitiated''        THEN sort_column := ''ppm_submitted_at'';
        WHEN ''closeoutLocation''         THEN sort_column := ''closeout_name'';
        WHEN ''branch''                   THEN sort_column := ''branch_out'';
        ELSE                              sort_column := ''ppm_submitted_at'';
    END CASE;

    sort_order := CASE WHEN sort_direction ILIKE ''desc'' THEN ''DESC'' ELSE ''ASC'' END;

    -- If sort and sort order is null, we default to submitted_at oldest -> latest

    RETURN QUERY EXECUTE format($FMT$
        WITH base AS (
            SELECT
                m.id,
                m.show,
                m.locator::TEXT                       AS locator,
                m.ppm_type::TEXT                      AS full_or_partial_ppm,
                m.locked_by                           AS locked_by,
                m.updated_at,
                m.lock_expires_at,
                m.sc_closeout_assigned_id             AS sc_closeout_assigned_id,
                m.counseling_transportation_office_id AS counseling_transportation_office_id,
                m.status::TEXT                        AS status,
                m.submitted_at::timestamptz           AS move_submitted_at,
                closeout_to.gbloc                     AS closeout_gbloc,
                json_build_object(
                    ''id'', o.id,
                    ''orders_type'', o.orders_type,
                    ''origin_duty_location_gbloc'', o.gbloc,
                    ''department_indicator'', o.department_indicator,
                    ''service_member'', json_build_object(
                        ''id'', sm.id,
                        ''first_name'', sm.first_name,
                        ''last_name'', sm.last_name,
                        ''edipi'', sm.edipi,
                        ''emplid'', sm.emplid,
                        ''affiliation'', sm.affiliation
                    ),
                    ''origin_duty_location'', json_build_object(
                        ''id'',   origin_dl.id,
                        ''name'', origin_dl.name,
                        ''address'', json_build_object(
                            ''street_address_1'', origin_dl_addr.street_address_1,
                            ''street_address_2'', origin_dl_addr.street_address_2,
                            ''city'',             origin_dl_addr.city,
                            ''state'',            origin_dl_addr.state,
                            ''postal_code'',      origin_dl_addr.postal_code
                        )
                    ),
                    ''new_duty_location'', json_build_object(
                        ''id'',   dest_dl.id,
                        ''name'', dest_dl.name,
                        ''address'', json_build_object(
                            ''street_address_1'', dest_dl_addr.street_address_1,
                            ''street_address_2'', dest_dl_addr.street_address_2,
                            ''city'',             dest_dl_addr.city,
                            ''state'',            dest_dl_addr.state,
                            ''postal_code'',      dest_dl_addr.postal_code
                        )
                    )
                )::JSONB AS orders,
                json_build_object(''name'', counseling_to.name, ''ID'', counseling_to.ID)::JSONB AS counseling_transportation_office,
                json_build_object(''name'', closeout_to.name, ''ID'', closeout_to.ID)::JSONB AS ppm_closeout_location,
                json_build_object(''first_name'', sc.first_name, ''last_name'', sc.last_name, ''id'', sc.id)::JSONB AS sc_assigned,
                COALESCE(ms_agg.mto_shipments, ''[]''::jsonb)          AS mto_shipments,
                COALESCE(ppm_agg.ppm_shipments, ''[]''::jsonb)         AS ppm_shipments,
                (sm.first_name || '' '' || sm.last_name)::TEXT  AS customer_name_out,
                origin_dl.name::TEXT                  AS origin_name,
                dest_dl.name::TEXT                    AS destination_name,
                counseling_to.name::TEXT              AS counseling_name,
                closeout_to.name::TEXT                AS closeout_name,
                (sc.first_name || '' '' || sc.last_name)::TEXT AS counselor_name,
                ppm_agg.latest_ppm_status::TEXT       AS ppm_status,
                ppm_agg.earliest_ppm_submitted_at::timestamptz AS ppm_submitted_at,
                sm.edipi                              AS edipi_out,
                sm.emplid                             AS emplid_out,
                sm.affiliation                        AS branch_out
            FROM moves m
            JOIN orders o                   ON o.id = m.orders_id
            JOIN service_members sm        ON sm.id = o.service_member_id
            JOIN duty_locations origin_dl  ON origin_dl.id = o.origin_duty_location_id
            LEFT JOIN duty_locations dest_dl ON dest_dl.id = o.new_duty_location_id
            LEFT JOIN transportation_offices counseling_to ON counseling_to.id = m.counseling_transportation_office_id
            LEFT JOIN transportation_offices closeout_to   ON closeout_to.id = m.closeout_office_id
            LEFT JOIN office_users sc       ON sc.id = m.sc_closeout_assigned_id
            LEFT JOIN addresses AS origin_dl_addr ON origin_dl.address_id = origin_dl_addr.id
            LEFT JOIN addresses AS dest_dl_addr ON dest_dl.address_id = dest_dl_addr.id
            LEFT JOIN LATERAL (
                SELECT json_agg(
                        json_build_object(
                            ''id'',     ms.id,
                            ''status'', ms.status
                        )
                        )::JSONB AS mto_shipments
                FROM   mto_shipments ms
                WHERE  ms.move_id = m.id
                ) ms_agg ON TRUE
            LEFT JOIN LATERAL (
                SELECT
                    json_agg(
                        json_build_object(
                        ''id'',           ps.id,
                        ''shipment_id'',  ps.shipment_id,
                        ''status'',       ps.status,
                        ''submitted_at'', ps.submitted_at
                        )
                    )::JSONB AS ppm_shipments,
                    MIN(ps.submitted_at) AS earliest_ppm_submitted_at,
                    MAX(ps.status) AS latest_ppm_status
                FROM mto_shipments ms2
                JOIN ppm_shipments ps ON ps.shipment_id = ms2.id
                WHERE ms2.move_id = m.id
                    -- Currently ps.status = ''NEEDS_CLOSEOUT'' is the only status that you can see in the queue
                    AND  ps.status   = ''NEEDS_CLOSEOUT''
                ) ppm_agg ON TRUE
                -- Filter out move entries that do not have a PPM with a closeout initiated at value (ppm submitted at = closeout initiated)
            WHERE m.show = TRUE AND ppm_agg.earliest_ppm_submitted_at IS NOT NULL
        ),
        filtered AS (
            SELECT * FROM base WHERE
              (
                $1 IS NULL
                -- Default query, match the office user GBLOC to the move closeout GBLOC
                -- unless the user GBLOC is of NAVY, TVCB (Marines), or USCG. These
                -- will use special filters
                OR (closeout_gbloc = $1
                    AND $1 NOT IN (''NAVY'', ''TVCB'',''USCG'')
                    AND branch_out NOT IN (''NAVY'',''MARINES'',''COAST_GUARD'')) -- Do not let these branches show up anywhere else outside of their dedicated PPM closeout offices
                -- Special closeout GBLOC handling
                OR ($1 = ''NAVY'' AND branch_out = ''NAVY'')        -- Navy -> Navy moves
                OR ($1 = ''TVCB'' AND branch_out = ''MARINES'')     -- TVCB -> Marine moves (USMC GBLOC is used for HHG moves. USMC uses TVCB closeout GBLOC explicitly for PPMs)
                OR ($1 = ''USCG'' AND branch_out = ''COAST_GUARD'') -- Coast Guard -> Coast Guard moves
                -- Again, USMC GBLOC does not go here! USMC is for HHG Marine moves, TVCB for PPM Marine moves
              )
              AND ($2  IS NULL OR customer_name_out ILIKE ''%%'' || $2 || ''%%'')
              AND ($3  IS NULL OR edipi_out   = $3)
              AND ($4  IS NULL OR emplid_out  = $4)
              AND ($5  IS NULL OR status = ANY($5))
              AND ($6  IS NULL OR locator ILIKE $6 || ''%%'')
              -- Loop over the ppm shipments agg and find the filter for submitted_at
              AND (
                $7 IS NULL OR
                EXISTS (
                SELECT 1
                FROM   jsonb_array_elements(ppm_shipments) elem
                WHERE  (elem->>''submitted_at'')::date = $7
                )
              )
              AND ($8  IS NULL OR branch_out = $8)
              AND ($9  IS NULL OR full_or_partial_ppm   = $9)
              AND ($10 IS NULL OR origin_name ILIKE ''%%'' || $10 || ''%%'')
              AND ($11 IS NULL OR counseling_name ILIKE ''%%'' || $11 || ''%%'')
              AND ($12 IS NULL OR destination_name ILIKE ''%%'' || $12 || ''%%'')
              AND ($13 IS NULL OR closeout_name ILIKE ''%%'' || $13 || ''%%'')
              AND ($14 OR (orders->>''orders_type'') <> ''SAFETY'')
              AND ($15 IS NULL OR counselor_name ILIKE ''%%'' || $15 || ''%%'')
        )
        SELECT
                id::UUID,
                show::BOOLEAN,
                locator::TEXT,
                full_or_partial_ppm::TEXT,
                (orders->>''id'')::UUID AS orders_id,
                locked_by::UUID,
                updated_at,
                lock_expires_at,
                sc_closeout_assigned_id::UUID,
                counseling_transportation_office_id::UUID,
                orders::JSONB,
                ppm_shipments::JSONB,
                counseling_transportation_office::JSONB,
                ppm_closeout_location::JSONB,
                sc_assigned::JSONB,
                mto_shipments::JSONB,
                status::TEXT,
                COUNT(*) OVER() AS total_count
        FROM   filtered
        ORDER  BY %I %s, id
        LIMIT  $16 OFFSET $17
    $FMT$, sort_column, sort_order)
    USING
      user_gbloc, customer_name_in, edipi, emplid, m_status, move_code,
      ppm_submitted_at_in, branch, moves_ppm_type_full_or_partial_ppm_filter,
      origin_duty_location, counseling_office, destination_duty_location,
      ppm_closeout_location_filter, has_safety_privilege, sc_assigned_user,
      per_page, offset_value;
END;
'
LANGUAGE plpgsql;