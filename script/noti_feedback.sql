-- Table: notification
CREATE TABLE notification (
    id SERIAL PRIMARY KEY,                         -- ID of the notification (Primary Key)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp of notification creation
    title TEXT NOT NULL,                           -- Title of the notification
    content TEXT NOT NULL,                         -- Content of the notification
    created_by UUID NOT NULL REFERENCES public.users(id) ON DELETE SET NULL -- ID of the admin who created the notification
);

-- Table: notification_tracking
CREATE TABLE notification_tracking (
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE, -- User ID (Foreign Key to users table)
    notification_id INTEGER NOT NULL REFERENCES notification(id) ON DELETE CASCADE, -- Notification ID (Foreign Key to notification table)
    status BOOLEAN NOT NULL,                      -- Read status of the notification
    PRIMARY KEY (user_id, notification_id)         -- Composite Primary Key
);

-- Table: feedback
CREATE TABLE feedback (
    id SERIAL PRIMARY KEY,                         -- ID of the feedback (Primary Key)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Timestamp when feedback was created
    content TEXT NOT NULL,                         -- Content of the feedback
    quiz_id INTEGER REFERENCES quiz(id) ON DELETE SET NULL, -- Quiz ID the feedback refers to (nullable)
    email TEXT,                                    -- Email of the sender (nullable)
    attaches TEXT                                  -- Links to attached images or files (nullable)
);
