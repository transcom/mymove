COPY public.re_contracts (id, code, name, created_at, updated_at) FROM stdin;
8ef44ca4-c589-4c39-93d8-3410a762ec6c	Pricing	Pricing	2020-06-15 19:07:59.702658	2020-06-15 19:07:59.702661
\.

COPY public.re_contract_years (id, contract_id, name, start_date, end_date, escalation, escalation_compounded, created_at, updated_at) FROM stdin;
9ca0c8d2-3b14-4a49-9709-1f710383642f	8ef44ca4-c589-4c39-93d8-3410a762ec6c	Base Period Year 3	2021-06-01	2022-05-31	1.01970	1.04071	2020-06-15 19:07:59.713313	2020-06-15 19:07:59.713314
\.

COPY public.re_domestic_service_areas (id, service_area, services_schedule, sit_pd_schedule, created_at, updated_at, contract_id) FROM stdin;
ea949687-431e-4c00-bbe7-646158471d4e	004	2	2	2020-06-15 19:08:03.986359	2020-06-15 19:08:03.98636	8ef44ca4-c589-4c39-93d8-3410a762ec6c
\.

COPY public.re_zip3s (id, zip3, domestic_service_area_id, created_at, updated_at, contract_id, rate_area_id, has_multiple_rate_areas, base_point_city, state) FROM stdin;
7eee5a4f-c457-49eb-bd6a-216fb00ab43c	321	ea949687-431e-4c00-bbe7-646158471d4e	2020-06-15 19:08:03.267183	2020-06-15 19:08:06.52456	8ef44ca4-c589-4c39-93d8-3410a762ec6c	\N	t	Crescent City	FL
\.
