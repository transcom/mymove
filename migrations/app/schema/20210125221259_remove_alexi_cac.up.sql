-- Remove akostibas CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='a5922b65a90304385b83f8a01983e8ad312c4e08c2c9322ebcf54cfd35e61fda';
