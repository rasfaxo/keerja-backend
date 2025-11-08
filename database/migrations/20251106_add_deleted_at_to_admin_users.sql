-- Migration: Add deleted_at column to admin_users table

ALTER TABLE admin_users
ADD COLUMN deleted_at TIMESTAMP;
