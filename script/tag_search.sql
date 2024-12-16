CREATE TABLE tag_search (
    id SERIAL PRIMARY KEY,                         -- ID of a tag
    title TEXT NOT NULL,                           -- Title of the tag
    is_shown BOOLEAN NOT NULL                      -- Whether the tag is visible
);

CREATE TABLE tag_position (
    id SERIAL PRIMARY KEY,                         -- ID of a tag position
    title TEXT NOT NULL,                           -- Title of the position
    position TEXT NOT NULL                         -- Display position of the tag on the interface
);
