package main

import (
	"fmt"
	"store"
	"strconv"
	"time"
)

var max = 10 //00000

func main() {

	db, err := store.Open("test.db")
	if err != nil {
		panic(err)
	}

	for i := 0; i < max; i++ {
		db.Put(strconv.Itoa(i), "hello how ai "+strconv.Itoa(i))
	}

	data, _, _ := db.Get("9")
	fmt.Println("data :", data)
	time.Sleep(time.Second * 10)

	for i := 0; i < max; i++ {
		a, _, _ := db.Get(strconv.Itoa(i))
		fmt.Println(a)
	}

	for i := 0; i < max; i++ {
		db.Put(strconv.Itoa(i), "")
	}
	data, _, _ = db.Get("9")
	fmt.Println("data :", data)

	for i := 0; i < max; i++ {
		a, _, _ := db.Get(strconv.Itoa(i))
		fmt.Println(a)
	}
	time.Sleep(time.Second * 20)
}
