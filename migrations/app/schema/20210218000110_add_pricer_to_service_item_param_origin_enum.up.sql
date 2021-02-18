-- Change origin from request_origin to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN origin TYPE VARCHAR(255);

-- Drop and create again request_origin enum:
DROP TYPE IF EXISTS service_item_param_origin;
CREATE TYPE service_item_param_origin AS ENUM (
    'PRIME',
    'SYSTEM',
    'PRICER'
    );

-- Revert origin from varchar to request_origin for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN origin TYPE service_item_param_origin USING (origin::service_item_param_origin);
