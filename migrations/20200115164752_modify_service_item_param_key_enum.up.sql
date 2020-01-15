-- https://stackoverflow.com/a/56376907

-- Change type from request_type to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE VARCHAR(255);

-- Drop and create again request_type enum:
DROP TYPE IF EXISTS service_item_param_type;
CREATE TYPE service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL',
    'PaymentServiceItemUUID'
    );

-- Revert type from varchar to request_type for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE service_item_param_type USING (type::service_item_param_type);
