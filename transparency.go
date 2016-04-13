package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
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
		n := merkle.FindLeaf(store, val)
		n.FindPath(store)
		res := make(map[string][]string)
		res["root_hash"] = []string{n.RootHash()}
		res["inclusion_proof"] = n.InclusionProof()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
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
	leaves = merkle.BuildTree(leaves)
	r := leaves[0].HashVal()
	fmt.Println(leaves)
	fmt.Println(r)
	save(leaves[0])

	http.HandleFunc("/verify/", verifyReq)
	http.ListenAndServe(":4000", nil)
}

func save(n *merkle.Node) {
	walk(n)
	walk(n)
}

func walk(n *merkle.Node) {
	n.Save(store)
	if n.L != nil {
		//fmt.Println("Left")
		//fmt.Println(n.L.ID)
		walk(n.L)
	}
	if n.R != nil {
		//fmt.Println("Right")
		//fmt.Println(n.R.ID)
		walk(n.R)
	} else {
		fmt.Println(n.InclusionProof())
	}
}

func loadLeaves() []*merkle.Node {
	var nodes []*merkle.Node
	dir := "./verified"
	files, _ := ioutil.ReadDir(dir)
	h := sha256.New()
	for _, f := range files {
		o, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer o.Close()
		if _, err := io.Copy(h, o); err != nil {
			log.Fatal(err)
		}
		nodes = append(nodes, &merkle.Node{Val: hex.EncodeToString(h.Sum(nil))})
		h.Reset()
	}
	return nodes
}
