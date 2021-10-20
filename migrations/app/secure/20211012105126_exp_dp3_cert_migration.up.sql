-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prd/stg/exp/demo
-- DO NOT include any sensitive data.

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
        '102628a2-b4a2-49a8-a850-cd0e4998f846',
        '205bf4f39b6100fad8733b1375edfcf6c5a4ee07a3fda7ca658ead9eaf53d8b3',
        'CN=api.exp.dp3.us',
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