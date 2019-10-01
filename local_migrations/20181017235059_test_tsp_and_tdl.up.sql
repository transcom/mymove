-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ----------------
-- Create a test TSP

INSERT INTO public.transportation_service_providers
	VALUES (
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',
		'TRS1',
	   	now(), now(),
		true
	);

-- ----------------
-- Create users for the TSP

INSERT INTO public.tsp_users
	VALUES (
		uuid_generate_v4(), NULL,
		'Joe', 'Schmoe', 'O',
		'joe.schmoe@example.com', '(555) 123-4567',
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',
		now(), now()
	);

-- ----------------
-- Create a test TDL

-- Create a Service Area for the test ZIP: 00000

INSERT INTO public.tariff400ng_zip3s
	VALUES (
		uuid_generate_v4(),
		'000',					-- ZIP3
		'Test Town', 'CA',		-- basepoint_city, state
		'000',					-- service_area
	   	'US00',					-- rate_area
		'0',					-- region
		now(), now()
	);

-- Create a TDL

INSERT INTO public.traffic_distribution_lists
	VALUES (
		'4a8e3a96-265b-4150-900e-f86eedb38b47',	-- UUID
		'US00',					-- source_rate_area
		'11',					-- destination_region (for ZIP 401xx, including Ft Knox)
		'D',					-- code_of_service
		now(), now()
	);

-- File rates in the TDL

INSERT INTO public.transportation_service_provider_performances
	VALUES (
		uuid_generate_v4(),
		'2018-01-01',							-- performance_period_start
		'2030-01-01',							-- performance_period_end (super far in the future)
		'4a8e3a96-265b-4150-900e-f86eedb38b47',	-- traffic_distribution_list_id
		NULL,									-- quality_band
		0,										-- offer_count
		100.0,									-- best_value_score
		'd7c0e4e0-ddcf-47b8-bdfd-6c0bce555b28',	-- transportation_service_provider_id
		now(), now(),							-- created_at, updated_at
		'2018-01-01',							-- rate_cycle_start
		'2030-01-01',							-- rate_cycle_end
		0.6,									-- linehaul_rate
		0.6										-- sit_rate
	);

