--
-- Data for Name: transportation_offices; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.transportation_offices VALUES ('ccf50409-9d03-4cac-a931-580649f1647a', 'aa899628-dabb-4724-8e4a-b4579c1550e0', 'Camp LeJeune (USMC)', '0f631e20-4edf-457d-8de8-431b4e17f5ae', 34.668,-77.327866, 'Monday; Tuesday; Wednesday and Friday 0730 -1600; Thursday 0730 - 1500', 'Walk-In Help; Briefings; Computer Lab; Appointments; QA Inspections', NULL, now(), now(), 'USMC');
--
-- Data for Name: office_phone_lines; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.office_phone_lines VALUES ('c20bbee8-72a8-4510-9dcb-20757e05ad36', 'ccf50409-9d03-4cac-a931-580649f1647a', '(910) 451-2377', NULL, false, 'voice', now(), now());
INSERT INTO public.office_phone_lines VALUES ('3f677f42-5191-4548-b50c-b94a3a48d99c', 'ccf50409-9d03-4cac-a931-580649f1647a', '750-4377', NULL, true, 'voice', now(), now());
--
-- Data for Name: office_emails; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.office_emails VALUES('33bec643-f01b-4ba4-be7a-e02a26f6240d', 'ccf50409-9d03-4cac-a931-580649f1647a', 'ppcig@usmc.mil', 'Customer Service',  now(), now())

