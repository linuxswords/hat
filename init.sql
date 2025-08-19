-- Create additional databases for testing
CREATE DATABASE hat_test;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE hat_development TO hat_user;
GRANT ALL PRIVILEGES ON DATABASE hat_test TO hat_user;