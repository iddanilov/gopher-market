-- +goose Up
CREATE TABLE balance
(
    user_id      VARCHAR(64) PRIMARY KEY NOT NULL,
    user_current INT                     NOT NULL,
    withdrawn    INT                     NOT NULL
);

CREATE TABLE withdrawals
(
    user_id INT PRIMARY KEY NOT NULL,
    order_number   VARCHAR(64)     NOT NULL,
    sum     VARCHAR(64)     NOT NULL,
    processed_at TIMESTAMP DEFAULT now()
);

CREATE UNIQUE INDEX tags_per_user_id
    ON orders (user_id);

-- +goose Down
DROP INDEX tags_per_user_id;
DROP TABLE balance;



type
Withdrawals struct {
	Id          int        `json:"-" db:"user_id"`
	Order       string     `json:"order" db:"order"`
	Sum         string     `json:"sum" db:"sum"`
	ProcessedAt *time.Time `json:"processed_at" db:"processed_at"`
}
