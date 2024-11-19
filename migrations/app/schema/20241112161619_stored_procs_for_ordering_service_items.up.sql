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


-- stored proc that creates auto-approved shipments based off of a shipment id
CREATE OR REPLACE PROCEDURE CreateApprovedServiceItemsForShipment(
    IN shipment_id UUID
)
AS $$
DECLARE
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
    SELECT ms.shipment_type, ms.market_code, ms.move_id, ms.pickup_address_id, ms.destination_address_id
    INTO s_type, m_code, move_id, pickup_address_id, destination_address_id
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF s_type IS NULL OR m_code IS NULL THEN
        RAISE EXCEPTION 'Shipment with ID % not found or missing required details.', shipment_id;
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
              AND rs.code = 'PODFSC'
              AND rsi.is_auto_approved = true
        LOOP
            BEGIN
                INSERT INTO mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    created_at,
                    updated_at
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    'APPROVED'::service_item_status,
                    NOW(),
                    NOW()
                );
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION 'Error creating PODFSC service item for shipment %: %', shipment_id, SQLERRM;
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
              AND rs.code = 'POEFSC'
              AND rsi.is_auto_approved = true
        LOOP
            BEGIN
                INSERT INTO mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    created_at,
                    updated_at
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    'APPROVED'::service_item_status,
                    NOW(),
                    NOW()
                );
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION 'Error creating POEFSC service item for shipment %: %', shipment_id, SQLERRM;
            END;
        END LOOP;
    ELSE
        RAISE EXCEPTION 'Invalid shipment direction for shipment %: Pickup is %CONUS, Destination is %CONUS.',
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
          AND rs.code NOT IN ('POEFSC', 'PODFSC')
    LOOP
        BEGIN
            INSERT INTO mto_service_items (
                mto_shipment_id,
                move_id,
                re_service_id,
                service_location,
                status,
                created_at,
                updated_at
            )
            VALUES (
                shipment_id,
                move_id,
                service_item.re_service_id,
                service_item.service_location,
                'APPROVED'::service_item_status,
                NOW(),
                NOW()
            );
        EXCEPTION
            WHEN OTHERS THEN
                RAISE EXCEPTION 'Error creating other service item for shipment %: %', shipment_id, SQLERRM;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE
OR REPLACE PROCEDURE CreateAccessorialServiceItems (
    IN shipment_id UUID,
    IN service_items JSONB[] -- Changed from TEXT[] to JSONB[]
) AS $$
DECLARE
s_type mto_shipment_type;
    m_code market_code_enum;
    move_id UUID;
    service_item RECORD;
    item JSONB;
BEGIN
    -- get the shipment type, market code, and move_id based on shipment_id
    SELECT ms.shipment_type, ms.market_code, ms.move_id
    INTO s_type, m_code, move_id
    FROM mto_shipments ms
    WHERE ms.id = shipment_id;

    IF s_type IS NULL OR m_code IS NULL THEN
        RAISE EXCEPTION 'Shipment with ID % not found or missing required details.', shipment_id;
    END IF;

    -- loop through each provided service item JSON object
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
              AND rs.code = (item->>'code')::text
              AND rsi.is_auto_approved = false
        LOOP
            BEGIN
                INSERT INTO mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    created_at,
                    updated_at,
                    sit_entry_date,
                    sit_customer_contacted
                    --and other sit related fields
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    'SUBMITTED'::service_item_status,
                    NOW(),
                    NOW(),
                    (item->>'sit_entry_date')::date,
                    (item->>'sit_customer_contacted')::date

                );
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION 'Error creating accessorial service item with code % for shipment %: %',
                                service_item.service_code, shipment_id, SQLERRM;
            END;
        END LOOP;
    END LOOP;
END;
$$ LANGUAGE plpgsql;