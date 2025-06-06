--B-22704 Michael Saki adding a function to update the progear and spousal progear actual weight totals in mto_shipments

DROP FUNCTION IF EXISTS update_actual_progear_weight_totals(uuid);

CREATE OR REPLACE FUNCTION update_actual_progear_weight_totals(ppm UUID)
RETURNS void
AS '
BEGIN
	    UPDATE mto_shipments
	    SET
	        actual_pro_gear_weight = NULLIF(t.self_weight, 0),
	        actual_spouse_pro_gear_weight = NULLIF(t.spouse_weight, 0)
	    FROM (
	        SELECT
	            SUM(CASE WHEN pw.belongs_to_self = TRUE THEN pw.weight::int ELSE 0 END) AS self_weight,
	            SUM(CASE WHEN pw.belongs_to_self = FALSE OR pw.belongs_to_self IS NULL THEN pw.weight::int ELSE 0 END) AS spouse_weight
	        FROM progear_weight_tickets pw
	        WHERE pw.ppm_shipment_id = ppm AND pw.deleted_at IS NULL
	    ) t
	    WHERE id = (SELECT shipment_id FROM ppm_shipments WHERE id = ppm);
END;
'
LANGUAGE plpgsql
