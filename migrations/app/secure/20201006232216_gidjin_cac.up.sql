
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
-- LOCAL MIGRATION SECURE DATA REMOVED
UPDATE client_certs
SET sha256_digest = 'fda2d1f3b265c280e20cea2e693fe86199f7777af802957cb9b1bbe3e6868338',
	subject = 'CN=gidjin,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = '4fea0eb1-0009-47a8-98f4-0a102ee53a4f';
