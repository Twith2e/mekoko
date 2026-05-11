CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    first_name TEXT NOT NULL, -- combination of first name and last name. might take middle name also, yet to decide.
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    date_of_birth DATE,
    phone_number TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);