--update duty loc name for NAS Whidbey Island
update duty_locations set name = 'NAS Whidbey Island, WA 98278', updated_at = now() where id = 'dac0aebc-87a4-475d-92f1-cddcfef7c607';

--remove duty loc Norcross, GA 30010
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '309ecf13-df38-4ef2-bc6c-506db0f46aba') THEN

		update orders set origin_duty_location_id = '0a6afb38-9466-4ccb-90ce-e31277fcc8e9', updated_at = now() where origin_duty_location_id = '309ecf13-df38-4ef2-bc6c-506db0f46aba';
		update orders set new_duty_location_id = '0a6afb38-9466-4ccb-90ce-e31277fcc8e9', updated_at = now() where new_duty_location_id = '309ecf13-df38-4ef2-bc6c-506db0f46aba';

		delete from duty_locations where id = '309ecf13-df38-4ef2-bc6c-506db0f46aba';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'd1a70c89-e8e4-49bb-9bcf-d19c4eeeedd4';

	END IF;

END $$;

--remove duty loc Virginia Beach, VA 23450
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '330b1d21-2b61-4d40-945e-855d3b37c922') THEN

		update orders set origin_duty_location_id = 'dd21a586-48e7-4b80-9dbf-1bd498ac50b0', updated_at = now() where origin_duty_location_id = '330b1d21-2b61-4d40-945e-855d3b37c922';
		update orders set new_duty_location_id = 'dd21a586-48e7-4b80-9dbf-1bd498ac50b0', updated_at = now() where new_duty_location_id = '330b1d21-2b61-4d40-945e-855d3b37c922';

		delete from duty_locations where id = '330b1d21-2b61-4d40-945e-855d3b37c922';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'e0f3d67a-1dd8-4bf2-8c84-c1630730b28d';

	END IF;

END $$;

--remove duty loc Seattle, WA 98175
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'd33381dc-f2fb-4ade-a8e3-6396d47d99c5') THEN

		update orders set origin_duty_location_id = '7f877bc1-f4cb-4261-8fdf-bf7bece554d7', updated_at = now() where origin_duty_location_id = 'd33381dc-f2fb-4ade-a8e3-6396d47d99c5';
		update orders set new_duty_location_id = '7f877bc1-f4cb-4261-8fdf-bf7bece554d7', updated_at = now() where new_duty_location_id = 'd33381dc-f2fb-4ade-a8e3-6396d47d99c5';

		delete from duty_locations where id = 'd33381dc-f2fb-4ade-a8e3-6396d47d99c5';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = '9ecd16a7-cb3d-4921-8db3-69d3829178c7';

	END IF;

END $$;

--remove duty loc Aurora, CO 80040
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'b037856b-c68e-4db8-8887-afe2ea278c6e') THEN

		update orders set origin_duty_location_id = 'fcefc02d-47ec-4e26-8a40-e6bc5759e22f', updated_at = now() where origin_duty_location_id = 'b037856b-c68e-4db8-8887-afe2ea278c6e';
		update orders set new_duty_location_id = 'fcefc02d-47ec-4e26-8a40-e6bc5759e22f', updated_at = now() where new_duty_location_id = 'b037856b-c68e-4db8-8887-afe2ea278c6e';

		delete from duty_locations where id = 'b037856b-c68e-4db8-8887-afe2ea278c6e';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'e602521f-46a0-42fb-8f27-811b77aaecb5';

	END IF;

END $$;

--remove duty loc Raleigh, NC 27602
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'b2eb208f-8acf-4683-88c8-d02013c272eb') THEN

		update orders set origin_duty_location_id = '369845b4-8996-4472-8acf-87c017aa11eb', updated_at = now() where origin_duty_location_id = 'b2eb208f-8acf-4683-88c8-d02013c272eb';
		update orders set new_duty_location_id = '369845b4-8996-4472-8acf-87c017aa11eb', updated_at = now() where new_duty_location_id = 'b2eb208f-8acf-4683-88c8-d02013c272eb';

		delete from duty_locations where id = 'b2eb208f-8acf-4683-88c8-d02013c272eb';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'bdcde56b-fb40-4726-a102-29107ce3d43c';

	END IF;

END $$;

--remove duty loc Washington, DC 20030
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '7f9be29b-14ac-46ff-ac4f-19c69fc9c264') THEN

		update orders set origin_duty_location_id = '3497f324-bc5d-46a4-ba6f-3fd861a1e776', updated_at = now() where origin_duty_location_id = '7f9be29b-14ac-46ff-ac4f-19c69fc9c264';
		update orders set new_duty_location_id = '3497f324-bc5d-46a4-ba6f-3fd861a1e776', updated_at = now() where new_duty_location_id = '7f9be29b-14ac-46ff-ac4f-19c69fc9c264';

		delete from duty_locations where id = '7f9be29b-14ac-46ff-ac4f-19c69fc9c264';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'fe416657-2515-4b54-8826-0a05cfd2400b';

	END IF;

END $$;

--remove duty loc Spanish Fort, AL 36577
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'b16720cf-39d3-4d2d-a77b-c1420fcec2b8') THEN

		update orders set origin_duty_location_id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc', updated_at = now() where origin_duty_location_id = 'b16720cf-39d3-4d2d-a77b-c1420fcec2b8';
		update orders set new_duty_location_id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc', updated_at = now() where new_duty_location_id = 'b16720cf-39d3-4d2d-a77b-c1420fcec2b8';

		delete from duty_locations where id = 'b16720cf-39d3-4d2d-a77b-c1420fcec2b8';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'f70cf6c4-59f2-45a2-a36a-261f30d9ccfe';

	END IF;

END $$;

--remove duty loc Austin, TX 78760
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '063f0e90-da3e-46cb-9a75-59b172a8337f') THEN

		update orders set origin_duty_location_id = '53210016-136d-4314-9494-bfe49ea428fa', updated_at = now() where origin_duty_location_id = '063f0e90-da3e-46cb-9a75-59b172a8337f';
		update orders set new_duty_location_id = '53210016-136d-4314-9494-bfe49ea428fa', updated_at = now() where new_duty_location_id = '063f0e90-da3e-46cb-9a75-59b172a8337f';

		delete from duty_locations where id = '063f0e90-da3e-46cb-9a75-59b172a8337f';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'f17462d0-4375-48d4-b068-300c6f8d508b';

	END IF;

END $$;

--remove duty loc Columbus, GA 31908
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = '8fc9c0c5-0dca-4217-ac4d-cc09bdeb4a97') THEN

		update orders set origin_duty_location_id = '7de8225d-6248-49c3-9708-2eb6cba9ac32', updated_at = now() where origin_duty_location_id = '8fc9c0c5-0dca-4217-ac4d-cc09bdeb4a97';
		update orders set new_duty_location_id = '7de8225d-6248-49c3-9708-2eb6cba9ac32', updated_at = now() where new_duty_location_id = '8fc9c0c5-0dca-4217-ac4d-cc09bdeb4a97';

		delete from duty_locations where id = '8fc9c0c5-0dca-4217-ac4d-cc09bdeb4a97';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'abe0e510-1abe-4135-86f1-fec6727ed220';

	END IF;

END $$;

--remove duty loc Las Vegas, NV 89136
DO $$
BEGIN

	IF EXISTS (SELECT 1 FROM duty_locations WHERE id = 'b865ba00-6949-4185-bf47-2587eb2666c6') THEN

		update orders set origin_duty_location_id = '35aaa898-6844-44f2-ae69-02b5b069c138', updated_at = now() where origin_duty_location_id = 'b865ba00-6949-4185-bf47-2587eb2666c6';
		update orders set new_duty_location_id = '35aaa898-6844-44f2-ae69-02b5b069c138', updated_at = now() where new_duty_location_id = 'b865ba00-6949-4185-bf47-2587eb2666c6';

		delete from duty_locations where id = 'b865ba00-6949-4185-bf47-2587eb2666c6';

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = 'aa990f46-6488-4764-a603-6908a464d76e';

	END IF;

END $$;