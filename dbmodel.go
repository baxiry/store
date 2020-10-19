package main

import (
	"database/sql"
	"fmt"
)

var (
	db  *sql.DB
	err error
)

func setdb() *sql.DB {
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/store?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	return db
}

// get all username
func getUsername(femail string) (string, string, string) {
	var name, email, password string
	err := db.QueryRow("SELECT username, email, password FROM users WHERE email = ?", femail).Scan(&name, &email, &password)
	if err != nil {
		fmt.Println(err.Error())
	}
	return name, email, password
}

func insertUser(user, pass, email, phon string) error {
	insert, err := db.Query(
		"INSERT INTO users(username, password, email, phon) VALUES ( ?, ?, ?, ? )",
		user, pass, email, phon)

	// if there is an error inserting, handle it
	if err != nil {
		return err
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return nil
}
