package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Product struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Image       string    `json:"image"`
	Featured    bool      `json:"featured"`
	Published   bool      `json:"published"`
	Category    Category  `json:"category"`
	Weighed     bool      `json:"weighed"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

func ProductExists(name string) bool {
	statement := "SELECT * FROM products WHERE name = $1"
	var product Product

	err := db.Get(&product, statement, name)

	return err != nil
}

func CreateProduct(name string, description string, price int, image string, categoryId string, weighed bool) ([]Product, error) {
	statement := "INSERT INTO products (id, name, description, price, image, featured, published, category, weighed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	tx := db.MustBegin()

	category, err := GetCategory(categoryId)
	if err != nil {
		return nil, err
	}

	newProduct := &Product{Id: uuid.NewV4().String(), Name: name, Description: description, Price: price, Image: image, Featured: false, Published: true, Category: *category, Weighed: weighed}

	if _, err := tx.Exec(statement, newProduct.Id, newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Image, false, true, newProduct.Category.Id, newProduct.Weighed); err != nil {
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

	updatedProducts, err := GetProducts()
	if err != nil {
		return nil, err
	}

	return updatedProducts, nil
}

func GetProducts() ([]Product, error) {
	var products []Product = make([]Product, 0)

	statement := "SELECT * FROM products ORDER BY created DESC"

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProduct(id string) (*Product, error) {
	statement := "SELECT * FROM products WHERE id = $1"
	var product Product

	err := db.Get(&product, statement, id)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetProductsByCategory(categoryId string) ([]Product, error) {
	var products []Product = make([]Product, 0)

	statement := "SELECT * FROM products WHERE category = $1"

	err := db.Select(&products, statement, categoryId)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (product *Product) Update(name string, description string, price int, image string, featured bool, published bool, categoryId string, weighed bool) ([]Product, error) {
	statement := "UPDATE products SET name = $1, description = $2, price = $3, image = $4, featured = $5, published = $6, category = $7, weighed = $8 WHERE id = $9"

	category, err := GetCategory(categoryId)
	if err != nil {
		return nil, err
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Image = image
	product.Featured = featured
	product.Published = published
	product.Category = *category
	product.Weighed = weighed

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, product.Name, product.Description, product.Price, product.Image, product.Featured, product.Published, product.Category.Id, product.Weighed); err != nil {
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

	updatedProducts, err := GetProducts()
	if err != nil {
		return nil, err
	}

	return updatedProducts, nil
}

func (product *Product) Delete() ([]Product, error) {
	statement := "DELETE FROM products WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, product.Id); err != nil {
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

	updatedProducts, err := GetProducts()
	if err != nil {
		return nil, err
	}

	return updatedProducts, nil
}
