package models

import (
	"fmt"
	"math"
	"slices"
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
		Date  time.Time
		Count int
	}

	var results []VisitCount = make([]VisitCount, 0)
	var whereStm string
	var horizonal []string
	var vertical []int

	switch timeframe {
	case L7:
		whereStm = "WHERE date >= NOW() - INTERVAL '7' DAY"
		for i := 6; i >= 0; i-- {
			horizonal = append(horizonal, time.Now().AddDate(0, 0, -i).Weekday().String())
			vertical = append(vertical, 0)
		}
	case PW:
		whereStm = "WHERE date >= DATE_TRUNC('week', NOW()) - INTERVAL '7 DAY' - INTERVAL '1 DAY' AND date < DATE_TRUNC('week', NOW())"
		offset := int(time.Now().Weekday())
		for i := 6 + offset; i >= offset; i-- {
			horizonal = append(horizonal, fmt.Sprint(time.Now().AddDate(0, 0, -i).Day()))
			vertical = append(vertical, 0)
		}
	case L30:
		whereStm = "WHERE date >= NOW() - INTERVAL '30' DAY"
		for i := 30; i > 0; i-- {
			if i%6 == 0 {
				horizonal = append(horizonal, fmt.Sprintf("%d - %d", time.Now().AddDate(0, 0, -i).Day(), time.Now().AddDate(0, 0, (-i)+6).Day()))
				vertical = append(vertical, 0)
			}
		}
	case PM:
		whereStm = "WHERE date >= DATE_TRUNC('month', NOW()) - INTERVAL '30 DAY' - INTERVAL '1 DAY' AND date < DATE_TRUNC('month', NOW())"
		offset := time.Now().Day()
		for i := 30 + offset; i > offset; i-- {
			if i%6 == 0 {
				horizonal = append(horizonal, fmt.Sprintf("%d - %d", time.Now().AddDate(0, 0, -i).Day(), time.Now().AddDate(0, 0, (-i)+6).Day()))
				vertical = append(vertical, 0)
			}
		}
	case L12:
		whereStm = "WHERE date >= NOW() - INTERVAL '12' MONTH"
		for i := 11; i >= 0; i-- {
			horizonal = append(horizonal, time.Now().AddDate(0, -i, 0).Month().String())
			vertical = append(vertical, 0)
		}
	case PY:
		whereStm = "WHERE date >= DATE_TRUNC('year', NOW()) - INTERVAL '12 MONTH' - INTERVAL '1 MONTH' AND date < DATE_TRUNC('year', NOW())"
		offset := int(time.Now().Month())
		for i := 11 + offset; i >= offset; i-- {
			horizonal = append(horizonal, time.Now().AddDate(0, -i, 0).Month().String())
			vertical = append(vertical, 0)
		}
	default:
		return Dataset{}, fmt.Errorf("Invalid timeframe: %v", timeframe)
	}

	var statement string
	if quality == UNIQUE {
		statement = `SELECT DATE(date) AS date, COUNT(DISTINCT ip) AS count FROM visits ` + whereStm + ` GROUP BY DATE(date) ORDER BY DATE(date) ASC`

	} else {
		statement = `SELECT DATE(date) AS date, COUNT(*) AS count FROM visits ` + whereStm + ` GROUP BY DATE(date) ORDER BY DATE(date) ASC`
	}

	err := db.Select(&results, statement)
	if err != nil {
		return Dataset{}, err
	}

	switch timeframe {
	case L7:
		for i := 0; i < len(horizonal); i++ {
			for j := 0; j < len(results); j++ {
				if results[j].Date.Weekday().String() == horizonal[i] {
					vertical[i] = results[j].Count
				}
			}
		}
	case PW:
		for i := 0; i < len(horizonal); i++ {
			for j := 0; j < len(results); j++ {
				if fmt.Sprint(results[j].Date.Day()) == horizonal[i] {
					vertical[i] = results[j].Count
				}
			}
		}
	case L30, PM:
		accumulator := 0
		vindex := 0
		for i := 30; i > 0; i-- {
			for j := 0; j < len(results); j++ {
				if results[j].Date.Day() == i {
					accumulator += results[j].Count
				}
			}

			if i%6 == 0 {
				vertical[vindex] = accumulator
				accumulator = 0
				vindex++
			}
		}
	case L12, PY:
		for i := 0; i < len(horizonal); i++ {
			for j := 0; j < len(results); j++ {
				if results[j].Date.Month().String() == horizonal[i] {
					vertical[i] = results[j].Count
				}
			}
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
	var avg float64

	statement := `SELECT AVG(duration) AS avg FROM visits`

	err := db.Get(&avg, statement)

	if err != nil {
		return 0, err
	}

	return int(math.Round(avg)), nil
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
	statement := `SELECT COUNT(*) AS count FROM visits WHERE views = 0`

	err := db.Get(&count, statement)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetDeviceOrigins() ([]DeviceOrigin, error) {

	type Extrapolator struct {
		Agent string
		Count int
	}
	var extrapolation []Extrapolator = make([]Extrapolator, 0)

	statement := `SELECT DISTINCT agent, Count(agent) as count FROM visits GROUP BY agent`

	err := db.Select(&extrapolation, statement)

	if err != nil {
		return nil, err
	}

	deviceOrigins := make([]DeviceOrigin, 0)

	for i := 0; i < len(extrapolation); i++ {
		ua := useragent.Parse(extrapolation[i].Agent)
		signature := ua.OS + " / " + ua.Name

		index := slices.IndexFunc(deviceOrigins, func(deviceOrigin DeviceOrigin) bool {
			return deviceOrigin.DeviceSignaure == signature
		})

		if index == -1 {
			deviceOrigins = append(deviceOrigins, DeviceOrigin{DeviceSignaure: signature, Count: extrapolation[i].Count})
		} else {
			deviceOrigins[index].Count += extrapolation[i].Count
		}
	}

	return deviceOrigins, nil
}

func GetVisitOrigins() ([]VisitOrigin, error) {
	var origins []VisitOrigin = make([]VisitOrigin, 0)

	statement := `SELECT DISTINCT sauce, Count(sauce) as count FROM visits GROUP BY sauce`

	err := db.Select(&origins, statement)
	if err != nil {
		return nil, err
	}

	return origins, nil
}
