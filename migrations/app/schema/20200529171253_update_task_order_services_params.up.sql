-- Change type from request_type to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE VARCHAR(255);

-- Drop and create again request_type enum:
DROP TYPE IF EXISTS service_item_param_type;
CREATE TYPE service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL',
    'PaymentServiceItemUUID',
    'TIMESTAMP'
    );

-- Revert type from varchar to request_type for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE service_item_param_type USING (type::service_item_param_type);

INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('958e43d9-a10c-4cf9-9737-4103f9d2de29','MTOAvailableToPrimeAt', 'Timestamp MTO was made available to prime', 'TIMESTAMP', 'SYSTEM', now(), now());

UPDATE service_params
SET service_item_param_key_id = (SELECT id FROM service_item_param_keys where key='MTOAvailableToPrimeAt')
WHERE service_id in (SELECT id FROM re_services WHERE code='MS' or code ='CS');