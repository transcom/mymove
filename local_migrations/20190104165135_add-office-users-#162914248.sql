-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Person1', 'Ima', NULL, 'ot1@test.test', '415-555-5555', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());
INSERT INTO public.office_users VALUES (uuid_generate_v4(), NULL, 'Person2', 'Ima', NULL, 'ot2@test.test', '415-555-5555', '0931a9dc-c1fd-444a-b138-6e1986b1714c', now(), now());