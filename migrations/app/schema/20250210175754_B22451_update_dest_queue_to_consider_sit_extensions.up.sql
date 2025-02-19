-- updating to consider sit extensions in PENDING status
CREATE OR REPLACE FUNCTION get_destination_queue(
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
    too_assigned_user TEXT DEFAULT NULL,
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
    too_assigned_id UUID,
    counseling_transportation_office_id UUID,
    orders JSONB,
    mto_shipments JSONB,
    counseling_transportation_office JSONB,
    too_assigned JSONB,
    total_count BIGINT
) AS $$
DECLARE
    sql_query TEXT;
    offset_value INTEGER;
    sort_column TEXT;
    sort_order TEXT;
BEGIN
    IF page < 1 THEN
        page := 1;
    END IF;

    IF per_page < 1 THEN
        per_page := 20;
    END IF;

    offset_value := (page - 1) * per_page;

    sql_query := '
        SELECT
            moves.id AS id,
            moves.show AS show,
            moves.locator::TEXT AS locator,
            moves.submitted_at::TIMESTAMP WITH TIME ZONE AS submitted_at,
            moves.orders_id AS orders_id,
            moves.status::TEXT AS status,
            moves.locked_by AS locked_by,
            moves.too_assigned_id AS too_assigned_id,
            moves.counseling_transportation_office_id AS counseling_transportation_office_id,
            json_build_object(
                ''id'', orders.id,
                ''origin_duty_location_gbloc'', orders.gbloc,
                ''service_member'', json_build_object(
                    ''id'', service_members.id,
                    ''first_name'', service_members.first_name,
                    ''last_name'', service_members.last_name,
                    ''edipi'', service_members.edipi,
                    ''emplid'', service_members.emplid,
                    ''affiliation'', service_members.affiliation
                ),
                ''origin_duty_location'', json_build_object(
                    ''name'', origin_duty_locations.name
                )
            )::JSONB AS orders,
            COALESCE(
                (
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
                    )
                    FROM (
                        SELECT DISTINCT ON (mto_shipments.id) mto_shipments.*
                        FROM mto_shipments
                        WHERE mto_shipments.move_id = moves.id
                    ) AS ms
                ),
                ''[]''
            )::JSONB AS mto_shipments,
            json_build_object(
                ''name'', counseling_offices.name
            )::JSONB AS counseling_transportation_office,
            json_build_object(
                ''first_name'', too_user.first_name,
                ''last_name'', too_user.last_name
            )::JSONB AS too_assigned,
            COUNT(*) OVER() AS total_count
        FROM moves
        JOIN orders ON moves.orders_id = orders.id
        JOIN mto_shipments ON mto_shipments.move_id = moves.id
        LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
        JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
        JOIN re_services ON mto_service_items.re_service_id = re_services.id
        JOIN service_members ON orders.service_member_id = service_members.id
        JOIN duty_locations AS new_duty_locations ON orders.new_duty_location_id = new_duty_locations.id
        JOIN duty_locations AS origin_duty_locations ON orders.origin_duty_location_id = origin_duty_locations.id
        LEFT JOIN office_users AS too_user ON moves.too_assigned_id = too_user.id
        LEFT JOIN office_users AS locked_user ON moves.locked_by = locked_user.id
        LEFT JOIN transportation_offices AS counseling_offices
            ON moves.counseling_transportation_office_id = counseling_offices.id
        LEFT JOIN shipment_address_updates ON shipment_address_updates.shipment_id = mto_shipments.id
        LEFT JOIN sit_extensions ON sit_extensions.mto_shipment_id = mto_shipments.id
        JOIN move_to_dest_gbloc ON move_to_dest_gbloc.move_id = moves.id
        WHERE moves.show = TRUE
    ';

    IF user_gbloc IS NOT NULL THEN
        sql_query := sql_query || ' AND move_to_dest_gbloc.gbloc = $1 ';
    END IF;

    IF customer_name IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            service_members.first_name || '' '' || service_members.last_name ILIKE ''%'' || $2 || ''%''
            OR service_members.last_name || '' '' || service_members.first_name ILIKE ''%'' || $2 || ''%''
        )';
    END IF;

    IF edipi IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.edipi ILIKE ''%'' || $3 || ''%'' ';
    END IF;

    IF emplid IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.emplid ILIKE ''%'' || $4 || ''%'' ';
    END IF;

    IF m_status IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.status = ANY($5) ';
    END IF;

    IF move_code IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.locator ILIKE ''%'' || $6 || ''%'' ';
    END IF;

    IF requested_move_date IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            mto_shipments.requested_pickup_date::DATE = $7::DATE
            OR ppm_shipments.expected_departure_date::DATE = $7::DATE
            OR (mto_shipments.shipment_type = ''HHG_OUTOF_NTS'' AND mto_shipments.requested_delivery_date::DATE = $7::DATE)
        )';
    END IF;

    IF date_submitted IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.submitted_at::DATE = $8::DATE ';
    END IF;

    IF branch IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.affiliation ILIKE ''%'' || $9 || ''%'' ';
    END IF;

    IF origin_duty_location IS NOT NULL THEN
        sql_query := sql_query || ' AND origin_duty_locations.name ILIKE ''%'' || $10 || ''%'' ';
    END IF;

    IF counseling_office IS NOT NULL THEN
        sql_query := sql_query || ' AND counseling_offices.name ILIKE ''%'' || $11 || ''%'' ';
    END IF;

    IF too_assigned_user IS NOT NULL THEN
        sql_query := sql_query || ' AND (too_user.first_name || '' '' || too_user.last_name) ILIKE ''%'' || $12 || ''%'' ';
    END IF;

    -- add destination queue-specific filters (pending dest address requests, pending dest SIT extension requests, dest SIT & dest shuttle service items)
    sql_query := sql_query || '
        AND (
            shipment_address_updates.status = ''REQUESTED''
            OR sit_extensions.status = ''PENDING''
            OR (
                mto_service_items.status = ''SUBMITTED''
                AND re_services.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSFSC'', ''DDSHUT'', ''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSFSC'', ''IDSHUT'')
            )
        )
    ';

    -- default sorting values if none are provided (move.id)
    sort_column := 'id';
    sort_order := 'ASC';

    IF sort IS NOT NULL THEN
        CASE sort
            WHEN 'locator' THEN sort_column := 'moves.locator';
            WHEN 'status' THEN sort_column := 'moves.status';
            WHEN 'customerName' THEN sort_column := 'service_members.last_name';
            WHEN 'edipi' THEN sort_column := 'service_members.edipi';
            WHEN 'emplid' THEN sort_column := 'service_members.emplid';
            WHEN 'requestedMoveDate' THEN sort_column := 'COALESCE(mto_shipments.requested_pickup_date, ppm_shipments.expected_departure_date, mto_shipments.requested_delivery_date)';
            WHEN 'appearedInTooAt' THEN sort_column := 'COALESCE(moves.submitted_at, moves.approvals_requested_at)';
            WHEN 'branch' THEN sort_column := 'service_members.affiliation';
            WHEN 'originDutyLocation' THEN sort_column := 'origin_duty_locations.name';
            WHEN 'counselingOffice' THEN sort_column := 'counseling_offices.name';
            WHEN 'assignedTo' THEN sort_column := 'too_user.last_name';
            ELSE
                sort_column := 'moves.id';
        END CASE;
    END IF;

    IF sort_direction IS NOT NULL THEN
        IF LOWER(sort_direction) = 'desc' THEN
            sort_order := 'DESC';
        ELSE
            sort_order := 'ASC';
        END IF;
    END IF;

    sql_query := sql_query || '
        GROUP BY
            moves.id,
            moves.show,
            moves.locator,
            moves.submitted_at,
            moves.orders_id,
            moves.status,
            moves.locked_by,
            moves.too_assigned_id,
            moves.counseling_transportation_office_id,
            mto_shipments.requested_pickup_date,
            mto_shipments.requested_delivery_date,
            ppm_shipments.expected_departure_date,
            orders.id,
            service_members.id,
            service_members.first_name,
            service_members.last_name,
            service_members.edipi,
            service_members.emplid,
            service_members.affiliation,
            origin_duty_locations.name,
            counseling_offices.name,
            too_user.first_name,
            too_user.last_name';
    sql_query := sql_query || format(' ORDER BY %s %s ', sort_column, sort_order);
    sql_query := sql_query || ' LIMIT $13 OFFSET $14 ';

    RETURN QUERY EXECUTE sql_query
    USING user_gbloc, customer_name, edipi, emplid, m_status, move_code, requested_move_date, date_submitted,
          branch, origin_duty_location, counseling_office, too_assigned_user, per_page, offset_value;

END;
$$ LANGUAGE plpgsql;
