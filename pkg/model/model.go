package model

import (
	"database/sql"
	"fmt"
	"log"
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

func ReformatDate(dateIn string) string {
	date, err := time.Parse("2006-01-02", dateIn)
	CheckErr(err)
	return fmt.Sprintf("%02d-%02d-%4d", date.Day(), date.Month(), date.Year())
}

func GetEndDate(dateIn string, duration int) string {
	date, err := time.Parse("2006-01-02", dateIn)
	CheckErr(err)
	endDate := date.AddDate(0, 0, duration)
	return fmt.Sprintf("%02d-%02d-%4d", endDate.Day(), endDate.Month(), endDate.Year())
}
