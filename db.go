package framework

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func MysqlConnect() (*sql.DB, error) {
	var db *sql.DB
	var err error
	if os.Getenv("mysqlGoogle") == "1" {
		connect := os.Getenv("CONNECTION")
		user := os.Getenv("USER")
		pass := os.Getenv("PASS")
		dbName := os.Getenv("DB")
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/%s",connect,user,pass,dbName))
	}else{
		db, err = sql.Open("mysql", os.Getenv("mysqlUser")+":"+
			os.Getenv("mysqlPwd")+"@"+os.Getenv("mysqlHost")+"/"+os.Getenv("mysqlDb"))
	}

	if (err) != nil {
		return nil, err
	}
	return db, nil
}
