--B-22660 Daniel Jordan added fetch_counseling_offices_for_oconus
--B-22914 Beth Grohmann moved to separate script and modify query for performance
CREATE OR REPLACE FUNCTION fetch_counseling_offices_for_oconus(p_duty_location_id UUID, p_gbloc_indicator TEXT)
RETURNS TABLE (id UUID, name TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
	with counseling_offices as (
			SELECT transportation_offices.id, transportation_offices.name, transportation_offices.address_id as counseling_address,
			  substring(a.postal_code, 1,3 ) as origin_zip, substring(a2.postal_code, 1,3 ) as dest_zip, j.code as gbloc
			FROM duty_locations
			JOIN addresses a on duty_locations.address_id = a.id
			JOIN v_locations v on (a.us_post_region_cities_id = v.uprc_id or v.uprc_id is null)
					and a.country_id = v.country_id
			JOIN re_oconus_rate_areas r on r.us_post_region_cities_id = v.uprc_id
			JOIN gbloc_aors on gbloc_aors.oconus_rate_area_id = r.id
			JOIN jppso_regions j on gbloc_aors.jppso_regions_id = j.id
			JOIN transportation_offices on j.code = transportation_offices.gbloc
			join addresses a2 on a2.id = transportation_offices.address_id
			WHERE duty_locations.provides_services_counseling = true and duty_locations.id = p_duty_location_id and j.code = p_gbloc_indicator
			    and transportation_offices.provides_ppm_closeout = true
			)
		SELECT counseling_offices.id, counseling_offices.name
			FROM counseling_offices
			JOIN addresses cnsl_address on counseling_offices.counseling_address = cnsl_address.id
			LEFT JOIN zip3_distances ON (
				(substring(cnsl_address.postal_code,1 ,3) = zip3_distances.to_zip3
				AND counseling_offices.origin_zip = zip3_distances.from_zip3)
				OR
				(substring(cnsl_address.postal_code,1 ,3) = zip3_distances.from_zip3
				AND counseling_offices.origin_zip = zip3_distances.to_zip3)
			)
			group by counseling_offices.id, counseling_offices.name, zip3_distances.distance_miles, counseling_offices.gbloc
			ORDER BY coalesce(zip3_distances.distance_miles,0) asc;
END;
$$;