-- name: GetContactInfoByEmailORPhone :many
SELECT id, email, phone_number, linked_id, created_at
FROM contact
WHERE
    NULLIF($1, '') IS NULL OR email = NULLIF($1, '')
    OR
    NULLIF($2, '') IS NULL OR phone_number = NULLIF($2, '');

-- name: InsertContactInfo :one
INSERT INTO contact
    (email, phone_number, linked_id, link_precedence)
VALUES
    ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: UpdateContactToSecondary :one
UPDATE contact
SET linked_id = $1, link_precedence = $2
WHERE
    ($3 IS NOT NULL AND email = $3)
    OR ($4 IS NOT NULL AND phone_number = $4)
RETURNING id;