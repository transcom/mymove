SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS allowable_weight integer;

COMMENT on COLUMN ppm_shipments.allowable_weight IS 'Combined allowable weight for all trips.';

UPDATE ppm_shipments
SET allowable_weight = summed_weights.summed_allowable_weight
FROM (
	SELECT ppm_shipment_id, SUM(full_weight - empty_weight) AS summed_allowable_weight FROM public.weight_tickets
	GROUP BY ppm_shipment_id
) AS summed_weights
WHERE ppm_shipments.id = summed_weights.ppm_shipment_id;

ALTER TABLE weight_tickets
    DROP COLUMN allowable_weight;
