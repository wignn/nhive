
CREATE DATABASE novelhive_users;
CREATE DATABASE novelhive_novels;
CREATE DATABASE novelhive_comments;
CREATE DATABASE novelhive_library;
CREATE DATABASE novelhive_notifications;

\c novelhive_comments;
CREATE EXTENSION IF NOT EXISTS ltree;

\c novelhive_users;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_novels;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_comments;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_library;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_notifications;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    payload JSONB,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read) WHERE is_read = FALSE;

CREATE TABLE IF NOT EXISTS fcm_tokens (
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, token)
);
