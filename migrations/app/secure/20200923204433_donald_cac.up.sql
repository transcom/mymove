-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.

-- Updating sha256 for new cac
UPDATE client_certs
    SET sha256_digest = 'fb364ff8130fd7c456a0c53809b78f80844cd75682253c9e53be7bcd7b5faaa9',
        subject = 'CN=donaldthai,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
        updated_at = NOW()
    WHERE id = 'de7605ec-2edd-4252-b176-27f0d3fe4b6f';