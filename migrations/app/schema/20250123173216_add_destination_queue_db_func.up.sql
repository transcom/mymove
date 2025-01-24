-- database function that returns a list of moves that have destination requests
-- this includes shipment address update requests & destination service items
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
    per_page INTEGER DEFAULT 20
)
RETURNS TABLE (
    id UUID,
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
BEGIN
    IF page < 1 THEN
        page := 1;
    END IF;

    IF per_page < 1 THEN
        per_page := 20;
    END IF;

    -- OFFSET for pagination
    offset_value := (page - 1) * per_page;

    sql_query := '
        SELECT
            moves.id AS id,
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
                json_agg(
                    json_build_object(
                        ''id'', mto_shipments.id,
                        ''shipment_type'', mto_shipments.shipment_type,
                        ''status'', mto_shipments.status,
                        ''requested_pickup_date'', TO_CHAR(mto_shipments.requested_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''scheduled_pickup_date'', TO_CHAR(mto_shipments.scheduled_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''requested_delivery_date'', TO_CHAR(mto_shipments.requested_delivery_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''approved_date'', TO_CHAR(mto_shipments.approved_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''prime_estimated_weight'', mto_shipments.prime_estimated_weight
                    )
                ) FILTER (WHERE mto_shipments.id IS NOT NULL),
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
        INNER JOIN orders ON moves.orders_id = orders.id
        LEFT JOIN mto_shipments ON mto_shipments.move_id = moves.id
        LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
        LEFT JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
        LEFT JOIN re_services ON mto_service_items.re_service_id = re_services.id
        LEFT JOIN service_members ON orders.service_member_id = service_members.id
        LEFT JOIN duty_locations AS new_duty_locations ON orders.new_duty_location_id = new_duty_locations.id
        LEFT JOIN duty_locations AS origin_duty_locations ON orders.origin_duty_location_id = origin_duty_locations.id
        LEFT JOIN office_users AS too_user ON moves.too_assigned_id = too_user.id
        LEFT JOIN office_users AS locked_user ON moves.locked_by = locked_user.id
        LEFT JOIN transportation_offices AS counseling_offices
            ON moves.counseling_transportation_office_id = counseling_offices.id
        LEFT JOIN shipment_address_updates ON shipment_address_updates.shipment_id = mto_shipments.id
        LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
        WHERE moves.show = TRUE
    ';

    -- adding conditionals for destination queue
    -- we only want to see moves that have destination requests (shipment address updates, destination service items in SUBMITTED status)
    sql_query := sql_query || '
        AND (
            shipment_address_updates.status = ''REQUESTED''
            OR (
                mto_service_items.status = ''SUBMITTED''
                AND re_services.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSHUT'', ''DDSFSC'', ''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSHUT'')
            )
        )
    ';

    -- this should always be passed in from the service object, but will nil check it anyway
    IF user_gbloc IS NOT NULL THEN
        sql_query := sql_query || ' AND move_to_gbloc.gbloc = ''%' || user_gbloc || '%'' ';
    END IF;

    IF customer_name IS NOT NULL AND customer_name <> '' THEN
        sql_query := sql_query || ' AND (
            service_members.first_name || '' '' || service_members.last_name ILIKE ''%' || customer_name || '%''
            OR service_members.last_name || '' '' || service_members.first_name ILIKE ''%' || customer_name || '%''
        ) ';
    END IF;

    IF edipi IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.edipi ILIKE ''%' || edipi || '%'' ';
    END IF;

    IF emplid IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.emplid ILIKE ''%' || emplid || '%'' ';
    END IF;

    IF m_status IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.status IN (SELECT unnest($1)) ';
    END IF;

    IF move_code IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.locator ILIKE ''%' || UPPER(move_code) || '%'' ';
    END IF;

    IF requested_move_date IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            mto_shipments.requested_pickup_date::DATE = ' || quote_literal(requested_move_date) || '::DATE
            OR ppm_shipments.expected_departure_date::DATE = ' || quote_literal(requested_move_date) || '::DATE
            OR (mto_shipments.shipment_type = ''HHG_OUTOF_NTS'' AND mto_shipments.requested_delivery_date::DATE = ' || quote_literal(requested_move_date) || '::DATE)
        )';
    END IF;

    IF date_submitted IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.submitted_at = ' || quote_literal(date_submitted) || ' ';
    END IF;

    IF branch IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.affiliation ILIKE ''%' || branch || '%'' ';
    END IF;

    IF origin_duty_location IS NOT NULL AND origin_duty_location <> '' THEN
        sql_query := sql_query || ' AND origin_duty_locations.name ILIKE ''%' || origin_duty_location || '%'' ';
    END IF;

    IF counseling_office IS NOT NULL THEN
        sql_query := sql_query || ' AND counseling_offices.name ILIKE ''%' || counseling_office || '%'' ';
    END IF;

    IF too_assigned_user IS NOT NULL THEN
        sql_query := sql_query || ' AND (too_user.first_name || '' '' || too_user.last_name) ILIKE ''%' || quote_literal(too_assigned_user) || '%'' ';
    END IF;

    sql_query := sql_query || '
        GROUP BY
            moves.id,
            moves.locator,
            moves.submitted_at,
            moves.orders_id,
            moves.status,
            moves.locked_by,
            moves.too_assigned_id,
            moves.counseling_transportation_office_id,
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
    sql_query := sql_query || ' ORDER BY moves.id ASC ';
    sql_query := sql_query || ' LIMIT ' || per_page || ' OFFSET ' || offset_value || ' ';

    RAISE NOTICE 'Query: %', sql_query;

    RETURN QUERY EXECUTE sql_query USING m_status;

END;
$$ LANGUAGE plpgsql;
