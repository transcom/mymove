-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.

-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '0eade0a51eb1ba4f17ce37cd3a6ee2a705eaa71dc0fb74ba1edd1218454930c0',
	subject = 'CN=reggieriser,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = '03191873-8cc5-4af1-9866-b4bf12842d54';
