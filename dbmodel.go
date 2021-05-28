package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

type Product struct {
	Id          int
	Title       string
	Description string
	Photo       string
	Photos      []string
	Price       string
}
func myProducts(owner string) []Product {                                                           
                                                                                                
    rows, err := db.Query("select id, title, description, photos, price from stores.products where owner = ?", owner)                      
    if err != nil {                                                                             
        fmt.Println(err)                                                                        
    }                                                                                           
    defer rows.Close() // ??
                                                                                                
    var products = []Product{}                                                                  
    var p = Product{}                                                                            
                                                                                                 
    // iterate over rows                                                                         
    for rows.Next() {                                                                            
        err = rows.Scan(&p.Id, &p.Title, &p.Description, &p.Photo, &p.Price)                    
        if err != nil {                                                                                                         
            fmt.Println("At myPorducts", err)
        }                                                                                                         

        products = append(products, p)                                                                            
                                                                                                                  
        fmt.Println(p)                                                                                             
    }                                                                                                             
    return products                                                                                          
}                      



func getProduct(id int) (Product, error) {
	var p Product
	var picts string
	err := db.QueryRow(
		"SELECT title, description, photos, price FROM stores.products WHERE id = ?",
		id).Scan(&p.Title, &p.Description, &picts, &p.Price)
	if err != nil {
		return p, err
	}

	// why not shold not close this connection ?

	list := strings.Split(picts, "];[")
	// TODO split return 2 item in some casess, is this a bug ?
	p.Photos = filter(list)
	//fmt.Println("product form db : ", p)
	return p, nil
}

// getCatigories get all photo name of catigories.
func getProductes(catigory string) ([]Product, error) {
	var p Product
	var picts string
	res, err := db.Query(
		"SELECT id, title, photos, price FROM stores.products WHERE catigory = ?", catigory)
	if err != nil {
		return nil, err
	}
	defer res.Close() // TODO I need understand this close in mariadb

	items := make([]Product, 0)
	for res.Next() {
		res.Scan(&p.Id, &p.Title, &picts, &p.Price)
		list := strings.Split(picts, "];[")
		// TODO split return 2 item in some casess, is this a bug ?
		p.Photo = list[0]
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
	defer insert.Close() // TODO why we need closeing this connection ?
    
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
		"mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local")
	if err != nil { // why no error when db is not runinig ?? 
        fmt.Println("run mysql server", err)
		// TODO report this error.

		// wehen db is stoped no error is return.
		// we expecte errore no database is runing

        // my be this error is fixed with panic ping pong bellow
	}
    
    if err = db.Ping(); err != nil {
        // TODO handle this error: dial tcp 127.0.0.1:3306: connect: connection refused
        fmt.Println("mybe database is not runing or error is: ", err)
        os.Exit(1)
    }
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
