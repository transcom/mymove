--B-22666   Tevin Adams    Adding foreign key fk_us_post_region_cities to the addresses table and populating us_post_region_cities_id where null

ALTER TABLE addresses
ADD CONSTRAINT fk_us_post_region_cities
FOREIGN KEY (us_post_region_cities_id) REFERENCES us_post_region_cities(id);
