DO $$
BEGIN
	
	--remove duty loc Johnston City, TN 37602
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687') THEN
	
		
		update orders set origin_duty_location_id = 'cd0c7325-15bb-45c7-a690-26c56c903ed7' where origin_duty_location_id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';
		update orders set new_duty_location_id = 'cd0c7325-15bb-45c7-a690-26c56c903ed7' where new_duty_location_id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';
		
		delete from duty_locations where id = 'd3a1be10-dcd7-4720-bcbe-7ba76d243687';
	
	END IF;

END $$;

DO $$
BEGIN
	
	--remove duty loc Oceanside, CA 92052
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de') THEN
	
		update orders set origin_duty_location_id = 'a6993e7b-4600-44b9-b288-04ca011143f0' where origin_duty_location_id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';
		update orders set new_duty_location_id = 'a6993e7b-4600-44b9-b288-04ca011143f0' where new_duty_location_id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';
		
		delete from duty_locations where id = '54ca99b7-3c2a-42b0-aa1a-ad071ac580de';
	
	END IF;

END $$;

DO $$
BEGIN

	--remove duty loc Albuquerque, NM 87103
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '2cc57072-19fa-438b-a44b-e349dff11763') THEN
	
		update orders set new_duty_location_id = '54acfb0e-222b-49eb-b94b-ccb00c6f529c' where new_duty_location_id = '2cc57072-19fa-438b-a44b-e349dff11763';
	
		delete from duty_locations where id = '2cc57072-19fa-438b-a44b-e349dff11763';

	END IF;

END $$;

DO $$
BEGIN
	
	--remove duty loc August, GA 30917
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '109ac405-47fb-4e1e-9efb-58290453ac09') THEN
	
		update orders set origin_duty_location_id = '595363c2-14ee-48e0-b318-e76ab0016453' where origin_duty_location_id = '109ac405-47fb-4e1e-9efb-58290453ac09';
		update orders set new_duty_location_id = '595363c2-14ee-48e0-b318-e76ab0016453' where new_duty_location_id = '109ac405-47fb-4e1e-9efb-58290453ac09';
		
		delete from duty_locations where id = '109ac405-47fb-4e1e-9efb-58290453ac09';
	
	END IF;

END $$;

DO $$
BEGIN
	
	--remove duty loc Frankfort, KY 40602
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'c7fadaa2-902f-4302-a7cd-108c525b96d4') THEN
	
		update orders set origin_duty_location_id = '1a973257-cd15-42a9-86be-a14796c014bc' where origin_duty_location_id = 'c7fadaa2-902f-4302-a7cd-108c525b96d4';
		update orders set new_duty_location_id = '1a973257-cd15-42a9-86be-a14796c014bc' where new_duty_location_id = 'c7fadaa2-902f-4302-a7cd-108c525b96d4';
		
		delete from duty_locations where id = 'c7fadaa2-902f-4302-a7cd-108c525b96d4';
	
	END IF;

END $$;

DO $$
BEGIN
	
	--remove duty loc Seattle, WA 98111
	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '2fb3e898-d6de-4be7-8576-7c7b10c2a706') THEN
	
		update orders set origin_duty_location_id = 'e7fdae4f-6be7-4264-99f8-03ee8541499c' where origin_duty_location_id = '2fb3e898-d6de-4be7-8576-7c7b10c2a706';
		update orders set new_duty_location_id = 'e7fdae4f-6be7-4264-99f8-03ee8541499c' where new_duty_location_id = '2fb3e898-d6de-4be7-8576-7c7b10c2a706';
		
		delete from duty_locations where id = '2fb3e898-d6de-4be7-8576-7c7b10c2a706';
	
	END IF;

END $$;

--add Joint Base Lewis McChord, WA 98438 duty location
DO $$
BEGIN

	IF NOT EXISTS (SELECT 1 FROM addresses WHERE id = '23d3140b-1ba2-400f-9d57-317034673c06') THEN

		INSERT INTO public.addresses
			(id, street_address_1, city, state, postal_code, created_at, updated_at, county, is_oconus, country_id, us_post_region_cities_id)
		VALUES('23d3140b-1ba2-400f-9d57-317034673c06'::uuid, 'n/a', 'JOINT BASE LEWIS MCCHORD', 'WA', '98438', now(),now(), 'PIERCE', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '81182dd4-1693-4b8d-9b6f-042bc4254019'::uuid);

	END IF;
	
	IF NOT EXISTS (SELECT 1 FROM duty_locations WHERE id = '109ac405-47fb-4e1e-9efb-58290453ac09') THEN

		INSERT INTO public.duty_locations
		(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
		VALUES('38fc6718-b80f-4761-a077-cfa62e414e27', 'Joint Base Lewis McChord, WA 98438', 'AIR_FORCE', '23d3140b-1ba2-400f-9d57-317034673c06'::uuid, now(), now(), '95abaeaa-452f-4fe0-9264-960cd2a15ccd', true);
		
	END IF;

	IF NOT EXISTS (SELECT 1 FROM duty_locations WHERE id = '693781f4-d011-4925-a492-aa1185f3f1fe') THEN

		INSERT INTO public.duty_locations
		(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
		VALUES('693781f4-d011-4925-a492-aa1185f3f1fe'::uuid, 'McChord AFB, WA 98438', 'AIR_FORCE', 'cc7894e3-148e-4e21-98df-37e45f0b2c9f'::uuid, now(), now(), '95abaeaa-452f-4fe0-9264-960cd2a15ccd', true);
		
	END IF;
	
END $$;

--associate duty loc Yuma, AZ 85365 to transportation office PPPO DMO MCAS Yuma - USMC
update duty_locations set transportation_office_id = '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4' where id = '9e94208a-881d-47bc-82c0-4f375471751e';

--update name for Alameda
update transportation_offices set name = 'PPSO Base Alameda - USCG' where id = '3fc4b408-1197-430a-a96a-24a5a1685b45';