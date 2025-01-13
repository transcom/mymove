DO
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'mto_service_item_type') THEN
        DROP TYPE mto_service_item_type
    END IF;

    CREATE TYPE  mto_service_item_type AS (
        id uuid,
        move_id uuid,
        mto_shipment_id uuid,
        re_service_code text,
        created_at timestamptz,
        updated_at timestamptz,
        reason text,
        pickup_postal_code text,
        description text,
        status public.service_item_status,
        rejection_reason text,
        approved_at timestamp,
        rejected_at timestamp,
        sit_postal_code text,
        sit_entry_date date,
        sit_departure_date date,
        sit_destination_final_address_id uuid,
        sit_origin_hhg_original_address_id uuid,
        sit_origin_hhg_actual_address_id uuid,
        estimated_weight int4,
        actual_weight int4,
        sit_destination_original_address_id uuid,
        sit_customer_contacted date,
        sit_requested_delivery date,
        requested_approvals_requested_status bool,
        customer_expense bool,
        customer_expense_reason text,
        sit_delivery_miles int4,
        pricing_estimate int4,
        standalone_crate bool,
        locked_price_cents int4,
        service_location public.service_location_enum,
        poe_location_id uuid,
        pod_location_id uuid
    );
END;

CREATE OR REPLACE PROCEDURE create_accessorial_service_items_for_shipment (
    IN shipment_id UUID,
    IN service_items mto_service_item_type[]
) AS '
DECLARE
    s_type mto_shipment_type;
    m_code market_code_enum;
    move_id UUID;
    service_item RECORD;
    item mto_service_item_type;
BEGIN
    -- get the shipment type, market code, and move_id based on shipment_id
    SELECT ms.shipment_type, ms.market_code, ms.move_id
    INTO s_type, m_code, move_id
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF s_type IS NULL OR m_code IS NULL THEN
        RAISE EXCEPTION ''Shipment with ID % not found or missing required details.'', shipment_id;
    END IF;

    -- loop through each provided service item  object
    FOREACH item IN ARRAY service_items
    LOOP
        FOR service_item IN
            SELECT rsi.id,
                   rs.id AS re_service_id,
                   rs.service_location,
                   rsi.is_auto_approved,
                   rs.code AS service_code
            FROM re_service_items rsi
            JOIN re_services rs ON rsi.service_id = rs.id
            WHERE rsi.shipment_type = s_type
              AND rsi.market_code = m_code
              AND rs.code = (item.re_service_code)
              AND rsi.is_auto_approved = false
        LOOP
            BEGIN
            IF NOT does_service_item_exist(service_item.re_service_id, shipment_id) THEN
                INSERT INTO mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    created_at,
                    updated_at,
                    sit_postal_code,
                    sit_entry_date,
                    sit_customer_contacted,
                    reason,
                    estimated_weight,
                    actual_weight,
                    pickup_postal_code,
                    description,
                    sit_destination_original_address_id,
                    sit_destination_final_address_id,
                    sit_requested_delivery,
                    sit_departure_date,
                    sit_origin_hhg_original_address_id,
                    sit_origin_hhg_actual_address_id,
                    customer_expense,
                    customer_expense_reason,
                    sit_delivery_miles,
                    standalone_crate
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    ''SUBMITTED''::service_item_status,
                    NOW(),
                    NOW(),
                    (item).sit_postal_code,
                    (item).sit_entry_date,
                    (item).sit_customer_contacted,
                    (item).reason,
                    (item).estimated_weight,
                    (item).actual_weight,
                    (item).pickup_postal_code,
                    (item).description,
                    (item).sit_destination_original_address_id,
                    (item).sit_destination_final_address_id,
                    (item).sit_requested_delivery,
                    (item).sit_departure_date,
                    (item).sit_origin_hhg_original_address_id,
                    (item).sit_origin_hhg_actual_address_id,
                    (item).customer_expense,
                    (item).customer_expense_reason,
                    (item).sit_delivery_miles,
                    (item).standalone_crate
                );
                END IF;
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION ''Error creating accessorial service item with code % for shipment %: %'',
                                service_item.service_code, shipment_id, SQLERRM;
            END;
        END LOOP;
    END LOOP;
END;
'
LANGUAGE plpgsql;