-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO public.tsp_users
	VALUES (
		uuid_generate_v4(), NULL,
		'Donkey', 'Kong', NULL,
		'd.kong@example.com', '(555) 342-4654',
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',
		now(), now()
	);

INSERT INTO public.tsp_users
	VALUES (
		uuid_generate_v4(), NULL,
		'Michael', 'Jackson', NULL,
		'michaelj@example.com', '(555) 432-8098',
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',
		now(), now()
	);

INSERT INTO public.tsp_users
	VALUES (
		uuid_generate_v4(), NULL,
		'Joey', 'Schabadoo', 'Joe-Joe',
		'j.jj.shabadoo@example.com', '(555) 348-9080',
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',
		now(), now()
	);
