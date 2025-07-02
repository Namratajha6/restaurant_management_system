ALTER TABLE user_address
ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS unique_primary_address_per_user
    ON user_address(user_id)
    WHERE is_primary = TRUE AND archived_at IS NULL;
