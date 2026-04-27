-- Init script to create separate databases for each service
-- This runs automatically when PostgreSQL container starts

CREATE DATABASE novelhive_users;
CREATE DATABASE novelhive_novels;
CREATE DATABASE novelhive_comments;
CREATE DATABASE novelhive_library;

-- Enable ltree extension for comments (threaded replies)
\c novelhive_comments;
CREATE EXTENSION IF NOT EXISTS ltree;

-- Enable uuid generation
\c novelhive_users;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_novels;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_comments;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c novelhive_library;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
