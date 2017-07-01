CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE images(
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  slug character varying(255) NOT NULL,
  title character varying(255) NOT NULL,
  url character varying(1024) NOT NULL,
  metadata jsonb,
  published_at timestamp with time zone DEFAULT now(),
  expired_at timestamp with time zone,
  publisher  character varying(255) NOT NULL,
  created_at timestamp with time zone DEFAULT now() NOT NULL,
  updated_at timestamp with time zone,
  restored_at timestamp with time zone,
  deleted_at timestamp with time zone
);

--
-- Name: images_pkey; Type: CONSTRAINT;
ALTER TABLE ONLY images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: images_slug_unique; Type: CONSTRAINT; 
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_slug_unique UNIQUE (slug);

--
-- Name: images_url_unique; Type: CONSTRAINT; Schema: public; Owner: intirf601_usr
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_url_unique UNIQUE (url);
