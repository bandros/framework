package framework

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func MongoConnect(ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	var driver = "mongodb"
	var mongoSrv = os.Getenv("mogoSrv")
	if mongoSrv == "true" {
		driver += "+srv"
	}
	var uri = driver + "://" + os.Getenv("mongoUser") + ":" +
		os.Getenv("mongoPwd") + "@" + os.Getenv("mongoHost")
	var opt = options.Client().ApplyURI(uri)
	var client, err = mongo.NewClient(opt)
	if err != nil {
		return nil, nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, client.Database(os.Getenv("mongoDb")), nil
}
