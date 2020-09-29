-- Remove rdhariwal CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='f2658a432de1f84813f4372d93873a328925df5a4fdc4ee948a025cb9520a83c';
