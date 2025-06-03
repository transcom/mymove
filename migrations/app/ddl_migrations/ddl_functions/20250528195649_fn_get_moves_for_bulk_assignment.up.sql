-- B-22968 Maria Traskowsky added get_moves_for_bulk_assignment, refactoring FetchMovesForBulkAssignmentTaskOrder() into db func
DROP FUNCTION IF EXISTS public.get_moves_for_bulk_assignment(v_gbloc);
CREATE OR REPLACE FUNCTION public.get_moves_for_bulk_assignment(v_gbloc TEXT) RETURNS TABLE (id UUID, earliest_date DATE) LANGUAGE plpgsql AS $function$ BEGIN RETURN QUERY
SELECT moves.id,
    MIN(
        LEAST(
            COALESCE(
                mto_shipments.requested_pickup_date,
                '9999-12-31'
            ),
            COALESCE(
                mto_shipments.requested_delivery_date,
                '9999-12-31'
            ),
            COALESCE(
                ppm_shipments.expected_departure_date,
                '9999-12-31'
            )
        )
    ) AS earliest_date
FROM moves
    INNER JOIN orders ON orders.id = moves.orders_id
    INNER JOIN service_members ON orders.service_member_id = service_members.id
    INNER JOIN mto_shipments ON mto_shipments.move_id = moves.id
    INNER JOIN duty_locations AS origin_dl ON orders.origin_duty_location_id = origin_dl.id
    LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    LEFT JOIN move_to_gbloc ON move_to_gbloc.move_id = moves.id
WHERE mto_shipments.deleted_at IS NULL
    AND moves.status IN (
        'APPROVALS REQUESTED',
        'SUBMITTED',
        'SERVICE COUNSELING COMPLETED'
    )
    AND moves.show = TRUE
    AND moves.too_assigned_id IS NULL
    AND orders.orders_type NOT IN ('BLUEBARK', 'WOUNDED_WARRIOR', 'SAFETY')
    AND (
        moves.ppm_type IS NULL
        OR moves.ppm_type = 'PARTIAL'
        OR (
            moves.ppm_type = 'FULL'
            AND origin_dl.provides_services_counseling = 'false'
        )
    )
    AND (
        v_gbloc IS NULL
        OR EXISTS (
            SELECT 1
            FROM mto_shipments ms
            WHERE ms.move_id = moves.id
                AND (
                    (
                        ms.shipment_type != 'HHG_OUTOF_NTS'
                        AND move_to_gbloc.gbloc = v_gbloc
                    )
                    OR (
                        ms.shipment_type = 'HHG_OUTOF_NTS'
                        AND orders.gbloc = v_gbloc
                    )
                )
                AND (
                    ms.status IN ('SUBMITTED', 'APPROVALS_REQUESTED')
                    OR (
                        ms.status = 'APPROVED'
                        AND (
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
                )
        )
    ) -- we want to omit shipments with ONLY destination queue-specific filters
    -- (pending dest address requests, pending dest SIT extension requests when there are dest SIT service items, submitted dest SIT & dest shuttle service items)
    AND NOT (
        (
            EXISTS (
                SELECT 1
                FROM mto_service_items msi
                    JOIN re_services rs ON msi.re_service_id = rs.id
                WHERE msi.mto_shipment_id IN (
                        SELECT ms.id
                        FROM mto_shipments ms
                        WHERE ms.move_id = moves.id
                    )
                    AND msi.status = 'SUBMITTED'
                    AND rs.code IN (
                        'DDFSIT',
                        'DDASIT',
                        'DDDSIT',
                        'DDSHUT',
                        'DDSFSC',
                        'IDFSIT',
                        'IDASIT',
                        'IDDSIT',
                        'IDSHUT',
                        'IDSFSC'
                    )
            )
            OR EXISTS (
                SELECT 1
                FROM shipment_address_updates sau
                WHERE sau.shipment_id IN (
                        SELECT ms.id
                        FROM mto_shipments ms
                        WHERE ms.move_id = moves.id
                    )
                    AND sau.status = 'REQUESTED'
            )
            OR (
                EXISTS (
                    SELECT 1
                    FROM sit_extensions se
                        JOIN mto_service_items msi ON se.mto_shipment_id = msi.mto_shipment_id
                        JOIN re_services rs ON msi.re_service_id = rs.id
                    WHERE se.mto_shipment_id IN (
                            SELECT ms.id
                            FROM mto_shipments ms
                            WHERE ms.move_id = moves.id
                        )
                        AND se.status = 'PENDING'
                        AND rs.code IN (
                            'DDFSIT',
                            'DDASIT',
                            'DDDSIT',
                            'DDSFSC',
                            'IDFSIT',
                            'IDASIT',
                            'IDDSIT',
                            'IDSFSC'
                        )
                )
                AND NOT EXISTS (
                    SELECT 1
                    FROM mto_service_items msi
                        JOIN re_services rs ON msi.re_service_id = rs.id
                    WHERE msi.mto_shipment_id IN (
                            SELECT ms.id
                            FROM mto_shipments ms
                            WHERE ms.move_id = moves.id
                        )
                        AND msi.status = 'SUBMITTED'
                        AND rs.code IN (
                            'ICRT',
                            'IUBPK',
                            'IOFSIT',
                            'IOASIT',
                            'IOPSIT',
                            'IOSHUT',
                            'IHUPK',
                            'IUCRT',
                            'DCRT',
                            'MS',
                            'CS',
                            'DOFSIT',
                            'DOASIT',
                            'DOPSIT',
                            'DOSFSC',
                            'IOSFSC',
                            'DUPK',
                            'DUCRT',
                            'DOSHUT',
                            'FSC',
                            'DMHF',
                            'DBTF',
                            'DBHF',
                            'IBTF',
                            'IBHF',
                            'DCRTSA',
                            'DLH',
                            'DOP',
                            'DPK',
                            'DSH',
                            'DNPK',
                            'INPK',
                            'UBP',
                            'ISLH',
                            'POEFSC',
                            'PODFSC',
                            'IHPK'
                        )
                )
            )
        )
        AND NOT (
            EXISTS (
                SELECT 1
                FROM mto_service_items msi
                    JOIN re_services rs ON msi.re_service_id = rs.id
                WHERE msi.mto_shipment_id IN (
                        SELECT ms.id
                        FROM mto_shipments ms
                        WHERE ms.move_id = moves.id
                    )
                    AND msi.status = 'SUBMITTED'
                    AND rs.code IN (
                        'ICRT',
                        'IUBPK',
                        'IOFSIT',
                        'IOASIT',
                        'IOPSIT',
                        'IOSHUT',
                        'IHUPK',
                        'IUCRT',
                        'DCRT',
                        'MS',
                        'CS',
                        'DOFSIT',
                        'DOASIT',
                        'DOPSIT',
                        'DOSFSC',
                        'IOSFSC',
                        'DUPK',
                        'DUCRT',
                        'DOSHUT',
                        'FSC',
                        'DMHF',
                        'DBTF',
                        'DBHF',
                        'IBTF',
                        'IBHF',
                        'DCRTSA',
                        'DLH',
                        'DOP',
                        'DPK',
                        'DSH',
                        'DNPK',
                        'INPK',
                        'UBP',
                        'ISLH',
                        'POEFSC',
                        'PODFSC',
                        'IHPK'
                    )
            )
        )
    )
GROUP BY moves.id
ORDER BY earliest_date;
END;
$function$;