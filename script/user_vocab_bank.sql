-- Table: user_vocab_category
CREATE TABLE user_vocab_category (
    id SERIAL PRIMARY KEY,                         -- ID of a vocabulary category (Primary Key)
    name TEXT NOT NULL,                            -- Name of the vocabulary category
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE, -- User ID (Foreign Key to users table)
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Timestamp of the last update
);

-- Table: user_vocab_bank
CREATE TABLE user_vocab_bank (
    id SERIAL PRIMARY KEY,                         -- ID of a vocabulary word (Primary Key)
    value TEXT NOT NULL,                           -- The vocabulary word
    word_class TEXT NOT NULL,                      -- Word class (e.g., noun, verb)
    meaning TEXT NOT NULL,                         -- Meaning of the word
    ipa TEXT NOT NULL,                             -- IPA transcription of the word
    example TEXT,                                  -- Example usage of the word
    note TEXT,                                     -- User's notes
    category INTEGER NOT NULL REFERENCES user_vocab_category(id) ON DELETE CASCADE, -- Vocabulary category ID
    status TEXT NOT NULL,                          -- Status of the word (e.g., learned, not learned)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Timestamp of when the word was added to the bank
);
