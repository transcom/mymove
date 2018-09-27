-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TNG Users
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Alyssa', 'Ogawa', NULL, 'alyssa@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Thadiun', 'Okona', NULL, 'thadiun@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Marie', 'Picard', NULL, 'marie@example.com', '(800) 555-1234', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());

-- DS9 users
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Miles', 'O''Brien', NULL, 'obrien@example.com', '(800) 555-5000', 'f50eb7f5-960a-46e8-aa64-6025b44132ab', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Julian', 'Bashir', NULL, 'bashir@example.com', '(800) 555-5000', 'f50eb7f5-960a-46e8-aa64-6025b44132ab', now(), now());
