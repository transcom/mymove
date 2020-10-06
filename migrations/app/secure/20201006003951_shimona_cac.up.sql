
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '46d6dc9b5624db0b0c319f13871d94f60e7fc43f212e76ae499ada702506539b',
	subject = 'CN=shimonacarvalho,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	updated_at = NOW()
WHERE id = '000c67fc-56dd-4732-9c48-b2cfb6e052ee';
