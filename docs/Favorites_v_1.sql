CREATE TABLE "asset_types"
(
    "name"       varchar PRIMARY KEY,
    "created_at" timestamp
);

CREATE TABLE "asset"
(
    "isin"       varchar PRIMARY KEY,
    "asset_type" varchar,
    "created_at" timestamp
);

CREATE TABLE "users"
(
    "id"         uuid PRIMARY KEY,
    "upk"        varchar,
    "created_at" timestamp
);

CREATE TABLE "favorites"
(
    "id"         integer PRIMARY KEY,
    "isin"       varchar,
    "user_id"    uuid,
    "asset_type" varchar,
    "deleted"    bool,
    "created_at" timestamp,
    "updated_at" timestamp
);

COMMENT ON COLUMN "favorites"."isin" IS 'Content of the post';

ALTER TABLE "favorites"
    ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "favorites"
    ADD FOREIGN KEY ("isin") REFERENCES "asset" ("isin");

ALTER TABLE "asset_types"
    ADD FOREIGN KEY ("name") REFERENCES "asset" ("asset_type");
