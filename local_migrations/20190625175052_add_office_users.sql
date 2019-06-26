-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

--Add CO MCLB ALBANY GA, transportation office and address
INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES ('5081b952-0154-4776-bce3-901e96a4774d', '814 Radford Blvd', 'Ste 20352', 'Albany', 'GA', '31704', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('d25be014-95ee-44f0-ac31-2f37160ed503', 'CO MCLB ALBANY GA', 'CNNQ', '297859bc-004d-4f11-9ac9-a657df596f25', 31.550341, -84.053009, now(), now());