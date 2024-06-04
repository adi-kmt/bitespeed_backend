CREATE TYPE link_precedence_enum AS ENUM ('primary', 'secondary');

CREATE TABLE IF NOT EXISTS contact (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100),
    phone_number VARCHAR(20),
    linked_id INTEGER,
    link_precedence link_precedence_enum NOT NULL DEFAULT 'primary',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP,

    UNIQUE (email, phone_number)
);

CREATE INDEX contact_email_phone ON contact (email, phone_number);