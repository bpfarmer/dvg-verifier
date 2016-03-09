package main

import (
	"database/sql"
	"fmt"
	"log"
	"transparency/merkle"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=gouser password=gouser dbname=transparency sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store := merkle.Store{DB: db}
	store.DropTables()
	store.AddTables()
	root := merkle.Node{Level: 0, Name: []byte("12345")}
	leaf := merkle.Node{Level: 1, Name: []byte("1234")}
	leaf.Parent = &root
	root.R = &leaf
	root.Save(&store)
	leaf.Save(&store)
	root.Save(&store)
	leaf.Save(&store)
	for i := 0; i < 4; i++ {
		merkle.ShiftDown(&leaf)
	}
}

func walk(n *merkle.Node) {
	if n.L != nil {
		fmt.Println("Walkin' Left")
		walk(n.L)
	} else if n.R != nil {
		fmt.Println("Walkin' Right")
		walk(n.R)
	}
}
