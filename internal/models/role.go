package models

import uuid "github.com/satori/go.uuid"

type Role struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func RoleExists(name string) bool {
	statement := "SELECT * FROM roles WHERE name = $1"
	var role Role

	err := db.Select(role, statement, name)
	return err != nil
}

func CreateRole(name string) (*Role, error) {
	statement := "INSERT INTO roles (id, name) VALUES ($1, $2)"

	role := &Role{Id: uuid.NewV4().String(), Name: name}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, role.Id, role.Name); err != nil {
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

	return role, nil
}

func GetRoles() ([]Role, error) {
	var roles []Role

	statement := "SELECT * FROM roles"

	err := db.Select(&roles, statement)

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func GetRoleById(id string) (*Role, error) {
	var role Role

	statement := "SELECT * FROM roles WHERE id = $1"

	err := db.Select(role, statement, id)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (role *Role) Delete() error {
	statement := "DELETE FROM roles WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, role.Id); err != nil {
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
