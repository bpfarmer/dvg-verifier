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
	log.Println("Transparency.verifyReq():val=" + val)
	res := make(map[string][]string)
	if len(addr) > 0 && len(val) > 0 {
		n := merkle.FindNode(store, val)
		log.Println("Transparency.verifyReq():node.Val=" + n.Val)
		if n == nil {
			res["error"] = []string{"Invalid"}
		} else {
			res["root_hash"] = []string{n.RootHash(store)}
			res["inclusion_proof"] = n.InclusionProof(store)
			//res["public_key"] = []string{hex.EncodeToString(pub[:])}
			//r, err := hex.DecodeString(n.RootHash(store))
			//if err != nil {
			//	log.Fatal(err)
			//}
			//s := ed25519.Sign(priv, r)[:]
			//fmt.Println(s)
			//res["signature"] = []string{hex.EncodeToString(s)}
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

// POST /add
func addReq(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	addr := req[1]
	val := req[2]
	log.Println("Transparency.addReq():val=" + val)
	res := make(map[string][]string)
	if len(addr) > 0 && len(val) > 0 {
		n := merkle.FindNode(store, val)
		log.Println("Transparency.addReq():node val=" + n.Val)
		if n == nil {
			res["error"] = []string{"Invalid"}
		} else {
			res["root_hash"] = []string{n.RootHash(store)}
			res["inclusion_proof"] = n.InclusionProof(store)
			res["public_key"] = []string{hex.EncodeToString(pub[:])}
			r, err := hex.DecodeString(n.RootHash(store))
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
	dbName := fmt.Sprintf("postgres://gouser:gouser@localhost/%s?sslmode=disable", "transparency")
	db, err := sql.Open("postgres", dbName) //"user=gouser password=gouser dbname=transparency sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store = &merkle.Store{DB: db}
	store.DropTables()
	store.AddTables()
	//pub, priv, err = ed25519.GenerateKey(rand.Reader)
	leaves := loadLeaves()
	addLoadedLeaves(leaves, store)
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/verify/", verifyReq)
	http.HandleFunc("/add/", addReq)
	http.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), nil)
}

func test(s *merkle.Store) {
	//Hashing some things and making them nodes
	log.Println("Hashing some things")
	nodes := [...]*merkle.Node{&merkle.Node{Val: "1"},
		&merkle.Node{Val: "2"},
		&merkle.Node{Val: "3"},
		&merkle.Node{Val: "4"},
		&merkle.Node{Val: "5"},
		&merkle.Node{Val: "6"},
		&merkle.Node{Val: "7"},
		&merkle.Node{Val: "8"},
		&merkle.Node{Val: "9"},
		&merkle.Node{Val: "10"},
		&merkle.Node{Val: "11"},
		&merkle.Node{Val: "12"},
		&merkle.Node{Val: "13"},
		&merkle.Node{Val: "14"},
		&merkle.Node{Val: "15"},
		&merkle.Node{Val: "16"},
		&merkle.Node{Val: "17"},
		&merkle.Node{Val: "18"}}

	//Store things in a Tree
	log.Println("Storing some things in a tree")
	tree := &merkle.Tree{}
	for _, n := range nodes {
		tree.AddLeaf(n, s)
	}

	/*
		log.Println("Counting Nodes at Level 0 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 0, s)))
		log.Println("Counting Nodes at Level 1 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 1, s)))
		log.Println("Counting Nodes at Level 2 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 2, s)))
		log.Println("Counting Nodes at Level 3 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 3, s)))
		log.Println("Counting Nodes at Level 4 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 4, s)))
		log.Println("Counting Nodes at Level 5 " + strconv.Itoa(tree.Root.CountNodesAtLevel(0, 5, s)))
	*/
	//Append to that Tree
	log.Println("Appending to that tree")
	//Remove from that Tree
	log.Println("Removing from that tree")
}

func addLoadedLeaves(leaves []*merkle.Node, s *merkle.Store) {
	tree := &merkle.Tree{}
	for _, n := range leaves {
		tree.AddLeaf(n, s)
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
