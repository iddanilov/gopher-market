-- +goose Up
CREATE TABLE orders
(
    order_number VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id      VARCHAR(64)             NOT NULL,
    status       VARCHAR(64)             NOT NULL,
    accrual      VARCHAR(64)             NOT NULL,
    uploaded_at  TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX tags_per_order_number_uploaded_at
    ON orders (order_number, uploaded_at);

-- +goose Down
DROP INDEX tags_per_order_number_uploaded_at;
DROP TABLE orders;