-- name: CreateFavorite :one
INSERT INTO favorites (
    user_id,
    place_id
) VALUES (
    $1, 
    $2
) RETURNING *;

-- name: ListFavorite :many
SELECT * FROM favorites
WHERE user_id = $1
ORDER BY created_at ASC
LIMIT $2
OFFSET $3;

-- name: RemoveFavorite :exec
DELETE FROM favorites
WHERE user_id = $1 AND place_id = $2;

-- name: GetFavorite :one
SELECT * FROM favorites
WHERE user_id = $1 AND place_id = $2;

