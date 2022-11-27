package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type MigrationsPostgres struct {
	db *sqlx.DB
}

func NewMigrationsPostgres(db *sqlx.DB) *MigrationsPostgres {
	return &MigrationsPostgres{db: db}
}

func (m *MigrationsPostgres) CreateUserTable(ctx context.Context) error {
	row, err := m.db.Query(`
select * from users;`,
	)
	if err != nil {
		if err.Error() == `pq: relation "users" does not exist` {
			_, err = m.db.ExecContext(ctx, `
CREATE TABLE users
(
    id            serial       not null unique,
    login         varchar(255) not null unique,
    password_hash varchar(255) not null
);`,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if row != nil {
		if row.Err() != nil {
			return err
		}
		defer row.Close()
	}
	log.Println("DB Create")

	return nil
}

func (m *MigrationsPostgres) CreateOrdersTable(ctx context.Context) error {
	row, err := m.db.Query(`
select * from orders;`,
	)
	if err != nil {
		if err.Error() == `pq: relation "orders" does not exist` {
			_, err = m.db.ExecContext(ctx, `
CREATE TABLE orders
(
    order_number VARCHAR(64) PRIMARY KEY NOT NULL,
    user_id      VARCHAR(64)             NOT NULL,
    status       VARCHAR(64)             NOT NULL,
    accrual      VARCHAR(64)             NOT NULL,
    uploaded_at  TIMESTAMP DEFAULT now()
);`,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if row != nil {
		if row.Err() != nil {
			return err
		}
		defer row.Close()
	}
	log.Println("DB Create")

	return nil
}

func (m *MigrationsPostgres) CreateBalanceTable(ctx context.Context) error {
	row, err := m.db.Query(`
select * from balance;`,
	)
	if err != nil {
		if err.Error() == `pq: relation "balance" does not exist` {
			_, err = m.db.ExecContext(ctx, `
CREATE TABLE balance
(
    user_id      VARCHAR(64) PRIMARY KEY NOT NULL,
    user_current INT                     NOT NULL,
    withdrawn    INT                     NOT NULL
);`,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if row != nil {
		if row.Err() != nil {
			return err
		}
		defer row.Close()
	}
	log.Println("DB Create")

	return nil
}

func (m *MigrationsPostgres) CreateWithdrawalsTable(ctx context.Context) error {
	row, err := m.db.Query(`
select * from withdrawals;`,
	)
	if err != nil {
		if err.Error() == `pq: relation "withdrawals" does not exist` {
			_, err = m.db.ExecContext(ctx, `
CREATE TABLE withdrawals
(
    user_id INT PRIMARY KEY NOT NULL,
    order_number   VARCHAR(64)     NOT NULL,
    sum     VARCHAR(64)     NOT NULL,
    processed_at TIMESTAMP DEFAULT now()
);`,
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if row != nil {
		if row.Err() != nil {
			return err
		}
		defer row.Close()
	}
	log.Println("DB Create")

	return nil
}
