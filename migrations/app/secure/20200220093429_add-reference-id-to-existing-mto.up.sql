-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.
LOCK TABLE move_task_orders IN SHARE MODE;

DO
$do$
    DECLARE
        current_move_task_order_id UUID;
        new_ref_id                 varchar(10);
        ref_id_count               int;
    BEGIN
        -- loop over all move task orders with reference_id null
        FOR current_move_task_order_id IN SELECT id FROM move_task_orders WHERE reference_id IS NULL
            LOOP
                LOOP
                    -- generate a random reference_id based on xxxx-xxxx
                    SELECT floor(random() * 9000 + 1000) || '-' || floor(random() * 9000 + 1000) INTO new_ref_id;

                    -- look up to see if reference_id is already being used by an MTO
                    SELECT COUNT(*) INTO ref_id_count FROM move_task_orders WHERE reference_id = new_ref_id;

                    -- if there are no collisions then break out of loop
                    IF ref_id_count = 0 THEN
                        EXIT;
                    END IF;
                END LOOP;

                UPDATE move_task_orders SET reference_id = new_ref_id WHERE id = current_move_task_order_id;
            END LOOP;
    END
$do$;
