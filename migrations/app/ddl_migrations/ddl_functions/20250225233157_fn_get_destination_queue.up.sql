-- B-22294 - Alex Lusk - Migrating get_destination_queue function to ddl_migrations
-- B-21824 - Samay Sofo replaced too_assigned_id with too_destination_assigned_id for destination assigned queue
-- B-21902 - Samay Sofo added has_safety_privilege parameter to filter out safety orders and also retrieved orders_type
-- B-22760 - Paul Stonebraker retrieve mto_service_items for the moves and delivery address update requests for the shipments
-- B-23545 - Daniel Jordan updating returns to use destination, filtering adjustments, removing gbloc return
-- B-23739 - Daniel Jordan updating returns to consider lock_expires_at
-- B-22759 - Paul Stonebraker add SIT extensions as part of the mto_shipments

-- database function that returns a list of moves that have destination requests
-- this includes shipment address update requests, destination SIT, & destination shuttle

DROP FUNCTION IF EXISTS get_destination_queue;
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
    new_duty_location TEXT DEFAULT NULL,
    counseling_office TEXT DEFAULT NULL,
    too_assigned_user TEXT DEFAULT NULL,
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
    lock_expires_at TIMESTAMP WITH TIME ZONE,
    too_destination_assigned_id UUID,
    counseling_transportation_office_id UUID,
    orders JSONB,
    mto_shipments JSONB,
    counseling_transportation_office JSONB,
    too_destination_assigned JSONB,
    mto_service_items JSONB,
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
            moves.lock_expires_at,
            moves.too_destination_assigned_id AS too_destination_assigned_id,
            moves.counseling_transportation_office_id AS counseling_transportation_office_id,
            json_build_object(
                ''id'', orders.id,
                ''orders_type'', orders.orders_type,
                ''service_member'', json_build_object(
                    ''id'', service_members.id,
                    ''first_name'', service_members.first_name,
                    ''last_name'', service_members.last_name,
                    ''edipi'', service_members.edipi,
                    ''emplid'', service_members.emplid,
                    ''affiliation'', service_members.affiliation
                ),
                ''new_duty_location'', json_build_object(
                    ''name'', new_duty_locations.name
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
                            ''prime_estimated_weight'', ms.prime_estimated_weight,
                            ''delivery_address_update'', json_build_object(
                                ''status'', ms.address_update_status
                            ),
                            ''sit_duration_updates'', (
                                SELECT json_agg(
                                    json_build_object(
                                        ''status'', se.status
                                    )
                                )
                                FROM sit_extensions se
                                LEFT JOIN mto_shipments ON mto_shipments.id = se.mto_shipment_id
                                LEFT JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
                                LEFT JOIN re_services ON mto_service_items.re_service_id = re_services.id
                                WHERE se.mto_shipment_id = ms.id AND re_services.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSFSC'', ''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSFSC'')
                            )
                        )
                    )
                    FROM (
                        SELECT DISTINCT ON (mto_shipments.id)
                            mto_shipments.id,
                            mto_shipments.shipment_type,
                            mto_shipments.status,
                            mto_shipments.requested_pickup_date,
                            mto_shipments.scheduled_pickup_date,
                            mto_shipments.approved_date,
                            mto_shipments.prime_estimated_weight,
                            shipment_address_updates.status as address_update_status
                        FROM mto_shipments
                        LEFT JOIN shipment_address_updates on shipment_address_updates.shipment_id = mto_shipments.id
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
                ''last_name'', too_user.last_name,
                ''id'', too_user.id
            )::JSONB AS too_destination_assigned,
            COALESCE(
                (
                    SELECT json_agg(
                        json_build_object(
                            ''id'', msi.id,
                            ''status'', msi.status,
                            ''re_service'', json_build_object(
                                ''code'', msi.code
                            )
                        )
                    )
                    FROM (
                        SELECT mto_service_items.id, mto_service_items.status, re_services.code
                        FROM mto_service_items
                        LEFT JOIN re_services on mto_service_items.re_service_id = re_services.id
                        WHERE mto_service_items.move_id = moves.id
                    ) as msi
                ),
                ''[]''
            )::JSONB as mto_service_items,
            COUNT(*) OVER() AS total_count
        FROM moves
        JOIN orders ON moves.orders_id = orders.id
        JOIN mto_shipments ON mto_shipments.move_id = moves.id
        LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
        JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
        JOIN re_services ON mto_service_items.re_service_id = re_services.id
        JOIN service_members ON orders.service_member_id = service_members.id
        JOIN duty_locations AS new_duty_locations ON orders.new_duty_location_id = new_duty_locations.id
        LEFT JOIN office_users AS too_user ON moves.too_destination_assigned_id = too_user.id
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

    IF new_duty_location IS NOT NULL THEN
        sql_query := sql_query || ' AND new_duty_locations.name ILIKE ''%'' || $10 || ''%'' ';
    END IF;

    IF counseling_office IS NOT NULL THEN
        sql_query := sql_query || ' AND counseling_offices.name ILIKE ''%'' || $11 || ''%'' ';
    END IF;

    IF too_assigned_user IS NOT NULL THEN
        sql_query := sql_query || ' AND (too_user.first_name || '' '' || too_user.last_name) ILIKE ''%'' || $12 || ''%'' ';
    END IF;

   -- filter out safety orders for users without safety privilege
   IF NOT has_safety_privilege THEN
    sql_query := sql_query || ' AND orders.orders_type != ''SAFETY'' ';
   END IF;

    -- add destination queue-specific filters (pending dest address requests, pending dest SIT extension requests when there are dest SIT service items, submitted dest SIT & dest shuttle service items)
    sql_query := sql_query || '
        AND (
            shipment_address_updates.status = ''REQUESTED''
            OR (
                sit_extensions.status = ''PENDING''
                AND re_services.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSFSC'', ''DDSHUT'', ''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSFSC'', ''IDSHUT'')
            )
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
            WHEN 'customerName' THEN sort_column := 'service_members.last_name, service_members.first_name';
            WHEN 'edipi' THEN sort_column := 'service_members.edipi';
            WHEN 'emplid' THEN sort_column := 'service_members.emplid';
            WHEN 'requestedMoveDate' THEN   sort_column := 'COALESCE(' || 'MIN(mto_shipments.requested_pickup_date),' || 'MIN(ppm_shipments.expected_departure_date),' || 'MIN(mto_shipments.requested_delivery_date)' || ')';
            WHEN 'appearedInTooAt' THEN sort_column := 'COALESCE(moves.submitted_at, moves.approvals_requested_at)';
            WHEN 'branch' THEN sort_column := 'service_members.affiliation';
            WHEN 'destinationDutyLocation' THEN sort_column := 'new_duty_locations.name';
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
            moves.lock_expires_at,
            moves.too_destination_assigned_id,
            moves.counseling_transportation_office_id,
            orders.id,
            service_members.id,
            service_members.first_name,
            service_members.last_name,
            service_members.edipi,
            service_members.emplid,
            service_members.affiliation,
            new_duty_locations.name,
            counseling_offices.name,
            too_user.first_name,
            too_user.last_name,
            too_user.id';

    -- handling ordering customer name by last, first
    IF sort = 'customerName' THEN
      sql_query := sql_query || format(
        ' ORDER BY service_members.last_name %s, service_members.first_name %s',
        sort_order, sort_order
      );
    ELSE
      sql_query := sql_query || format(
        ' ORDER BY %s %s',
        sort_column, sort_order
      );
    END IF;

    IF sort_column <> 'moves.locator' OR sort <> 'customerName' OR sort <> 'requestedMoveDate' THEN
      sql_query := sql_query || ', moves.locator ASC';
    END IF;
    sql_query := sql_query || ' LIMIT $14 OFFSET $15 ';

    RETURN QUERY EXECUTE sql_query
    USING user_gbloc, customer_name, edipi, emplid, m_status, move_code, requested_move_date, date_submitted,
          branch, new_duty_location, counseling_office, too_assigned_user, has_safety_privilege, per_page, offset_value;

END;
$$ LANGUAGE plpgsql;
