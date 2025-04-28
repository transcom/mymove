-- B-22945 Paul Stonebraker remove actual postal code columns from ppm_shipments table
ALTER TABLE ppm_shipments
    DROP COLUMN IF EXISTS actual_pickup_postal_code,
    DROP COLUMN IF EXISTS actual_destination_postal_code;