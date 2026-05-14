CREATE TABLE password_reset_attempts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    attempt_type TEXT NOT NULL CHECK (attempt_type IN ('reset')),
    token_expires_at TIMESTAMPTZ NOT NULL,
    token_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_password_reset_attempts_user
        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);