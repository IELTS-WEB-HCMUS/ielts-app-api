CREATE TABLE type (
    id SERIAL PRIMARY KEY,                         -- ID of the type
    title TEXT NOT NULL                           -- Title of the type
);

-- Table: quiz
CREATE TABLE quiz (
    id SERIAL PRIMARY KEY,                         -- ID of a quiz (Primary Key)
    status TEXT NOT NULL,                          -- Status of the quiz ('published' or 'unpublished')
    user_created UUID NOT NULL REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the creator
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Quiz creation timestamp
    user_updated UUID REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the editor
    date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last updated timestamp
    type INTEGER NOT NULL REFERENCES type(id) ON DELETE SET NULL, -- Quiz type (listening/reading)
    content TEXT,                                  -- Content for reading quizzes
    title TEXT NOT NULL,                           -- Title of the quiz
    time INTEGER NOT NULL,                         -- Time limit for the quiz (in minutes)
    listening_file TEXT,                           -- Link to the listening file for listening quizzes
    level INTEGER NOT NULL,                        -- Band level of the quiz
    vote_count INTEGER NOT NULL DEFAULT 0,         -- Number of upvotes
    description TEXT NOT NULL,                     -- Description of the quiz
    thumbnail TEXT,                                -- Thumbnail image link
    mode INTEGER NOT NULL,                         -- Quiz mode (0 for single, 1 for full test)
    is_public BOOLEAN NOT NULL DEFAULT TRUE        -- Whether the quiz is public
);

-- Table: part
CREATE TABLE part (
    id SERIAL PRIMARY KEY,                         -- ID of a part (Primary Key)
    user_created UUID NOT NULL REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the creator
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Part creation timestamp
    user_updated UUID REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the editor
    date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last updated timestamp
    title TEXT NOT NULL,                           -- Title of the part
    content TEXT NOT NULL,                         -- Content of the part
    description TEXT NOT NULL,                     -- Description of the part
    question_count INTEGER NOT NULL DEFAULT 0,     -- Number of questions in the part
    type INTEGER NOT NULL REFERENCES type(id) ON DELETE SET NULL, -- Part type (listening/reading)
    level INTEGER,                                 -- Difficulty level (1, 2, 3 for reading)
    quiz_id INTEGER NOT NULL REFERENCES quiz(id) ON DELETE CASCADE -- Quiz ID the part belongs to
);

-- Table: quiz_part
CREATE TABLE quiz_part (
    quiz_id INTEGER NOT NULL REFERENCES quiz(id) ON DELETE CASCADE, -- Quiz ID (Foreign Key)
    part_id INTEGER NOT NULL REFERENCES part(id) ON DELETE CASCADE, -- Part ID (Foreign Key)
    sort INTEGER NOT NULL,                          -- Order of the part in the quiz
    PRIMARY KEY (quiz_id, part_id)                  -- Composite Primary Key
);


-- Table: question
CREATE TABLE question (
    id SERIAL PRIMARY KEY,                         -- ID of a question (Primary Key)
    user_created UUID NOT NULL REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the creator
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Question creation timestamp
    user_updated UUID REFERENCES public.users(id) ON DELETE SET NULL, -- ID of the editor
    date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Last updated timestamp
    content TEXT NOT NULL,                         -- Content of the question
    type TEXT NOT NULL,                            -- Question type (e.g., SINGLE-SELECTION, MULTIPLE, FILL-IN-THE-BLANK)
    single_choice_radio JSON,                      -- Answers for SINGLE-RADIO type questions
    selection JSON,                                -- Answers for SINGLE-SELECTION type questions
    multiple_choice JSON,                          -- Answers for MULTIPLE type questions
    gap_fill_in_blank JSON,                        -- Answers for FILL-IN-THE-BLANK type questions
    selection_option JSON,                         -- Options for SINGLE-SELECTION type questions
    "order" INTEGER ,                      -- Order of the question in a part
    explain JSON NOT NULL,                         -- Explanations for answers
    question_type TEXT NOT NULL,                   -- Question type (e.g., IELTS-related type)
    part_id INTEGER NOT NULL REFERENCES part(id) ON DELETE CASCADE -- Part ID the question belongs to
);


-- Table: quiz_tag_search
CREATE TABLE quiz_tag_search (
    quiz_id INTEGER NOT NULL REFERENCES quiz(id) ON DELETE CASCADE, -- Quiz ID
    tag_search_id INTEGER NOT NULL REFERENCES tag_search(id) ON DELETE CASCADE, -- Tag ID
    PRIMARY KEY (quiz_id, tag_search_id)         -- Composite Primary Key
);


--DONE--