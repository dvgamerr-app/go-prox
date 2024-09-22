-- +goose Up
-- +goose StatementBegin

-- CREATE TYPE manga_type AS ENUM ('manga', 'manhwa', 'doujin');
-- CREATE TYPE manga_lang AS ENUM ('TH', 'EN', 'CH', 'JP', 'KR');

-- CREATE TABLE IF NOT EXISTS "stash"."cinema_showing" (
--   "s_name" varchar(120) NOT NULL,
--   "s_bind" varchar(120),
--   "s_display" varchar(200) NOT NULL,
--   "t_release" timestamptz NOT NULL,
--   "s_genre" varchar(40) NOT NULL,
--   "n_week" int2 NOT NULL,
--   "n_year" int4 NOT NULL,
--   "n_time" int4 NOT NULL DEFAULT 0,
--   "s_url" text NOT NULL,
--   "s_cover" text NOT NULL,
--   "o_theater" jsonb NOT NULL DEFAULT '[]'::jsonb
-- );

-- CREATE TABLE IF NOT EXISTS "stash"."manga_collection" (
--   "s_name" varchar(255) NOT NULL,
--   "s_title" text NOT NULL,
--   "s_url" text NOT NULL,
--   "s_thumbnail" varchar(200) NOT NULL,
--   "b_translate" boolean NOT NULL DEFAULT 'f',
--   "e_type" manga_type NOT NULL DEFAULT 'manga',
--   "e_lang" manga_lang NOT NULL DEFAULT 'TH',
--   "n_total" int4 NOT NULL DEFAULT 0,
--   "o_image" jsonb NOT NULL DEFAULT '[]'
-- );



CREATE SCHEMA IF NOT EXISTS "comic";

CREATE TABLE "comic"."manga" (
  "id" serial PRIMARY KEY,
  "s_name" varchar(500) NOT NULL,
  "s_romanji" varchar(500),
  "s_kanji" varchar(500),
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "comic"."manga_chapter" (
  "manga_id" int4 NOT NULL,
  "s_chapter" varchar(5) NOT NULL,
  "n_no" int4 NOT NULL,
  "s_name" varchar(50),
  "b_deleted" boolean DEFAULT false,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
 
ALTER TABLE "comic"."manga_chapter" ADD FOREIGN KEY ("manga_id") REFERENCES  "comic"."manga" ("id");

ALTER TABLE "comic"."manga_chapter" ADD CONSTRAINT uq_manga_chapter UNIQUE ("manga_id", "s_chapter", "n_no");

CREATE INDEX "idx_manga_chapter__deleted" ON "comic"."manga_chapter" USING BTREE ("manga_id", "s_chapter");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "comic"."manga_chapter" CASCADE;
DROP TABLE "comic"."manga" CASCADE;
-- +goose StatementEnd
