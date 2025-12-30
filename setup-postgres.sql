-- PostgreSQL Setup Script for Ritel-App
-- Run this with: psql -U postgres -f setup-postgres.sql

-- Drop existing database and user if exists (clean start)
DROP DATABASE IF EXISTS ritel_db;
DROP USER IF EXISTS ritel;

-- Create database
CREATE DATABASE ritel_db;

-- Create user
CREATE USER ritel WITH PASSWORD 'ritel123';

-- Grant privileges on database
GRANT ALL PRIVILEGES ON DATABASE ritel_db TO ritel;

-- Connect to the database
\c ritel_db

-- Grant schema permissions
GRANT ALL ON SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ritel;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ritel;

-- Set default privileges for future tables and sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ritel;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ritel;

-- Display success message
\echo '✅ Database setup completed!'
\echo 'Database: ritel_db'
\echo 'User: ritel'
\echo 'Password: ritel123'
