package postgres

import (
	"fmt"
	"github.com/gopher-market/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (r *OrdersPostgres) LoadOrder(userID int, orderID string) error {
	query := fmt.Sprintf("INSERT INTO orders(order_number, user_id, status) values ($1, $2, $3)")
	log.Println(query)

	_, err := r.db.Exec(query, orderID, userID, "NEW")
	if err != nil {
		log.Println("Can't Update Order: ", err)
	}
	return nil
}

func (r *OrdersPostgres) GetOrders(userID int) ([]models.Order, error) {
	var result []models.Order
	fmt.Println(userID)
	query := fmt.Sprintf(`
		SELECT order_number, user_id, status, accrual, uploaded_at 
		FROM orders
    	WHERE user_id = $1`)
	if err := r.db.Select(&result, query, strconv.Itoa(userID)); err != nil {
		return nil, err
	}
	fmt.Println(result)
	return result, nil

}
