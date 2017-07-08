-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE EXTENSION IF NOT EXISTS hstore WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS artists (
    id integer PRIMARY KEY,
    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    urls text[] DEFAULT '{}'::text[],
    tags text[] DEFAULT '{}'::text[]
);

CREATE INDEX IF NOT EXISTS index_artists_on_tags ON artists USING gin (tags);
CREATE INDEX IF NOT EXISTS index_artists_on_urls ON artists USING gin (urls);

CREATE TABLE IF NOT EXISTS categories (
    id integer PRIMARY KEY,
    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    tags text[] DEFAULT '{}'::text[],
    ignored_tags text[] DEFAULT '{}'::text[],
    product_id integer
);

CREATE INDEX IF NOT EXISTS index_categories_on_tags ON categories USING gin (tags);
CREATE INDEX IF NOT EXISTS index_categories_on_ignored_tags ON categories USING gin (ignored_tags);

CREATE TABLE IF NOT EXISTS designs (
    id integer PRIMARY KEY,
    artist_id integer NOT NULL references artists(id) ON DELETE CASCADE,
    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    description text DEFAULT '' NOT NULL,
    tags text[] DEFAULT '{}'::text[],
    category_tags text[] DEFAULT '{}'::text[],
    mature boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS category_designs (
    id integer PRIMARY KEY,
    category_id integer references categories(id) ON DELETE CASCADE,
    design_id integer references designs(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS index_designs_on_tags ON designs USING gin (tags);
CREATE INDEX IF NOT EXISTS index_designs_on_category_tags ON designs USING gin (category_tags);

CREATE TABLE IF NOT EXISTS sites (
    id integer PRIMARY KEY,
    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    domain_name text NOT NULL,
    affiliate_url text,
    deal_scraper boolean DEFAULT false NOT NULL,
    full_scraper boolean DEFAULT false NOT NULL,
    active boolean DEFAULT false NOT NULL,
    display_order integer DEFAULT 99 NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id integer PRIMARY KEY,
    design_id integer NOT NULL references designs(id) ON DELETE CASCADE,
    site_id integer NOT NULL references sites(id) ON DELETE CASCADE,
    slug text NOT NULL UNIQUE,
    url text NOT NULL,
    active boolean DEFAULT false NOT NULL,
    deal boolean DEFAULT false NOT NULL,
    last_chance boolean DEFAULT false NOT NULL,
    tags text[] DEFAULT '{}'::text[],
    prices hstore DEFAULT ''::hstore NOT NULL,
    expires_at timestamp without time zone,
    active_at timestamp without time zone,
    image_background text DEFAULT '#000000'::text NOT NULL,
    image_updated_at date DEFAULT '2010-01-01'::date NOT NULL
);

CREATE INDEX IF NOT EXISTS index_products_on_active_and_deal ON products USING btree (active, deal);
CREATE UNIQUE INDEX IF NOT EXISTS index_products_on_site_id_and_design_id ON products USING btree (site_id, design_id);
CREATE INDEX IF NOT EXISTS index_products_on_tags ON products USING gin (tags);

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    email text NOT NULL UNIQUE,
    admin boolean DEFAULT false NOT NULL,
    api_access boolean DEFAULT false NOT NULL,
    api_token text NOT NULL,
    encrypted_password text NOT NULL
);

CREATE INDEX IF NOT EXISTS index_users_on_api_access_and_api_token ON users USING btree (api_access, api_token);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
