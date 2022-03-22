-- Remove jacquelineIO's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest = 'f1c5a949a12dbd8ffa37faf0e4f4a6930a97534993d290b102117b6094e69d3b';

-- Remove jaynawallaceTRUSS's CAC using unique sha256 digest of the client cert
DELETE FROM client_certs WHERE sha256_digest = '58c0fc8d1bee1cbad735294079d785c798c603978238630b4d4b5321b469f3db';
