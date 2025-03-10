--B-22692 Brian Manley added prime_acknowledge_moves_shipments
CREATE OR REPLACE PROCEDURE prime_acknowledge_moves_shipments(json_data jsonb)
LANGUAGE plpgsql
AS '
DECLARE
    v_move_record jsonb;
    v_shipment_record jsonb;
    v_log_message text;
    v_move_id UUID;
    v_move_prime_acknowledged_at TIMESTAMP;
    v_shipment_id UUID;
    v_shipment_prime_acknowledged_at TIMESTAMP;
    v_move_exists BOOLEAN;
    v_shipment_exists BOOLEAN;
BEGIN
    RAISE NOTICE ''Starting prime_acknowledge_moves_shipments procedure'';

    -- Check if the input JSON is null or empty
    IF json_data IS NULL OR json_data = ''[]''::jsonb THEN
        RAISE WARNING ''Input JSON is null or empty'';
        RETURN;
    END IF;

    -- Loop through each move in the JSON array
    FOR v_move_record IN SELECT jsonb_array_elements(json_data)
    LOOP
        -- Check if the move exists in the moves table
        v_move_id := NULLIF(v_move_record->>''id'', '''')::UUID;
        SELECT EXISTS(SELECT 1 FROM moves WHERE id = v_move_id) INTO v_move_exists;
        IF NOT v_move_exists THEN
            v_log_message := format(''Move with id %s does not exist in the moves table'', v_move_id);
            RAISE WARNING ''%'', v_log_message;
            CONTINUE;
        END IF;

        v_move_prime_acknowledged_at := NULLIF(v_move_record->>''primeAcknowledgedAt'', '''')::TIMESTAMP;
        IF v_move_prime_acknowledged_at IS NOT NULL THEN
            -- Update the moves table only if the existing prime_acknowledged_at value is NULL
            UPDATE moves
            SET prime_acknowledged_at = v_move_prime_acknowledged_at
            WHERE id = v_move_id
              AND prime_acknowledged_at IS NULL;

            -- Check if the update affected any rows
            IF NOT FOUND THEN
                v_log_message := format(''Move not updated (prime_acknowledged_at already set): %s'', v_move_id);
                RAISE WARNING ''%'', v_log_message;
            ELSE
                v_log_message := format(''Successfully updated moves.prime_acknowledged_at value to %s for id %s'', v_move_prime_acknowledged_at, v_move_id);
                RAISE NOTICE ''%'', v_log_message;
            END IF;
        END IF;

        -- Check if mtoShipments exists and is an array
        IF jsonb_typeof(v_move_record->''mtoShipments'') = ''array'' THEN
            -- Loop through each shipment in the mtoShipments array
            FOR v_shipment_record IN SELECT jsonb_array_elements(v_move_record->''mtoShipments'')
            LOOP
                v_shipment_id := NULLIF(v_shipment_record->>''id'', '''')::UUID;
                SELECT EXISTS(SELECT 1 FROM mto_shipments WHERE id = v_shipment_id and move_id = v_move_id) INTO v_shipment_exists;
                IF NOT v_shipment_exists THEN
                    v_log_message := format(''Shipment with id %s and move_id %s does not exist in the mto_shipments table'', v_shipment_id, v_move_id);
                    RAISE WARNING ''%'', v_log_message;
                    CONTINUE;
                END IF;

                v_shipment_prime_acknowledged_at := NULLIF(v_shipment_record->>''primeAcknowledgedAt'', '''')::TIMESTAMP;
                IF v_shipment_prime_acknowledged_at IS NOT NULL THEN
                    -- Update the mto_shipments table only if the existing prime_acknowledged_at column is NULL
                    UPDATE mto_shipments
                    SET prime_acknowledged_at = v_shipment_prime_acknowledged_at
                    WHERE id = v_shipment_id
                        AND prime_acknowledged_at IS NULL;

                    -- Check if the update affected any rows
                    IF NOT FOUND THEN
                        v_log_message := format(''Shipment not updated (prime_acknowledged_at already set): %s'', v_shipment_id);
                        RAISE WARNING ''%'', v_log_message;
                    ELSE
                        v_log_message := format(''Successfully updated mto_shipments.prime_acknowledged_at value to %s for id %s'', v_shipment_prime_acknowledged_at, v_shipment_id);
                        RAISE NOTICE ''%'', v_log_message;
                    END IF;
                END IF;
            END LOOP;
        END IF;
    END LOOP;
END;
';