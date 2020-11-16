
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = 'f1c5a949a12dbd8ffa37faf0e4f4a6930a97534993d290b102117b6094e69d3b',
	subject = 'CN=jacquelineIO,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = 'cfaea7ca-33de-40c8-b4cb-3e6f90f315e8';
