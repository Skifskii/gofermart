-- Создание таблицы заказов
CREATE TABLE orders (
    number VARCHAR(20) PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    accrual NUMERIC(10, 2),
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id)
)