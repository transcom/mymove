-- Make new nullable fields (we'll set them to non-nullable after we provide values).
ALTER TABLE payment_requests
    ADD COLUMN payment_request_number text,
    ADD COLUMN sequence_number integer;

-- Generate the values for those fields for any existing payment requests.
-- Lock the table to prevent any concurrent updates that may affect our sequencing.
LOCK TABLE payment_requests IN SHARE MODE;

DO
$do$
    DECLARE
        current_move_task_order_id uuid;
        current_payment_request_id uuid;
        mto_reference_id           varchar(255);
        new_sequence_number        int;
    BEGIN
        -- Get all distinct task order IDs associated with payment requests
        FOR current_move_task_order_id IN SELECT DISTINCT move_task_order_id FROM payment_requests
            LOOP
                -- Initialize the reference ID and sequence number for the payment requests for
                -- that task order
                SELECT mto.reference_id
                INTO mto_reference_id
                FROM move_task_orders mto
                WHERE id = current_move_task_order_id;

                new_sequence_number := 1;

                -- For each payment request of that task order, set the payment_request_number and
                -- sequence_number
                FOR current_payment_request_id IN SELECT id
                                                  FROM payment_requests
                                                  WHERE move_task_order_id = current_move_task_order_id
                                                  ORDER BY created_at
                    LOOP
                        UPDATE payment_requests
                        SET payment_request_number = mto_reference_id || '-' || new_sequence_number,
                            sequence_number        = new_sequence_number,
                            updated_at             = now()
                        WHERE id = current_payment_request_id;

                        new_sequence_number := new_sequence_number + 1;
                    END LOOP;
            END LOOP;
    END
$do$;

-- Now make the columns not null and establish the unique constraints.
ALTER TABLE payment_requests
    ALTER COLUMN payment_request_number SET NOT NULL,
    ALTER COLUMN sequence_number SET NOT NULL,
    ADD CONSTRAINT payment_requests_payment_request_number_unique_key UNIQUE (payment_request_number),
    ADD CONSTRAINT payment_requests_sequence_number_unique_key UNIQUE (move_task_order_id, sequence_number);
