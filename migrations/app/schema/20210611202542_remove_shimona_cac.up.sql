-- Remove shimona CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='46d6dc9b5624db0b0c319f13871d94f60e7fc43f212e76ae499ada702506539b';
