-- Table: payment
CREATE TABLE payment (
    id SERIAL PRIMARY KEY,                         -- ID of the transaction (Primary Key)
    amount INTEGER NOT NULL,                       -- Amount of the transaction
    status TEXT NOT NULL,                          -- Status of the transaction (e.g., success or failure)
    type TEXT NOT NULL,                            -- Type of transaction (e.g., premium upgrade, vocabulary credits)
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE, -- User ID (Foreign Key to users table)
    transaction_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Timestamp of the transaction
);
