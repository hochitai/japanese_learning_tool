--
-- PostgreSQL database dump
--

-- Dumped from database version 16.1
-- Dumped by pg_dump version 16.0

-- Started on 2024-06-15 21:52:39

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
-- TOC entry 4863 (class 1262 OID 81934)
-- Name: learning; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE learning WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'English_United States.1252';


ALTER DATABASE learning OWNER TO postgres;

\connect learning

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
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 4864 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 220 (class 1259 OID 98329)
-- Name: favorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.favorites (
    id integer NOT NULL,
    word_id integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.favorites OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 98328)
-- Name: favorites_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.favorites_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.favorites_id_seq OWNER TO postgres;

--
-- TOC entry 4865 (class 0 OID 0)
-- Dependencies: 219
-- Name: favorites_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.favorites_id_seq OWNED BY public.favorites.id;


--
-- TOC entry 218 (class 1259 OID 81945)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    name text NOT NULL,
    salt text NOT NULL,
    permission text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 81944)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 4866 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 216 (class 1259 OID 81936)
-- Name: words; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.words (
    id integer NOT NULL,
    characters text NOT NULL,
    pronunciation text NOT NULL,
    meaning text,
    level text NOT NULL
);


ALTER TABLE public.words OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 81935)
-- Name: words_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.words_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.words_id_seq OWNER TO postgres;

--
-- TOC entry 4867 (class 0 OID 0)
-- Dependencies: 215
-- Name: words_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.words_id_seq OWNED BY public.words.id;


--
-- TOC entry 4700 (class 2604 OID 98332)
-- Name: favorites id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites ALTER COLUMN id SET DEFAULT nextval('public.favorites_id_seq'::regclass);


--
-- TOC entry 4699 (class 2604 OID 81948)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 4698 (class 2604 OID 81939)
-- Name: words id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.words ALTER COLUMN id SET DEFAULT nextval('public.words_id_seq'::regclass);


--
-- TOC entry 4857 (class 0 OID 98329)
-- Dependencies: 220
-- Data for Name: favorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.favorites (id, word_id, user_id) VALUES (3, 10, 8);
INSERT INTO public.favorites (id, word_id, user_id) VALUES (4, 11, 6);


--
-- TOC entry 4855 (class 0 OID 81945)
-- Dependencies: 218
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, username, password, name, salt, permission) VALUES (5, 'johnwhite', 'de30fa3dc34b8e6de3010a6995cc078d', 'Jack', 'FeHjxDnPPphoNSuc', 'user');
INSERT INTO public.users (id, username, password, name, salt, permission) VALUES (6, 'jack', '5eaaa58d6276fcc001b4d6fdd11bdcd5', 'jack', 'ZPVwQy5AUgbwqcEO', 'user');
INSERT INTO public.users (id, username, password, name, salt, permission) VALUES (8, 'john', 'ebfc8482e083b6b753fb208f4e902c01', 'john', 'aene8PJarHxYp3mO', 'admin');


--
-- TOC entry 4853 (class 0 OID 81936)
-- Dependencies: 216
-- Data for Name: words; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (1, 'ありがとう', 'arigatoo', 'Cảm ơn', 'public');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (2, 'どこ', 'doko', 'ở đâu', 'public');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (3, 'そう', 'soo', 'ra thế', 'public');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (4, 'はい', 'hai', 'vâng', 'public');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (5, 'ある', 'aru', 'có', 'public');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (10, 'じゃね', 'ja ne', 'Goodbye', 'private');
INSERT INTO public.words (id, characters, pronunciation, meaning, level) VALUES (11, 'すごい', 'sugoi', 'tuyệt vời', 'private');


--
-- TOC entry 4868 (class 0 OID 0)
-- Dependencies: 219
-- Name: favorites_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.favorites_id_seq', 4, true);


--
-- TOC entry 4869 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 8, true);


--
-- TOC entry 4870 (class 0 OID 0)
-- Dependencies: 215
-- Name: words_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.words_id_seq', 11, true);


--
-- TOC entry 4704 (class 2606 OID 90137)
-- Name: users constraint_username; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT constraint_username UNIQUE (username);


--
-- TOC entry 4708 (class 2606 OID 98334)
-- Name: favorites favorites_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT favorites_pkey PRIMARY KEY (id);


--
-- TOC entry 4706 (class 2606 OID 81952)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4702 (class 2606 OID 81943)
-- Name: words words_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.words
    ADD CONSTRAINT words_pkey PRIMARY KEY (id);


-- Completed on 2024-06-15 21:52:39

--
-- PostgreSQL database dump complete
--

