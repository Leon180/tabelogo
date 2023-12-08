-- name: CreatePlace :one
INSERT INTO places (
  google_id,
  tw_display_name,
  tw_formatted_address,
  tw_weekday_descriptions,
  administrative_area_level_1,
  country,
  google_map_uri,
  international_phone_number,
  lat,
  lng,
  primary_type,
  rating,
  types,
  user_rating_count,
  website_uri,
  place_version
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 1
) RETURNING *;

-- name: GetPlaceByGoogleId :one
SELECT * FROM places
WHERE google_id = $1 LIMIT 1;

-- name: UpdatePlace :one
UPDATE places SET
    tw_display_name = COALESCE(sqlc.narg(tw_display_name), tw_display_name),
    tw_formatted_address = COALESCE(sqlc.narg(tw_formatted_address), tw_formatted_address),
    tw_weekday_descriptions = COALESCE(sqlc.narg(tw_weekday_descriptions), tw_weekday_descriptions),
    administrative_area_level_1 = COALESCE(sqlc.narg(administrative_area_level_1), administrative_area_level_1),
    country = COALESCE(sqlc.narg(country), country),
    google_map_uri = COALESCE(sqlc.narg(google_map_uri), google_map_uri),
    international_phone_number = COALESCE(sqlc.narg(international_phone_number), international_phone_number),
    lat = COALESCE(sqlc.narg(lat), lat),
    lng = COALESCE(sqlc.narg(lng), lng),
    primary_type = COALESCE(sqlc.narg(primary_type), primary_type),
    rating = COALESCE(sqlc.narg(rating), rating),
    types = COALESCE(sqlc.narg(types), types),
    user_rating_count = COALESCE(sqlc.narg(user_rating_count), user_rating_count),
    website_uri = COALESCE(sqlc.narg(website_uri), website_uri),
    place_version = (place_version + 1)
WHERE google_id = sqlc.arg(google_id) AND place_version = sqlc.arg(place_version)
RETURNING *;

-- name: DeletePlace :exec
DELETE FROM places
WHERE google_id = $1 AND place_version = $2;
