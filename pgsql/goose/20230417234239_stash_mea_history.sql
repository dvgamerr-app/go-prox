-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS "stash";
CREATE TABLE "stash"."mea_history" (
	"t_monthy" timestamptz NOT NULL,
	"t_added" timestamptz NOT NULL,
	"n_unit" int8 NOT NULL DEFAULT 0,
	"n_baht" float8 NOT NULL DEFAULT 0,
	"t_created" timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stash"."mea_history";
-- +goose StatementEnd
