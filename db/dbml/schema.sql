-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-11-23T15:12:23.077Z

CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "places" (
  "place_id" bigserial PRIMARY KEY,
  "google_id" varchar,
  "tw_display_name" varchar,
  "jp_display_name" varchar,
  "primary_type" varchar,
  "rating" numeric(2,1),
  "user_rating_count" int,
  "jp_formatted_address" varchar,
  "en_city" varchar,
  "jp_district" varchar,
  "international_phone_number" varchar,
  "tw_weekday_descriptions" varchar[],
  "accessibility_options" varchar[],
  "google_map_uri" varchar,
  "website_uri" varchar,
  "photos_name" varchar[],
  "types" varchar[],
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "favorites" (
  "user_id" bigserial,
  "place_id" bigserial,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "favorites" ADD FOREIGN KEY ("place_id") REFERENCES "places" ("place_id");