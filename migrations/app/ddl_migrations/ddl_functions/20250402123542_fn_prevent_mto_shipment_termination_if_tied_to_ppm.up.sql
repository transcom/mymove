-- b-22910 cam
-- add a new trigger to protect mto_shipments from termination if they're
-- tied to a PPM
CREATE OR REPLACE FUNCTION prevent_mto_shipment_termination_if_tied_to_ppm()
RETURNS TRIGGER AS $$
DECLARE
  v_ppm_exists BOOLEAN;
BEGIN
  -- Only block if we're trying to set the status to TERMINATED_FOR_CAUSE
  IF NEW.status = 'TERMINATED_FOR_CAUSE' THEN
    -- Lookup to see if a PPM is tied to this shipment
    SELECT EXISTS (
      SELECT 1
      FROM ppm_shipments
      WHERE shipment_id = NEW.id
    ) INTO v_ppm_exists;

    IF v_ppm_exists THEN
      RAISE EXCEPTION 'Cannot update mto_shipments row: Cannot set status to TERMINATED_FOR_CAUSE: shipment is associated with a PPM and PPMs do not qualify for termination';
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS mto_shipments_prevent_termination_if_tied_to_ppm ON mto_shipments;
CREATE TRIGGER mto_shipments_prevent_termination_if_tied_to_ppm
BEFORE UPDATE ON mto_shipments
FOR EACH ROW
EXECUTE FUNCTION prevent_mto_shipment_termination_if_tied_to_ppm();