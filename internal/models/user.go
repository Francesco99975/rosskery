package models

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DbUser struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func UserExists(email string) bool {
	statement := "SELECT * FROM users WHERE email = $1"
	var user DbUser

	err := db.Select(user, statement, email)
	return err != nil
}


func CreateUser(u *User, password string, roleId string) (*User, error) {
	statement := "INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)"
	roleStatement := "INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password: %v", err)
	}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, u.Id, u.Username, u.Email, hashedPassword); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if _, err := tx.Exec(roleStatement, u.Id, roleId); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	newUser, err := GetUserById(u.Id)
	if err != nil {
		return nil, err
	}

	return newUser.ToUser()
}

func GetUserFromEmail(email string) (*DbUser, error) {
	var user DbUser
	statement := "SELECT id, username, email, password, created, updated FROM users WHERE email = $1"

	err := db.Select(&user, statement, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserById(id string) (*DbUser, error) {
	var user DbUser
	statement := "SELECT id, username, email, password, created, updated FROM users WHERE id = $1"

	err := db.Select(&user, statement, id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserFromUsername(username string) (*DbUser, error) {
	var user DbUser
	statement := "SELECT id, username, email, password, created, updated FROM users WHERE username = $1"

	err := db.Select(&user, statement, username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]*DbUser, error) {
	var users []*DbUser
	statement := "SELECT id, username, email, password, created, updated FROM users"

	err := db.Select(&users, statement)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersByRole(role Role) ([]*DbUser, error) {
	var users []*DbUser
	statement := "SELECT id, username, email, password, created, updated FROM users WHERE role = $1"

	err := db.Select(&users, statement, role)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (user *DbUser) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	 })

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func (user *DbUser) UpdatePassword(oldPassword, password string) error {
	statement := "UPDATE users SET password = $1 WHERE id = $2"

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("Wrong Password. Unauthorized: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("Error while hashing password: %v", err)
	}

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, hashedPassword, user.Id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = tx.Commit()

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (user *DbUser) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("Wrong Password. Unauthorized: %v", err)
	}
	return nil
}

func (user *DbUser) Update(data *User) error {
	statement := "UPDATE users SET username = $1, email = $2 WHERE id = $5"

	user.Username = data.Username
	user.Email = data.Email

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, user.Username, user.Email, user.Id); err != nil {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err := tx.Commit()

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (user *DbUser) UpdateRole(role Role) error {
	statement := "UPDATE users_roles SET roleid = $1 WHERE id = $2"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, role.Id, user.Id); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err := tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (user *DbUser) Delete() error {
	statement := "DELETE FROM users WHERE id = $1"

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, user.Id); err != nil {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err := tx.Commit()

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}

func (u *DbUser) ToUser() (*User, error) {
	role := Role{}

	statement := "SELECT roleid FROM users_roles WHERE id = $1"

	err := db.Select(&role, statement, u.Id)

	if err != nil {
		return nil, err
	}

	return &User{
		Id: u.Id,
		Username: u.Username,
		Email: u.Email,
		Role: role,
	}, nil
}

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role Role `json:"role"`
}
