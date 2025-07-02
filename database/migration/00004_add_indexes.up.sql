-- for users table
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE archived_at IS NULL;

-- Add indexes for restaurant table
CREATE INDEX IF NOT EXISTS idx_restaurant_location ON restaurant(latitude, longitude) WHERE archived_at IS NULL;

-- Add indexes for dishes table
CREATE INDEX IF NOT EXISTS idx_dishes_restaurant_id ON dishes(restaurant_id) WHERE archived_at IS NULL;

