-- Создание таблицы пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    balance NUMERIC(15, 2) DEFAULT 0.00,
    withdrawn NUMERIC(15, 2) DEFAULT 0.00
);

-- Индекс для быстрого поиска по логину
CREATE INDEX idx_users_login ON users(login);