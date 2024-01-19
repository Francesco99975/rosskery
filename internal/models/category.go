package models

import uuid "github.com/satori/go.uuid"

type Category struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func CategoryExists(name string) bool {
	statement := "SELECT * FROM categories WHERE name = $1"
	var category Category

	err := db.Select(category, statement, name)
	return err != nil
}

func CreateCategory(name string) (*Category, error) {
	statement := "INSERT INTO categories (id, name) VALUES ($1, $2)"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, uuid.NewV4().String(), name); err != nil {
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

	return &Category{Id: uuid.NewV4().String(), Name: name}, nil
}


func GetCategories() ([]Category, error) {
	var categories []Category

	statement := "SELECT * FROM categories"

	err := db.Select(&categories, statement)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategory(id string) (*Category, error) {
	var category Category

	statement := "SELECT * FROM categories WHERE id = $1"

	err := db.Select(&category, statement, id)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (category *Category) Update(name string) error {
	statement := "UPDATE categories SET name = $1 WHERE id = $2"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, name, category.Id); err != nil {
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

func (category *Category) Delete() error {
	statement := "DELETE FROM categories WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, category.Id); err != nil {
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
