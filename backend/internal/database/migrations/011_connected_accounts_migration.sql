ALTER TABLE connected_accounts ADD COLUMN IF NOT EXISTS did TEXT,
ADD COLUMN IF NOT EXISTS pds_url TEXT,ADD COLUMN IF NOT EXISTS password TEXT;

DO $$
BEGIN
     IF EXISTS(
        SELECT 1 FROM information_schema.columns
        WHERE table_name='connected_accounts'AND column_name='refersh_token_encrypted'
     )THEN ALTER TABLE connected_accounts RENAME COLUMN refersh_token_encrypted TO refresh_token_encrypted;
     END IF;
END $$;


ALTER TABLE connected_accounts DROP CONSTRAINT IF EXISTS connected_accounts_user_id_platform_id_platform_account_id_key;


CREATE UNIQUE INDEX IF NOT EXISTS idx_connected_accounts_platform_unique ON connected_accounts(user_id,platform_id,COALESCE(did,handle));

CREATE INDEX IF NOT EXISTS idx_connected_accounts_did ON connected_accounts(did) WHERE did IS NOT NULL
