--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.9
-- Dumped by pg_dump version 9.6.9

-- Started on 2018-10-29 15:15:48 -04

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 7 (class 2615 OID 48323)
-- Name: callcenter; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA callcenter;


ALTER SCHEMA callcenter OWNER TO postgres;

--
-- TOC entry 515 (class 1255 OID 49183)
-- Name: consultdatacallcenter(text, text); Type: FUNCTION; Schema: callcenter; Owner: postgres
--

CREATE FUNCTION callcenter.consultdatacallcenter(date_from text, date_to text) RETURNS TABLE(result json)
    LANGUAGE plpgsql
    AS $_$
BEGIN
	IF TRIM($1) != '' AND TRIM($2) != '' THEN
		RETURN QUERY SELECT array_to_json(array_agg(row_to_json(t)))
		FROM
		(SELECT date_ride, driver, origin, destination, rider, tks FROM callcenter.callcenterdata WHERE date_ride BETWEEN to_date($1,'YYYY-MM-DD') AND to_date($2,'YYYY-MM-DD')) t;
	END IF;
	IF TRIM($1) = '' AND TRIM($2) = '' THEN
		RETURN QUERY SELECT array_to_json(array_agg(row_to_json(t)))
		FROM
		(SELECT date_ride, driver, origin, destination, rider, tks FROM callcenter.callcenterdata) t;
	END IF;
	IF TRIM($1) != '' AND TRIM($2) = '' THEN
		RETURN QUERY SELECT array_to_json(array_agg(row_to_json(t)))
		FROM
		(SELECT date_ride, driver, origin, destination, rider, tks FROM callcenter.callcenterdata WHERE date_ride >= to_date($1,'YYYY-MM-DD')) t;
	END IF;
	IF TRIM($1) = '' AND TRIM($2) != '' THEN
		RETURN QUERY SELECT array_to_json(array_agg(row_to_json(t)))
		FROM
		(SELECT date_ride, driver, origin, destination, rider, tks FROM callcenter.callcenterdata WHERE date_ride <= to_date($2,'YYYY-MM-DD')) t;
	END IF;			
END;

$_$;


ALTER FUNCTION callcenter.consultdatacallcenter(date_from text, date_to text) OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 411 (class 1259 OID 48876)
-- Name: callcenterdata; Type: TABLE; Schema: callcenter; Owner: postgres
--

CREATE TABLE callcenter.callcenterdata (
    id integer NOT NULL,
    date_ride date,
    driver character varying(100),
    origin character varying(100),
    destination character varying(100),
    rider character varying(100),
    tks character varying(100)
);


ALTER TABLE callcenter.callcenterdata OWNER TO postgres;

--
-- TOC entry 410 (class 1259 OID 48874)
-- Name: callcenterdata_id_seq; Type: SEQUENCE; Schema: callcenter; Owner: postgres
--

CREATE SEQUENCE callcenter.callcenterdata_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE callcenter.callcenterdata_id_seq OWNER TO postgres;

--
-- TOC entry 2913 (class 0 OID 0)
-- Dependencies: 410
-- Name: callcenterdata_id_seq; Type: SEQUENCE OWNED BY; Schema: callcenter; Owner: postgres
--

ALTER SEQUENCE callcenter.callcenterdata_id_seq OWNED BY callcenter.callcenterdata.id;


--
-- TOC entry 412 (class 1259 OID 49109)
-- Name: userdata; Type: TABLE; Schema: callcenter; Owner: postgres
--

CREATE TABLE callcenter.userdata (
    driver text,
    uuid text
);


ALTER TABLE callcenter.userdata OWNER TO postgres;

--
-- TOC entry 2788 (class 2604 OID 48879)
-- Name: callcenterdata id; Type: DEFAULT; Schema: callcenter; Owner: postgres
--

ALTER TABLE ONLY callcenter.callcenterdata ALTER COLUMN id SET DEFAULT nextval('callcenter.callcenterdata_id_seq'::regclass);


--
-- TOC entry 2790 (class 2606 OID 48884)
-- Name: callcenterdata callcenterdata_id; Type: CONSTRAINT; Schema: callcenter; Owner: postgres
--

ALTER TABLE ONLY callcenter.callcenterdata
    ADD CONSTRAINT callcenterdata_id PRIMARY KEY (id);


-- Completed on 2018-10-29 15:15:49 -04

--
-- PostgreSQL database dump complete
--

