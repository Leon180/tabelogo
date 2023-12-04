-- name: CreateFavorite :one
INSERT INTO favorites (
    user_id,
    place_id
) VALUES (
    $1, 
    $2
) RETURNING *;

-- name: ListFavoritesByCreateTime :many
SELECT
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.place_id = places.place_id
WHERE user_id = $1
ORDER BY favorites.created_at ASC
LIMIT $2
OFFSET $3;

-- name: ListFavoritesByCountry :many
SELECT
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.place_id = places.place_id
WHERE user_id = $1 AND country = $2
ORDER BY favorites.created_at ASC
LIMIT $3
OFFSET $4;

-- name: ListFavoritesByCountrAndRegion :many
SELECT
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
    favorites.created_at,
    favorites.updated_at,
    favorite_id
FROM favorites JOIN places ON favorites.place_id = places.place_id
WHERE user_id = $1 AND country = $2 AND administrative_area_level_1 = $3
ORDER BY favorites.created_at ASC
LIMIT $4
OFFSET $5;

-- name: GetCountryList :many
SELECT DISTINCT country FROM favorites JOIN places ON favorites.place_id = places.place_id
WHERE user_id = $1
ORDER BY country ASC;

-- name: GetRegionList :many
SELECT DISTINCT administrative_area_level_1 FROM favorites JOIN places ON favorites.place_id = places.place_id
WHERE user_id = $1 AND country = $2
ORDER BY administrative_area_level_1 ASC;

-- name: RemoveFavorite :exec
DELETE FROM favorites
WHERE user_id = $1 AND place_id = $2;

-- name: GetFavorite :one
SELECT * FROM favorites
WHERE user_id = $1 AND place_id = $2;

