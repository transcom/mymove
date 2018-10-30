-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Remove users
DELETE FROM public.tsp_users WHERE email = 'joe.schmoe@example.com' AND
transportation_service_provider_id = 'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28';

-- Add new users
INSERT INTO public.tsp_users VALUES (uuid_generate_v4(), NULL, 'Moe', 'Howard', NULL, 'm.howard@example.com', '(555) 123-4566', 'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',now(), now());
