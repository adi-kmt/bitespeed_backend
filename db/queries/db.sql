-- name: GetContactInfoByEmailORPhone :many
SELECT id, email, phone_number, linked_id, created_at
FROM contact
WHERE
 (email = $1 OR $1 = '') OR (phone_number = $2 OR $2 = '');

-- name: InsertContactInfo :one
INSERT INTO contact
    (email, phone_number, linked_id, link_precedence)
VALUES
    ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING id;

-- name: UpdateContactToSecondary :exec
UPDATE contact
SET linked_id = $1, link_precedence = 'secondary'
WHERE
    id = $2;