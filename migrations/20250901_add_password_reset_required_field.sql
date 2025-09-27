-- Add is_password_reset_required column to users table
ALTER TABLE users 
ADD COLUMN is_password_reset_required BOOLEAN NOT NULL DEFAULT FALSE;
