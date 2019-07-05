ALTER TABLE service_members ADD COLUMN require_access_code boolean;

-- Updates require_access_code to false
DO $$
BEGIN
	FOR rec IN SELECT id FROM service_members
	LOOP
		UPDATE service_members SET require_access_code=false WHERE id=rec.id;
	END LOOP;
END $$;

ALTER TABLE service_members ALTER COLUMN require_access_code NOT NULL;
