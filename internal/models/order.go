package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Purchase struct {
	Id       string    `json:"id"`
	Product  Product   `json:"product"`
	Quantity int       `json:"quantity"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
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
	var purchases []Purchase = make([]Purchase, 0)

	statement := "SELECT * FROM purchases"

	err := db.Select(&purchases, statement)
	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func GetOrderPurchases(orderId string) ([]Purchase, error) {
	var purchases []Purchase = make([]Purchase, 0)

	statement := "SELECT * FROM purchases WHERE orderid = $1"

	err := db.Select(&purchases, statement, orderId)

	if err != nil {
		return nil, err
	}

	return purchases, nil
}

func GetProductPurchases(productId string) ([]Purchase, error) {
	var purchases []Purchase = make([]Purchase, 0)

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

	err := db.Get(&purchase, statement, id)

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

type PaymentMethod string

const (
	CASH   PaymentMethod = "cash"
	STRIPE PaymentMethod = "stripe"
	PAYPAL PaymentMethod = "paypal"
)

var PaymentMethods = []PaymentMethod{CASH, STRIPE, PAYPAL}

func ParsePaymentMethod(method string) PaymentMethod {
	switch method {
	case "cash":
		return CASH
	case "stripe":
		return STRIPE
	case "paypal":
		return PAYPAL
	default:
		return CASH
	}
}

type Order struct {
	Id         string     `json:"id"`
	Customer   Customer   `json:"customer"`
	Purchases  []Purchase `json:"purchases"`
	Pickuptime time.Time  `json:"pickuptime"`
	Fulfilled  bool       `json:"fulfilled"`
	Method     string     `json:"method"`
	Created    time.Time  `json:"created"`
	Updated    time.Time  `json:"updated"`
}

func CreateOrder(customerId string, pickuptime time.Time, items []PurchasedItem, method PaymentMethod) (*Order, error) {
	statement := "INSERT INTO orders (id, customer, pickuptime, fulfilled, method) VALUES ($1, $2, $3, $4, $5)"

	customer, err := GetCustomer(customerId)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{Id: uuid.NewV4().String(), Customer: *customer, Pickuptime: pickuptime, Purchases: make([]Purchase, len(items)), Fulfilled: false, Method: string(method)}

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
	var orders []Order = make([]Order, 0)

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
	var orders []Order = make([]Order, 0)

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
	var orders []Order = make([]Order, 0)

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
	var orders []Order = make([]Order, 0)

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

	err := db.Get(&order, statement, id)
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

type RankedOrder struct {
	Id       string    `json:"id"`
	Cost     int       `json:"cost"`
	Customer string    `json:"customer_name"`
	Created  time.Time `json:"created"`
}

type RankedSeller struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Sold     int    `json:"sold"`
}

type RankedGainer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Gained   int    `json:"gained"`
}

type FinancesResponse struct {
	OrdersAmount    int `json:"orders_amount"`    // All orders made
	OutstandingCash int `json:"outstanding_cash"` // Unpaid cash orders
	PendingMoney    int `json:"pending_money"`    // Paid online but still unfulfilled
	Gains           int `json:"gains"`            // All money from fulfilled orders
	Total           int `json:"total"`            // Total money registred under every order made

	OrdersData   Dataset `json:"orders_data"`
	MonetaryData Dataset `json:"monetary_data"`

	PreferredMethodData []Dataset `json:"preferred_method_data"`

	FilledPie Pie `json:"filled_pie"`
	MethodPie Pie `json:"method_pie"`

	RankedOrders   []RankedOrder  `json:"ranked_orders"`
	ToppedSellers  []RankedSeller `json:"topped_sellers"`
	ToppedGainers  []RankedGainer `json:"topped_gainers"`
	FloppedSellers []RankedSeller `json:"flopped_sellers"`
	FloppedGainers []RankedGainer `json:"flopped_gainers"`
}

func GetOrdersAmount() (int, error) {
	var amount int
	statement := "SELECT COUNT(*) FROM orders"

	err := db.Get(&amount, statement)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func GetOutstandingCash() (int, error) {
	var outstanding int

	statement := `SELECT COALESCE(SUM(products.price * purchases.quantity), 0) AS outstanding
								FROM orders
								JOIN purchases ON orders.id = purchases.orderid
								JOIN products ON purchases.productid = products.id
								WHERE orders.fulfilled = false
								AND orders.method = 'cash'`

	err := db.Get(&outstanding, statement)
	if err != nil {
		return 0, err
	}

	return outstanding, nil
}

func GetPendingMoney() (int, error) {
	var pending int

	statement := `SELECT COALESCE(SUM(total_cost), 0) AS pending
								FROM (
										SELECT COALESCE(SUM(p.quantity * pr.price), 0) AS total_cost
										FROM orders o
										JOIN purchases p ON o.id = p.orderid
										JOIN products pr ON p.productid = pr.id
										WHERE o.fulfilled = false
										AND o.method != 'cash'
										GROUP BY o.id
								) AS order_totals`

	err := db.Get(&pending, statement)
	if err != nil {
		return 0, err
	}

	return pending, nil
}

func GetGains() (int, error) {
	var gains int

	statement := `SELECT COALESCE(SUM(total_cost), 0) AS gains
								FROM (
										SELECT COALESCE(SUM(p.quantity * pr.price), 0) AS total_cost
										FROM orders o
										JOIN purchases p ON o.id = p.orderid
										JOIN products pr ON p.productid = pr.id
										WHERE o.fulfilled = true
										GROUP BY o.id
								) AS order_totals`

	err := db.Get(&gains, statement)
	if err != nil {
		return 0, err
	}

	return gains, nil
}

func GetTotalFromOrders() (int, error) {
	var total int

	statement := `SELECT COALESCE(SUM(total_cost), 0) AS total
								FROM (
										SELECT COALESCE(SUM(p.quantity * pr.price), 0) AS total_cost
										FROM orders o
										JOIN purchases p ON o.id = p.orderid
										JOIN products pr ON p.productid = pr.id
										GROUP BY o.id
								) AS order_totals`

	err := db.Get(&total, statement)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetOrdersData(timeframe Timeframe, method PaymentMethod, fulfilled bool) (Dataset, error) {

	var results []Count = make([]Count, 0)
	var whereStm string

	horizontal, err := GetHorizonalDataAndQueryByTimeframe("orders.created", timeframe, &whereStm)
	if err != nil {
		return Dataset{}, err
	}

	switch method {
	case CASH:
		whereStm += " AND method = 'cash'"
	default:
		whereStm += " AND method != 'cash'"
	}

	if !fulfilled {
		whereStm += " AND fulfilled = false"
	} else {
		whereStm += " AND fulfilled = true"
	}

	statement := `SELECT DATE(created) as date, COUNT(*) as count FROM orders ` + whereStm + ` GROUP BY created ORDER BY created ASC`

	err = db.Select(&results, statement)
	if err != nil {
		return Dataset{}, err
	}

	vertical, err := ComputeVertical(results, horizontal, timeframe)
	if err != nil {
		return Dataset{}, err
	}

	return Dataset{Horizontal: horizontal, Vertical: vertical}, nil
}

func GetMonetaryData(timeframe Timeframe, method PaymentMethod, fulfilled bool) (Dataset, error) {
	var results []Count = make([]Count, 0)
	var whereStm string

	horizontal, err := GetHorizonalDataAndQueryByTimeframe("orders.created", timeframe, &whereStm)
	if err != nil {
		return Dataset{}, err
	}

	switch method {
	case CASH:
		whereStm += " AND method = 'cash'"
	default:
		whereStm += " AND method != 'cash'"
	}

	if !fulfilled {
		whereStm += " AND fulfilled = false"
	} else {
		whereStm += " AND fulfilled = true"
	}

	statement := `SELECT DATE(orders.created) as date, COALESCE(SUM(products.price * purchases.quantity), 0) as count FROM orders JOIN purchases ON orders.id = purchases.orderid JOIN products ON purchases.productid = products.id ` + whereStm + `  GROUP BY orders.created ORDER BY orders.created ASC`

	err = db.Select(&results, statement)

	if err != nil {
		return Dataset{}, err
	}

	vertical, err := ComputeVertical(results, horizontal, timeframe)
	if err != nil {
		return Dataset{}, err
	}

	return Dataset{Horizontal: horizontal, Vertical: vertical}, nil
}

func GetPreferredMethodData(timeframe Timeframe, fulfilled bool) ([]Dataset, error) {
	var results []Dataset = make([]Dataset, 0)
	var whereStm string

	horizontal, err := GetHorizonalDataAndQueryByTimeframe("created", timeframe, &whereStm)

	if err != nil {
		return nil, err
	}

	if !fulfilled {
		whereStm += " AND fulfilled = false"
	} else {
		whereStm += " AND fulfilled = true"
	}

	for method := range PaymentMethods {
		var result []Count = make([]Count, 0)
		whereStm2 := whereStm + ` AND method = '` + string(PaymentMethods[method]) + `' `

		statement := `SELECT
							DATE(created) AS date,
							COUNT(*) AS count
							FROM
									orders ` + whereStm2 + `
							GROUP BY
									date, method
							ORDER BY
									date DESC, method`

		err = db.Select(&result, statement)

		if err != nil {
			return nil, err
		}

		vertical, err := ComputeVertical(result, horizontal, timeframe)
		if err != nil {
			return nil, err
		}

		results = append(results, Dataset{Horizontal: horizontal, Vertical: vertical})

	}

	return results, nil
}

func GetTopSellers() ([]RankedSeller, error) {
	var results []RankedSeller = make([]RankedSeller, 0)

	statement := `SELECT
								pr.id AS id,
								pr.name AS name,
								pr.category AS category,
								COALESCE(SUM(p.quantity), 0) AS sold
								FROM
										products pr
								JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, pr.category
								ORDER BY
										sold DESC
								LIMIT 10`

	err := db.Select(&results, statement)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetTopOrders() ([]RankedOrder, error) {
	var results []RankedOrder = make([]RankedOrder, 0)

	statement := `SELECT
								o.id AS id,
								COALESCE(SUM(p.quantity * pr.price), 0) AS cost,
								c.fullname AS customer,
								o.created AS created
								FROM
										orders o
								JOIN
										customers c ON o.customer = c.id
								JOIN
										purchases p ON o.id = p.orderid
								JOIN
										products pr ON p.productid = pr.id
								GROUP BY
										o.id, c.fullname, o.created
								ORDER BY
										cost DESC
								LIMIT 10`

	err := db.Select(&results, statement)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetTopGainers() ([]RankedGainer, error) {
	var results []RankedGainer = make([]RankedGainer, 0)

	statement := `SELECT
								pr.id AS id,
								pr.name AS name,
								pr.category AS category,
								COALESCE(SUM(p.quantity * pr.price), 0) AS gained
								FROM
										products pr
								JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, pr.category
								ORDER BY
										gained DESC
								LIMIT 10`

	err := db.Select(&results, statement)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetFlopSellers() ([]RankedSeller, error) {
	var results []RankedSeller = make([]RankedSeller, 0)

	statement := `SELECT
								pr.id AS id,
								pr.name AS name,
								pr.category AS category,
								COALESCE(SUM(p.quantity), 0) AS sold
								FROM
										products pr
								JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, pr.category
								ORDER BY
										sold ASC
								LIMIT 10`

	err := db.Select(&results, statement)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetFlopGainers() ([]RankedGainer, error) {
	var results []RankedGainer = make([]RankedGainer, 0)

	statement := `SELECT
								pr.id AS id,
								pr.name AS name,
								pr.category AS category,
								COALESCE(SUM(p.quantity * pr.price)) AS gained
								FROM
										products pr
								JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, pr.category
								ORDER BY
										gained ASC
								LIMIT 10`

	err := db.Select(&results, statement)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetFilledPie() (Pie, error) {
	type OrderState struct {
		Filled      float64
		Unfulfilled float64
	}

	var results OrderState

	statement := `SELECT
								ROUND(COALESCE(COUNT(CASE WHEN fulfilled THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 0), 2) AS filled,
								ROUND(COALESCE(COUNT(CASE WHEN NOT fulfilled THEN 1 END) * 100.0 / NULLIF(COUNT(*), 0), 0), 2) AS unfulfilled
								FROM orders`

	err := db.Get(&results, statement)
	if err != nil {
		return Pie{}, err
	}

	return Pie{Title: "Orders State", Items: []PieItem{{Label: "Fulfilled", Value: results.Filled}, {Label: "Unfulfilled", Value: results.Unfulfilled}}}, nil
}

func GetMethodsPie() (Pie, error) {

	var results []PieItem = make([]PieItem, 0)

	statement := `SELECT
								method AS label,
								ROUND(COALESCE(COUNT(*) * 100.0 / NULLIF((SELECT COUNT(*) FROM orders), 0), 0), 2) AS value
								FROM
										orders
								GROUP BY
										method`

	err := db.Select(&results, statement)
	if err != nil {
		return Pie{}, err
	}

	return Pie{Title: "Used Payment Methods", Items: results}, nil
}
