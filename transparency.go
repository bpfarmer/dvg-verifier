package main

import (
	"crypto/rand"
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

	"github.com/agl/ed25519"
	_ "github.com/lib/pq"
)

var store *merkle.Store
var pub *[32]byte
var priv *[64]byte

// /verify/stanford.edu/asdfasdf9as8d7f0as98df7as
func verifyReq(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	addr := req[1]
	val := req[2]
	res := make(map[string][]string)
	if len(addr) > 0 && len(val) > 0 {
		n := merkle.FindLeaf(store, val)
		if n == nil {
			res["error"] = []string{"Invalid"}
		} else {
			n.FindPath(store)
			res["root_hash"] = []string{n.RootHash()}
			res["inclusion_proof"] = n.InclusionProof()
			res["public_key"] = []string{hex.EncodeToString(pub[:])}
			r, err := hex.DecodeString(n.RootHash())
			if err != nil {
				log.Fatal(err)
			}
			s := ed25519.Sign(priv, r)[:]
			//fmt.Println(s)
			res["signature"] = []string{hex.EncodeToString(s)}
		}
	}
	js, err := json.Marshal(res)
	if err != nil || len(val) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	db, err := sql.Open("postgres", "postgres://gouser:gouser@localhost/transparency?sslmode=disable") //"user=gouser password=gouser dbname=transparency sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store = &merkle.Store{DB: db}
	store.DropTables()
	store.AddTables()
	pub, priv, err = ed25519.GenerateKey(rand.Reader)
	//fmt.Println(pub)
	//fmt.Println(priv)
	leaves := loadLeaves()
	//fmt.Println(len(leaves))
	leaves = merkle.BuildTree(leaves)
	r := leaves[0].HashVal()
	//fmt.Println(leaves)
	fmt.Println(r)
	save(leaves[0])

	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/verify/", verifyReq)
	http.Handle("/", fs)
	http.ListenAndServe(":4000", nil)
}

func save(n *merkle.Node) {
	walk(n)
	walk(n)
}

func walk(n *merkle.Node) {
	n.Save(store)
	if n.L != nil {
		walk(n.L)
	}
	if n.R != nil {
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
	nodes = append(nodes, &merkle.Node{Val: "98cfd7226e2670eafa1f06e370d97be8047c3031e3785ac9438bfdb32e1e4041"})
	return nodes
}
