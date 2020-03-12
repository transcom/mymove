ALTER TABLE client_certs
    ALTER COLUMN sha256_digest SET NOT NULL,
    ALTER COLUMN subject SET NOT NULL,
    ALTER COLUMN allow_dps_auth_api SET NOT NULL,
    ALTER COLUMN allow_orders_api SET NOT NULL;
