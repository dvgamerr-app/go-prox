-- +goose Up
-- +goose StatementBegin

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
