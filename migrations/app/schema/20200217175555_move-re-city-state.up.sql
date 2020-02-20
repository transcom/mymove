-- Add nullable (for now) city/state columns to the re_zip3s table
ALTER TABLE re_zip3s
    ADD COLUMN base_point_city varchar(80),
    ADD COLUMN state varchar(80);

-- Populate city/state with data from the re_domestic_service_areas table
UPDATE re_zip3s z3
SET base_point_city = dsa.base_point_city,
    state           = dsa.state
FROM re_domestic_service_areas dsa
WHERE z3.domestic_service_area_id = dsa.id;

-- Make the columns not nullable
ALTER TABLE re_zip3s
    ALTER COLUMN base_point_city SET NOT NULL,
    ALTER COLUMN state SET NOT NULL;

-- Drop the (now) old city/state columns from the re_domestic_service_areas table
ALTER TABLE re_domestic_service_areas
    DROP COLUMN base_point_city,
    DROP COLUMN state;
