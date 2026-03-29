-- +goose Up
INSERT INTO users (id, email, password_hash, role, display_name, is_private, created_at)
VALUES (
    gen_random_uuid(),
    'curator_general@gmail.com',
    '$2a$10$KoaV52VqHrTSjbiDWtPqL.uFMMeSpHqXrNVhZMXIfIl/uSbtbSLua',
    'curator',
    'Главный Админ',
    false,
    NOW()
);

-- +goose Down
DELETE FROM users WHERE email = 'curator_general@gmail.com';
