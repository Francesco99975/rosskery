package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Purchase struct {
	Id string `json:"id"`
	Product Product `json:"product"`
	Quantity int `json:"quantity"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func CreatePurchase(productId string, quantity int) (*Purchase, error) {
	statement := "INSERT INTO purchases (id, productid, quantity) VALUES ($1, $2, $3)"

	product, err := GetProduct(productId)
	if err != nil {
		return nil, err
	}

	newPurchase := &Purchase{Id: uuid.NewV4().String(), Product: *product, Quantity: quantity}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, newPurchase.Id, newPurchase.Product.Id, newPurchase.Quantity); err != nil {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	return newPurchase, nil
}

func GetPurchases() ([]Purchase, error) {
	var purchases []Purchase

	statement := "SELECT * FROM purchases"

	err := db.Select(&purchases, statement)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func GetOrderPurchases(orderId string) ([]Purchase, error) {
	var purchases []Purchase

	statement := "SELECT * FROM purchases WHERE orderid = $1"

	err := db.Select(&purchases, statement, orderId)

	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func GetProductPurchases(productId string) ([]Purchase, error) {
	var purchases []Purchase

	statement := "SELECT * FROM purchases WHERE productid = $1"

	err := db.Select(&purchases, statement, productId)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func GetPurchase(id string) (*Purchase, error) {
	var purchase Purchase

	statement := "SELECT * FROM purchases WHERE id = $1"

	err := db.Select(purchase, statement, id)

	if err != nil {
		return nil, err
	}

	return &purchase, nil
}

func (p *Purchase) Delete() (*Purchase, error) {
	statement := "DELETE FROM purchases WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, p.Id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	return p, nil
}

type Order struct {
	Id string `json:"id"`
	Customer Customer `json:"customer"`
	Purchases []Purchase `json:"purchases"`
	Pickuptime time.Time `json:"pickuptime"`
	Fulfilled bool `json:"fulfilled"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func CreateOrder(customerId string, pickuptime time.Time, items []PurchasedItem) (*Order, error) {
	statement := "INSERT INTO orders (id, customer, pickuptime, fulfilled) VALUES ($1, $2, $3, $4)"

	customer, err := GetCustomer(customerId)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{Id: uuid.NewV4().String(), Customer: *customer, Pickuptime: pickuptime, Purchases: make([]Purchase, len(items)), Fulfilled: false}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, newOrder.Id, newOrder.Customer.Id, newOrder.Pickuptime, newOrder.Fulfilled); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	for i, item := range items {
		purchase, err := CreatePurchase(item.Product.Id, item.Quantity)
		if err != nil {
			return nil, err
		}

		newOrder.Purchases[i] = *purchase
	}

	return newOrder, nil
}

func GetOrders() ([]Order, error) {
	var orders []Order

	statement := "SELECT * FROM orders"

	err := db.Select(&orders, statement)

	for _, order := range orders {
		order.Purchases, err = GetOrderPurchases(order.Id)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetCustomerOrders(customerId string) ([]Order, error) {
	var orders []Order

	statement := "SELECT * FROM orders WHERE customer = $1"

	err := db.Select(&orders, statement, customerId)

	for _, order := range orders {
		order.Purchases, err = GetOrderPurchases(order.Id)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetFulfilledOrders() ([]Order, error) {
	var orders []Order

	statement := "SELECT * FROM orders WHERE fulfilled = $1"

	err := db.Select(&orders, statement, true)

	for _, order := range orders {
		order.Purchases, err = GetOrderPurchases(order.Id)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetUnfulfilledOrders() ([]Order, error) {
	var orders []Order

	statement := "SELECT * FROM orders WHERE fulfilled = $1"

	err := db.Select(&orders, statement, false)

	for _, order := range orders {
		order.Purchases, err = GetOrderPurchases(order.Id)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrder(id string) (*Order, error) {
	var order Order

	statement := "SELECT * FROM orders WHERE id = $1"

	err := db.Select(order, statement, id)
	if err != nil {
		return nil, err
	}

	order.Purchases, err = GetOrderPurchases(order.Id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *Order) Fulfill() (*Order, error) {
	statement := "UPDATE orders SET fulfilled = $1 WHERE id = $2"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, true, o.Id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	return o, nil
}

func (o *Order) CalculateTotal() int {
	total := 0
	for _, purchase := range o.Purchases {
		total += purchase.Product.Price * purchase.Quantity
	}
	return total
}

func (o *Order) Delete() error {
	statement := "DELETE FROM orders WHERE id = $1"

	for _, purchase := range o.Purchases {
		_, err := purchase.Delete()
		if err != nil {
			return err
		}
	}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, o.Id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}
