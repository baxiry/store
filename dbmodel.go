package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	db  *sql.DB
	err error
)

// getCatigories get all photo name of catigories.
func getCatigories(catigory string) ([]string, error) {
	var photos string
	res, err := db.Query(
		"SELECT photos FROM stores.products WHERE catigory = ?", catigory)
	if err != nil {
		return nil, err
	}
	defer res.Close() // TODO I need understand this close in mariadb

	items := make([]string, 0)
	for res.Next() {
		res.Scan(&photos)
		list := strings.Split(photos, "];[")
		// TODO split return 2 item in some casess, is this a bug ?

		items = append(items, list...)
		// TODO we need just avatar photo
	}
	return filter(items), nil
}

func insertProduct(owner, title, catigory, details, picts string, price int) error {
	insert, err := db.Query(
		"INSERT INTO stores.products(owner, title, catigory, description, price, photos) VALUES ( ?, ?, ?, ?, ?, ?)",
		owner, title, catigory, details, price, picts)

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

func setdb() *sql.DB {
	db, err = sql.Open(
		"mysql", "root:@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("new db connection")
	return db
}

// some tools
func filter(slc []string) []string {
	res := make([]string, 0)
	for _, v := range slc {
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}
