--
-- PostgreSQL database dump
--

-- Dumped from database version 16.0
-- Dumped by pg_dump version 16.0

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: appointments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.appointments (
    id integer NOT NULL,
    timetable_id bigint NOT NULL,
    pacient_username character varying(255) NOT NULL,
    "time" timestamp with time zone NOT NULL
);


ALTER TABLE public.appointments OWNER TO postgres;

--
-- Name: appointments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.appointments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.appointments_id_seq OWNER TO postgres;

--
-- Name: appointments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.appointments_id_seq OWNED BY public.appointments.id;


--
-- Name: history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.history (
    id integer NOT NULL,
    date timestamp with time zone NOT NULL,
    pacient_uuid character varying(36) NOT NULL,
    hospital_uuid character varying(36) NOT NULL,
    doctor_uuid character varying(36) NOT NULL,
    room character varying(255) NOT NULL,
    data character varying(255) NOT NULL
);


ALTER TABLE public.history OWNER TO postgres;

--
-- Name: history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.history_id_seq OWNER TO postgres;

--
-- Name: history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.history_id_seq OWNED BY public.history.id;


--
-- Name: hospital; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hospital (
    uuid character varying(36) NOT NULL,
    name character varying(255) NOT NULL,
    address character varying(255) NOT NULL,
    contact_phone character varying(255) NOT NULL
);


ALTER TABLE public.hospital OWNER TO postgres;

--
-- Name: hospital_rooms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.hospital_rooms (
    id integer NOT NULL,
    hospital_uuid character varying(36) NOT NULL,
    room character varying(255) NOT NULL
);


ALTER TABLE public.hospital_rooms OWNER TO postgres;

--
-- Name: hospital_rooms_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.hospital_rooms_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.hospital_rooms_id_seq OWNER TO postgres;

--
-- Name: hospital_rooms_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.hospital_rooms_id_seq OWNED BY public.hospital_rooms.id;


--
-- Name: timetable; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.timetable (
    id integer NOT NULL,
    hospital_uuid character varying(36) NOT NULL,
    doctor_uuid character varying(36) NOT NULL,
    "from" timestamp with time zone NOT NULL,
    "to" timestamp with time zone NOT NULL,
    room character varying(255) NOT NULL
);


ALTER TABLE public.timetable OWNER TO postgres;

--
-- Name: timetable_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.timetable_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.timetable_id_seq OWNER TO postgres;

--
-- Name: timetable_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.timetable_id_seq OWNED BY public.timetable.id;


--
-- Name: user_and_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_and_roles (
    id integer NOT NULL,
    user_uuid character varying(36) NOT NULL,
    role character varying(255) DEFAULT 'user'::character varying NOT NULL
);


ALTER TABLE public.user_and_roles OWNER TO postgres;

--
-- Name: user_and_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_and_roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_and_roles_id_seq OWNER TO postgres;

--
-- Name: user_and_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_and_roles_id_seq OWNED BY public.user_and_roles.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    uuid character varying(36) NOT NULL,
    username character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    password character varying(255) NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: appointments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments ALTER COLUMN id SET DEFAULT nextval('public.appointments_id_seq'::regclass);


--
-- Name: history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.history ALTER COLUMN id SET DEFAULT nextval('public.history_id_seq'::regclass);


--
-- Name: hospital_rooms id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hospital_rooms ALTER COLUMN id SET DEFAULT nextval('public.hospital_rooms_id_seq'::regclass);


--
-- Name: timetable id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timetable ALTER COLUMN id SET DEFAULT nextval('public.timetable_id_seq'::regclass);


--
-- Name: user_and_roles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_roles ALTER COLUMN id SET DEFAULT nextval('public.user_and_roles_id_seq'::regclass);


--
-- Data for Name: appointments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.appointments (id, timetable_id, pacient_username, "time") FROM stdin;
\.


--
-- Data for Name: history; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.history (id, date, pacient_uuid, hospital_uuid, doctor_uuid, room, data) FROM stdin;
\.


--
-- Data for Name: hospital; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.hospital (uuid, name, address, contact_phone) FROM stdin;
\.


--
-- Data for Name: hospital_rooms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.hospital_rooms (id, hospital_uuid, room) FROM stdin;
\.


--
-- Data for Name: timetable; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.timetable (id, hospital_uuid, doctor_uuid, "from", "to", room) FROM stdin;
\.


--
-- Data for Name: user_and_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_and_roles (id, user_uuid, role) FROM stdin;
4	a9f0a3b0-f680-4e21-9a1e-3f48d6fb52cd	user
1	f8534269-2b59-46e4-b962-82a217f47c66	admin
2	39098b85-00ac-42af-b28c-8a0a7cde44ef	manager
3	b07a0c0f-bf38-4217-ba5a-a6727ddd7437	doctor
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (uuid, username, first_name, last_name, password) FROM stdin;
f8534269-2b59-46e4-b962-82a217f47c66	admin	Admin	Adminov	admin
39098b85-00ac-42af-b28c-8a0a7cde44ef	manager	Manager	Managerov	manager
b07a0c0f-bf38-4217-ba5a-a6727ddd7437	doctor	Doctor	Doctorov	doctor
a9f0a3b0-f680-4e21-9a1e-3f48d6fb52cd	user	User	Userov	user
\.


--
-- Name: appointments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.appointments_id_seq', 1, false);


--
-- Name: history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.history_id_seq', 1, false);


--
-- Name: hospital_rooms_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.hospital_rooms_id_seq', 1, false);


--
-- Name: timetable_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.timetable_id_seq', 1, false);


--
-- Name: user_and_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_and_roles_id_seq', 4, true);


--
-- Name: appointments appointments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_pkey PRIMARY KEY (id);


--
-- Name: history history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_pkey PRIMARY KEY (id);


--
-- Name: hospital hospital_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hospital
    ADD CONSTRAINT hospital_name_key UNIQUE (name);


--
-- Name: hospital hospital_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hospital
    ADD CONSTRAINT hospital_pkey PRIMARY KEY (uuid);


--
-- Name: hospital_rooms hospital_rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hospital_rooms
    ADD CONSTRAINT hospital_rooms_pkey PRIMARY KEY (id);


--
-- Name: timetable timetable_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timetable
    ADD CONSTRAINT timetable_pkey PRIMARY KEY (id);


--
-- Name: user_and_roles user_and_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_roles
    ADD CONSTRAINT user_and_roles_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (uuid);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: appointments appointments_pacient_username_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_pacient_username_fkey FOREIGN KEY (pacient_username) REFERENCES public.users(username) ON DELETE CASCADE;


--
-- Name: appointments appointments_timetable_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.appointments
    ADD CONSTRAINT appointments_timetable_id_fkey FOREIGN KEY (timetable_id) REFERENCES public.timetable(id) ON DELETE CASCADE;


--
-- Name: history history_pacient_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.history
    ADD CONSTRAINT history_pacient_uuid_fkey FOREIGN KEY (pacient_uuid) REFERENCES public.users(uuid) ON DELETE CASCADE;


--
-- Name: hospital_rooms hospital_rooms_hospital_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.hospital_rooms
    ADD CONSTRAINT hospital_rooms_hospital_uuid_fkey FOREIGN KEY (hospital_uuid) REFERENCES public.hospital(uuid) ON DELETE CASCADE;


--
-- Name: timetable timetable_doctor_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timetable
    ADD CONSTRAINT timetable_doctor_uuid_fkey FOREIGN KEY (doctor_uuid) REFERENCES public.users(uuid) ON DELETE CASCADE;


--
-- Name: timetable timetable_hospital_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.timetable
    ADD CONSTRAINT timetable_hospital_uuid_fkey FOREIGN KEY (hospital_uuid) REFERENCES public.hospital(uuid) ON DELETE CASCADE;


--
-- Name: user_and_roles user_and_roles_user_uuid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_and_roles
    ADD CONSTRAINT user_and_roles_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES public.users(uuid) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

