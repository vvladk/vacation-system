package model

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
	"vsystem/config"

	_ "github.com/mattn/go-sqlite3"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	dbConn, err := sql.Open("sqlite3", config.SqlConf.Path+config.SqlConf.DbName)
	CheckErr(err)
	err = dbConn.Ping()
	CheckErr(err)
	log.Println("Connection to DB has been created")
	return dbConn
}
func RoundDays(in float64) float64 {
	return math.Round(in*100) / 100
}

func ReformatDate(dateIn string) string {
	date, err := time.Parse("2006-01-02", dateIn)
	CheckErr(err)
	return fmt.Sprintf("%02d-%02d-%4d", date.Day(), date.Month(), date.Year())
}

func GetEndDate(dateIn string, duration float64) string {
	var endDate time.Time
	date, err := time.Parse("2006-01-02", dateIn)
	if duration < 1.0 {
		endDate = date
	} else {
		CheckErr(err)
		endDate = date.AddDate(0, 0, int(duration))
	}
	return fmt.Sprintf("%02d-%02d-%4d", endDate.Day(), endDate.Month(), endDate.Year())
}
func GetEndDate4Cal(dateIn string, duration float64) string {
	var endDate time.Time
	date, err := time.Parse("2006-01-02", dateIn)
	CheckErr(err)
	endDate = date.AddDate(0, 0, int(duration+1))
	return fmt.Sprintf("%04d-%02d-%02d", endDate.Year(), endDate.Month(), endDate.Day())
}

func GetDuration(startDate, endDate string) int {
	workHours := 9
	startWorkDay := 9
	endWOrkDay := 18
	offSet := 15

	sDate, err := time.Parse("2006-01-02", startDate)
	sDate = sDate.Add(time.Hour * time.Duration(startWorkDay))
	CheckErr(err)

	nextDate := sDate

	eDate, err := time.Parse("2006-01-02", endDate)
	eDate = eDate.Add(time.Hour * time.Duration(endWOrkDay))
	CheckErr(err)
	duration := 0

	if eDate.Weekday() == time.Saturday {
		eDate = eDate.Add(-time.Hour * time.Duration(24))
	}

	if eDate.Weekday() == time.Sunday {
		eDate = eDate.Add(-time.Hour * time.Duration(48))
	}

	for {
		if (nextDate.Weekday() == time.Saturday) || (nextDate.Weekday() == time.Sunday) {
			nextDate = nextDate.Add(time.Hour * time.Duration(24))
			continue
		}
		nextDate = nextDate.Add(time.Hour * time.Duration(workHours))
		duration++
		log.Println(eDate, nextDate)
		if nextDate.Equal(eDate) {
			break
		}
		nextDate = nextDate.Add(time.Hour * time.Duration(offSet))
	}

	return duration - GetExtraDays(startDate, endDate)
}

// Return number of extra day between 2 dates
func GetExtraDays(sDate, eDate string) int {
	db := GetDB()
	defer db.Close()
	var days int
	sqlS := `SELECT COUNT(*) FROM extra_days WHERE extra_day BETWEEN ? AND ?`

	db.QueryRow(sqlS, sDate, eDate).Scan(&days)

	return days
}
