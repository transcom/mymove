-- Change type to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE VARCHAR(255);

-- Drop and create again service_item_param_type enum:
DROP TYPE IF EXISTS service_item_param_type;
CREATE TYPE service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL',
    'TIMESTAMP',
    'PaymentServiceItemUUID',
    'BOOLEAN'
    );

-- Revert type from varchar to service_item_param_type for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE service_item_param_type USING (type::service_item_param_type);
