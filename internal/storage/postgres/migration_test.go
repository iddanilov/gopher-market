package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestMigrationsPostgres_CreateBalanceTable(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MigrationsPostgres{
				db: tt.fields.db,
			}
			if err := m.CreateBalanceTable(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CreateBalanceTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMigrationsPostgres_CreateOrdersTable(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MigrationsPostgres{
				db: tt.fields.db,
			}
			if err := m.CreateOrdersTable(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrdersTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMigrationsPostgres_CreateUserTable(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MigrationsPostgres{
				db: tt.fields.db,
			}
			if err := m.CreateUserTable(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMigrationsPostgres_CreateWithdrawalsTable(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MigrationsPostgres{
				db: tt.fields.db,
			}
			if err := m.CreateWithdrawalsTable(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CreateWithdrawalsTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
