package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/Francesco99975/rosskery/internal/helpers"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"

	uuid "github.com/satori/go.uuid"
)

type DbPurchase struct {
	Id        string    `json:"id"`
	ProductId string    `json:"product_id" db:"productid"`
	Quantity  int       `json:"quantity"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type Purchase struct {
	Id       string    `json:"id"`
	Product  Product   `json:"product"`
	Quantity int       `json:"quantity"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

func (dbp *DbPurchase) ConvertToProduct(product Product) *Purchase {
	return &Purchase{
		Id:       dbp.Id,
		Product:  product,
		Quantity: dbp.Quantity,
		Created:  dbp.Created,
		Updated:  dbp.Updated,
	}
}

func CreatePurchase(tx *sqlx.Tx, orderId string, productId string, quantity int) (*Purchase, error) {
	statement := "INSERT INTO purchases (id, productid, quantity, orderid) VALUES ($1, $2, $3, $4)"

	product, err := GetProduct(productId)
	if err != nil {
		return nil, fmt.Errorf("error getting product while submitting purchase: %s", err)
	}

	newPurchase := &Purchase{Id: uuid.NewV4().String(), Product: *product, Quantity: quantity}

	if _, err := tx.Exec(statement, newPurchase.Id, newPurchase.Product.Id, newPurchase.Quantity, orderId); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, fmt.Errorf("error inserting purchase: %s", err)
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
	var purchases []DbPurchase = make([]DbPurchase, 0)

	statement := "SELECT id, productid, quantity, created, updated FROM purchases WHERE orderid = $1"

	err := db.Select(&purchases, statement, orderId)

	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(purchases, func(dbp DbPurchase) Purchase {
		product, err := GetProduct(dbp.ProductId)
		if err != nil {
			panic(err)
		}
		return *dbp.ConvertToProduct(*product)
	}), nil
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

func GetColorForMethod(method PaymentMethod) int {
	switch method {
	case CASH:
		return 0x22BB6F
	case STRIPE:
		return 0xD22371
	case PAYPAL:
		return 0x0D3575
	default:
		return 0x22BB6F
	}
}

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

type DbOrder struct {
	Id         string    `json:"id"`
	CustomerId string    `json:"customer_id" db:"customer"`
	Pickuptime time.Time `json:"pickuptime"`
	Fulfilled  bool      `json:"fulfilled"`
	Method     string    `json:"method"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
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

func (dbp *DbOrder) ConvertToOrder(customer Customer, purchases []Purchase) *Order {
	return &Order{
		Id:         dbp.Id,
		Customer:   customer,
		Purchases:  purchases,
		Pickuptime: dbp.Pickuptime,
		Fulfilled:  dbp.Fulfilled,
		Method:     dbp.Method,
		Created:    dbp.Created,
		Updated:    dbp.Updated,
	}
}

func CreateOrder(customerId string, pickuptime time.Time, items []PurchasedItem, method PaymentMethod) (*Order, error) {
	statement := "INSERT INTO orders (id, customer, pickuptime, fulfilled, method) VALUES ($1, $2, $3, $4, $5)"

	customer, err := GetDbCustomer(customerId)
	if err != nil {
		return nil, err
	}
	tx := db.MustBegin()

	newOrder := &Order{Id: uuid.NewV4().String(), Customer: *(*customer).ConvertToCustomer(time.Time{}, 0), Pickuptime: pickuptime, Purchases: make([]Purchase, len(items)), Fulfilled: false, Method: string(method)}

	if _, err = tx.Exec(statement, newOrder.Id, newOrder.Customer.Id, newOrder.Pickuptime, newOrder.Fulfilled, newOrder.Method); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	for i, item := range items {
		purchase, err := CreatePurchase(tx, newOrder.Id, item.ProductId, item.Quantity)
		if err != nil {
			return nil, err
		}

		newOrder.Purchases[i] = *purchase
	}

	if err = tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error rolling back transaction: %v", rollbackErr)
		}
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	createdOrder, err := GetOrder(newOrder.Id)
	if err != nil {
		return nil, err
	}

	log.Debugf("Created order %v", createdOrder)

	return createdOrder, nil
}

func GetOrders() ([]Order, error) {
	var db_orders []DbOrder = make([]DbOrder, 0)
	var orders []Order = make([]Order, 0)

	statement := "SELECT * FROM orders"

	err := db.Select(&db_orders, statement)

	for _, db_order := range db_orders {
		customer, err := GetCustomer(db_order.CustomerId)
		if err != nil {
			return nil, err
		}
		purchases, err := GetOrderPurchases(db_order.Id)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *db_order.ConvertToOrder(*customer, purchases))
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrder(id string) (*Order, error) {
	var db_order DbOrder

	statement := "SELECT * FROM orders WHERE id = $1"

	err := db.Get(&db_order, statement, id)
	if err != nil {
		return nil, err
	}
	customer, err := GetCustomer(db_order.CustomerId)
	if err != nil {
		return nil, err
	}

	purchases, err := GetOrderPurchases(db_order.Id)
	if err != nil {
		return nil, err
	}

	return db_order.ConvertToOrder(*customer, purchases), nil
}

func GetOverbooked() (string, error) {
	type Overbooked struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var overbooked []Overbooked = make([]Overbooked, 0)
	statement := `SELECT
									DATE(pickuptime) AS date,
									COUNT(*) AS count
								FROM
										orders
								GROUP BY
										DATE(pickuptime)
								HAVING
										COUNT(*) > 3
								ORDER BY
										date;`
	err := db.Select(&overbooked, statement)
	if err != nil {
		return "", err
	}

	var dates []string = make([]string, 0)

	for _, val := range overbooked {
		dates = append(dates, val.Date.Format("2006-01-02"))
	}

	return strings.Join(dates, ","), nil
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

	updatedOrder, err := GetOrder(o.Id)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (o *Order) CalculateTotal() int {
	total := 0
	for _, purchase := range o.Purchases {
		total += purchase.Product.Price * purchase.Quantity
	}
	return total
}

func (o *Order) Delete() ([]Order, error) {
	statement := "DELETE FROM orders WHERE id = $1"

	for _, purchase := range o.Purchases {
		_, err := purchase.Delete()
		if err != nil {
			return nil, err
		}
	}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, o.Id); err != nil {
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

	updatedOrders, err := GetOrders()
	if err != nil {
		return nil, err
	}

	return updatedOrders, nil
}

type RankedOrder struct {
	Id       string    `json:"id"`
	Cost     int       `json:"cost"`
	Customer string    `json:"customer"`
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

type FinancesStats struct {
	OrdersAmount    int `json:"orders_amount"`    // All orders made
	OutstandingCash int `json:"outstanding_cash"` // Unpaid cash orders
	PendingMoney    int `json:"pending_money"`    // Paid online but still unfulfilled
	Gains           int `json:"gains"`            // All money from fulfilled orders
	Total           int `json:"total"`
}

type OrdersStandingsResponse struct {
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

	statement := `SELECT COALESCE(
            ROUND(SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
						)),
            0
        ) AS outstanding
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

	statement := `SELECT COALESCE(ROUND(SUM(total_cost)), 0) AS pending
								FROM (
										SELECT COALESCE(
            SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
            ),
            0
        ) AS total_cost
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

	statement := `SELECT COALESCE(ROUND(SUM(total_cost), 0)) AS gains
								FROM (
										SELECT COALESCE(
            SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
            ),
            0
        ) AS total_cost
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

	statement := `SELECT COALESCE(ROUND(SUM(total_cost), 0)) AS total
								FROM (
										SELECT COALESCE(
            SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
            ),
            0
        ) AS total_cost
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

	statement := `SELECT DATE(orders.created) as date, COALESCE(
            SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
            ),
            0
        ) as count FROM orders JOIN purchases ON orders.id = purchases.orderid JOIN products ON purchases.productid = products.id ` + whereStm + `  GROUP BY orders.created ORDER BY orders.created ASC`

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

		results = append(results, Dataset{Topic: string(PaymentMethods[method]), Horizontal: horizontal, Vertical: vertical, Color: GetColorForMethod(PaymentMethods[method])})

	}

	return results, nil
}

func GetTopSellers() ([]RankedSeller, error) {
	var results []RankedSeller = make([]RankedSeller, 0)

	statement := `SELECT
										pr.id AS id,
										pr.name AS name,
										cat.name AS category,
										COALESCE(SUM(p.quantity), 0) AS sold
								FROM
										products pr
								JOIN
										categories cat ON pr.category = cat.id
								LEFT JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, cat.name
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
								COALESCE(
            ROUND(SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
						)),
            0
        ) AS cost,
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
										cat.name AS category,
										COALESCE(
            ROUND(SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
						)),
            0
        ) AS gained
								FROM
										products pr
								JOIN
										categories cat ON pr.category = cat.id
								LEFT JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, cat.name
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
								cat.name AS category,
								COALESCE(SUM(p.quantity), 0) AS sold
								FROM
										products pr
								JOIN
										categories cat ON pr.category = cat.id
								LEFT JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, cat.name
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
								cat.name AS category,
								COALESCE(
            ROUND(SUM(
                CASE
                    WHEN pr.weighed = true THEN (p.quantity / 10.0) * pr.price
                    ELSE p.quantity * pr.price
                END
						)),
            0
        ) AS gained
								FROM
										products pr
								JOIN
										categories cat ON pr.category = cat.id
								LEFT JOIN
										purchases p ON pr.id = p.productid
								GROUP BY
										pr.id, pr.name, cat.name
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

	return Pie{Title: "Orders State", Items: []PieItem{{Label: "Fulfilled", Value: results.Filled, Color: 0x00FF00}, {Label: "Unfulfilled", Value: results.Unfulfilled, Color: 0xFF0000}}}, nil
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

	for i := 0; i < len(results); i++ {
		results[i].Color = GetColorForMethod(ParsePaymentMethod(results[i].Label))
	}

	return Pie{Title: "Used Payment Methods", Items: results}, nil
}
