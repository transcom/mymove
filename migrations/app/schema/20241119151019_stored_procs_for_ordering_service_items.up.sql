CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE mto_service_items
  ALTER COLUMN id SET DEFAULT uuid_generate_v4();


-- creating function to get address is_oconus
CREATE OR REPLACE FUNCTION get_is_oconus(address_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    is_oconus BOOLEAN;
BEGIN
    SELECT a.is_oconus
    INTO is_oconus
    FROM addresses a
    WHERE a.id = address_id;

    RETURN is_oconus;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE EXCEPTION 'Address with ID % not found', address_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION does_service_item_exist(
    service_id UUID,
    shipment_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM mto_service_items
        WHERE re_service_id = service_id
        AND mto_shipment_id = shipment_id
    ) THEN
        RAISE EXCEPTION 'Service item already exists for service_id % and shipment_id %', service_id, shipment_id;
    END IF;
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- stored proc that creates auto-approved service items based off of a shipment id
CREATE OR REPLACE PROCEDURE create_approved_service_items_for_shipment(
    IN shipment_id UUID
)
AS '
DECLARE
    s_status mto_shipment_status;
    s_type mto_shipment_type;
    m_code market_code_enum;
    move_id UUID;
    pickup_address_id UUID;
    destination_address_id UUID;
    is_pickup_oconus BOOLEAN;
    is_destination_oconus BOOLEAN;
    service_item RECORD;
BEGIN
    -- get shipment type, market code, move_id, and address IDs based on shipment_id
    SELECT ms.shipment_type, ms.market_code, ms.move_id, ms.pickup_address_id, ms.destination_address_id, ms.status
    INTO s_type, m_code, move_id, pickup_address_id, destination_address_id, s_status
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF s_type IS NULL OR m_code IS NULL THEN
        RAISE EXCEPTION ''Shipment with ID % not found or missing required details.'', shipment_id;
    END IF;

    IF s_status  IN (''APPROVED'') THEN
        RAISE EXCEPTION ''Shipment with ID % is already in APPROVED status'', shipment_id;
    END IF;

    -- get the is_oconus values for both pickup and destination addresses - this determines POD/POE creation
    is_pickup_oconus := get_is_oconus(pickup_address_id);
    is_destination_oconus := get_is_oconus(destination_address_id);

    -- determine which service item to create based on shipment direction
    -- first create the direction-specific service item (POEFSC or PODFSC)
    IF is_pickup_oconus AND NOT is_destination_oconus THEN
        -- Shipment is OCONUS to CONUS, create PODFSC item
        FOR service_item IN
            SELECT rsi.id,
                   rs.id AS re_service_id,
                   rs.service_location,
                   rsi.is_auto_approved
            FROM re_service_items rsi
            JOIN re_services rs ON rsi.service_id = rs.id
            WHERE rsi.shipment_type = s_type
              AND rsi.market_code = m_code
              AND rs.code = ''PODFSC''
              AND rsi.is_auto_approved = true
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
                    approved_at
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    ''APPROVED''::service_item_status,
                    NOW(),
                    NOW(),
                    NOW()
                );
                END IF;
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION ''Error creating PODFSC service item for shipment %: %'', shipment_id, SQLERRM;
            END;
        END LOOP;
    ELSIF NOT is_pickup_oconus AND is_destination_oconus THEN
        -- Shipment is CONUS to OCONUS, create POEFSC item
        FOR service_item IN
            SELECT rsi.id,
                   rs.id AS re_service_id,
                   rs.service_location,
                   rsi.is_auto_approved
            FROM re_service_items rsi
            JOIN re_services rs ON rsi.service_id = rs.id
            WHERE rsi.shipment_type = s_type
              AND rsi.market_code = m_code
              AND rs.code = ''POEFSC''
              AND rsi.is_auto_approved = true
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
                    approved_at
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    ''APPROVED''::service_item_status,
                    NOW(),
                    NOW(),
                    NOW()
                );
                END IF;
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION ''Error creating POEFSC service item for shipment %: %'', shipment_id, SQLERRM;
            END;
        END LOOP;
    ELSE
        RAISE EXCEPTION ''Invalid shipment direction for shipment %: Pickup is %CONUS, Destination is %CONUS.'',
                         shipment_id, is_pickup_oconus, is_destination_oconus;
    END IF;

    -- create all other auto-approved service items, filtering out the POEFSC or PODFSC service items
    FOR service_item IN
        SELECT rsi.id,
               rs.id AS re_service_id,
               rs.service_location,
               rsi.is_auto_approved
        FROM re_service_items rsi
        JOIN re_services rs ON rsi.service_id = rs.id
        WHERE rsi.shipment_type = s_type
          AND rsi.market_code = m_code
          AND rsi.is_auto_approved = true
          AND rs.code NOT IN (''POEFSC'', ''PODFSC'')
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
                approved_at
            )
            VALUES (
                shipment_id,
                move_id,
                service_item.re_service_id,
                service_item.service_location,
                ''APPROVED''::service_item_status,
                NOW(),
                NOW(),
                NOW()
            );
            End IF;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE EXCEPTION ''Error creating other service item for shipment %: %'', shipment_id, SQLERRM;
        END;
    END LOOP;
END;
'
LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'mto_service_item_type') THEN
        CREATE TYPE  mto_service_item_type  AS (
    id uuid,
            move_id uuid,
            mto_shipment_id uuid,
            re_service_id uuid,
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
END IF;
END
$$;


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

    IF s_type <> item.shipment_type THEN
    RAISE EXCEPTION ''Shipment type mismatch. Expected %, but got %.'', s_type, item.shipment_type;
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
              AND rs.id = (item.re_service_id)
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