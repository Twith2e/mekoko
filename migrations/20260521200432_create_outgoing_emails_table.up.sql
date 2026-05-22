CREATE TABLE outgoing_emails (
    id BIGSERIAL PRIMARY KEY,
    public_id TEXT NOT NULL UNIQUE,
    recipient BIGINT NOT NULL,
    message_id TEXT,
    subject TEXT NOT NULL,
    reason_for_failure TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    email_struct jsonb NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'successful', 'failed')),
    last_retry_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT fk_outgoing_emails_recipient
        FOREIGN KEY(recipient) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX ON outgoing_emails(status, retry_count);