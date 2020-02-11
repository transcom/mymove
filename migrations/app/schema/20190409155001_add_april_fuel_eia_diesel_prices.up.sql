CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO public.fuel_eia_diesel_prices (id, pub_date, rate_start_date, rate_end_date,
										   eia_price_per_gallon_millicents, baseline_rate, created_at, updated_at)
VALUES (uuid_generate_v4(), '2019-04-01', '2019-04-15', '2019-05-14', 307800, 5, now(), now());
