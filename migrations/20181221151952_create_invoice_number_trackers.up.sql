CREATE TABLE invoice_number_trackers (
    standard_carrier_alpha_code text NOT NULL,
    year integer NOT NULL,
    sequence_number integer NOT NULL,
    PRIMARY KEY(standard_carrier_alpha_code, year)
);

-- Fix any invoices already in the system that may have placeholder invoice numbers.
DO
  $do$
    DECLARE
      invoice                          RECORD;
      scac                             TEXT;
      processed_shipment_invoice_count INT;
      shipment_year                    INT;
      shipment_two_digit_year          VARCHAR(2);
      base_invoice_number              VARCHAR(255);
      new_sequence_number              INT;
    BEGIN

      -- Go through all the invoices
      FOR invoice IN SELECT * FROM invoices ORDER BY created_at
        LOOP
          -- Get the SCAC for the invoice
          SELECT tsp.standard_carrier_alpha_code INTO scac
          FROM invoices i
                 INNER JOIN shipments s on i.shipment_id = s.id
                 INNER JOIN shipment_offers so on s.id = so.shipment_id
                 INNER JOIN transportation_service_provider_performances tspp
                            on so.transportation_service_provider_performance_id = tspp.id
                 INNER JOIN transportation_service_providers tsp on tspp.transportation_service_provider_id = tsp.id
          WHERE i.shipment_id = invoice.shipment_id
            AND so.accepted = TRUE;

          -- No need to update this row if already using SCAC.
          CONTINUE WHEN invoice.invoice_number LIKE scac || '%';

          -- Do we already have correct invoice numbers for this shipment?
          -- They may have already been processed in this loop.
          SELECT count(*) INTO processed_shipment_invoice_count
          FROM invoices i
                 INNER JOIN shipments s on i.shipment_id = s.id
          WHERE i.shipment_id = invoice.shipment_id
            AND invoice_number LIKE scac || '%';

          IF processed_shipment_invoice_count = 0 THEN
            -- Get shipment's created at date.
            SELECT EXTRACT(YEAR FROM created_at), to_char(created_at, 'YY') INTO shipment_year, shipment_two_digit_year
            FROM shipments s
            WHERE s.ID = invoice.shipment_id;

            -- Insert/update sequence number.
            INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
            VALUES (scac, shipment_year, 1)
            ON CONFLICT (standard_carrier_alpha_code, year)
               DO
                 UPDATE
                 SET sequence_number = trackers.sequence_number + 1
                 WHERE trackers.standard_carrier_alpha_code = scac AND trackers.year = shipment_year
               RETURNING sequence_number INTO new_sequence_number;

            -- Update baseline invoice number.
            UPDATE invoices
            SET invoice_number = scac || shipment_two_digit_year || to_char(new_sequence_number, 'fm0000')
            WHERE id = invoice.id;
          ELSE
            -- Get baseline invoice number.
            SELECT invoice_number INTO base_invoice_number
            FROM invoices i
                   INNER JOIN shipments s on i.shipment_id = s.id
            WHERE i.shipment_id = invoice.shipment_id
              AND invoice_number LIKE scac || '%'
              AND invoice_number NOT LIKE '%-%';

            -- Update invoice number for subsequent invoices.
            UPDATE invoices
            SET invoice_number = base_invoice_number || '-' || to_char(processed_shipment_invoice_count, 'fm00')
            WHERE id = invoice.id;
          END IF;
        END LOOP;

    END
    $do$;

-- Add a unique constraint for the invoice number.
ALTER TABLE invoices ADD CONSTRAINT unique_invoice_number UNIQUE (invoice_number);
