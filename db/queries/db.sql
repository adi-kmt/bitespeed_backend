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
SET linked_id = $1, link_precedence = 'secondary'
WHERE
    NULLIF($2, '') IS NULL OR email = NULLIF($2, '')
    OR
    NULLIF($3, '') IS NULL OR phone_number = NULLIF($3, '')
RETURNING id;