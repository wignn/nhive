
CREATE DATABASE novelhive_users;
CREATE DATABASE novelhive_novels;
CREATE DATABASE novelhive_comments;
CREATE DATABASE novelhive_library;

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
