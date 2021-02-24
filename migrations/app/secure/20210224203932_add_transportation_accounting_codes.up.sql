-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.

-- Based this on 20201208215206_load_zip3_distances.up.sql

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Temporarily drop indexes to speed up load

DROP INDEX public.transportation_accounting_codes_tac_idx;

--
-- Data for Name: transportation_accounting_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

-- NOTE these are fake TAC values that were created based on most of the rules for TAC values.

COPY public.transportation_accounting_codes (id, tac) FROM stdin;
e170982f-e356-4699-9e78-0658a7c957eb	E01A
0c0cfb1f-ffbb-4004-b089-04da8f2a16b9	E11A
4e845ac6-2fe3-4904-b499-ee38b15566e4	E12A
ec3fc4ce-ad2d-45d2-aaa8-401de3eb82a3	E13A
1ddf6bed-7691-49e0-b66a-747e17dc71b9	E14A
55849130-451e-479b-abc3-990d592ce34d	E15A
c7a6b4bb-6714-4f4b-88a0-56c6bcd1e9c5	E16A
b49b139e-662f-4628-b9cc-de5be12a3a68	E17A
89518b5d-5976-4145-a75d-6994f2121098	E18A
f9d3b41e-5f50-47e1-b93a-d55830148ee1	E19A
5c055341-2995-40b5-a54b-830044208903	1111
ce6a23ad-ae74-4ba1-b2cc-645d3da6160a	2222
5aaac38b-7d1b-4974-ad1a-3bfd737c0cab	3333
09c69c2e-403c-497c-9e5b-e88fede5a700	4444
4929db2d-957f-4dd1-b2d2-43e059571a6a	5555
b4dccee3-3065-4f02-b05e-9ea5920f378a	6666
8c8de24f-9896-45bc-8faa-b4582457c515	7777
11360093-bef0-4746-95e9-c232a7b36458	8888
e3566ec6-04df-4d02-add6-5db39b5ec4e6	9999
\.

CREATE INDEX transportation_accounting_codes_tac_idx ON public.transportation_accounting_codes USING btree(tac);
