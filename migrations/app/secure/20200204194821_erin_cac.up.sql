-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
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
        'a8744ea1-38ee-458b-ba28-e8d04b4f76c3',
        'd72d8acfa6c5e8e6aa036f69a817a000db8387dfa5ed4397d1aa636b82949330',
        'CN=tinyels,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
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