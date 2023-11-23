-- name: CreatePlace :one
INSERT INTO places (
    google_id,
    tw_display_name,
    jp_display_name,
    primary_type,
    rating,
    user_rating_count,
    jp_formatted_address,
    en_city,
    jp_district,
    international_phone_number,
    tw_weekday_descriptions,
    accessibility_options,
    google_map_uri,
    website_uri,
    photos_name,
    types
) VALUES (
    $1, 
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16
) RETURNING *;

-- name: GetPlaceById :one
SELECT * FROM places
WHERE place_id = $1 LIMIT 1;

-- name: GetPlaceByGoogleId :one
SELECT * FROM places
WHERE google_id = $1 LIMIT 1;

-- name: UpdatePlace :one
UPDATE places SET
    google_id = COALESCE(sqlc.narg(google_id), google_id),
    tw_display_name = COALESCE(sqlc.narg(tw_display_name), tw_display_name),
    jp_display_name = COALESCE(sqlc.narg(jp_display_name), jp_display_name),
    primary_type = COALESCE(sqlc.narg(primary_type), primary_type),
    rating = COALESCE(sqlc.narg(rating), rating),
    user_rating_count = COALESCE(sqlc.narg(user_rating_count), user_rating_count),
    jp_formatted_address = COALESCE(sqlc.narg(jp_formatted_address), jp_formatted_address),
    en_city = COALESCE(sqlc.narg(en_city), en_city),
    jp_district = COALESCE(sqlc.narg(jp_district), jp_district),
    international_phone_number = COALESCE(sqlc.narg(international_phone_number), international_phone_number),
    tw_weekday_descriptions = COALESCE(sqlc.narg(tw_weekday_descriptions), tw_weekday_descriptions),
    accessibility_options = COALESCE(sqlc.narg(accessibility_options), accessibility_options),
    google_map_uri = COALESCE(sqlc.narg(google_map_uri), google_map_uri),
    website_uri = COALESCE(sqlc.narg(website_uri), website_uri),
    photos_name = COALESCE(sqlc.narg(photos_name), photos_name),
    types = COALESCE(sqlc.narg(types), types)
WHERE place_id = sqlc.arg(place_id)
RETURNING *;

-- name: DeletePlace :exec
DELETE FROM places
WHERE place_id = $1;
