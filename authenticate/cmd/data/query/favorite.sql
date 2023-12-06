-- name: CreateFavorite :one
INSERT INTO favorites (
    user_id,
    google_id
) VALUES (
    $1, 
    $2
) RETURNING *;

-- name: ListFavoritesByCreateTime :many
SELECT
    favorites.google_id,
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.google_id = places.google_id
WHERE user_id = $1
ORDER BY favorites.created_at ASC
LIMIT $2
OFFSET $3;

-- name: ListFavoritesByCountry :many
SELECT
    favorites.google_id,
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.google_id = places.google_id
WHERE user_id = $1 AND country = $2
ORDER BY favorites.created_at ASC
LIMIT $3
OFFSET $4;

-- name: ListFavoritesByCountrAndRegion :many
SELECT
    favorites.google_id,
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.google_id = places.google_id
WHERE user_id = $1 AND country = $2 AND administrative_area_level_1 = $3
ORDER BY favorites.created_at ASC
LIMIT $4
OFFSET $5;

-- name: GetCountryList :many
SELECT DISTINCT country FROM favorites JOIN places ON favorites.google_id = places.google_id
WHERE user_id = $1
ORDER BY country ASC;

-- name: GetRegionList :many
SELECT DISTINCT administrative_area_level_1 FROM favorites JOIN places ON favorites.google_id = places.google_id
WHERE user_id = $1 AND country = $2
ORDER BY administrative_area_level_1 ASC;

-- name: RemoveFavorite :exec
DELETE FROM favorites
WHERE user_id = $1 AND google_id = $2;

-- name: GetFavorite :one
SELECT * FROM favorites
WHERE user_id = $1 AND google_id = $2;

