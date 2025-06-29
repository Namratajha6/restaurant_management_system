
CREATE TYPE role_type AS ENUM ('admin', 'sub_admin', 'user');


CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     name TEXT NOT NULL,
                                     email TEXT NOT NULL UNIQUE,
                                     password TEXT NOT NULL,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                     archived_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS user_role (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         user_id UUID REFERENCES users(id) NOT NULL,
                                         role_type role_type NOT NULL,
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                         archived_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS user_address (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            user_id UUID REFERENCES users(id) NOT NULL,
                                            address TEXT NOT NULL,
                                            latitude DOUBLE PRECISION,
                                            longitude DOUBLE PRECISION,
                                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                            archived_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS user_session (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            user_id UUID REFERENCES users(id) NOT NULL,
                                            refresh_token TEXT NOT NULL,
                                            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                            archived_at TIMESTAMP WITH TIME ZONE
);


CREATE TABLE IF NOT EXISTS restaurant (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                          name TEXT NOT NULL,
                                          address TEXT NOT NULL,
                                          latitude DOUBLE PRECISION,
                                          longitude DOUBLE PRECISION,
                                          created_by UUID REFERENCES users(id) NOT NULL,
                                          rating NUMERIC(2,1) NOT NULL CHECK (rating >= 0.0 AND rating <= 5.0),
                                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                          archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS dishes (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      restaurant_id UUID REFERENCES restaurant(id) NOT NULL,
                                      name TEXT NOT NULL,
                                      description TEXT,
                                      price NUMERIC(10,2),
                                      created_by UUID REFERENCES users(id) NOT NULL,
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                      archived_at TIMESTAMP WITH TIME ZONE
);

