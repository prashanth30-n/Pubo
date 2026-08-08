ALTER TABLE connected_accounts
    ADD COLUMN IF NOT EXISTS did TEXT,
    ADD COLUMN IF NOT EXISTS pds_url TEXT,
    ADD COLUMN IF NOT EXISTS password TEXT;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'connected_accounts'
          AND column_name = 'refersh_token_encrypted'
    ) THEN
        ALTER TABLE connected_accounts
            RENAME COLUMN refersh_token_encrypted TO refresh_token_encrypted;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS post_platform_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    connected_account_id UUID NOT NULL REFERENCES connected_accounts(id) ON DELETE CASCADE,
    platform_code TEXT NOT NULL,
    external_post_id TEXT,
    status post_status NOT NULL DEFAULT 'processing',
    failure_reason TEXT,
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_connected_accounts_did
    ON connected_accounts(did)
    WHERE did IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_post_platform_status_post_id
    ON post_platform_status(post_id);
