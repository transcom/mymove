CREATE OR REPLACE FUNCTION get_destination_queue(
    move_code TEXT DEFAULT NULL,                -- Search parameter for Move Code
    input_move_id UUID DEFAULT NULL,            -- Search parameter for Move ID
    page INTEGER DEFAULT 1,                     -- Page number for pagination
    per_page INTEGER DEFAULT 20,                -- Number of results per page
    branch TEXT DEFAULT NULL,                   -- Filter: service_member.affiliation
    edipi TEXT DEFAULT NULL,                    -- Filter: service_member.edipi
    emplid TEXT DEFAULT NULL,                   -- Filter: service_member.emplid
    customer_name TEXT DEFAULT NULL,            -- Filter: service_member.first_name + last_name
    destination_duty_location TEXT DEFAULT NULL,-- Filter: orders.new_duty_location_id.name
    origin_duty_location TEXT DEFAULT NULL,     -- Filter: orders.origin_duty_location_id.name
    origin_gbloc TEXT DEFAULT NULL,             -- Filter: move.counseling_office_transportation_office.gbloc
    submitted_at TIMESTAMP DEFAULT NULL,        -- Filter: moves.submitted_at
    appeared_in_too_at TIMESTAMP DEFAULT NULL,  -- Filter: moves.appeared_in_too_at
    too_assigned_user TEXT DEFAULT NULL         -- Filter: moves.too_assigned_id -> office_users.first_name + last_name
)
RETURNS TABLE (
    move_id UUID,
    locator TEXT,
    orders_id UUID,
    available_to_prime_at TIMESTAMP WITH TIME ZONE,
    show BOOLEAN,
    total_count BIGINT
) AS $$
DECLARE
    sql_query TEXT;
    offset_value INTEGER;
BEGIN
    -- OFFSET for pagination
    offset_value := (page - 1) * per_page;

    sql_query := '
        SELECT moves.id AS move_id,
               moves.locator::TEXT AS locator,
               moves.orders_id,
               moves.available_to_prime_at,
               moves.show,
               COUNT(*) OVER() AS total_count
        FROM moves
        INNER JOIN orders ON moves.orders_id = orders.id
        INNER JOIN service_members ON orders.service_member_id = service_members.id
        LEFT JOIN duty_locations AS new_duty_locations ON orders.new_duty_location_id = new_duty_locations.id
        LEFT JOIN duty_locations AS origin_duty_locations ON orders.origin_duty_location_id = origin_duty_locations.id
        LEFT JOIN office_users ON moves.too_assigned_id = office_users.id
        LEFT JOIN transportation_offices AS counseling_offices
            ON moves.counseling_transportation_office_id = counseling_offices.id
        WHERE moves.available_to_prime_at IS NOT NULL
          AND moves.show = TRUE
    ';

    -- add dynamic filters based on provided parameters
    IF move_code IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.locator ILIKE ''%' || UPPER(move_code) || '%'' ';
    END IF;

    IF input_move_id IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.id = ''' || input_move_id || ''' ';
    END IF;

    IF branch IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.affiliation ILIKE ''%' || branch || '%'' ';
    END IF;

    IF edipi IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.edipi ILIKE ''%' || edipi || '%'' ';
    END IF;

    IF emplid IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.emplid ILIKE ''%' || emplid || '%'' ';
    END IF;

    IF customer_name IS NOT NULL THEN
        sql_query := sql_query || ' AND (service_members.first_name || '' '' || service_members.last_name) ILIKE ''%' || customer_name || '%'' ';
    END IF;

    IF destination_duty_location IS NOT NULL THEN
        sql_query := sql_query || ' AND new_duty_locations.name ILIKE ''%' || destination_duty_location || '%'' ';
    END IF;

    IF origin_duty_location IS NOT NULL THEN
        sql_query := sql_query || ' AND origin_duty_locations.name ILIKE ''%' || origin_duty_location || '%'' ';
    END IF;

    IF origin_gbloc IS NOT NULL THEN
        sql_query := sql_query || ' AND counseling_offices.gbloc ILIKE ''%' || origin_gbloc || '%'' ';
    END IF;

    IF submitted_at IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.submitted_at = ''' || submitted_at || ''' ';
    END IF;

    IF appeared_in_too_at IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.appeared_in_too_at = ''' || appeared_in_too_at || ''' ';
    END IF;

    IF too_assigned_user IS NOT NULL THEN
        sql_query := sql_query || ' AND (office_users.first_name || '' '' || office_users.last_name) ILIKE ''%' || too_assigned_user || '%'' ';
    END IF;

    sql_query := sql_query || ' ORDER BY moves.id ASC ';
    sql_query := sql_query || ' LIMIT ' || per_page || ' OFFSET ' || offset_value || ' ';

    RETURN QUERY EXECUTE sql_query;

END;
$$ LANGUAGE plpgsql;