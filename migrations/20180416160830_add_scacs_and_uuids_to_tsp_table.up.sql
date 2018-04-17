-- Remove (for now) constraint requiring TSP name field to be not null
ALTER TABLE ONLY transportation_service_providers ALTER COLUMN name DROP not null;

-- Import SCACs and UUIDs into existing transportation_service_providers table
INSERT INTO transportation_service_providers (id, standard_carrier_alpha_code, created_at, updated_at)
SELECT id, scac, CURRENT_TIMESTAMP as created_at, CURRENT_TIMESTAMP as updated_at
FROM transportation_service_provider_data;

DROP TABLE transportation_service_provider_data;
