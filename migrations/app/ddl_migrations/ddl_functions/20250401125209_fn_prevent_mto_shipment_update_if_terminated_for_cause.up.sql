-- b-22910 cam
-- add a new trigger to protect mto_shipment data from modification if it is terminated for cause
-- this is an extra precaution if an update operation makes it past the server
CREATE OR REPLACE FUNCTION prevent_update_if_terminated_for_cause()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.status = 'TERMINATED_FOR_CAUSE' THEN
    RAISE EXCEPTION 'Cannot update mto_shipments row: shipment status is TERMINATED_FOR_CAUSE and is protected from update operations';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS mto_shipments_prevent_update_if_terminated ON mto_shipments;
CREATE TRIGGER mto_shipments_prevent_update_if_terminated
BEFORE UPDATE ON mto_shipments
FOR EACH ROW
EXECUTE FUNCTION prevent_update_if_terminated_for_cause();