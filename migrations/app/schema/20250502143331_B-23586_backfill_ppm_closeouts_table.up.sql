DO $$
DECLARE
  ppm_id uuid;
BEGIN
  FOR ppm_id IN
    SELECT id
    FROM ppm_shipments
    WHERE status = 'CLOSEOUT_COMPLETE'
  LOOP
    CALL calculate_ppm_closeout(ppm_id, true);
  END LOOP;
END;
$$;
