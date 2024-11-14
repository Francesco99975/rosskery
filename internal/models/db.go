package models

import (
	"fmt"
	"log"
	"os"
	"time"

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

func GetNewTx() *sqlx.Tx {
	tx := db.MustBegin()

	return tx
}

type Count struct {
	Date  time.Time
	Count int
}

func GetHorizonalDataAndQueryByTimeframe(tableDate string, timeframe Timeframe, timeframeQuery *string) ([]string, error) {
	var horizontal []string

	switch timeframe {
	case L7:
		*timeframeQuery = "WHERE " + tableDate + " >= NOW() - INTERVAL '7' DAY"
		for i := 6; i >= 0; i-- {
			horizontal = append(horizontal, time.Now().AddDate(0, 0, -i).Weekday().String())
		}
	case PW:
		*timeframeQuery = "WHERE " + tableDate + " >= DATE_TRUNC('week', NOW()) - INTERVAL '7 DAY' - INTERVAL '1 DAY' AND " + tableDate + " < DATE_TRUNC('week', NOW())"
		offset := int(time.Now().Weekday())
		for i := 6 + offset; i >= offset; i-- {
			horizontal = append(horizontal, fmt.Sprint(time.Now().AddDate(0, 0, -i).Day()))
		}
	case L30:
		*timeframeQuery = "WHERE " + tableDate + " >= NOW() - INTERVAL '30' DAY"
		for i := 30; i > 0; i-- {
			if i%6 == 0 {
				horizontal = append(horizontal, fmt.Sprintf("%d - %d", time.Now().AddDate(0, 0, -i).Day(), time.Now().AddDate(0, 0, (-i)+6).Day()))
			}
		}
	case PM:
		*timeframeQuery = "WHERE " + tableDate + " >= DATE_TRUNC('month', NOW()) - INTERVAL '30 DAY' - INTERVAL '1 DAY' AND " + tableDate + " < DATE_TRUNC('month', NOW())"
		offset := time.Now().Day()
		for i := 30 + offset; i > offset; i-- {
			if i%6 == 0 {
				horizontal = append(horizontal, fmt.Sprintf("%d - %d", time.Now().AddDate(0, 0, -i).Day(), time.Now().AddDate(0, 0, (-i)+6).Day()))
			}
		}
	case L12:
		*timeframeQuery = "WHERE " + tableDate + " >= NOW() - INTERVAL '12' MONTH"
		for i := 11; i >= 0; i-- {
			horizontal = append(horizontal, time.Now().AddDate(0, -i, 0).Month().String())
		}
	case PY:
		*timeframeQuery = "WHERE " + tableDate + " >= DATE_TRUNC('year', NOW()) - INTERVAL '12 MONTH' - INTERVAL '1 MONTH' AND " + tableDate + " < DATE_TRUNC('year', NOW())"
		offset := int(time.Now().Month())
		for i := 11 + offset; i >= offset; i-- {
			horizontal = append(horizontal, time.Now().AddDate(0, -i, 0).Month().String())
		}
	default:
		return nil, fmt.Errorf("invalid timeframe: %v", timeframe)
	}

	return horizontal, nil
}

func ComputeVertical(results []Count, horizontal []string, timeframe Timeframe) ([]int, error) {
	var vertical []int = make([]int, len(horizontal))

	switch timeframe {
	case L7:
		for i := 0; i < len(horizontal); i++ {
			for j := 0; j < len(results); j++ {
				if results[j].Date.Weekday().String() == horizontal[i] {
					vertical[i] = results[j].Count
				}
			}
		}
	case PW:
		for i := 0; i < len(horizontal); i++ {
			for j := 0; j < len(results); j++ {
				if fmt.Sprint(results[j].Date.Day()) == horizontal[i] {
					vertical[i] = results[j].Count
				}
			}
		}
	case L30, PM:
		accumulator := 0
		vindex := -1
		for i := 30; i > -6; i-- {

			for j := 0; j < len(results); j++ {
				current := time.Now().AddDate(0, 0, -i)
				resdate := results[j].Date

				if resdate.Year() == current.Year() && resdate.Month() == current.Month() && resdate.Day() == current.Day() {
					accumulator += results[j].Count
				}
			}

			if i%6 == 0 {
				if vindex > -1 {
					vertical[vindex] = accumulator
				}
				accumulator = 0
				vindex++
			}

		}
	case L12, PY:
		for i := 0; i < len(horizontal); i++ {
			for j := 0; j < len(results); j++ {
				if results[j].Date.Month().String() == horizontal[i] {

					vertical[i] += results[j].Count
				}
			}
		}
	default:
		return nil, fmt.Errorf("invalid timeframe: %v", timeframe)
	}

	return vertical, nil
}
