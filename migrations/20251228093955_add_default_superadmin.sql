-- +goose Up
-- +goose StatementBegin
-- Check if a superadmin already exists to avoid duplicates
INSERT INTO users (id, email, password_hash, role, created_at, updated_at)
SELECT 
  '01JH5F8K1H9J6Z3V5QY8K2W1Z9', 
  'superadmin@x.com', 
  '$2a$10$u09/NZ1vo9ByVamR5hMpM.4cgMJiczR5uAUz7D60jyZdwauUOkfCW', 
  'superadmin', 
  NOW(), 
  NOW()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'superadmin@x.com'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE email = 'superadmin@x.com';
-- +goose StatementEnd
