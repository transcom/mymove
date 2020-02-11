-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ABC users
UPDATE public.office_users SET first_name = 'test44', email = 'test44@example.com' WHERE email = 'test43@example.com';

-- DEF users
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Gibbons', 'Peter', NULL, 'gibbons@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Lumbergh', 'Bill', NULL, 'lumbergh@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Waddams', 'Milton', NULL, 'waddams@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Smykowski', 'Tom', NULL, 'smykowski@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
