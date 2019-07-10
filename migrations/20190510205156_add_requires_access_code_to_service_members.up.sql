ALTER TABLE service_members ADD COLUMN requires_access_code boolean;

-- Updates requires_access_code to false
DO $$
DECLARE
	rec RECORD;
BEGIN
	FOR rec IN SELECT id FROM service_members
	LOOP
		UPDATE service_members SET requires_access_code=false WHERE id=rec.id;
	END LOOP;
END $$;

ALTER TABLE service_members ALTER COLUMN requires_access_code SET NOT NULL;
