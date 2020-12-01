-- Remove cgilmer CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='15697aff85d93a2f04259618b66f09f9976887b3f2ee206395c42e4897aefe61';
