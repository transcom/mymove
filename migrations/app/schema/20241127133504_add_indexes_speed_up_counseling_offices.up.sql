--adding indexes to help speed up counseling_offices dropdown
CREATE INDEX IF NOT EXISTS idx_addresses_postal_code ON addresses(postal_code);
CREATE INDEX IF NOT EXISTS idx_addresses_us_post_region_cities_id ON addresses(us_post_region_cities_id);
CREATE INDEX IF NOT EXISTS idx_addresses_country_id ON addresses(country_id);
CREATE INDEX IF NOT EXISTS idx_duty_locations_provides_services_counseling ON duty_locations(provides_services_counseling);
CREATE INDEX IF NOT EXISTS idx_re_us_post_regions_uspr_zip_id ON re_us_post_regions(uspr_zip_id);
CREATE INDEX IF NOT EXISTS idx_zip3_distances_from_zip3_to_zip3 ON zip3_distances(from_zip3, to_zip3);
CREATE INDEX IF NOT EXISTS idx_zip3_distances_to_zip3_from_zip3 ON zip3_distances(to_zip3, from_zip3);
CREATE INDEX IF NOT EXISTS idx_re_oconus_rate_areas_us_post_region_cities_id ON re_oconus_rate_areas(us_post_region_cities_id);
CREATE INDEX IF NOT EXISTS idx_gbloc_aors_oconus_rate_area_id ON gbloc_aors(oconus_rate_area_id);
CREATE INDEX IF NOT EXISTS idx_gbloc_aors_jppso_regions_id ON gbloc_aors(jppso_regions_id);
CREATE INDEX IF NOT EXISTS idx_transportation_offices_provides_ppm_closeout ON transportation_offices(provides_ppm_closeout);