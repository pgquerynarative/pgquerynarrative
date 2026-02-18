-- ⚠️  SECURITY WARNING: This file uses default passwords for DEVELOPMENT ONLY.
-- 
-- In production:
--   1. Use environment variables or secrets management
--   2. Set passwords via: psql -v readonly_pwd='...' -v app_pwd='...' -f create-roles.sql
--   3. Or use: CREATE ROLE ... LOGIN PASSWORD :'password_from_env';
--   4. Never commit production passwords to version control
--
-- Create roles if they don't exist
-- Default passwords are: pgquerynarrative_readonly, pgquerynarrative_app
-- CHANGE THESE IN PRODUCTION!
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'pgquerynarrative_readonly') THEN
        CREATE ROLE pgquerynarrative_readonly LOGIN PASSWORD 'pgquerynarrative_readonly';
    END IF;
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'pgquerynarrative_app') THEN
        CREATE ROLE pgquerynarrative_app LOGIN PASSWORD 'pgquerynarrative_app';
    END IF;
END
$$;

-- Grant permissions to readonly role
GRANT USAGE ON SCHEMA demo TO pgquerynarrative_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA demo TO pgquerynarrative_readonly;
ALTER DEFAULT PRIVILEGES IN SCHEMA demo GRANT SELECT ON TABLES TO pgquerynarrative_readonly;

-- Grant permissions to app role
GRANT USAGE ON SCHEMA app TO pgquerynarrative_app;
GRANT ALL ON ALL TABLES IN SCHEMA app TO pgquerynarrative_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA app GRANT ALL ON TABLES TO pgquerynarrative_app;
