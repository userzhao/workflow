package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	MySqlDB *sql.DB
	CasSqlDB *sql.DB
)

func init() {
	var err error
	CasSqlDB, err = sql.Open("mysql", "yu.zhang:@(192.168.162.108:3306)/cas_23?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = CasSqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	MySqlDB, err = sql.Open("mysql", "root:123@(127.0.0.1:3306)/go?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = MySqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

