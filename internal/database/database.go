package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	username string
	password string
	host     string
	port     string
	database string
	db       *sql.DB
}

func NewDatabase(username string, password string, host string, port string, database string) *Database {
	return &Database{username: username, password: password, host: host, port: port, database: database}
}

func (d *Database) Connect() error {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.username, d.password, d.host, d.port, d.database)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Database) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}
	return nil
}
