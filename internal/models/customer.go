package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Customer struct {
	Id string `json:"id"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func CustomerExists(email string) bool {
	statement := "SELECT * FROM customers WHERE email = $1"
	var customer Customer

	err := db.Get(&customer, statement, email)

	return err != nil
}


func CreateCustomer(fullname string, email string, address string, phone string) (*Customer, error) {
	statement := "INSERT INTO customers (id, fullname, email, address, phone) VALUES ($1, $2, $3, $4, $5)"

	c := &Customer{Id: uuid.NewV4().String(), Fullname: fullname, Email: email, Address: address, Phone: phone}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, c.Id, c.Fullname, c.Email, c.Address, c.Phone); err != nil {
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

	return c, nil
}

func GetCustomers() ([]Customer, error) {
	var customers []Customer = make([]Customer, 0)

	statement := "SELECT * FROM customers"

	err := db.Select(&customers, statement)

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func GetCustomer(id string) (*Customer, error) {
	var customer Customer

	statement := "SELECT * FROM customers WHERE id = $1"

	err := db.Get(&customer, statement, id)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func GetCustomerByEmail(email string) (*Customer, error) {
	var customer Customer

	statement := "SELECT * FROM customers WHERE email = $1"

	err := db.Get(&customer, statement, email)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (customer *Customer) Update(fullname string, email string, address string, phone string) error {
	statement := "UPDATE customers SET fullname = $1, email = $2, address = $3, phone = $4 WHERE id = $5"

	customer.Fullname = fullname
	customer.Email = email
	customer.Address = address
	customer.Phone = phone

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, customer.Fullname, customer.Email, customer.Address, customer.Phone, customer.Id); err != nil {
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

func (customer *Customer) Delete() error {
	statement := "DELETE FROM customers WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, customer.Id); err != nil {
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
