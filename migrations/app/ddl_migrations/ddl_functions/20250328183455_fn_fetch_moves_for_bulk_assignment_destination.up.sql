--B-22151 Jonathan Spight  added fetch_moves_for_bulk_assignment_destination
CREATE OR REPLACE FUNCTION public.fetch_moves_for_bulk_assignment_destination(
    v_gbloc text
)
RETURNS TABLE (ID uuid, earliest_date date)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT
        moves.id,
        MIN(
            LEAST(
                COALESCE(mto_shipments.requested_pickup_date, '9999-12-31'::date),
                COALESCE(mto_shipments.requested_delivery_date, '9999-12-31'::date),
                COALESCE(ppm_shipments.expected_departure_date, '9999-12-31'::date)
            )
        ) AS earliest_date
    FROM moves
    INNER JOIN orders ON orders.id = moves.orders_id
    INNER JOIN service_members ON orders.service_member_id = service_members.id
    INNER JOIN mto_shipments ON mto_shipments.move_id = moves.id
    JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
    JOIN re_services  ON mto_service_items.re_service_id = re_services.id
    INNER JOIN duty_locations AS origin_dl ON orders.origin_duty_location_id = origin_dl.id
    LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
    LEFT JOIN shipment_address_updates ON shipment_address_updates.shipment_id = mto_shipments.id
    LEFT JOIN sit_extensions ON sit_extensions.mto_shipment_id = mto_shipments.id
    JOIN move_to_dest_gbloc ON move_to_dest_gbloc.move_id = moves.id
    WHERE
        mto_shipments.deleted_at IS NULL
        AND moves.status IN (
            'APPROVALS REQUESTED'
        )
        AND moves.show = TRUE
        AND moves.too_destination_assigned_id IS NULL
        AND orders.orders_type NOT IN (
            'BLUEBARK',
            'WOUNDED_WARRIOR',
            'SAFETY'
        )
        AND (
            moves.ppm_type IS NULL
            OR (
                moves.ppm_type = 'PARTIAL'
                OR (
                    moves.ppm_type = 'FULL'
                    AND origin_dl.provides_services_counseling = 'false'
                )
            )
        )
        AND (
            shipment_address_updates.status = 'REQUESTED'
            OR (
                sit_extensions.status = 'PENDING'
                AND re_services.code IN (
                    'DDFSIT', 'DDASIT', 'DDDSIT', 'DDSFSC', 'DDSHUT',
                    'IDFSIT', 'IDASIT', 'IDDSIT', 'IDSFSC', 'IDSHUT'
                )
            )
            OR (
                mto_service_items.status = 'SUBMITTED'
                AND re_services.code IN (
                    'DDFSIT', 'DDASIT', 'DDDSIT', 'DDSHUT', 'DDSFSC',
                    'IDFSIT', 'IDASIT', 'IDDSIT', 'IDSHUT', 'IDSFSC'
                )
            )
        )
        AND move_to_dest_gbloc.gbloc  = v_gbloc
    GROUP BY moves.id
    ORDER BY earliest_date;
END
$function$;