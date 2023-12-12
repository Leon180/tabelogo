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
  "jp_display_name" varchar NOT NULL DEFAULT '',
  "tw_formatted_address" varchar NOT NULL,
  "tw_weekday_descriptions" varchar[] NOT NULL DEFAULT '{}',
  "administrative_area_level_1" varchar NOT NULL DEFAULT '',
  "country" varchar NOT NULL DEFAULT '',
  "google_map_uri" varchar NOT NULL,
  "international_phone_number" varchar NOT NULL DEFAULT '',
  "lat" numeric NOT NULL,
  "lng" numeric NOT NULL,
  "primary_type" varchar NOT NULL DEFAULT '',
  "rating" numeric NOT NULL DEFAULT 0,
  "types" varchar[] NOT NULL DEFAULT '{}',
  "user_rating_count" integer NOT NULL DEFAULT 0,
  "website_uri" varchar NOT NULL DEFAULT '',
  "place_version" integer NOT NULL DEFAULT 1,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "favorites" (
  "is_favorite" boolean NOT NULL DEFAULT true,
  "user_email" varchar NOT NULL,
  "google_id" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("user_email", "google_id")
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "places" ("google_id");

ALTER TABLE "favorites" ADD FOREIGN KEY ("user_email") REFERENCES "users" ("email");

ALTER TABLE "favorites" ADD FOREIGN KEY ("google_id") REFERENCES "places" ("google_id");
