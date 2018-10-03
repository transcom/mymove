-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TOS Users
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Leonard', 'McCoy', NULL, 'bones@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Janice', 'Rand', NULL, 'randj@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
