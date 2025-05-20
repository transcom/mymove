--B-22704 Michael Saki adding a function to update the progear and spousal progear actual weight totals in mto_shipments

CREATE OR REPLACE FUNCTION update_actual_progear_weight_totals(ppm UUID)
RETURNS void
LANGUAGE plpgsql AS $$
BEGIN
    update mto_shipments
        SET actual_pro_gear_weight = NULLIF((
            SELECT coalesce(sum(pw.weight::int), 0)
            FROM progear_weight_tickets pw
            WHERE
                pw.ppm_shipment_id = ppm
                AND pw.deleted_at IS NULL
                AND pw.belongs_to_self = true
        ), 0),
        actual_spouse_pro_gear_weight = NULLIF((
            SELECT coalesce(sum(pw.weight::int), 0)
            FROM progear_weight_tickets pw
            WHERE
                pw.ppm_shipment_id = ppm
                AND pw.deleted_at IS NULL
                AND (pw.belongs_to_self = FALSE OR pw.belongs_to_self IS NULL)
        ), 0)
    WHERE id = (
        SELECT shipment_id
        FROM ppm_shipments
        WHERE id = ppm
    );
END;
$$;
