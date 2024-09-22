-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS "notice";

CREATE TYPE "notice"."notify" AS ENUM (
  'telegram',
  'slack',
  'msteam',
  'line',
  'line-notify',
  'workplace',
  'email',
  'webhook',
  'native'
);

CREATE TABLE "notice"."provider" (
  "id" serial PRIMARY KEY,
  "s_liff_id" varchar(33) NOT NULL,
  "e_type" "notice"."notify"NOT NULL,
  "o_param" jsonb NOT NULL,
  "b_deleted" boolean DEFAULT false,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "notice"."room" (
  "id" serial PRIMARY KEY,
  "notice_provider_id" int4 NOT NULL,
  "o_param" jsonb NOT NULL,
  "b_deleted" boolean DEFAULT false,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "notice"."section" (
  "id" serial PRIMARY KEY,
  "s_liff_id" varchar(33) NOT NULL,
  "s_name" varchar(20) NOT NULL,
  "n_uuid" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "t_deleted" timestamp WITH TIME ZONE DEFAULT NULL,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE "notice"."subscriber" (
  "notice_section_id" int4 NOT NULL,
  "notice_room_id" int4 NOT NULL,
  "t_deleted" timestamp WITH TIME ZONE DEFAULT NULL,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "notice"."event" (
  "notice_room_id" int4 NOT NULL,
  "o_payload" jsonb NOT NULL,
  "b_sended" boolean DEFAULT false,
  "t_created" timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE "notice"."room" ADD FOREIGN KEY ("notice_provider_id") REFERENCES "notice"."provider" ("id");
ALTER TABLE "notice"."subscriber" ADD FOREIGN KEY ("notice_section_id") REFERENCES "notice"."section" ("id");
ALTER TABLE "notice"."subscriber" ADD FOREIGN KEY ("notice_room_id") REFERENCES "notice"."room" ("id");
ALTER TABLE "notice"."event" ADD FOREIGN KEY ("notice_room_id") REFERENCES "notice"."room" ("id");

ALTER TABLE "notice"."section" ADD CONSTRAINT uq_notice_section UNIQUE ("s_liff_id", "s_name");

CREATE INDEX "idx_section__deleted" ON "notice"."section" USING BTREE ("t_deleted");
CREATE INDEX "idx_provider__deleted" ON "notice"."provider" USING BTREE ("b_deleted");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "notice"."subscriber";
DROP TABLE "notice"."event";
DROP TABLE "notice"."section";
DROP TABLE "notice"."room";
DROP TABLE "notice"."provider";
DROP TYPE "notice"."notify";
-- +goose StatementEnd
