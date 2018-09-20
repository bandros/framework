package framework

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)


func MysqlConnect() (*sql.DB, error) {
	db, err := sql.Open("mysql",os.Getenv("mysqlUser")+":"+
	os.Getenv("mysqlPwd")+"@"+os.Getenv("mysqlHost")+"/"+os.Getenv("mysqlDb"))
	if(err) != nil {
		return  nil, err
	}
	return db,nil
}

