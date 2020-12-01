-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.

-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '2475a3de7da3082c5b3ff364cd8c567906206cba4c50541b5f9f0b6972284fb3',
    subject = 'CN=ledeep,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
    updated_at = NOW()
WHERE id = '9928caf4-072c-438a-8a5e-07b213bc1826';