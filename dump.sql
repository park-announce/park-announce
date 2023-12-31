--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3 (Debian 15.3-1.pgdg120+1)
-- Dumped by pg_dump version 15.3 (Debian 15.3-1.pgdg120+1)

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

--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: pa_corporation_locations; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporation_locations (
    id character varying(50) NOT NULL,
    corporation_id character varying(50) NOT NULL,
    geog public.geography NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    available_location_count integer DEFAULT 0 NOT NULL,
    name character varying(100) NOT NULL
);


ALTER TABLE public.pa_corporation_locations OWNER TO park_announce;

--
-- Name: pa_corporation_prices; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporation_prices (
    id character varying(50) NOT NULL,
    corporation_location_id character varying(50) NOT NULL,
    price double precision NOT NULL,
    currency character varying(10) NOT NULL,
    price_description character varying(50) NOT NULL
);


ALTER TABLE public.pa_corporation_prices OWNER TO park_announce;

--
-- Name: pa_corporation_roles; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporation_roles (
    id character varying(50) NOT NULL,
    name character varying(200) NOT NULL
);


ALTER TABLE public.pa_corporation_roles OWNER TO park_announce;

--
-- Name: pa_corporation_users; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporation_users (
    id character varying(50) NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
    email character varying(200) NOT NULL,
    password character varying(200) NOT NULL,
    corporation_id character varying(50) NOT NULL,
    role_id character varying(50) NOT NULL
);


ALTER TABLE public.pa_corporation_users OWNER TO park_announce;

--
-- Name: pa_corporations; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporations (
    id character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    status smallint DEFAULT 1 NOT NULL
);


ALTER TABLE public.pa_corporations OWNER TO park_announce;

--
-- Name: pa_locations; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_locations (
    id character varying(50) NOT NULL,
    geog public.geography,
    status smallint DEFAULT 0 NOT NULL,
    owner_id character varying(50) NOT NULL,
    assigned_user_id character varying(50),
    scheduled_available_time bigint,
    location_type smallint DEFAULT 0 NOT NULL
);


ALTER TABLE public.pa_locations OWNER TO park_announce;

--
-- Name: pa_user_passwords; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_user_passwords (
    id character varying(50) NOT NULL,
    user_id character varying(50) NOT NULL,
    password character varying(100) NOT NULL
);


ALTER TABLE public.pa_user_passwords OWNER TO park_announce;

--
-- Name: pa_users; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_users (
    id character varying(50) NOT NULL,
    email character varying(200) NOT NULL,
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    mobile_phone character varying(20) NOT NULL,
    city_code smallint NOT NULL
);


ALTER TABLE public.pa_users OWNER TO park_announce;

--
-- Data for Name: pa_corporation_locations; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_corporation_locations (id, corporation_id, geog, status, available_location_count, name) FROM stdin;
229c5221-0dd6-40a4-b0de-cd63424afeXX	932b7062-6435-4be1-83d1-b37f9d3f0448	0101000020E6100000426D937699043D4094A31074D98C4440	0	12	ispark - uskudar
\.


--
-- Data for Name: pa_corporation_prices; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_corporation_prices (id, corporation_location_id, price, currency, price_description) FROM stdin;
229c5221-0dd6-40a4-b0de-cd63424afeYY	229c5221-0dd6-40a4-b0de-cd63424afeXX	20	TL	1 saat
\.


--
-- Data for Name: pa_corporation_roles; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_corporation_roles (id, name) FROM stdin;
0d811a85-c53b-4d5e-9c61-97d9607259e1	admin
2b2a8cc3-471f-44cd-ac46-fd595c091c60	user
\.


--
-- Data for Name: pa_corporation_users; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_corporation_users (id, status, email, password, corporation_id, role_id) FROM stdin;
932b7062-6435-4be1-83d1-b37f9d3f0333	1	resulguldibi@gmail.com	$2a$10$E8rQ34gcb/80PT.c1o.WXu22AEcf7BPTqHBsWWfmi1dw9bjNleSAu	932b7062-6435-4be1-83d1-b37f9d3f0448	0d811a85-c53b-4d5e-9c61-97d9607259e1
09bcbb19-1915-446f-858e-2e357cd31e58	1	fatihberksoz@gmail.com	$2a$10$zmpuE3xTwNoLwvPBmf2yQuVPt8GOrekioqaPolEFE5aHEccRPayG6	932b7062-6435-4be1-83d1-b37f9d3f0448	2b2a8cc3-471f-44cd-ac46-fd595c091c60
\.


--
-- Data for Name: pa_corporations; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_corporations (id, name, status) FROM stdin;
932b7062-6435-4be1-83d1-b37f9d3f0448	ispark	1
\.


--
-- Data for Name: pa_locations; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_locations (id, geog, status, owner_id, assigned_user_id, scheduled_available_time, location_type) FROM stdin;
d2c7a5d9-8ae2-4cc7-a14a-4075b308f505	0101000020E61000004C58BBAC7E033D40A06EE384E48C4440	0	932b7062-6435-4be1-83d1-b37f9d3f0448	\N	\N	0
9bb18f75-70a8-4d97-90d9-3b1b69c35561	0101000020E61000005C051BAE24033D40A059FEC1068D4440	0	932b7062-6435-4be1-83d1-b37f9d3f0448	\N	\N	0
\.


--
-- Data for Name: pa_user_passwords; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_user_passwords (id, user_id, password) FROM stdin;
7779a639-61c9-4645-b5ec-e55005d2a4f9	3472ec31-f624-46c5-a4ce-02604680da52	$2a$10$q9.0Fazp0qNFMbNMQr4ZZOyTrqYG63WZPiX0Hh3cZpp20ea6HwTO6
\.


--
-- Data for Name: pa_users; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.pa_users (id, email, first_name, last_name, status, mobile_phone, city_code) FROM stdin;
f6d89519-fd2d-4b29-be94-8559b725c84b	2@gmail.com	kamil	tayyare	1	905542891616	34
3472ec31-f624-46c5-a4ce-02604680da52	2@gmail.com	kamil	tayyare	1	905542891616	34
\.


--
-- Data for Name: spatial_ref_sys; Type: TABLE DATA; Schema: public; Owner: park_announce
--

COPY public.spatial_ref_sys (srid, auth_name, auth_srid, srtext, proj4text) FROM stdin;
\.


--
-- Name: pa_corporation_locations pa_corporation_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_locations
    ADD CONSTRAINT pa_corporation_locations_pkey PRIMARY KEY (id);


--
-- Name: pa_corporation_prices pa_corporation_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_prices
    ADD CONSTRAINT pa_corporation_prices_pkey PRIMARY KEY (id);


--
-- Name: pa_corporation_roles pa_corporation_roles_name_key; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_roles
    ADD CONSTRAINT pa_corporation_roles_name_key UNIQUE (name);


--
-- Name: pa_corporation_roles pa_corporation_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_roles
    ADD CONSTRAINT pa_corporation_roles_pkey PRIMARY KEY (id);


--
-- Name: pa_corporation_users pa_corporation_users_email_key; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_users
    ADD CONSTRAINT pa_corporation_users_email_key UNIQUE (email);


--
-- Name: pa_corporation_users pa_corporation_users_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_users
    ADD CONSTRAINT pa_corporation_users_pkey PRIMARY KEY (id);


--
-- Name: pa_corporations pa_corporations_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporations
    ADD CONSTRAINT pa_corporations_pkey PRIMARY KEY (id);


--
-- Name: pa_locations pa_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_locations
    ADD CONSTRAINT pa_locations_pkey PRIMARY KEY (id);


--
-- Name: pa_user_passwords pa_user_passwords_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_user_passwords
    ADD CONSTRAINT pa_user_passwords_pkey PRIMARY KEY (id);


--
-- Name: pa_users pa_users_pkey; Type: CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_users
    ADD CONSTRAINT pa_users_pkey PRIMARY KEY (id);


--
-- Name: pa_locations_geog_idx; Type: INDEX; Schema: public; Owner: park_announce
--

CREATE INDEX pa_locations_geog_idx ON public.pa_locations USING gist (geog);


--
-- Name: pa_corporation_locations pa_corporation_locations_corporation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_locations
    ADD CONSTRAINT pa_corporation_locations_corporation_id_fkey FOREIGN KEY (corporation_id) REFERENCES public.pa_corporations(id) NOT VALID;


--
-- Name: pa_corporation_prices pa_corporation_prices_corporation_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_prices
    ADD CONSTRAINT pa_corporation_prices_corporation_location_id_fkey FOREIGN KEY (corporation_location_id) REFERENCES public.pa_corporation_locations(id) NOT VALID;


--
-- Name: pa_corporation_users pa_corporation_users_corporation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_users
    ADD CONSTRAINT pa_corporation_users_corporation_id_fkey FOREIGN KEY (corporation_id) REFERENCES public.pa_corporations(id) NOT VALID;


--
-- Name: pa_corporation_users pa_corporation_users_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_corporation_users
    ADD CONSTRAINT pa_corporation_users_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.pa_corporation_roles(id) NOT VALID;


--
-- Name: pa_user_passwords pa_user_passwords_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: park_announce
--

ALTER TABLE ONLY public.pa_user_passwords
    ADD CONSTRAINT pa_user_passwords_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.pa_users(id);


--
-- PostgreSQL database dump complete
--

