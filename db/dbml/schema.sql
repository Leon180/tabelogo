-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-12-04T05:31:05.975Z

CREATE TABLE "users" (
  "user_id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "active" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "places" (
  "place_id" bigserial PRIMARY KEY,
  "google_id" varchar UNIQUE NOT NULL,
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
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "favorites" (
  "favorite_id" bigserial PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "place_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "favorites" ADD FOREIGN KEY ("place_id") REFERENCES "places" ("place_id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");
