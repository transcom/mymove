--B-22666   Tevin Adams    Adding foreign key fk_us_post_region_cities to the addresses table

DO $$
BEGIN
    IF NOT EXISTS (SELECT constraint_name
        FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
        Where table_name = 'addresses' and constraint_name = 'fk_us_post_region_cities') then
        ALTER TABLE addresses
            ADD CONSTRAINT fk_us_post_region_cities
            FOREIGN KEY (us_post_region_cities_id) REFERENCES us_post_region_cities(id);
    END IF;
END $$;
