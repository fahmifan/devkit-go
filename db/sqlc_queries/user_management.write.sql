-- name: SaveUser :one
INSERT INTO users (id, "name", email, "password", "role", active, created_at, updated_at) 
VALUES (@id, @name, @email, @password, @role, @active, @created_at, @updated_at)
ON CONFLICT (id) DO UPDATE SET 
    "name" = EXCLUDED."name",
    email = EXCLUDED.email,
    "password" = EXCLUDED."password",
    "role" = EXCLUDED."role",
    active = EXCLUDED.active,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at
RETURNING id;