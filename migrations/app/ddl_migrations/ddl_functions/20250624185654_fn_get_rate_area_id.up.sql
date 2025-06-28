-- B-23853 Beth Grohmann Initial check-in, added additional logic to pull from re_zip5_rate_areas

drop function if exists get_rate_area_id;

CREATE OR REPLACE FUNCTION get_rate_area_id(
    address_id UUID,
    service_item_id UUID,
    c_id uuid,
    OUT o_rate_area_id UUID
)
RETURNS UUID AS $$
DECLARE
    is_oconus BOOLEAN;
    zip3_value TEXT;
	zip5_value TEXT;
BEGIN
    is_oconus := get_is_oconus(address_id);

    IF is_oconus THEN
        -- re_oconus_rate_areas if is_oconus is TRUE
        SELECT ro.rate_area_id
        INTO o_rate_area_id
        FROM addresses a
        JOIN re_oconus_rate_areas ro
        ON a.us_post_region_cities_id = ro.us_post_region_cities_id
        JOIN re_rate_areas ra ON ro.rate_area_id = ra.id
        WHERE a.id = address_id
            AND ra.contract_id = c_id;
    ELSE
        -- re_zip3s if is_oconus is FALSE
        SELECT rupr.zip3, rupr.uspr_zip_id
        INTO zip3_value, zip5_value
        FROM addresses a
        JOIN us_post_region_cities uprc
        ON a.us_post_region_cities_id = uprc.id
        JOIN re_us_post_regions rupr
        ON uprc.us_post_regions_id = rupr.id
        WHERE a.id = address_id;

        -- use the zip3 value to find the rate_area_id in re_zip3s
       SELECT rz.rate_area_id
        INTO o_rate_area_id
        FROM re_zip3s rz
        JOIN re_rate_areas ra
        ON rz.rate_area_id = ra.id
        WHERE rz.zip3 = zip3_value
            AND ra.contract_id = c_id;

        IF o_rate_area_id IS NULL THEN

			--if o_rate_area_id is null check zip5
			SELECT rz.rate_area_id
	          INTO o_rate_area_id
			  FROM re_zip5_rate_areas rz
			  JOIN re_rate_areas ra
			    ON rz.rate_area_id = ra.id
	         WHERE rz.zip5 = zip5_value
	           AND ra.contract_id = c_id;

	    END IF;
    END IF;

    -- Raise an exception if no rate area is found
    IF o_rate_area_id IS NULL THEN
        RAISE EXCEPTION 'Rate area not found for address % for service item ID %', address_id, service_item_id;
    END IF;
END;
$$ LANGUAGE plpgsql;