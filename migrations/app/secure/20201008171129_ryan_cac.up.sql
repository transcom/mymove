
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = 'b8f9f403dc38cffe9866cb19b8383dbcaaecd41e92bf4c21ee7f8db94faf4d11',
    subject = 'CN=ryan-koch,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
    updated_at = NOW()
WHERE id = '1e7998d0-3145-4252-b293-7a6a3d52cb32';
