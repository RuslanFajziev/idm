--- employee ---
CREATE TABLE IF NOT EXISTS "employee"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz not null,
    "updated_at" timestamptz,
	primary key ("id")
);

--- role ---
CREATE TABLE IF NOT EXISTS "role"
(
    "id" bigint GENERATED ALWAYS AS IDENTITY,
    "name" text not null,
    "create_at" timestamptz not null,
    "updated_at" timestamptz,
	primary key ("id")
);