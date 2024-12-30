CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    name VARCHAR(255) NOT NULL,                    
    public_id TEXT                                 
);

CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- UUID primary key
    first_name VARCHAR(255),                       -- First name of the user
    last_name VARCHAR(255),                        -- Last name of the user
    fullname VARCHAR(255),                         -- Full name of the user
    email VARCHAR(255) UNIQUE NOT NULL,            -- User's email, must be unique
    password VARCHAR(255),                         -- Encrypted password (nullable for social login users)
    avatar TEXT,                                   -- Avatar image ID
    status VARCHAR(50) DEFAULT 'active',           -- Status (e.g., active, inactive)
    role UUID REFERENCES roles(id),                -- Foreign key to roles table (role UUID)
    token VARCHAR(255),                            -- Token for verification or authentication
    provider VARCHAR(255),                         -- Social login provider (e.g., google, facebook)
    email_notifications BOOLEAN DEFAULT TRUE,      -- Email notifications flag    
    is_active BOOLEAN DEFAULT TRUE,                -- Flag for whether the user is active
    date_created TIMESTAMPTZ DEFAULT NOW(),        -- Account creation timestamp
    vocab_usage_count INT                          -- Vocabulary usage count
);

CREATE TABLE public.student_target (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),   -- UUID primary key
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Foreign key to the users table
    target_study_duration INT,                       -- Study duration (in minutes or hours)
    target_reading FLOAT,                            -- Reading score target
    target_listening FLOAT,                          -- Listening score target
    target_speaking FLOAT,                           -- Speaking score target
    target_writing FLOAT,                            -- Writing score target
    next_exam_date TIMESTAMPTZ                       -- Timestamp for the next exam date
);


--DONE--