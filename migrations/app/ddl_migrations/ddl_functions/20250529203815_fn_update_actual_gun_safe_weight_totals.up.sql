-- B-23373 Brooklyn Welsh - Create function to update actual_gun_safe_weight column in mto_shipments table, used whenever a gun safe weight ticket is added or deleted. Accounts for all tickets linked to the shipment
DROP FUNCTION IF EXISTS update_actual_gun_safe_weight_totals(uuid);

CREATE OR REPLACE FUNCTION update_actual_gun_safe_weight_totals(ppm UUID)
RETURNS void
AS
$$
BEGIN
	    UPDATE mto_shipments
	    SET
	        actual_gun_safe_weight = NULLIF(t.weight, 0)
	    FROM (
	        SELECT
	            SUM(gw.weight::int) AS weight
	        FROM gunsafe_weight_tickets gw
	        WHERE gw.ppm_shipment_id = ppm AND gw.deleted_at IS NULL
	    ) t
	    WHERE id = (SELECT shipment_id FROM ppm_shipments WHERE id = ppm);
END;
$$
LANGUAGE plpgsql;
