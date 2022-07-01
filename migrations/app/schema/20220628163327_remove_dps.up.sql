-- DPS is no longer needed
DROP TABLE dps_users;

ALTER TABLE client_certs DROP COLUMN allow_dps_auth_api;
