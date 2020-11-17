package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

func setdb() *sql.DB {
	db, err = sql.Open(
		"mysql", "root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("new db connection")
	return db
}
func insertProduct(owner, title, details, picts string, price int) error {
	insert, err := db.Query(
		"INSERT INTO stores.products(owner, title, description, price, photos) VALUES ( ?, ?, ?, ?, ?)",
		owner, title, details, price, picts)

	// if there is an error inserting, handle it
	if err != nil {
		return err
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return nil
}

// get all username
func getUsername(femail string) (string, string, string) {
	var name, email, password string
	err := db.QueryRow(
		"SELECT username, email, password FROM stores.users WHERE email = ?",
		femail).Scan(&name, &email, &password)
	if err != nil {
		fmt.Println("no result or", err.Error())
	}
	return name, email, password
}

func insertUser(user, pass, email, phon string) error {
	insert, err := db.Query(
		"INSERT INTO stores.users(username, password, email, phon) VALUES ( ?, ?, ?, ? )",
		user, pass, email, phon)

	// if there is an error inserting, handle it
	if err != nil {
		return err
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	return nil
}
