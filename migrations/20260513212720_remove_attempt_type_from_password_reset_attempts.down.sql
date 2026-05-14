ALTER TABLE password_reset_attempts
ADD attempt_type TEXT NOT NULL CHECK (attempt_type IN ('reset'));