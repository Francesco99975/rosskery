package models

import (
	"mime/multipart"
	"time"

	"github.com/Francesco99975/rosskery/internal/helpers"
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

type DbProduct struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int       `json:"price"`
	Image        string    `json:"image"`
	Featured     bool      `json:"featured"`
	Published    bool      `json:"published"`
	CategoryId   string    `json:"category_id" db:"category_id"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Weighed      bool      `json:"weighed"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

func (dbp *DbProduct) ConvertToProduct() *Product {
	return &Product{
		Id:          dbp.Id,
		Name:        dbp.Name,
		Description: dbp.Description,
		Price:       dbp.Price,
		Image:       dbp.Image,
		Featured:    dbp.Featured,
		Published:   dbp.Published,
		Category:    Category{Id: dbp.CategoryId, Name: dbp.CategoryName},
		Weighed:     dbp.Weighed,
		Created:     dbp.Created,
		Updated:     dbp.Updated,
	}
}

func ProductExists(name string) bool {
	statement := `SELECT p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
									FROM products
									JOIN categories c ON p.category = c.id
									WHERE id = $1`
	var product DbProduct

	err := db.Get(&product, statement, name)

	return err != nil
}

func CreateProduct(name string, description string, price int, file *multipart.FileHeader, categoryId string, weighed bool) ([]Product, error) {
	statement := "INSERT INTO products(id, name, description, price, image, featured, published, category, weighed) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	tx := db.MustBegin()

	category, err := GetCategory(categoryId)
	if err != nil {
		return nil, err
	}

	newProduct := &Product{Id: uuid.NewV4().String(), Name: name, Description: description, Price: price, Featured: false, Published: true, Category: *category, Weighed: weighed}

	imageUrl, err := helpers.ImageUpload(file, "products", newProduct.Id)
	if err != nil {
		return nil, err
	}

	newProduct.Image = imageUrl

	if _, err := tx.Exec(statement, newProduct.Id, newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Image, newProduct.Featured, newProduct.Published, newProduct.Category.Id, newProduct.Weighed); err != nil {
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
	var products []DbProduct = make([]DbProduct, 0)

	statement := `SELECT
									p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
								FROM products p
								JOIN categories c ON p.category = c.id
								ORDER BY created DESC`

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(products, func(dbp DbProduct) Product {
		return *dbp.ConvertToProduct()
	}), nil
}

func GetPublishedProducts() ([]Product, error) {
	var products []DbProduct = make([]DbProduct, 0)
	statement := `SELECT
									p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
								FROM products p
								JOIN categories c ON p.category = c.id
								WHERE p.published = true
								ORDER BY created DESC`

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(products, func(dbp DbProduct) Product {
		return *dbp.ConvertToProduct()
	}), nil
}

func GetFeaturedProducts() ([]Product, error) {
	var products []DbProduct = make([]DbProduct, 0)
	statement := `SELECT
									p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
								FROM products p
								JOIN categories c ON p.category = c.id
								WHERE p.featured = true AND p.published = true
								ORDER BY created DESC`

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(products, func(dbp DbProduct) Product {
		return *dbp.ConvertToProduct()
	}), nil
}

func GetNewArrivals() ([]Product, error) {
	var products []DbProduct = make([]DbProduct, 0)
	statement := `SELECT
									p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
								FROM products p
								JOIN categories c ON p.category = c.id
								WHERE p.published = true AND p.created > NOW() - INTERVAL '1 WEEK'
								ORDER BY created DESC`

	err := db.Select(&products, statement)
	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(products, func(dbp DbProduct) Product {
		return *dbp.ConvertToProduct()
	}), nil
}

func GetProduct(id string) (*Product, error) {
	statement := `SELECT p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
									FROM products p
									JOIN categories c ON p.category = c.id
									WHERE p.id = $1`
	var product DbProduct

	err := db.Get(&product, statement, id)
	if err != nil {
		return nil, err
	}

	return product.ConvertToProduct(), nil
}

func GetProductsByCategory(categoryId string) ([]Product, error) {
	var products []DbProduct = make([]DbProduct, 0)

	statement := `SELECT
									p.id AS id,
									p.name AS name,
									p.description AS description,
									p.price AS price,
									p.image AS image,
									p.featured AS featured,
									p.published AS published,
									p.weighed AS weighed,
									p.created AS created,
									p.updated AS updated,
									c.id AS category_id,
									c.name AS category_name
								FROM products p
								JOIN categories c ON p.category = c.id
								WHERE p.category = $1
								ORDER BY created DESC`

	err := db.Select(&products, statement, categoryId)

	if err != nil {
		return nil, err
	}

	return helpers.MapSlice(products, func(dbp DbProduct) Product {
		return *dbp.ConvertToProduct()
	}), nil
}

func (product *Product) Update(name string, description string, price int, file *multipart.FileHeader, featured bool, published bool, categoryId string, weighed bool) ([]Product, error) {
	statement := "UPDATE products SET name = $1, description = $2, price = $3, image = $4, featured = $5, published = $6, category = $7, weighed = $8 WHERE id = $9"

	category, err := GetCategory(categoryId)
	if err != nil {
		return nil, err
	}
	var imageUrl string
	if file != nil {
		imageUrl, err = helpers.ImageUpload(file, "products", product.Id)
		if err != nil {
			return nil, err
		}
	} else {
		imageUrl = product.Image
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Image = imageUrl
	product.Featured = featured
	product.Published = published
	product.Category = *category
	product.Weighed = weighed

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, product.Name, product.Description, product.Price, product.Image, product.Featured, product.Published, product.Category.Id, product.Weighed, product.Id); err != nil {
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

	if err := helpers.DeleteImage("products", product.Id); err != nil {
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
