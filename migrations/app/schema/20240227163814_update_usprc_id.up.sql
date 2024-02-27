CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE us_post_region_cities
ALTER COLUMN id SET DEFAULT uuid_generate_v4();
