
--
-- TOC entry 199 (class 1259 OID 16423)
-- Name: operator; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.operator (
    id uuid NOT NULL,
    name character varying NOT NULL,
    organization_id uuid NOT NULL
);


ALTER TABLE public.operator OWNER TO postgres;

--
-- TOC entry 196 (class 1259 OID 16385)
-- Name: organization; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.organization (
    id uuid NOT NULL,
    name character varying NOT NULL,
    api_key character varying NOT NULL
);


ALTER TABLE public.organization OWNER TO postgres;

--
-- TOC entry 204 (class 1259 OID 16515)
-- Name: run; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.run (
    id uuid NOT NULL,
    name character varying NOT NULL,
    start_time bigint NOT NULL,
    end_time bigint NOT NULL,
    session_id uuid,
    thing_id uuid NOT NULL,
    CONSTRAINT start_time_less_than_end_time CHECK ((start_time < end_time))
);


ALTER TABLE public.run OWNER TO postgres;

--
-- TOC entry 205 (class 1259 OID 16609)
-- Name: run_comment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.run_comment (
    id uuid NOT NULL,
    run_id uuid NOT NULL,
    user_id uuid NOT NULL,
    last_update bigint NOT NULL,
    content character varying
);


ALTER TABLE public.run_comment OWNER TO postgres;

--
-- TOC entry 198 (class 1259 OID 16409)
-- Name: sensor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sensor (
    id uuid NOT NULL,
    small_id integer NOT NULL,
    type character varying NOT NULL,
    last_update bigint NOT NULL,
    category character varying,
    name character varying NOT NULL,
    frequency integer NOT NULL,
    unit character varying,
    can_id bigint NOT NULL,
    disabled boolean,
    thing_id uuid NOT NULL,
    upper_calibration double precision,
    lower_calibration double precision,
    conversion_multiplier double precision,
    upper_warning double precision,
    lower_warning double precision,
    upper_danger double precision,
    lower_danger double precision,
    upper_bound double precision,
    lower_bound double precision,
    significance double precision,
    CONSTRAINT "smallId_valid_value_check" CHECK (((small_id >= 0) AND (small_id <= 255)))
);


ALTER TABLE public.sensor OWNER TO postgres;

--
-- TOC entry 201 (class 1259 OID 16463)
-- Name: session; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.session (
    id uuid NOT NULL,
    name character varying NOT NULL,
    description character varying,
    start_date bigint NOT NULL,
    end_date bigint NOT NULL,
    thing_id uuid NOT NULL,
    CONSTRAINT start_date_less_than_end_date CHECK ((start_date < end_date))
);


ALTER TABLE public.session OWNER TO postgres;

--
-- TOC entry 203 (class 1259 OID 16497)
-- Name: session_comment; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.session_comment (
    id uuid NOT NULL,
    session_id uuid NOT NULL,
    user_id uuid NOT NULL,
    last_update bigint NOT NULL,
    content character varying
);


ALTER TABLE public.session_comment OWNER TO postgres;

--
-- TOC entry 197 (class 1259 OID 16393)
-- Name: thing; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.thing (
    id uuid NOT NULL,
    name character varying NOT NULL,
    organization_id uuid NOT NULL
);


ALTER TABLE public.thing OWNER TO postgres;

--
-- TOC entry 200 (class 1259 OID 16436)
-- Name: thing_operator; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.thing_operator (
    id uuid NOT NULL,
    operator_id uuid NOT NULL,
    thing_id uuid NOT NULL
);


ALTER TABLE public.thing_operator OWNER TO postgres;

--
-- TOC entry 202 (class 1259 OID 16484)
-- Name: user; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."user" (
    id uuid NOT NULL,
    display_name character varying NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL,
    organization_id uuid,
    role character varying NOT NULL
);


ALTER TABLE public."user" OWNER TO postgres;

--
-- TOC entry 2840 (class 2606 OID 16430)
-- Name: operator operator_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operator
    ADD CONSTRAINT operator_pkey PRIMARY KEY (id);


--
-- TOC entry 2831 (class 2606 OID 16392)
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- TOC entry 2860 (class 2606 OID 16616)
-- Name: run_comment run_comment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run_comment
    ADD CONSTRAINT run_comment_pkey PRIMARY KEY (id);


--
-- TOC entry 2858 (class 2606 OID 16523)
-- Name: run run_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run
    ADD CONSTRAINT run_pkey PRIMARY KEY (id);


--
-- TOC entry 2836 (class 2606 OID 16417)
-- Name: sensor sensor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sensor
    ADD CONSTRAINT sensor_pkey PRIMARY KEY (id);


--
-- TOC entry 2855 (class 2606 OID 16504)
-- Name: session_comment session_comment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session_comment
    ADD CONSTRAINT session_comment_pkey PRIMARY KEY (id);


--
-- TOC entry 2848 (class 2606 OID 16470)
-- Name: session session_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (id);


--
-- TOC entry 2845 (class 2606 OID 16440)
-- Name: thing_operator thing_operator_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thing_operator
    ADD CONSTRAINT thing_operator_pkey PRIMARY KEY (id);


--
-- TOC entry 2833 (class 2606 OID 16400)
-- Name: thing thing_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thing
    ADD CONSTRAINT thing_pkey PRIMARY KEY (id);


--
-- TOC entry 2851 (class 2606 OID 16491)
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- TOC entry 2841 (class 1259 OID 16586)
-- Name: fki_o; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_o ON public.thing_operator USING btree (thing_id);


--
-- TOC entry 2842 (class 1259 OID 16580)
-- Name: fki_operator_id_fk1; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_operator_id_fk1 ON public.thing_operator USING btree (operator_id);


--
-- TOC entry 2837 (class 1259 OID 16456)
-- Name: fki_organization_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_organization_id ON public.operator USING btree (id);


--
-- TOC entry 2838 (class 1259 OID 16546)
-- Name: fki_organization_id_fk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_organization_id_fk ON public.operator USING btree (organization_id);


--
-- TOC entry 2856 (class 1259 OID 16552)
-- Name: fki_session_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_session_id ON public.run USING btree (session_id);


--
-- TOC entry 2834 (class 1259 OID 16462)
-- Name: fki_thing_id_fk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_thing_id_fk ON public.sensor USING btree (id);


--
-- TOC entry 2843 (class 1259 OID 16597)
-- Name: fki_thing_id_fk1; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_thing_id_fk1 ON public.thing_operator USING btree (thing_id);


--
-- TOC entry 2846 (class 1259 OID 16540)
-- Name: fki_thing_id_fl; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_thing_id_fl ON public.session USING btree (thing_id);


--
-- TOC entry 2852 (class 1259 OID 16568)
-- Name: fki_u; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_u ON public.session_comment USING btree (session_id);


--
-- TOC entry 2853 (class 1259 OID 16574)
-- Name: fki_user_id_fk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_user_id_fk ON public.session_comment USING btree (user_id);


--
-- TOC entry 2849 (class 1259 OID 16608)
-- Name: fki_user_id_org_id_fk; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX fki_user_id_org_id_fk ON public."user" USING btree (organization_id);


--
-- TOC entry 2865 (class 2606 OID 16575)
-- Name: thing_operator operator_id_fk1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thing_operator
    ADD CONSTRAINT operator_id_fk1 FOREIGN KEY (operator_id) REFERENCES public.operator(id) ON DELETE CASCADE;


--
-- TOC entry 2861 (class 2606 OID 16401)
-- Name: thing organizationId_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thing
    ADD CONSTRAINT "organizationId_fk" FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- TOC entry 2863 (class 2606 OID 16541)
-- Name: operator organization_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operator
    ADD CONSTRAINT organization_id_fk FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- TOC entry 2872 (class 2606 OID 16617)
-- Name: run_comment run_id_comment_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run_comment
    ADD CONSTRAINT run_id_comment_fk FOREIGN KEY (run_id) REFERENCES public.run(id) ON DELETE CASCADE;


--
-- TOC entry 2870 (class 2606 OID 16547)
-- Name: run session_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run
    ADD CONSTRAINT session_id_fk FOREIGN KEY (session_id) REFERENCES public.session(id) ON DELETE SET NULL;


--
-- TOC entry 2868 (class 2606 OID 16563)
-- Name: session_comment session_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session_comment
    ADD CONSTRAINT session_id_fk FOREIGN KEY (session_id) REFERENCES public.session(id) ON DELETE CASCADE;


--
-- TOC entry 2866 (class 2606 OID 16535)
-- Name: session thing_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT thing_id_fk FOREIGN KEY (thing_id) REFERENCES public.thing(id) ON DELETE CASCADE;


--
-- TOC entry 2871 (class 2606 OID 16553)
-- Name: run thing_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run
    ADD CONSTRAINT thing_id_fk FOREIGN KEY (thing_id) REFERENCES public.thing(id) ON DELETE CASCADE;


--
-- TOC entry 2862 (class 2606 OID 16558)
-- Name: sensor thing_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sensor
    ADD CONSTRAINT thing_id_fk FOREIGN KEY (thing_id) REFERENCES public.thing(id) ON DELETE CASCADE;


--
-- TOC entry 2864 (class 2606 OID 16592)
-- Name: thing_operator thing_id_fk1; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.thing_operator
    ADD CONSTRAINT thing_id_fk1 FOREIGN KEY (thing_id) REFERENCES public.thing(id) ON DELETE CASCADE;


--
-- TOC entry 2873 (class 2606 OID 16622)
-- Name: run_comment user_id_comment_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.run_comment
    ADD CONSTRAINT user_id_comment_fk FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- TOC entry 2869 (class 2606 OID 16569)
-- Name: session_comment user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.session_comment
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE;


--
-- TOC entry 2867 (class 2606 OID 16603)
-- Name: user user_id_org_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_id_org_id_fk FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE SET NULL;


-- Completed on 2022-05-21 17:33:23 MDT

--
-- PostgreSQL database dump complete
--

