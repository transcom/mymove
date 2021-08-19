-- Remove erin CAC using unique sha256 digest of client cert
DELETE FROM client_certs WHERE sha256_digest='d72d8acfa6c5e8e6aa036f69a817a000db8387dfa5ed4397d1aa636b82949330';
