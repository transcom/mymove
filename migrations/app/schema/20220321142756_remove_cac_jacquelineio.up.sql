-- Remove jacquelineIO's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest = 'f1c5a949a12dbd8ffa37faf0e4f4a6930a97534993d290b102117b6094e69d3b';
