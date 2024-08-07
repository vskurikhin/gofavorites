CREATE TABLE "asset_types"
(
    "name"       varchar PRIMARY KEY,
    "deleted"    bool,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "assets"
(
    "isin"       varchar PRIMARY KEY,
    "asset_type" varchar NOT NULL,
    "deleted"    bool,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "users"
(
    "upk"        varchar PRIMARY KEY,
    "deleted"    bool,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE TABLE "favorites"
(
    "id"         uuid PRIMARY KEY,
    "isin"       varchar NOT NULL,
    "user_upk"   varchar NOT NULL,
    "version"    bigint,
    "deleted"    bool,
    "created_at" timestamp,
    "updated_at" timestamp
);

CREATE UNIQUE INDEX ON "favorites" ("isin", "user_upk");

COMMENT ON COLUMN "asset_types"."name" IS 'Type of asset, example bonds, stocks, funds ... etc';

COMMENT ON COLUMN "assets"."isin" IS 'International Securities Identification Numbers';

COMMENT ON COLUMN "users"."upk" IS 'User primary key';

COMMENT ON COLUMN "favorites"."version" IS 'CBUL version';

ALTER TABLE "favorites"
    ADD FOREIGN KEY ("user_upk") REFERENCES "users" ("upk");

ALTER TABLE "favorites"
    ADD FOREIGN KEY ("isin") REFERENCES "assets" ("isin");

ALTER TABLE "assets"
    ADD FOREIGN KEY ("asset_type") REFERENCES "asset_types" ("name");
