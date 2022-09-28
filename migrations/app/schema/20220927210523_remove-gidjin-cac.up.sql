-- Remove gidjin CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='31d9903c22119796ebb7ea04321ae35ebd6265015aae260ff7662e66de0c6e1d';
