package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type ConnSettings struct {
	user string
	pass string
	name string
	host string
	dbms string
	port string
}

func GetConnSettings() ConnSettings {
	return ConnSettings{
		user: os.Getenv("DB_USER"),
		pass: os.Getenv("DB_PASS"),
		name: os.Getenv("DB_NAME"),
		host: os.Getenv("DB_HOST"),
		dbms: os.Getenv("DB_TYPE"),
		port: os.Getenv("DB_PORT"),
	}
}

func GetInstance() (*sql.DB, error) {
	stngs := GetConnSettings()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", stngs.user, stngs.pass, stngs.host, stngs.port, stngs.name)
	db, err := sql.Open(stngs.dbms, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
