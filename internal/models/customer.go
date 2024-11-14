package models

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DbCustomer struct {
	Id       string    `json:"id"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type Customer struct {
	Id          string    `json:"id"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Created     time.Time `json:"created"`
	LastOrdered time.Time `json:"last_ordered" db:"last_ordered"`
	TotalSpent  int       `json:"total_spent" db:"total_spent"`
}

func (dbp *DbCustomer) ConvertToCustomer(lastOrdered time.Time, totalSpent int) *Customer {
	return &Customer{
		Id:          dbp.Id,
		Fullname:    dbp.Fullname,
		Email:       dbp.Email,
		Address:     dbp.Address,
		Phone:       dbp.Phone,
		Created:     dbp.Created,
		LastOrdered: lastOrdered,
		TotalSpent:  totalSpent,
	}
}

func CustomerExists(email string) (bool, error) {
	var exists bool

	existQuery := `SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)`

	if err := db.Get(&exists, existQuery, email); err != nil {
		return false, err
	}

	return exists, nil
}

func CreateCustomer(fullname string, email string, address string, phone string) (*DbCustomer, error) {
	statement := "INSERT INTO customers (id, fullname, email, address, phone) VALUES ($1, $2, $3, $4, $5)"

	tx := db.MustBegin()

	c := &DbCustomer{Id: uuid.NewV4().String(), Fullname: fullname, Email: email, Address: address, Phone: phone}

	if _, err := tx.Exec(statement, c.Id, c.Fullname, c.Email, c.Address, c.Phone); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error rolling back transaction: %v", rollbackErr)
		}
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return c, nil
}

func GetCustomers() ([]Customer, error) {
	var customers []Customer = make([]Customer, 0)

	statement := `SELECT
									c.id as id,
									c.fullname as fullname,
									c.email as email,
									c.address as address,
									c.phone as phone,
									c.created as created,
									MAX(o.created) AS last_ordered,
									COALESCE(SUM(p.quantity * pr.price), 0) AS total_spent
								FROM
										customers c
								LEFT JOIN
										orders o ON c.id = o.customer
								LEFT JOIN
										purchases p ON o.id = p.orderid
								LEFT JOIN
										products pr ON p.productid = pr.id
								GROUP BY
										c.id, c.fullname, c.email
								ORDER BY
										c.fullname ASC`

	err := db.Select(&customers, statement)

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func GetCustomer(id string) (*Customer, error) {
	var customer Customer

	statement := `SELECT
									c.id as id,
									c.fullname as fullname,
									c.email as email,
									c.address as address,
									c.phone as phone,
									c.created as created,
									MAX(o.created) AS last_ordered,
									COALESCE(SUM(p.quantity * pr.price), 0) AS total_spent
								FROM
										customers c
								LEFT JOIN
										orders o ON c.id = o.customer
								LEFT JOIN
										purchases p ON o.id = p.orderid
								LEFT JOIN
										products pr ON p.productid = pr.id
								WHERE c.id = $1
								GROUP BY
										c.id, c.fullname, c.email
								ORDER BY
										c.fullname ASC`

	err := db.Get(&customer, statement, id)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func GetDbCustomer(id string) (*DbCustomer, error) {
	var customer DbCustomer

	statement := "SELECT * FROM customers WHERE id = $1"

	err := db.Get(&customer, statement, id)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func GetCustomerByEmail(email string) (*DbCustomer, error) {
	var customer DbCustomer

	statement := "SELECT * FROM customers WHERE email = $1"

	err := db.Get(&customer, statement, email)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (customer *DbCustomer) Update(fullname string, email string, address string, phone string) error {
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
			return fmt.Errorf("error rolling back transaction: %v", rollbackErr)
		}
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (customer *Customer) Delete() ([]Customer, error) {
	statement := "DELETE FROM customers WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, customer.Id); err != nil {
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

	updatedCustomers, err := GetCustomers()
	if err != nil {
		return nil, err
	}

	return updatedCustomers, nil
}

type Spender struct {
	Id       string `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Spent    int    `json:"spent"`
}

type CustomersStats struct {
	TotalCustomers int `json:"total_customers"`

	CustomersData []Dataset `json:"customers_data"`

	TopSpenders []Spender `json:"top_spenders"`
}

func GetAllCustomersAmount() (int, error) {
	statement := "SELECT COUNT(*) FROM customers"

	var totalCustomers int

	err := db.Get(&totalCustomers, statement)
	if err != nil {
		return 0, err
	}

	return totalCustomers, nil
}

func GetCustomersData(timeframe Timeframe) ([]Dataset, error) {
	var newResults []Count = make([]Count, 0)
	var oldResults []Count = make([]Count, 0)
	var havingStm string

	horizontal, err := GetHorizonalDataAndQueryByTimeframe("orders.created", timeframe, &havingStm)
	if err != nil {
		return nil, err
	}

	havingStm = strings.Replace(havingStm, "WHERE", "HAVING COUNT(*) = 1 AND ", 1)

	newCustomersStatement := `SELECT
														DATE(orders.created) AS date,
														COUNT(DISTINCT customer) AS count
														FROM
																orders
														GROUP BY
																orders.created ` + havingStm + ` ORDER BY orders.created;`
	err = db.Select(&newResults, newCustomersStatement)
	if err != nil {
		return nil, err
	}

	newVertical, err := ComputeVertical(newResults, horizontal, timeframe)
	if err != nil {
		return nil, err
	}

	havingStm = strings.Replace(havingStm, "=", ">", 1)

	oldCustomersStatement := `SELECT
														DATE(orders.created) AS date,
														COUNT(DISTINCT customer) AS count
														FROM
																orders
														GROUP BY
																orders.created ` + havingStm + ` ORDER BY orders.created;`

	err = db.Select(&oldResults, oldCustomersStatement)
	if err != nil {
		return nil, err
	}

	oldVertical, err := ComputeVertical(oldResults, horizontal, timeframe)
	if err != nil {
		return nil, err
	}

	return []Dataset{{Topic: "New Customers", Horizontal: horizontal, Vertical: newVertical, Color: 0x1CE2D4}, {Topic: "Old Customers", Horizontal: horizontal, Vertical: oldVertical, Color: 0xCFE410}}, nil
}

func GetTopSpenders() ([]Spender, error) {
	var spenders []Spender = make([]Spender, 0)

	statement := `SELECT
								c.id AS id,
								c.fullname AS fullname,
								c.email AS email,
								COALESCE(SUM(p.quantity * pr.price), 0) AS spent
								FROM
										customers c
								LEFT JOIN
										orders o ON c.id = o.customer
								LEFT JOIN
										purchases p ON o.id = p.orderid
								LEFT JOIN
										products pr ON p.productid = pr.id
								GROUP BY
										c.id, c.fullname, c.email
								ORDER BY
										spent DESC
								LIMIT 10`

	err := db.Select(&spenders, statement)
	if err != nil {
		return nil, err
	}

	return spenders, nil
}
