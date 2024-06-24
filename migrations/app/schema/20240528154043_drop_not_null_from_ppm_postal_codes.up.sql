ALTER TABLE ppm_shipments
ALTER COLUMN pickup_postal_code DROP NOT NULL,
ALTER COLUMN destination_postal_code DROP NOT NULL;