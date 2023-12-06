-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-11-23T15:12:23.077Z

CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "places" (
  "google_id" varchar PRIMARY KEY NOT NULL,
  "tw_display_name" varchar NOT NULL,
  "tw_formatted_address" varchar NOT NULL,
  "tw_weekday_descriptions" varchar[],
  "administrative_area_level_1" varchar,
  "country" varchar,
  "google_map_uri" varchar NOT NULL,
  "international_phone_number" varchar,
  "lat" numeric NOT NULL,
  "lng" numeric NOT NULL,
  "primary_type" varchar,
  "rating" numeric,
  "types" varchar[],
  "user_rating_count" integer,
  "website_uri" varchar,
  "place_version" integer NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "favorites" (
  "favorite_id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "google_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "places" ("google_id");

CREATE INDEX ON "favorites" ("user_id");

CREATE INDEX ON "favorites" ("google_id");

ALTER TABLE "favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "favorites" ADD FOREIGN KEY ("google_id") REFERENCES "places" ("google_id");
