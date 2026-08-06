CREATE TABLE IF NOT EXISTS post_platform_status(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    connected_account_id UUID NOT NULL REFERENCES connected_accounts(id) ON DELETE CASCADE
    platform_code  TEXT NOT NULL,
    external_post_id TEXT,
    status post_status NOT NULL DEFAULT 'processing',
    failure_reason TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
