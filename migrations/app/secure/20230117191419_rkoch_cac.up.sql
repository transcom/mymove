-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '8e2a89d0e578ba80cde13b726be2ac0b12f3ca4519a27d2f793c3f1f6bfabc03',
	subject = 'CN=ryan-koch,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = '1e7998d0-3145-4252-b293-7a6a3d52cb32';
