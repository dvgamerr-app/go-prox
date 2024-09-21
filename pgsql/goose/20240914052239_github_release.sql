-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS "stash";

CREATE TABLE "stash"."github_release" (
	"owner" varchar(50) NOT NULL,
	"name" varchar(50) NOT NULL,
	"tag_name" varchar(10) NOT NULL,
	"published" TIMESTAMP WITH TIME ZONE,
	"url" varchar(500) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stash"."github_release" CASCADE;
-- +goose StatementEnd
