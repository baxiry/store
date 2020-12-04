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

type Product struct {
	Photos []string
	Title  string
}

// getCatigories get all photo name of catigories.
func getProductes(catigory string) ([]Product, error) {
	var p Product
	var picts string
	res, err := db.Query(
		"SELECT title, photos FROM stores.products WHERE catigory = ?", catigory)
	if err != nil {
		return nil, err
	}
	defer res.Close() // TODO I need understand this close in mariadb

	items := make([]Product, 0)
	for res.Next() {
		res.Scan(&p.Title, &picts)
		list := strings.Split(picts, "];[")
		// TODO split return 2 item in some casess, is this a bug ?

		p.Photos = filter(list)

		items = append(items, p)
		// TODO we need just avatar photo

	}
	return items, nil
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
	fmt.Println("new db connection")
	fmt.Println(err)
	// TODO report this error.
	// wehen db is stoped no error is return.
	// we expecte errore no database is runing

	return db
}

// some tools
func filter(slc []string) []string {
	res := make([]string, 0)
	for _, v := range slc {
		if v != "" {
			res = append(res, v) // TODO this need improve fo performence
		}
	}
	return res
}
