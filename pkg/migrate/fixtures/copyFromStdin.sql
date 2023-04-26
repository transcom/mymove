--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5 (Debian 10.5-2.pgdg90+1)
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

-- comment
-- blah
-- blah2
-- blah blah blah blah
COPY public.re_services (id, code, name, created_at, updated_at, priority) FROM stdin;
10000012-2c32-4529-ad8a-131df722cb17	12	Twelve	2020-03-23 16:31:50.313853	2020-03-23 16:31:50.313853	1
10000013-ef6e-45b1-9d3d-8a89e46af743	13	Thirteen	2020-03-23 16:31:50.313853	2020-03-23 16:31:50.313853	2
\.


--
-- PostgreSQL database dump complete
--
