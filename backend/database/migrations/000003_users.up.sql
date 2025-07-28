CREATE TYPE TypeUser AS ENUM ('Sender', 'Carrier', 'Receiver');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    address VARCHAR(255) NOT NULL UNIQUE,
    type TypeUser NOT NULL
);

CREATE INDEX idx_users_id ON users(id);