-- Remove sandy's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='0beb55298a9bf146e53213266fc45fcc2b83c22ab2721aea9c2bce0a9ccc4acd';
-- Remove lynzt's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest='38789424b4a63082a8f911ba2015b7c414ad50e77de94ea0b251ad7d70ff7ed2';
