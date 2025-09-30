-- name: GetOriginalUrlByCode :one
SELECT * FROM links WHERE short_code = ? AND (expires_at > ? OR expires_at IS NULL);

-- name: CreateShortLink :exec
INSERT INTO links (
    id,tenant_id,user_id,short_code,original_url,created_at,expires_at
) VALUES (
    ?,?,?,?,?,?,?
)