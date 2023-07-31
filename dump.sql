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
    available_location_count integer DEFAULT 0 NOT NULL
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
-- Name: pa_corporation_users; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_corporation_users (
    id character varying(50) NOT NULL,
    status smallint DEFAULT 1,
    email character varying(200) NOT NULL,
    password character varying(10) NOT NULL,
    corporation_id character varying(50) NOT NULL
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
-- Name: pa_users; Type: TABLE; Schema: public; Owner: park_announce
--

CREATE TABLE public.pa_users (
    id character varying(50) NOT NULL,
    email character varying(200) NOT NULL,
    password character varying(10)
);


ALTER TABLE public.pa_users OWNER TO park_announce;

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
-- PostgreSQL database dump complete
--

