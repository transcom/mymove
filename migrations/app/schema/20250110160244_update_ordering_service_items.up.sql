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
    -- collect the service items into a temporary table for sorting
    CREATE TEMPORARY TABLE temp_mto_service_items (
        mto_shipment_id uuid,
        move_id uuid,
        re_service_id uuid,
        service_location service_location_enum,
        status service_item_status,
        sort text
    ) ON COMMIT DROP;

    -- first create the direction-specific service item (POEFSC or PODFSC)
    IF is_pickup_oconus AND NOT is_destination_oconus THEN
        -- Shipment is OCONUS to CONUS, create PODFSC item
        FOR service_item IN
            SELECT rsi.id,
                   rs.id AS re_service_id,
                   rs.service_location,
                   rsi.is_auto_approved,
                   rsi.sort
            FROM re_service_items rsi
            JOIN re_services rs ON rsi.service_id = rs.id
            WHERE rsi.shipment_type = s_type
              AND rsi.market_code = m_code
              AND rs.code = ''PODFSC''
              AND rsi.is_auto_approved = true
        LOOP
            BEGIN
            IF NOT does_service_item_exist(service_item.re_service_id, shipment_id) THEN
                INSERT INTO temp_mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    sort
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    ''APPROVED''::service_item_status,
                    service_item.sort
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
                   rsi.is_auto_approved,
                   rsi.sort
            FROM re_service_items rsi
            JOIN re_services rs ON rsi.service_id = rs.id
            WHERE rsi.shipment_type = s_type
              AND rsi.market_code = m_code
              AND rs.code = ''POEFSC''
              AND rsi.is_auto_approved = true
        LOOP
            BEGIN
           IF NOT does_service_item_exist(service_item.re_service_id, shipment_id) THEN
                INSERT INTO temp_mto_service_items (
                    mto_shipment_id,
                    move_id,
                    re_service_id,
                    service_location,
                    status,
                    sort
                )
                VALUES (
                    shipment_id,
                    move_id,
                    service_item.re_service_id,
                    service_item.service_location,
                    ''APPROVED''::service_item_status,
                    service_item.sort
                );
                END IF;
            EXCEPTION
                WHEN OTHERS THEN
                    RAISE EXCEPTION ''Error creating POEFSC service item for shipment %: %'', shipment_id, SQLERRM;
            END;
        END LOOP;
    END IF;

    -- create all other auto-approved service items, filtering out the POEFSC or PODFSC service items
    FOR service_item IN
        SELECT rsi.id,
               rs.id AS re_service_id,
               rs.service_location,
               rsi.is_auto_approved,
               rsi.sort
        FROM re_service_items rsi
        JOIN re_services rs ON rsi.service_id = rs.id
        WHERE rsi.shipment_type = s_type
          AND rsi.market_code = m_code
          AND rsi.is_auto_approved = true
          AND rs.code NOT IN (''POEFSC'', ''PODFSC'')
    LOOP
        BEGIN
         IF NOT does_service_item_exist(service_item.re_service_id, shipment_id) THEN
            INSERT INTO temp_mto_service_items (
                mto_shipment_id,
                move_id,
                re_service_id,
                service_location,
                status,
                sort
            )
            VALUES (
                shipment_id,
                move_id,
                service_item.re_service_id,
                service_item.service_location,
                ''APPROVED''::service_item_status,
                service_item.sort
            );
            End IF;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE EXCEPTION ''Error creating other service item for shipment %: %'', shipment_id, SQLERRM;
        END;
    END LOOP;

    -- Insert the mto_service_items in order
    FOR service_item IN
        SELECT tmsi.mto_shipment_id,
               tmsi.move_id,
               tmsi.re_service_id,
               tmsi.service_location,
               tmsi.status
        FROM temp_mto_service_items tmsi
        ORDER BY sort
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
                approved_at
            )
            VALUES (
                service_item.mto_shipment_id,
                service_item.move_id,
                service_item.re_service_id,
                service_item.service_location,
                service_item.status,
                NOW(),
                NOW(),
                NOW()
            );
        EXCEPTION
            WHEN OTHERS THEN
               RAISE EXCEPTION ''Error creating service items from temp table for shipment %: %'', shipment_id, SQLERRM;
        END;
    END LOOP;
END;
'
LANGUAGE plpgsql;
