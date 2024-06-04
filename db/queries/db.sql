-- name: GetContactInfoByEmailORPhone :many
SELECT id, email, phone_number, linked_id, created_at
FROM contact
WHERE
    ($1 IS NOT NULL AND email = $1)
    OR ($2 IS NOT NULL AND phone_number = $2);

-- name: InsertContactInfo :exec
INSERT INTO contact
    (email, phone_number, linked_id, link_precedence)
VALUES
    ($1, $2, $3, $4);

-- name: UpdateContactToSecondary :exec
UPDATE contact
SET linked_id = $1, link_precedence = $2
WHERE
    ($3 IS NOT NULL AND email = $3)
    OR ($4 IS NOT NULL AND phone_number = $4);