package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"transparency/merkle"

	_ "github.com/lib/pq"
)

var store *merkle.Store

// /verify/stanford.edu/asdfasdf9as8d7f0as98df7as
func verifyReq(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	addr := req[1]
	val := req[2]
	if len(addr) > 0 && len(val) > 0 {

	}
	fmt.Fprintf(w, "Request: %s!", r.URL.Path[1:])
}

func main() {
	db, err := sql.Open("postgres", "user=gouser password=gouser dbname=transparency sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store = &merkle.Store{DB: db}
	store.DropTables()
	store.AddTables()

	leaves := loadLeaves()
	fmt.Println(leaves)
	walk(merkle.BuildTree(leaves)[0])
	walk(merkle.BuildTree(leaves)[0])

	http.HandleFunc("/verify/", verifyReq)
	http.ListenAndServe(":4000", nil)
}

func walk(n *merkle.Node) {
	//fmt.Println(n.ID)
	if n.L != nil {
		//fmt.Println("Left")
		//fmt.Println(n.L.ID)
		n.L.Save(store)
		walk(n.L)
	}
	if n.R != nil {
		//fmt.Println("Right")
		//fmt.Println(n.R.ID)
		n.R.Save(store)
		walk(n.R)
	}
}

func loadLeaves() []*merkle.Node {
	var nodes []*merkle.Node
	dir := "./verified"
	files, _ := ioutil.ReadDir(dir)
	hasher := sha256.New()
	for _, f := range files {
		fmt.Println(f.Name())
		o, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer o.Close()
		if _, err := io.Copy(hasher, o); err != nil {
			log.Fatal(err)
		}
		fmt.Println(hex.EncodeToString(hasher.Sum(nil)))
		nodes = append(nodes, &merkle.Node{Val: hasher.Sum(nil)})
		hasher.Reset()
	}
	hasher.Write(nodes[0].Val)
	hasher.Write(nodes[1].Val)
	fmt.Println(hex.EncodeToString(hasher.Sum(nil)))
	return nodes
}
