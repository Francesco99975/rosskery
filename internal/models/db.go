package models

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB


func Setup(dsn string) {
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
        log.Fatalln(err)
  }

	err = db.Ping()
	if err != nil {
        log.Fatalln(err)
  }

	schema, err := os.ReadFile("sql/init.sql")
	if err != nil {
        log.Fatalln(err)
  }

	db.MustExec(string(schema))

	var count int

	rows, err := db.Query(`SELECT COUNT(*) AS count
													FROM users u
													JOIN users_roles ur ON u.id = ur.userid
													JOIN roles r ON ur.roleid = r.id
													WHERE r.role = 'admin';`)

	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if count == 0 {
		userId := uuid.NewV4()
		username := os.Getenv("ADMIN_USERNAME")
		email := os.Getenv("ADMIN_EMAIL")
		password := os.Getenv("ADMIN_PASSWORD")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			log.Fatalln(err)
		}

		statement := "INSERT INTO users(id, username, email, password) VALUES($1, $2, $3, $4);"

		_, err = db.Exec(statement, userId.String(), username, email, hashedPassword)

		if err != nil {
			log.Fatalln(err)
		}

		var id int

		rolesRows, err := db.Query(`SELECT id FROM roles WHERE role = 'admin';`)

		if err != nil {
			log.Fatalln(err)
		}

		for rolesRows.Next() {
			err = rolesRows.Scan(&id)
			if err != nil {
				log.Fatalln(err)
			}
		}

		statementUR := "INSERT INTO users_roles(userid, roleid) VALUES($1, $2);"

		_, err = db.Exec(statementUR, userId.String(), id)

		if err != nil {
			log.Fatalln(err)
		}

	}

}
