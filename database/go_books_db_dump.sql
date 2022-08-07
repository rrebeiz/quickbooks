--
-- PostgreSQL database dump
--

-- Dumped from database version 14.0 (Debian 14.0-1.pgdg110+1)
-- Dumped by pg_dump version 14.3

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
-- Name: authors; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.authors (
    id bigint NOT NULL,
    author_name character varying(512) NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.authors OWNER TO devuser;

--
-- Name: authors_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.authors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authors_id_seq OWNER TO devuser;

--
-- Name: authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.authors_id_seq OWNED BY public.authors.id;


--
-- Name: books; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.books (
    id bigint NOT NULL,
    title character varying(512) NOT NULL,
    author_id integer NOT NULL,
    publication_year integer NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    slug character varying(512) NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.books OWNER TO devuser;

--
-- Name: books_genres; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.books_genres (
    id bigint NOT NULL,
    book_id integer NOT NULL,
    genre_id integer NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.books_genres OWNER TO devuser;

--
-- Name: books_genres_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.books_genres_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_genres_id_seq OWNER TO devuser;

--
-- Name: books_genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.books_genres_id_seq OWNED BY public.books_genres.id;


--
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.books_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO devuser;

--
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- Name: genres; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.genres (
    id bigint NOT NULL,
    genre_name character varying(255) NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.genres OWNER TO devuser;

--
-- Name: genres_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.genres_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.genres_id_seq OWNER TO devuser;

--
-- Name: genres_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.genres_id_seq OWNED BY public.genres.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.tokens (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    email text NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    expiry timestamp(0) with time zone NOT NULL
);


ALTER TABLE public.tokens OWNER TO devuser;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tokens_id_seq OWNER TO devuser;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp(0) with time zone DEFAULT now() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    password_hash bytea NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    account_type character varying(255) DEFAULT 'user'::character varying NOT NULL,
    updated_at timestamp(0) with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: authors id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.authors ALTER COLUMN id SET DEFAULT nextval('public.authors_id_seq'::regclass);


--
-- Name: books id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- Name: books_genres id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books_genres ALTER COLUMN id SET DEFAULT nextval('public.books_genres_id_seq'::regclass);


--
-- Name: genres id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.genres ALTER COLUMN id SET DEFAULT nextval('public.genres_id_seq'::regclass);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: authors; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.authors (id, author_name, created_at, updated_at, version) FROM stdin;
1	James Clear	2022-08-01 07:03:15+00	2022-08-01 07:03:15+00	1
2	Frank Herbert	2022-08-01 07:18:53+00	2022-08-01 07:18:53+00	1
\.


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.books (id, title, author_id, publication_year, created_at, updated_at, slug, description) FROM stdin;
1	Atomic Habits	1	1990	2022-08-01 07:11:35+00	2022-08-01 07:11:35+00	atomic-habits	An atomic habit is a regular practice or routine that is not only small and easy to do but is also the source of incredible power; a component of the system of compound growth. Bad habits repeat themselves again and again not because you do not want to change, but because you have the wrong system for change.
2	Dune	2	1965	2022-08-01 07:32:52+00	2022-08-01 07:32:52+00	dune	set in the distant future amidst a feudal interstellar society in which various noble houses control planetary fiefs. It tells the story of young Paul Atreides, whose family accepts the stewardship of the planet Arrakis
3	Dune Messiah	2	1969	2022-08-01 07:34:15+00	2022-08-01 07:34:15+00	dune-messiah	Dune Messiah continuesthe story of Paul Atreides, better known—and feared—as the man christened Muad Dib. As Emperor of the known universe, he possesses more power than a single man was ever meant to wield
4	Children of Dune	2	1976	2022-08-01 07:35:24+00	2022-08-01 07:35:24+00	children-of-dune	After Paul walks into the desert, Leto II and Ghanima, his children, become the center of politics; this leads to a series of events that shapes the Dune universe.
5	God Emperor of Dune	2	1981	2022-08-01 07:36:42+00	2022-08-01 07:36:42+00	god-emperor-of-dune	Leto II Atreides, the God Emperor, has ruled the universe as a tyrant for 3,500 years after becoming a hybrid of human and giant sandworm in Children of Dune.
6	Heretics of Dune	2	1984	2022-08-01 07:37:36+00	2022-08-01 07:37:36+00	heretics-of-dune	Heretics of Dune tells of life 3,500 years after the death of the Tyrant, Leto II, as ferocious "Honored Matres" stream in from the "Scattering," and the Bene Gesserit work to unite a ghola, Duncan Idaho, and a Fremen girl, Sheeana, who has the power to command worms.
7	Chapterhouse: Dune	2	1985	2022-08-01 07:39:00+00	2022-08-01 07:39:00+00	chapterhouse:-dune	A direct follow-up to Heretics of Dune, the novel chronicles the continued struggles of the Bene Gesserit Sisterhood against the violent Honored Matres, who are succeeding in their bid to seize control of the universe and destroy the factions and planets that oppose them.
8	Test Book	2	2022	2022-08-01 14:50:31+00	2022-08-02 05:03:00+00	test-book	A test book
\.


--
-- Data for Name: books_genres; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.books_genres (id, book_id, genre_id, created_at, updated_at) FROM stdin;
2	1	9	2022-08-02 04:17:02+00	2022-08-02 04:17:02+00
16	8	4	2022-08-02 05:03:00+00	2022-08-02 05:03:00+00
17	8	5	2022-08-02 05:03:00+00	2022-08-02 05:03:00+00
18	2	1	2022-08-05 14:33:45+00	2022-08-05 14:33:45+00
19	3	1	2022-08-05 14:33:48+00	2022-08-05 14:33:48+00
20	4	1	2022-08-05 14:33:50+00	2022-08-05 14:33:50+00
21	5	1	2022-08-05 14:33:56+00	2022-08-05 14:33:56+00
22	6	1	2022-08-05 14:33:58+00	2022-08-05 14:33:58+00
23	7	1	2022-08-05 14:34:01+00	2022-08-05 14:34:01+00
24	8	1	2022-08-05 14:34:03+00	2022-08-05 14:34:03+00
\.


--
-- Data for Name: genres; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.genres (id, genre_name, created_at, updated_at) FROM stdin;
1	Science Fiction	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
2	Fantasy	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
3	Romance	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
4	Thriller	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
5	Mystery	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
6	Horror	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
7	Classic	2022-08-01 07:06:13+00	2022-08-01 07:06:13+00
9	Self-help	2022-08-02 04:16:38+00	2022-08-02 04:16:38+00
\.


--
-- Data for Name: tokens; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.tokens (id, user_id, email, token, token_hash, created_at, updated_at, expiry) FROM stdin;
79	1	admin@test.com	46RSVYVIJLNKMDLXU6QAZ5PZJA	\\x5e6e85b86f3967b5a28abcdcaffbdbcf8a7d4fdc8e6fd7ecd4e86bf2f6ddbfdf	2022-08-07 06:41:02+00	2022-08-07 06:41:02+00	2022-08-08 06:41:02+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.users (id, created_at, name, email, password_hash, version, account_type, updated_at) FROM stdin;
1	2022-07-25 14:53:01+00	Admin	admin@test.com	\\x24326124313224776b557478435775386d323934746759314f7877742e635362567a575866415639516c52656c35364b3167705238366e7353573147	3	admin	2022-08-05 14:47:22+00
3	2022-07-27 15:08:32+00	Test2	test2@test.com	\\x243261243132243735374177647a354a736959326e555172446a734c4f70505042546c66436a782e317564712f3333307033575a3853746572727757	6	user	2022-08-05 14:47:22+00
14	2022-07-29 08:37:53+00	Test3	test3@test.com	\\x243261243132246d416b4269572e363861706b395a6b70414b684e372e635230714a5a3139474c39686e6a5a30796750476b317366586c6370412f79	1	user	2022-08-05 14:47:22+00
15	2022-08-03 04:01:09+00	test4	test4@test.com	\\x243261243132244b5863737541366f6d74706541516a58622f62716f4f72453931546b776a772f5a4b6a75476f544a38392e37543838674637672f36	1	user	2022-08-05 14:47:22+00
\.


--
-- Name: authors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.authors_id_seq', 2, true);


--
-- Name: books_genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.books_genres_id_seq', 25, true);


--
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.books_id_seq', 9, true);


--
-- Name: genres_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.genres_id_seq', 9, true);


--
-- Name: tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.tokens_id_seq', 80, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.users_id_seq', 15, true);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (id);


--
-- Name: books_genres books_genres_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books_genres
    ADD CONSTRAINT books_genres_pkey PRIMARY KEY (id);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: genres genres_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.genres
    ADD CONSTRAINT genres_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: books_title_idx; Type: INDEX; Schema: public; Owner: devuser
--

CREATE INDEX books_title_idx ON public.books USING gin (to_tsvector('simple'::regconfig, (title)::text));


--
-- Name: books books_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.authors(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: books_genres books_genres_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books_genres
    ADD CONSTRAINT books_genres_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: books_genres books_genres_genre_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.books_genres
    ADD CONSTRAINT books_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: tokens tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

