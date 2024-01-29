-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
INSERT INTO users (
    id,
    okta_email,
    created_at,
    updated_at)
VALUES (
    '87fc5974-fbc9-4719-a3e2-b609647478d7',
    '25b64f60444878e22c3cbfbbfdeb6e3e38832ade1c9704a7bd906b709c15bf38' || '@api.move.mil',
    now(),
    now());

INSERT INTO users_roles (
    id,
    role_id,
    user_id,
    created_at,
    updated_at)
VALUES (
    uuid_generate_v4(),
    (SELECT id FROM roles WHERE role_type = 'prime'),
    '87fc5974-fbc9-4719-a3e2-b609647478d7',
    now(),
    now());

INSERT INTO public.client_certs (
    id,
    sha256_digest,
    subject,
    user_id,
    allow_orders_api,
    allow_prime,
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
    allow_navy_orders_write)
VALUES (
    '3a80db0d-a204-49f9-a9b2-359f57378e01',
    '25b64f60444878e22c3cbfbbfdeb6e3e38832ade1c9704a7bd906b709c15bf38',
    'C=US, O=U.S. Government, OU=ECA, OU=IdenTrust, OU=MOVEHQ INC., CN=mmb.gov.uat.homesafeconnect.com',
    '87fc5974-fbc9-4719-a3e2-b609647478d7',
    true,
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
    true);
