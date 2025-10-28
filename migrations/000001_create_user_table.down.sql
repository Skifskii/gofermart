-- Откат индекса по логину
DROP INDEX IF EXISTS idx_users_login;

-- Откат создания таблицы пользователей
DROP TABLE IF EXISTS users;