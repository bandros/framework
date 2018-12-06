package framework

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func MysqlConnect() (*sql.DB, error) {
	var db *sql.DB
	var err error
	if os.Getenv("MYSQL_CONNECTION") != "" {
		db, err = sql.Open("mysql", os.Getenv("MYSQL_CONNECTION"))
	} else {
		db, err = sql.Open("mysql", os.Getenv("mysqlUser")+":"+
			os.Getenv("mysqlPwd")+"@"+os.Getenv("mysqlHost")+"/"+os.Getenv("mysqlDb"))
	}

	if (err) != nil {
		return nil, err
	}
	return db, nil
}
