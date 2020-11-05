-- Remove mr337 CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='6f684e8ed9f697a5b6099f31f5346042db6339ecd8029000df8924369defb069';
