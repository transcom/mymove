-- with the deprecation of the createSITAddressUpdateRequest endpoint
-- sit_address_updates will no longer be used
DROP TABLE IF EXISTS sit_address_updates;

-- also going to drop the status enum created in 20230504204015_creating_SIT_destination_address_update_table.up.sql
DROP TYPE IF EXISTS sit_address_update_status;