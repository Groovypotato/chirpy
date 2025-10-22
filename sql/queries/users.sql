-- name: CreateUser :one
INSERT INTO users(id,created_at,updated_at,email,hashed_password)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: InsertTestUser :one
INSERT INTO users (id,
created_at,
updated_at,
email,
hashed_password
) 
VALUES ("e364ea76-dd6d-4d05-93f0-b321b876ff68"
,"2025-10-22 11:40:20.883407",
"2025-10-22 11:40:20.883407",
"gleasoncr@gmail.com",
"$argon2id$v=19$m=131072,t=4,p=16$3zfzL14P9fWukKMNUiSPLA$Pyhb9RzdwqnDMmxjNGO3ZmK5i4hiYX6icQQQW3wqNNg")
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUsers :many
SELECT *
FROM users;

-- name: DeleteUsers :exec
DELETE FROM users;