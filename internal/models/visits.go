package models

import "time"

type Visit struct {
	Id string
	Ip string
	Views int
	Duration int
	Sauce string
	Agent string
	Date time.Time
}

func (v *Visit) Archive() error {
	v.Duration = int(time.Since(v.Date).Milliseconds())

	statement := `INSERT INTO visits(id, ip, views, duration, sauce, agent) VALUES($1, $2, $3, $4, $5, $6);`

	tx := db.MustBegin()

	if _, err := tx.Exec(statement, v.Id, v.Ip, v.Views, v.Duration, v.Sauce, v.Agent); err != nil {
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

func GetVisits() ([]Visit, error) {
	var visits []Visit

	statement := `SELECT * FROM visits`

	err := db.Select(&visits, statement)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetVisitsFromDate(date time.Time) ([]Visit, error) {
	var visits []Visit

	statement := `SELECT * FROM visits WHERE date >= $1 ORDER BY date ASC`

	err := db.Select(&visits, statement, date)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetVisitsFromIp(ip string) ([]Visit, error) {
	var visits []Visit

	statement := `SELECT * FROM visits WHERE ip = $1 ORDER BY date ASC`

	err := db.Select(&visits, statement, ip)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func CountUniqueIps() (int, error) {
	var count int

	statement := `SELECT COUNT(DISTINCT ip) FROM visits`

	err := db.Get(&count, statement)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetVisitsByMostViews() ([]Visit, error) {
	statement := `SELECT * FROM visits ORDER BY views DESC`

	var visits []Visit

	err := db.Select(&visits, statement)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetVisitsByMostDuration() ([]Visit, error) {
	statement := `SELECT * FROM visits ORDER BY duration DESC`

	var visits []Visit

	err := db.Select(&visits, statement)
	if err != nil {
		return nil, err
	}

	return visits, nil
}



