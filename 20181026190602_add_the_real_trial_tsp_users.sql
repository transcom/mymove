CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--------
-- Users for Planetary Van Lines

-- Add new users
INSERT INTO public.tsp_users VALUES (uuid_generate_v4(), NULL, 'Moe', 'Howard', NULL, 'm.howard@example.com', '(555) 123-4566', 'c71bdb14-ed86-4c92-bf06-93c0865f5070',now(), now());

--------
-- Users for Green Chip LLC dba D'Lux Moving and Storage

-- Add new users
INSERT INTO public.tsp_users VALUES (uuid_generate_v4(), NULL, 'Larry', 'Fine', NULL, 'larry.fine@example.com', '(555) 983-1235', 'b98d3deb-abe9-4609-8d6e-36b2c50873c0',now(), now());

--------
-- Users for Security Storage Company of Washington LLC

-- Add new users
INSERT INTO public.tsp_users VALUES (uuid_generate_v4(), NULL, 'Curly', 'Howard', NULL, 'c.howard@example.com', '(555) 954-2342', 'b6f06674-1b6b-4b93-9ec6-293d5d846876',now(), now());
