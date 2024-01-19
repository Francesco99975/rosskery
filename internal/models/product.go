package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price int `json:"price"`
	Image string `json:"image"`
	Featured bool `json:"featured"`
	Published bool `json:"published"`
	Category Category `json:"category"`
	Weighed bool `json:"weighed"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func ProductExists(name string) bool {
	statement := "SELECT * FROM products WHERE name = $1"
	var product Product

	err := db.Select(product, statement, name)

	return err != nil
}

func CreateProduct(name string, description string, price int, image string, featured bool, published bool, categoryId string, weighed bool) (*Product, error) {
	statement := "INSERT INTO products (id, name, description, price, image, featured, published, category, weighed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, uuid.NewV4().String(), name, description, price, image, featured, published, categoryId, weighed); err != nil {
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

	category, err := GetCategory(categoryId)
	if err != nil {
		return nil, err
	}

	return &Product{Id: uuid.NewV4().String(), Name: name, Description: description, Price: price, Image: image, Featured: featured, Published: published, Category: *category, Weighed: weighed}, nil
}

func GetProducts() ([]Product, error) {
	var products []Product

	statement := "SELECT * FROM products"

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProduct(id string) (*Product, error) {
	statement := "SELECT * FROM products WHERE id = $1"
	var product Product

	err := db.Select(product, statement, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetProductsByCategory(categoryId string) ([]Product, error) {
	var products []Product

	statement := "SELECT * FROM products WHERE category = $1"

	err := db.Select(&products, statement, categoryId)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (product *Product) Update(name string, description string, price int, image string, featured bool, published bool, categoryId string, weighed bool) error {
	statement := "UPDATE products SET name = $1, description = $2, price = $3, image = $4, featured = $5, published = $6, category = $7, weighed = $8 WHERE id = $9"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, name, description, price, image, featured, published, categoryId, weighed); err != nil {
		errr := err
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return errr
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (product *Product) Delete() error {
	statement := "DELETE FROM products WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, product.Id); err != nil {
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
