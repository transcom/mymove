-- Import SCACs and UUIDs into existing transportation_service_providers table
SELECT
  id,
  scac
  INTO transportation_service_providers
  FROM transportation_service_provider_data;

DROP TABLE transportation_service_provider_data;
