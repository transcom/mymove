ALTER TABLE client_certs ADD COLUMN IF NOT EXISTS allow_pptas bool DEFAULT false NOT NULL;

COMMENT ON COLUMN client_certs.allow_pptas IS 'Indicates whether or not the cert grants access to the PPTAS API';
