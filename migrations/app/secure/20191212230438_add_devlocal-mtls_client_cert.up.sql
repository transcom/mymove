-- This migration allows a snake oil cert to have read/write access to all orders and the prime API.
-- IT SHOULD ONLY BE RUN LOCALLY
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
INSERT INTO public.client_certs
    (
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
VALUES
    (
        '190b1e07-eef8-445a-9696-5a2b49ee488d',
        '2c0c1fc67a294443292a9e71de0c71cc374fe310e8073f8cdc15510f6b0ef4db',
        '/C=US/ST=DC/L=Washington/O=Truss/OU=AppClientTLS/CN=devlocal',
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