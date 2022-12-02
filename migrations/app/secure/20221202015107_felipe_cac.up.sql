-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prd/stg/exp/demo
-- DO NOT include any sensitive data.

-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '074382e25c702f11aa88d92f29a8b9ee3187d654f47258bd5a4a5a864d3d4b74',
	subject = 'CN=felipe-lee,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = '5eda0836-2910-4280-8f37-f129e3859f2a';
