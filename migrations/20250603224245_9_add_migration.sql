-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "employee"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz DEFAULT now(),
    "update_at" timestamptz DEFAULT now(),

    primary key ("id")
);

CREATE TABLE IF NOT EXISTS "role"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz DEFAULT now(),
    "update_at" timestamptz DEFAULT now(),

    CONSTRAINT "pk_role_name" PRIMARY KEY ("name")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "employee";
DROP TABLE "role";
-- +goose StatementEnd
