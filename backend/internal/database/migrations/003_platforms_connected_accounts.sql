CREATE TABLE IF NOT EXISTS platforms(
    id SMALLSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO platforms(code,display_name)
VALUES
('x','X'),
('linkedin','LinkedIn'),
('bluesky','Bluesky')
ON CONFLICT (code) DO NOTHING;


CREATE TABLE IF NOT EXISTS connected_accounts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES pubo_users(clerk_user_id) on DELETE CASCADE,
    platform_id SMALLINT NOT NULL REFERENCES platforms(id),
    -- platform_account_id TEXT NOT NULL this has been removed through another migration,
    handle TEXT, 
    display_name TEXT,
    avatar_url TEXT,
    access_token_encrypted TEXT,
    refersh_token_encrypted TEXT,
    token_expires_at TIMESTAMP,
    scopes TEXT [],
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_synced_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id,platform_id,platform_account_id)
);

CREATE INDEX IF NOT EXISTS idx_connected_accounts_user_id ON connected_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_connected_accounts_platform_id ON connected_accounts(platform_id);
CREATE INDEX IF NOT EXISTS idx_connected_accounts_user_active ON connected_accounts(user_id,is_active);
