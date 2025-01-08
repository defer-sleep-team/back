--
-- PostgreSQL database dump
--

-- Dumped from database version 11.22 (Debian 11.22-4.pgdg110+1)
-- Dumped by pg_dump version 16.4

-- Started on 2024-09-29 13:49:47

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

DROP DATABASE IF EXISTS aether;
--
-- TOC entry 3101 (class 1262 OID 25632)
-- Name: aether; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE aether WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE aether OWNER TO postgres;

\connect aether

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
-- TOC entry 6 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- TOC entry 3103 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 220 (class 1255 OID 33834)
-- Name: after_post_insert(); Type: FUNCTION; Schema: public; Owner: dbuser
--

CREATE FUNCTION public.after_post_insert() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Вставка записи в таблицу post_stats
    INSERT INTO post_stats (post_id, likes, views) VALUES (NEW.id, 0, 0);

    -- Вставка записи в таблицу ratios
    INSERT INTO ratios (post_id, ratio) VALUES (NEW.id, 10000);

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.after_post_insert() OWNER TO dbuser;

SET default_tablespace = '';

--
-- TOC entry 217 (class 1259 OID 25819)
-- Name: comments; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.comments (
    id integer NOT NULL,
    content character varying(2048) NOT NULL,
    reg_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.comments OWNER TO dbuser;

--
-- TOC entry 216 (class 1259 OID 25817)
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.comments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.comments_id_seq OWNER TO dbuser;

--
-- TOC entry 3105 (class 0 OID 0)
-- Dependencies: 216
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;


--
-- TOC entry 206 (class 1259 OID 25721)
-- Name: plan_posts; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.plan_posts (
    subscription_plan_id integer,
    post_id integer
);


ALTER TABLE public.plan_posts OWNER TO dbuser;

--
-- TOC entry 218 (class 1259 OID 25829)
-- Name: post_comments; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.post_comments (
    post_id integer,
    comment_id integer
);


ALTER TABLE public.post_comments OWNER TO dbuser;

--
-- TOC entry 208 (class 1259 OID 25736)
-- Name: post_images; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.post_images (
    id integer NOT NULL,
    post_id integer,
    image_url character varying(255) NOT NULL
);


ALTER TABLE public.post_images OWNER TO dbuser;

--
-- TOC entry 207 (class 1259 OID 25734)
-- Name: post_images_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.post_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.post_images_id_seq OWNER TO dbuser;

--
-- TOC entry 3106 (class 0 OID 0)
-- Dependencies: 207
-- Name: post_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.post_images_id_seq OWNED BY public.post_images.id;


--
-- TOC entry 213 (class 1259 OID 25783)
-- Name: post_likes; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.post_likes (
    post_id integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.post_likes OWNER TO dbuser;

--
-- TOC entry 214 (class 1259 OID 25798)
-- Name: post_stats; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.post_stats (
    post_id integer,
    likes integer DEFAULT 0,
    views integer DEFAULT 0
);


ALTER TABLE public.post_stats OWNER TO dbuser;

--
-- TOC entry 211 (class 1259 OID 25757)
-- Name: post_tags; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.post_tags (
    post_id integer,
    tag_id integer
);


ALTER TABLE public.post_tags OWNER TO dbuser;

--
-- TOC entry 204 (class 1259 OID 25696)
-- Name: posts; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.posts (
    id integer NOT NULL,
    description character varying(2048),
    is_private boolean DEFAULT false,
    is_nsfw boolean DEFAULT false,
    reg_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.posts OWNER TO dbuser;

--
-- TOC entry 203 (class 1259 OID 25694)
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.posts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.posts_id_seq OWNER TO dbuser;

--
-- TOC entry 3107 (class 0 OID 0)
-- Dependencies: 203
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- TOC entry 215 (class 1259 OID 25808)
-- Name: ratios; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.ratios (
    post_id integer,
    ratio integer DEFAULT 10000
);


ALTER TABLE public.ratios OWNER TO dbuser;

--
-- TOC entry 200 (class 1259 OID 25657)
-- Name: subscription_plans; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.subscription_plans (
    id integer NOT NULL,
    user_id integer,
    name character varying(50) NOT NULL,
    price numeric(10,2) NOT NULL
);


ALTER TABLE public.subscription_plans OWNER TO dbuser;

--
-- TOC entry 199 (class 1259 OID 25655)
-- Name: subscription_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.subscription_plans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.subscription_plans_id_seq OWNER TO dbuser;

--
-- TOC entry 3108 (class 0 OID 0)
-- Dependencies: 199
-- Name: subscription_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.subscription_plans_id_seq OWNED BY public.subscription_plans.id;


--
-- TOC entry 210 (class 1259 OID 25749)
-- Name: tags; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.tags (
    id integer NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.tags OWNER TO dbuser;

--
-- TOC entry 209 (class 1259 OID 25747)
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tags_id_seq OWNER TO dbuser;

--
-- TOC entry 3109 (class 0 OID 0)
-- Dependencies: 209
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- TOC entry 219 (class 1259 OID 25842)
-- Name: user_comments; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_comments (
    user_id integer,
    comment_id integer
);


ALTER TABLE public.user_comments OWNER TO dbuser;

--
-- TOC entry 202 (class 1259 OID 25681)
-- Name: user_followers; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_followers (
    follower_id integer,
    followee_id integer
);


ALTER TABLE public.user_followers OWNER TO dbuser;

--
-- TOC entry 198 (class 1259 OID 25646)
-- Name: user_ips; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_ips (
    user_id integer,
    ip_address character varying(45) NOT NULL,
    last_login timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_ips OWNER TO dbuser;

--
-- TOC entry 205 (class 1259 OID 25708)
-- Name: user_posts; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_posts (
    user_id integer,
    post_id integer
);


ALTER TABLE public.user_posts OWNER TO dbuser;

--
-- TOC entry 201 (class 1259 OID 25668)
-- Name: user_subscriptions; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_subscriptions (
    user_id integer,
    subscription_plan_id integer
);


ALTER TABLE public.user_subscriptions OWNER TO dbuser;

--
-- TOC entry 212 (class 1259 OID 25770)
-- Name: user_tags; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.user_tags (
    user_id integer,
    tag_id integer
);


ALTER TABLE public.user_tags OWNER TO dbuser;

--
-- TOC entry 197 (class 1259 OID 25635)
-- Name: users; Type: TABLE; Schema: public; Owner: dbuser
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100),
    password character varying(100) NOT NULL,
    avatar character varying(40) DEFAULT '1.jpg'::character varying,
    bio character varying(1024) DEFAULT 'Это твоё описание'::character varying,
    privilege_level integer DEFAULT 0,
    payments character varying(55),
    block boolean DEFAULT false
);


ALTER TABLE public.users OWNER TO dbuser;

--
-- TOC entry 196 (class 1259 OID 25633)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: dbuser
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO dbuser;

--
-- TOC entry 3110 (class 0 OID 0)
-- Dependencies: 196
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: dbuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 2923 (class 2604 OID 25822)
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);


--
-- TOC entry 2918 (class 2604 OID 25739)
-- Name: post_images id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_images ALTER COLUMN id SET DEFAULT nextval('public.post_images_id_seq'::regclass);


--
-- TOC entry 2914 (class 2604 OID 25699)
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- TOC entry 2913 (class 2604 OID 25660)
-- Name: subscription_plans id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.subscription_plans ALTER COLUMN id SET DEFAULT nextval('public.subscription_plans_id_seq'::regclass);


--
-- TOC entry 2919 (class 2604 OID 25752)
-- Name: tags id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- TOC entry 2907 (class 2604 OID 25638)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 2926 (class 2606 OID 25863)
-- Name: ratios ck_ratios_ratio; Type: CHECK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE public.ratios
    ADD CONSTRAINT ck_ratios_ratio CHECK ((ratio >= 0)) NOT VALID;


--
-- TOC entry 2934 (class 2606 OID 25862)
-- Name: subscription_plans ck_sub_name; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT ck_sub_name UNIQUE (name);


--
-- TOC entry 2925 (class 2606 OID 25864)
-- Name: user_followers ck_user_follow_to_self; Type: CHECK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE public.user_followers
    ADD CONSTRAINT ck_user_follow_to_self CHECK ((follower_id <> followee_id)) NOT VALID;


--
-- TOC entry 2932 (class 2606 OID 25868)
-- Name: user_ips ck_user_ips; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_ips
    ADD CONSTRAINT ck_user_ips UNIQUE (ip_address);


--
-- TOC entry 2950 (class 2606 OID 25828)
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- TOC entry 2942 (class 2606 OID 25741)
-- Name: post_images post_images_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_images
    ADD CONSTRAINT post_images_pkey PRIMARY KEY (id);


--
-- TOC entry 2948 (class 2606 OID 25787)
-- Name: post_likes post_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT post_likes_pkey PRIMARY KEY (post_id, user_id);


--
-- TOC entry 2940 (class 2606 OID 25707)
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- TOC entry 2936 (class 2606 OID 25662)
-- Name: subscription_plans subscription_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);


--
-- TOC entry 2944 (class 2606 OID 25756)
-- Name: tags tags_name_key; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_name_key UNIQUE (name);


--
-- TOC entry 2946 (class 2606 OID 25754)
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- TOC entry 2938 (class 2606 OID 25858)
-- Name: user_followers user_followers_follower_id_followee_id_key; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_followers
    ADD CONSTRAINT user_followers_follower_id_followee_id_key UNIQUE (follower_id, followee_id);


--
-- TOC entry 2928 (class 2606 OID 33833)
-- Name: users username; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT username UNIQUE (username);


--
-- TOC entry 2930 (class 2606 OID 25645)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 2974 (class 2620 OID 33835)
-- Name: posts post_insert_trigger; Type: TRIGGER; Schema: public; Owner: dbuser
--

CREATE TRIGGER post_insert_trigger AFTER INSERT ON public.posts FOR EACH ROW EXECUTE PROCEDURE public.after_post_insert();


--
-- TOC entry 2959 (class 2606 OID 25729)
-- Name: plan_posts plan_posts_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.plan_posts
    ADD CONSTRAINT plan_posts_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2960 (class 2606 OID 25724)
-- Name: plan_posts plan_posts_subscription_plan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.plan_posts
    ADD CONSTRAINT plan_posts_subscription_plan_id_fkey FOREIGN KEY (subscription_plan_id) REFERENCES public.subscription_plans(id);


--
-- TOC entry 2970 (class 2606 OID 25837)
-- Name: post_comments post_comments_comment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_comments
    ADD CONSTRAINT post_comments_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES public.comments(id);


--
-- TOC entry 2971 (class 2606 OID 25832)
-- Name: post_comments post_comments_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_comments
    ADD CONSTRAINT post_comments_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2961 (class 2606 OID 25742)
-- Name: post_images post_images_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_images
    ADD CONSTRAINT post_images_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2966 (class 2606 OID 25788)
-- Name: post_likes post_likes_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT post_likes_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2967 (class 2606 OID 25793)
-- Name: post_likes post_likes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_likes
    ADD CONSTRAINT post_likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2968 (class 2606 OID 25803)
-- Name: post_stats post_stats_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_stats
    ADD CONSTRAINT post_stats_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2962 (class 2606 OID 25760)
-- Name: post_tags post_tags_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2963 (class 2606 OID 25765)
-- Name: post_tags post_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- TOC entry 2969 (class 2606 OID 25812)
-- Name: ratios ratios_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.ratios
    ADD CONSTRAINT ratios_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2952 (class 2606 OID 25663)
-- Name: subscription_plans subscription_plans_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2972 (class 2606 OID 25850)
-- Name: user_comments user_comments_comment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_comments
    ADD CONSTRAINT user_comments_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES public.comments(id);


--
-- TOC entry 2973 (class 2606 OID 25845)
-- Name: user_comments user_comments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_comments
    ADD CONSTRAINT user_comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2955 (class 2606 OID 25689)
-- Name: user_followers user_followers_followee_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_followers
    ADD CONSTRAINT user_followers_followee_id_fkey FOREIGN KEY (followee_id) REFERENCES public.users(id);


--
-- TOC entry 2956 (class 2606 OID 25684)
-- Name: user_followers user_followers_follower_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_followers
    ADD CONSTRAINT user_followers_follower_id_fkey FOREIGN KEY (follower_id) REFERENCES public.users(id);


--
-- TOC entry 2951 (class 2606 OID 25650)
-- Name: user_ips user_ips_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_ips
    ADD CONSTRAINT user_ips_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2957 (class 2606 OID 25716)
-- Name: user_posts user_posts_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_posts
    ADD CONSTRAINT user_posts_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- TOC entry 2958 (class 2606 OID 25711)
-- Name: user_posts user_posts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_posts
    ADD CONSTRAINT user_posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2953 (class 2606 OID 25676)
-- Name: user_subscriptions user_subscriptions_subscription_plan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_subscription_plan_id_fkey FOREIGN KEY (subscription_plan_id) REFERENCES public.subscription_plans(id);


--
-- TOC entry 2954 (class 2606 OID 25671)
-- Name: user_subscriptions user_subscriptions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 2964 (class 2606 OID 25778)
-- Name: user_tags user_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_tags
    ADD CONSTRAINT user_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id);


--
-- TOC entry 2965 (class 2606 OID 25773)
-- Name: user_tags user_tags_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: dbuser
--

ALTER TABLE ONLY public.user_tags
    ADD CONSTRAINT user_tags_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3102 (class 0 OID 0)
-- Dependencies: 3101
-- Name: DATABASE aether; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON DATABASE aether TO dbuser;


--
-- TOC entry 3104 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE USAGE ON SCHEMA public FROM PUBLIC;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2024-09-29 13:49:52

--
-- PostgreSQL database dump complete
--

