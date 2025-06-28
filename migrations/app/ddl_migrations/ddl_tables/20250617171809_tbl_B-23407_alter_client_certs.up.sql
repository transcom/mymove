-- B-23407 Anthony Mann adding PPTAS Affiliation column for client certs

ALTER TABLE client_certs
ADD COLUMN IF NOT EXISTS pptas_affiliation TEXT NULL;

COMMENT ON COLUMN client_certs.pptas_affiliation IS 'Indicates the authorized affiliation for PPTAS API';