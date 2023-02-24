-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prd/stg/exp/demo
-- DO NOT include any sensitive data.
UPDATE client_certs
SET sha256_digest = 'e8426c603b1bb240ceb4cc629df81819def2f2aba55ad3f9fecd9c615fe5d6ae',
	updated_at = NOW()
WHERE id = 'b8e767d6-fb38-4602-b7f1-6bead16ca7e1';
