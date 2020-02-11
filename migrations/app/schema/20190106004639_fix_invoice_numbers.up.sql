-- Lock both the invoice_number_trackers and invoices tables to prevent concurrent updates.  Just
-- trying to prevent another invoice number from landing in the middle of this transaction.
LOCK TABLE invoice_number_trackers, invoices IN SHARE MODE;

-- Delete all current invoice number trackers just in case something slipped in before this migration
-- runs.  We're regenerating all invoice numbers, so it should be OK to reset these sequence numbers.
DELETE
FROM invoice_number_trackers;

-- Fix any invoices already in the system since they have placeholder invoice numbers.
DO $do$
  DECLARE
    current_shipment_id     UUID;
    current_invoice_id      UUID;
    scac                    TEXT;
    invoice_count           INT;
    shipment_year           INT;
    shipment_two_digit_year VARCHAR(2);
    base_invoice_number     VARCHAR(255);
    new_sequence_number     INT;
    target_invoice_number   VARCHAR(255);
  BEGIN
    -- Get all distinct shipment IDs currently associated with invoices, ordered by earliest invoice creation.
    FOR current_shipment_id IN SELECT DISTINCT shipment_id, MIN(created_at) as min_created_at
                               FROM invoices
                               GROUP BY shipment_id
                               ORDER BY MIN(created_at)
      LOOP
        scac := NULL;

        -- Get the SCAC and creation year for the current_shipment_id
        SELECT tsp.standard_carrier_alpha_code,
               EXTRACT(YEAR FROM s.created_at),
               to_char(s.created_at, 'YY')
               INTO scac, shipment_year, shipment_two_digit_year
        FROM shipments s
               INNER JOIN shipment_offers so ON s.id = so.shipment_id
               INNER JOIN transportation_service_provider_performances tspp
                          ON so.transportation_service_provider_performance_id = tspp.id
               INNER JOIN transportation_service_providers tsp ON tspp.transportation_service_provider_id = tsp.id
        WHERE s.id = current_shipment_id
          AND so.accepted = TRUE
        ORDER BY so.created_at
        LIMIT 1;

        -- If scac is null, that probably means there was no accepted shipment offer for some strange reason.
        -- Raise exception; migration will likely fail anyway trying to apply the invoice number unique constraint.
        IF scac IS NULL THEN
          RAISE EXCEPTION 'Shipment ID % has no accepted shipment offer, so unable to generate proper invoice number.', current_shipment_id;
        END IF;

        -- Get all invoice records for that shipment, ordered by creation date.
        invoice_count := 0;
        base_invoice_number := NULL;
        FOR current_invoice_id IN SELECT id FROM invoices WHERE shipment_id = current_shipment_id ORDER BY created_at
          LOOP
            IF invoice_count = 0 THEN
              -- Set the first invoice number to the baseline invoice number.

              -- Insert/update sequence number.
              INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
              VALUES (scac, shipment_year, 1)
              ON CONFLICT (
                 standard_carrier_alpha_code,
                 year)
                 DO
                   UPDATE
                   SET sequence_number = trackers.sequence_number + 1
                   WHERE trackers.standard_carrier_alpha_code = scac AND trackers.year = shipment_year
                 RETURNING sequence_number INTO new_sequence_number;

              base_invoice_number := scac || shipment_two_digit_year || to_char(new_sequence_number, 'fm0000');
              target_invoice_number := base_invoice_number;
            ELSE
              -- Set subsequent invoice numbers to the baseline number suffixed by "-01", "-02", etc.
              target_invoice_number := base_invoice_number || '-' || to_char(invoice_count, 'fm00');
            END IF;

            -- Update the invoice_number for this invoice to the target number determined above.
            UPDATE invoices
            SET invoice_number = target_invoice_number,
                updated_at     = now()
            WHERE id = current_invoice_id;

            invoice_count := invoice_count + 1;
          END LOOP;
      END LOOP;
  END $do$;

-- Now add a unique constraint for the invoice number and set the shipment_id to non-nullable.
ALTER TABLE invoices
  ADD CONSTRAINT unique_invoice_number UNIQUE (invoice_number),
  ALTER COLUMN shipment_id SET NOT NULL;
