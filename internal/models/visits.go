package models

import (
	"fmt"
	"time"

	"github.com/mileusna/useragent"
)

type Visit struct {
	Id       string    `json:"id"`
	Ip       string    `json:"ip"`
	Views    int       `json:"views"`
	Duration int       `json:"duration"`
	Sauce    string    `json:"sauce"`
	Agent    string    `json:"agent"`
	Date     time.Time `json:"date"`
}

type VisitQuality string

const (
	ALL    VisitQuality = "all"
	UNIQUE VisitQuality = "unique"
)

func ParseVisitQuality(str string) VisitQuality {
	switch str {
	case "all":
		return ALL
	case "unique":
		return UNIQUE
	}
	return ALL
}

type Timeframe string

const (
	L7  Timeframe = "l7"  // by day name
	PW  Timeframe = "pw"  // by day number
	L30 Timeframe = "l30" // by day in 5 fractions
	PM  Timeframe = "pm"  // by day in 5 fractions
	L12 Timeframe = "l12" // by month
	PY  Timeframe = "py"  // by month
)

func ParseTimeframe(str string) Timeframe {
	switch str {
	case "l7":
		return L7
	case "pw":
		return PW
	case "l30":
		return L30
	case "pm":
		return PM
	case "l12":
		return L12
	case "py":
		return PY
	}
	return L7
}

type VisitOrigin struct {
	Sauce string `json:"sauce"` // referrer
	Count int    `json:"count"`
}

type DeviceOrigin struct {
	DeviceSignaure string `json:"device_signature"` // user-agent (OS+Browser)
	Count          int    `json:"count"`
}

type VisitsResponse struct {
	Current             int            `json:"current"`
	TotalViews          int            `json:"total_views"`
	BounceRate          string         `json:"bounce_rate"`
	AvgVisitDuration    int            `json:"avg_visit_duration"` //seconds
	TotalVisits         int            `json:"total_visits"`
	TotalUniqueVisitors int            `json:"total_unique_visitors"`
	VisitOrigins        []VisitOrigin  `json:"visit_origins"`
	DeviceOrigins       []DeviceOrigin `json:"device_origins"`
	Data                Dataset        `json:"data"`
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
	var visits []Visit = make([]Visit, 0)

	statement := `SELECT * FROM visits`

	err := db.Select(&visits, statement)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetVisitsByQualityAndTimeframe(quality VisitQuality, timeframe Timeframe) (Dataset, error) {
	type VisitCount struct {
		date  time.Time
		count int
	}

	var resl7 []VisitCount
	var whereStm string
	var horizonal []string
	var vertical []int

	switch timeframe {
	case L7:
		whereStm = "WHERE date >= NOW() - INTERVAL '7' DAY"
	case PW:
		whereStm = "WHERE date >= DATE_TRUNC('week', NOW()) - INTERVAL '7 DAY' - INTERVAL '1 DAY' AND date < DATE_TRUNC('week', NOW())"
	case L30:
		whereStm = "WHERE date >= NOW() - INTERVAL '30' DAY"
	case PM:
		whereStm = "WHERE date >= DATE_TRUNC('month', NOW()) - INTERVAL '30 DAY' - INTERVAL '1 DAY' AND date < DATE_TRUNC('month', NOW())"
	case L12:
		whereStm = "WHERE date >= NOW() - INTERVAL '12' MONTH"
	case PY:
		whereStm = "WHERE date >= DATE_TRUNC('year', NOW()) - INTERVAL '12 MONTH' - INTERVAL '1 MONTH' AND date < DATE_TRUNC('year', NOW()"
	default:
		return Dataset{}, fmt.Errorf("Invalid timeframe: %v", timeframe)
	}

	var statement string
	if quality == UNIQUE {
		statement = `SELECT date, COUNT(DISTINCT ip) AS count FROM visits ` + whereStm + ` GROUP BY date ORDER BY date ASC`

	} else {
		statement = `SELECT date, COUNT(*) AS count FROM visits ` + whereStm + ` GROUP BY date ORDER BY date ASC`
	}

	err := db.Select(&resl7, statement)
	if err != nil {
		return Dataset{}, err
	}

	accumulator := 0
	for i := 0; i < len(resl7); i++ {
		switch timeframe {
		case L7, PW:
			horizonal = append(horizonal, resl7[i].date.Weekday().String())
			vertical = append(vertical, resl7[i].count)
		case L30, PM:
			accumulator += resl7[i].count
			if i%6 == 0 {
				horizonal = append(horizonal, fmt.Sprintf("%d - %d", resl7[i].date.AddDate(0, 0, -6).Day(), resl7[i].date.Day()))
				vertical = append(vertical, accumulator)
				accumulator = 0
			}
		case L12, PY:
			horizonal = append(horizonal, resl7[i].date.Month().String())
			vertical = append(vertical, resl7[i].count)
		}
	}

	return Dataset{Horizontal: horizonal, Vertical: vertical}, nil
}

func GetVisitsFromIp(ip string) ([]Visit, error) {
	var visits []Visit = make([]Visit, 0)

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

func CountTotalViews() (int, error) {
	var count int
	statement := `SELECT SUM(views) FROM visits`

	err := db.Get(&count, statement)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetVisitsByMostViews() ([]Visit, error) {
	statement := `SELECT * FROM visits ORDER BY views DESC`

	var visits []Visit = make([]Visit, 0)

	err := db.Select(&visits, statement)

	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetAverageVisitDuration() (int, error) {
	var avg int

	statement := `SELECT AVG(duration) FROM visits`

	err := db.Get(&avg, statement)

	if err != nil {
		return 0, err
	}

	return avg, nil
}

func GetVisitsByMostDuration() ([]Visit, error) {
	statement := `SELECT * FROM visits ORDER BY duration DESC`

	var visits []Visit = make([]Visit, 0)

	err := db.Select(&visits, statement)
	if err != nil {
		return nil, err
	}

	return visits, nil
}

func GetVisitsWithZeroViews() (int, error) {
	var count int
	statement := `SELECT COUNT(*) FROM visits WHERE views = 0`

	err := db.Get(&count, statement)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetDeviceOrigins() ([]DeviceOrigin, error) {

	type Extrapolator struct {
		agent string
		count int
	}
	var extrapolation []Extrapolator = make([]Extrapolator, 0)

	statement := `SELECT DISTINCT agent, Count(DISTINCT agent) as count FROM visits`

	err := db.Select(&extrapolation, statement)

	if err != nil {
		return nil, err
	}

	deviceOrigins := make([]DeviceOrigin, 0)

	for i := 0; i < len(extrapolation); i++ {
		ua := useragent.Parse(extrapolation[i].agent)
		deviceOrigins = append(deviceOrigins, DeviceOrigin{DeviceSignaure: ua.OS + " / " + ua.Name, Count: extrapolation[i].count})
	}

	return deviceOrigins, nil
}

func GetVisitOrigins() ([]VisitOrigin, error) {
	var origins []VisitOrigin = make([]VisitOrigin, 0)

	statement := `SELECT DISTINCT origin as sauce, Count(DISTINCT origin) as count FROM visits`

	err := db.Select(&origins, statement)
	if err != nil {
		return nil, err
	}

	return origins, nil
}
