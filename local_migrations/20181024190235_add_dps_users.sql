-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'aileen@dds.mil', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
INSERT INTO public.dps_users VALUES (uuid_generate_v4(), 'test@example.com', now(), now());
