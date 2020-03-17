
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
INSERT INTO public.client_certs (
    id,
    sha256_digest,
    subject,
    allow_dps_auth_api,
    allow_orders_api,
    created_at,
    updated_at,
    allow_air_force_orders_read,
    allow_air_force_orders_write,
    allow_army_orders_read,
    allow_army_orders_write,
    allow_coast_guard_orders_read,
    allow_coast_guard_orders_write,
    allow_marine_corps_orders_read,
    allow_marine_corps_orders_write,
    allow_navy_orders_read,
    allow_navy_orders_write,
    allow_prime)
VALUES (
           '078a8da8-3d9f-4c96-82da-b36e71317c08',
           'fcee9e07caf20dbcc2c795652c232d8b554e01e999941f23167d31e80a3ca330',
           'CN=kimallen,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
           false,
           true,
           now(),
           now(),
           true,
           true,
           true,
           true,
           true,
           true,
           true,
           true,
           true,
           true,
           true);


