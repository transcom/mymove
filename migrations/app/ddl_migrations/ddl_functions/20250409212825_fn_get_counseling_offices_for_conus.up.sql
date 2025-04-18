--B-22660 Daniel Jordan added fetch_counseling_offices_for_conus
--B-22914 Beth Grohmann moved to separate script and modify query for performance
CREATE OR REPLACE FUNCTION fetch_counseling_offices_for_conus(p_duty_location_id UUID)
RETURNS TABLE (id UUID, name TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    with counseling_offices as (
                SELECT transportation_offices.id, transportation_offices.name, transportation_offices.address_id as counseling_address, substring(addresses.postal_code, 1,3 ) as pickup_zip
                        FROM postal_code_to_gblocs
                        JOIN addresses on postal_code_to_gblocs.postal_code = addresses.postal_code
                        JOIN duty_locations on addresses.id = duty_locations.address_id
                        JOIN transportation_offices on postal_code_to_gblocs.gbloc = transportation_offices.gbloc
                        WHERE duty_locations.provides_services_counseling = true and duty_locations.id = p_duty_location_id
                )
        SELECT counseling_offices.id, counseling_offices.name
                FROM counseling_offices
                JOIN duty_locations duty_locations2 on counseling_offices.id = duty_locations2.transportation_office_id
                JOIN addresses on counseling_offices.counseling_address = addresses.id
                JOIN re_us_post_regions on addresses.postal_code = re_us_post_regions.uspr_zip_id
                LEFT JOIN zip3_distances ON (
		                (re_us_post_regions.zip3 = zip3_distances.to_zip3
		            AND counseling_offices.pickup_zip = zip3_distances.from_zip3)
		                OR
		                (re_us_post_regions.zip3 = zip3_distances.from_zip3
		            AND counseling_offices.pickup_zip = zip3_distances.to_zip3)
		        )
                WHERE duty_locations2.provides_services_counseling = true
        group by counseling_offices.id, counseling_offices.name, zip3_distances.distance_miles
                ORDER BY coalesce(zip3_distances.distance_miles,0), counseling_offices.name asc;
END;
$$;