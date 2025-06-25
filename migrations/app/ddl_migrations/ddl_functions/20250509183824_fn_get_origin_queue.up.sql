-- B-23540 - Daniel Jordan - initial function creation for TOO origin queue refactor into db func
-- B-23739 - Daniel Jordan - updating returns to consider lock_expires_at
-- B-23767  Daniel Jordan - updating query to exclude FULL PPM types that provide SC and null PPM types
-- B-22712 -- Paul Stonebraker - add move data for excess weight, amended orders; attach diversions and SIT extensions to mto shipments
-- B-23582 - Paul Stonebraker - update to handle new task order queue specific assignment column


DROP FUNCTION IF EXISTS get_origin_queue;
CREATE OR REPLACE FUNCTION get_origin_queue(
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
    too_task_order_assigned_id UUID,
    counseling_transportation_office_id UUID,
    orders JSONB,
    mto_shipments JSONB,
    mto_service_items JSONB,
    counseling_transportation_office JSONB,
    too_task_order_assigned JSONB,
    excess_weight_qualified_at TIMESTAMP WITH TIME ZONE,
    excess_weight_acknowledged_at TIMESTAMP WITH TIME ZONE,
    excess_unaccompanied_baggage_weight_qualified_at TIMESTAMP WITH TIME ZONE,
    excess_unaccompanied_baggage_weight_acknowledged_at TIMESTAMP WITH TIME ZONE,
    total_count BIGINT
) AS $$
DECLARE
    sql_query TEXT;
    offset_value INTEGER;
    sort_column TEXT;
    sort_order TEXT;
    total_count BIGINT;
BEGIN
    IF page < 1 THEN
        page := 1;
    END IF;
    IF per_page < 1 THEN
        per_page := 20;
    END IF;
    offset_value := (page - 1) * per_page;

    IF sort IS NOT NULL THEN
        CASE sort
            WHEN 'locator' THEN sort_column := 'base.locator';
            WHEN 'status' THEN sort_column := 'CASE base.status
                WHEN ''SERVICE COUNSELING COMPLETED'' THEN 1
                WHEN ''SUBMITTED'' THEN 2
                WHEN ''NEEDS SERVICE COUNSELING'' THEN 3
                WHEN ''APPROVALS REQUESTED'' THEN 4
                WHEN ''APPROVED'' THEN 5
                ELSE 99 END';
            WHEN 'customerName' THEN sort_column := 'base.sm_last_name, base.sm_first_name';
            WHEN 'edipi' THEN sort_column := 'base.sm_edipi';
            WHEN 'emplid' THEN sort_column := 'base.sm_emplid';
            WHEN 'requestedMoveDate' THEN sort_column := 'LEAST(base.earliest_requested_pickup_date, base.earliest_expected_departure_date, base.earliest_requested_delivery_date)';
            WHEN 'appearedInTooAt' THEN sort_column := 'GREATEST(base.submitted_at, base.service_counseling_completed_at, base.approvals_requested_at)';
            WHEN 'branch' THEN sort_column := 'base.sm_affiliation';
            WHEN 'originDutyLocation' THEN sort_column := 'base.origin_duty_location_name';
            WHEN 'counselingOffice' THEN sort_column := 'base.counseling_office_name';
            WHEN 'assignedTo' THEN sort_column := 'base.too_user_last_name';
            ELSE sort_column := 'base.id';
        END CASE;
    END IF;

    IF sort_direction IS NOT NULL THEN
        IF LOWER(sort_direction) = 'desc' THEN
            sort_order := 'DESC';
        ELSE
            sort_order := 'ASC';
        END IF;
    END IF;

    IF sort_column IS NULL THEN
        sort_column := 'CASE base.status
            WHEN ''APPROVALS REQUESTED'' THEN 1
            WHEN ''SUBMITTED'' THEN 2
            WHEN ''SERVICE COUNSELING COMPLETED'' THEN 3
            ELSE 99 END';
    END IF;

    IF sort_order IS NULL THEN
        sort_order := 'ASC';
    END IF;

    sql_query := '
    WITH base AS (
        SELECT
            moves.id,
            moves.show,
            moves.locator::TEXT AS locator,
            moves.submitted_at,
            moves.orders_id,
            moves.status,
            moves.locked_by,
            moves.lock_expires_at,
            moves.too_task_order_assigned_id,
            moves.counseling_transportation_office_id,
            moves.service_counseling_completed_at,
            moves.approvals_requested_at,
            moves.excess_weight_qualified_at,
            moves.excess_weight_acknowledged_at,
            moves.excess_unaccompanied_baggage_weight_qualified_at,
            moves.excess_unaccompanied_baggage_weight_acknowledged_at,
            orders.id AS orders_id_inner,
            orders.orders_type,
            orders.department_indicator AS orders_department_indicator,
            orders.gbloc,
            orders.uploaded_amended_orders_id,
            orders.amended_orders_acknowledged_at,
            service_members.id AS sm_id,
            service_members.first_name AS sm_first_name,
            service_members.last_name AS sm_last_name,
            service_members.edipi AS sm_edipi,
            service_members.emplid AS sm_emplid,
            service_members.affiliation AS sm_affiliation,
            origin_duty_locations.id AS origin_duty_location_id,
            origin_duty_locations.name AS origin_duty_location_name,
            addr.street_address_1 AS origin_duty_location_street_address_1,
            addr.street_address_2 AS origin_duty_location_street_address_2,
            addr.city AS origin_duty_location_city,
            addr.state AS origin_duty_location_state,
            addr.postal_code AS origin_duty_location_postal_code,
            counseling_offices.name AS counseling_office_name,
            too_user.first_name AS too_user_first_name,
            too_user.last_name AS too_user_last_name,
            too_user.id AS too_user_id,
            shipments.mto_shipments,
            shipments.earliest_requested_pickup_date,
            shipments.earliest_requested_delivery_date,
            ppm_dates.earliest_expected_departure_date,
            service_items.mto_service_items
        FROM moves
        JOIN orders ON moves.orders_id = orders.id
        JOIN service_members ON orders.service_member_id = service_members.id
        JOIN duty_locations AS origin_duty_locations ON orders.origin_duty_location_id = origin_duty_locations.id
        LEFT JOIN addresses AS addr ON origin_duty_locations.address_id = addr.id
        LEFT JOIN office_users AS too_user ON moves.too_task_order_assigned_id = too_user.id
        LEFT JOIN transportation_offices AS counseling_offices ON moves.counseling_transportation_office_id = counseling_offices.id
        JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
        LEFT JOIN LATERAL (
            SELECT
                json_agg(
                    json_build_object(
                        ''id'', ms.id,
                        ''shipment_type'', ms.shipment_type,
                        ''status'', ms.status,
                        ''requested_pickup_date'', TO_CHAR(ms.requested_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''scheduled_pickup_date'', TO_CHAR(ms.scheduled_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''requested_delivery_date'', TO_CHAR(ms.requested_delivery_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''approved_date'', TO_CHAR(ms.approved_date, ''YYYY-MM-DD"T00:00:00Z"''),
                        ''prime_estimated_weight'', ms.prime_estimated_weight,
                        ''ppm_shipment'', CASE
                            WHEN ppm.id IS NOT NULL THEN json_build_object(
                                ''expected_departure_date'', TO_CHAR(ppm.expected_departure_date, ''YYYY-MM-DD"T00:00:00Z"'')
                            )
                            ELSE NULL
                        END,
                        ''diversion'', ms.diversion,
                        ''sit_duration_updates'', (
                            SELECT json_agg(
                                json_build_object(
                                    ''status'', se.status
                                )
                            )
                            FROM sit_extensions se
                            LEFT JOIN mto_shipments ON mto_shipments.id = se.mto_shipment_id
                            LEFT JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
                            RIGHT JOIN re_services ON mto_service_items.re_service_id = re_services.id
                            WHERE se.mto_shipment_id = ms.id AND re_services.code IN (''DOFSIT'', ''DOASIT'', ''DOPSIT'', ''DOSFSC'', ''IOFSIT'', ''IOASIT'', ''IOPSIT'', ''IOSFSC'')
                        )
                    )
                )::JSONB AS mto_shipments,
                MIN(ms.requested_pickup_date) AS earliest_requested_pickup_date,
                MIN(ms.requested_delivery_date) AS earliest_requested_delivery_date
            FROM mto_shipments ms
            LEFT JOIN ppm_shipments ppm ON ppm.shipment_id = ms.id
            WHERE ms.move_id = moves.id
        ) shipments ON TRUE
        LEFT JOIN LATERAL (
            SELECT MIN(ppm.expected_departure_date) AS earliest_expected_departure_date
            FROM ppm_shipments ppm
            JOIN mto_shipments ms ON ppm.shipment_id = ms.id
            WHERE ms.move_id = moves.id
        ) ppm_dates ON TRUE
        LEFT JOIN LATERAL (
            SELECT
                json_agg(
                    json_build_object(
                        ''id'', si.id,
                        ''status'', si.status,
                        ''re_service'', json_build_object(
                            ''code'', rs.code
                        )
                    )
                )::JSONB AS mto_service_items
            FROM mto_service_items si
            JOIN re_services rs ON si.re_service_id = rs.id
            WHERE si.mto_shipment_id IN (
                SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id
            )
        ) service_items ON TRUE
        WHERE moves.show = TRUE
        AND (moves.ppm_type IS NULL OR moves.ppm_type = ''PARTIAL'' OR (moves.ppm_type = ''FULL'' AND origin_duty_locations.provides_services_counseling = false)) ';

    IF user_gbloc IS NOT NULL THEN
        sql_query := sql_query || '
        AND EXISTS (
            SELECT 1 FROM mto_shipments ms
            WHERE ms.move_id = moves.id
            AND (
                (ms.shipment_type != ''HHG_OUTOF_NTS'' AND move_to_gbloc.gbloc = $1)
                OR (ms.shipment_type = ''HHG_OUTOF_NTS'' AND orders.gbloc = $1)
            )
            AND (ms.status IN (''SUBMITTED'',''APPROVALS_REQUESTED'')
        	    OR (ms.status = ''APPROVED''
        			AND
                		(
                        	moves.excess_weight_qualified_at IS NOT NULL
                        	AND moves.excess_weight_acknowledged_at IS NULL
		                )
        		        OR (
                	        moves.excess_unaccompanied_baggage_weight_qualified_at IS NOT NULL
                    	    AND moves.excess_unaccompanied_baggage_weight_acknowledged_at IS NULL
                		)
                        OR (
                            orders.uploaded_amended_orders_id IS NOT NULL
                            AND orders.amended_orders_acknowledged_at IS NULL
                        )
        		)
    		)
        )';
    END IF;

    IF customer_name IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            service_members.first_name || '' '' || service_members.last_name ILIKE ''%'' || $2 || ''%''
            OR service_members.last_name || '' '' || service_members.first_name ILIKE ''%'' || $2 || ''%''
            OR service_members.last_name || '', '' || service_members.first_name ILIKE ''%'' || $2 || ''%''
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
        sql_query := sql_query || '
        AND EXISTS (
            SELECT 1
            FROM mto_shipments ms2
            LEFT JOIN ppm_shipments ppm2 ON ppm2.shipment_id = ms2.id
            WHERE ms2.move_id = moves.id
              AND (
                   ms2.requested_pickup_date::DATE = $7::DATE
                    OR ppm2.expected_departure_date::DATE = $7::DATE
                    OR (ms2.shipment_type = ''HHG_OUTOF_NTS''
                        AND ms2.requested_delivery_date::DATE = $7::DATE)
              )
        )';
    END IF;

    IF date_submitted IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            moves.submitted_at::DATE = $8::DATE OR
            moves.service_counseling_completed_at::DATE = $8::DATE OR
            moves.approvals_requested_at::DATE = $8::DATE) ';
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
        sql_query := sql_query || ' AND (
            too_user.first_name || '' '' || too_user.last_name ILIKE ''%'' || $12 || ''%''
            OR too_user.last_name || '' '' || too_user.first_name ILIKE ''%'' || $12 || ''%''
        )';
    END IF;

    IF NOT has_safety_privilege THEN
        sql_query := sql_query || ' AND orders_type != ''SAFETY'' ';
    END IF;

    -- we want to omit shipments with ONLY destination queue-specific filters
    -- (pending dest address requests, pending dest SIT extension requests when there are dest SIT service items, submitted dest SIT & dest shuttle service items)
    sql_query := sql_query || '
            AND NOT (
					(
						EXISTS (
							SELECT 1
							FROM mto_service_items msi
							JOIN re_services rs ON msi.re_service_id = rs.id
							WHERE msi.mto_shipment_id IN (SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id)
							AND msi.status = ''SUBMITTED''
							AND rs.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSHUT'', ''DDSFSC'',
											''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSHUT'', ''IDSFSC'')
						)
						OR EXISTS (
							SELECT 1
							FROM shipment_address_updates sau
							WHERE sau.shipment_id IN (SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id)
							AND sau.status = ''REQUESTED''
						)
						OR (
							EXISTS (
								SELECT 1
								FROM sit_extensions se
								JOIN mto_service_items msi ON se.mto_shipment_id = msi.mto_shipment_id
								JOIN re_services rs ON msi.re_service_id = rs.id
								WHERE se.mto_shipment_id IN (SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id)
								AND se.status = ''PENDING''
								AND rs.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSFSC'',
												''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSFSC'')
							)
							AND NOT EXISTS (
								SELECT 1
								FROM mto_service_items msi
								JOIN re_services rs ON msi.re_service_id = rs.id
								WHERE msi.mto_shipment_id IN (SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id)
								AND msi.status = ''SUBMITTED''
								AND rs.code IN (''ICRT'', ''IUBPK'', ''IOFSIT'', ''IOASIT'', ''IOPSIT'', ''IOSHUT'',
												''IHUPK'', ''IUCRT'', ''DCRT'', ''MS'', ''CS'', ''DOFSIT'', ''DOASIT'',
												''DOPSIT'', ''DOSFSC'', ''IOSFSC'', ''DUPK'', ''DUCRT'', ''DOSHUT'',
												''FSC'', ''DMHF'', ''DBTF'', ''DBHF'', ''IBTF'', ''IBHF'', ''DCRTSA'',
												''DLH'', ''DOP'', ''DPK'', ''DSH'', ''DNPK'', ''INPK'', ''UBP'',
												''ISLH'', ''POEFSC'', ''PODFSC'', ''IHPK'')
							)
						)
					)
					AND NOT (
						EXISTS (
							SELECT 1
							FROM mto_service_items msi
							JOIN re_services rs ON msi.re_service_id = rs.id
							WHERE msi.mto_shipment_id IN (SELECT ms.id FROM mto_shipments ms WHERE ms.move_id = moves.id)
							AND msi.status = ''SUBMITTED''
							AND rs.code IN (''ICRT'', ''IUBPK'', ''IOFSIT'', ''IOASIT'', ''IOPSIT'', ''IOSHUT'',
											''IHUPK'', ''IUCRT'', ''DCRT'', ''MS'', ''CS'', ''DOFSIT'', ''DOASIT'',
											''DOPSIT'', ''DOSFSC'', ''IOSFSC'', ''DUPK'', ''DUCRT'', ''DOSHUT'',
											''FSC'', ''DMHF'', ''DBTF'', ''DBHF'', ''IBTF'', ''IBHF'', ''DCRTSA'',
											''DLH'', ''DOP'', ''DPK'', ''DSH'', ''DNPK'', ''INPK'', ''UBP'',
											''ISLH'', ''POEFSC'', ''PODFSC'', ''IHPK'')
						)
					)
				)
            ';

    sql_query := sql_query || ' )
    SELECT
        id::UUID,
        show::BOOLEAN,
        locator::TEXT,
        submitted_at::TIMESTAMP WITH TIME ZONE AS submitted_at,
        orders_id_inner::UUID AS orders_id,
        status::TEXT,
        locked_by::UUID AS locked_by,
        lock_expires_at::TIMESTAMP WITH TIME ZONE AS lock_expires_at,
        too_task_order_assigned_id::UUID AS too_task_order_assigned_id,
        counseling_transportation_office_id::UUID AS counseling_transportation_office_id,
        json_build_object(
            ''id'', orders_id_inner,
            ''orders_type'', orders_type,
            ''origin_duty_location_gbloc'', gbloc,
            ''department_indicator'', orders_department_indicator,
            ''service_member'', json_build_object(
                ''id'', sm_id,
                ''first_name'', sm_first_name,
                ''last_name'', sm_last_name,
                ''edipi'', sm_edipi,
                ''emplid'', sm_emplid,
                ''affiliation'', sm_affiliation
            ),
            ''origin_duty_location'', json_build_object(
                ''id'',   origin_duty_location_id,
                ''name'', origin_duty_location_name,
                ''address'', json_build_object(
                    ''street_address_1'', origin_duty_location_street_address_1,
                    ''street_address_2'', origin_duty_location_street_address_2,
                    ''city'',             origin_duty_location_city,
                    ''state'',            origin_duty_location_state,
                    ''postal_code'',      origin_duty_location_postal_code
                )
            ),
            ''uploaded_amended_orders_id'', uploaded_amended_orders_id,
            ''amended_orders_acknowledged_at'', amended_orders_acknowledged_at::TIMESTAMP WITH TIME ZONE
        )::JSONB AS orders,
        COALESCE(mto_shipments, ''[]''::JSONB) AS mto_shipments,
        COALESCE(mto_service_items, ''[]''::JSONB) AS mto_service_items,
        json_build_object(''name'', counseling_office_name)::JSONB AS counseling_transportation_office,
        json_build_object(''first_name'', too_user_first_name, ''last_name'', too_user_last_name, ''id'', too_user_id)::JSONB AS too_task_order_assigned,
        excess_weight_qualified_at::TIMESTAMP WITH TIME ZONE,
        excess_weight_acknowledged_at::TIMESTAMP WITH TIME ZONE,
        excess_unaccompanied_baggage_weight_qualified_at::TIMESTAMP WITH TIME ZONE,
        excess_unaccompanied_baggage_weight_acknowledged_at::TIMESTAMP WITH TIME ZONE,
        COUNT(*) OVER() AS total_count
        FROM base ';

    IF sort = 'customerName' THEN
        sql_query := sql_query || format(
            ' ORDER BY sm_last_name %s, sm_first_name %s, locator ASC ',
            sort_order, sort_order
        );
    ELSE
        sql_query := sql_query || ' ORDER BY ' || sort_column || ' ' || sort_order || ', locator ASC ';
    END IF;

    sql_query := sql_query || ' LIMIT $14 OFFSET $15 ';

    RETURN QUERY EXECUTE sql_query
    USING user_gbloc, customer_name, edipi, emplid, m_status, move_code, requested_move_date, date_submitted,
          branch, origin_duty_location, counseling_office, too_assigned_user, has_safety_privilege, per_page, offset_value;
END;
$$ LANGUAGE plpgsql;
