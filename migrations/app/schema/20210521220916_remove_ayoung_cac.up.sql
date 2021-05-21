-- Remove abbyoung CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='68e7ce8ec4681f4eb1cbebf2f2e3e25eb1edb4675689087b66efabd0b04b330a';
