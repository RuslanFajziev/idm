-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "employee"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz not null,
    "updated_at" timestamptz,
	primary key ("id")
);

CREATE TABLE IF NOT EXISTS "role"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz not null,
    "updated_at" timestamptz,
	primary key ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "employee";
DROP TABLE "role";
-- +goose StatementEnd
