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
	"net/http/httputil"
	"os"
	"strings"
	"transparency/merkle"

	_ "github.com/lib/pq"
)

var port string
var store *merkle.Store
var authToken string
var pub *[32]byte
var priv *[64]byte

func main() {
	log.Println("main():setting environment variables")
	port = fmt.Sprintf(":%s", os.Args[1])
	log.Println("main():port=" + port)
	authToken = os.Args[2]
	log.Println("main():authToken=" + authToken)

	dbName := fmt.Sprintf("postgres://gouser:gouser@localhost/%s?sslmode=disable", "transparency")
	log.Println("main():setting up db connection to dbName=" + dbName)
	db, err := sql.Open("postgres", dbName) //"user=gouser password=gouser dbname=transparency sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store = &merkle.Store{DB: db}
	store.AddTables()
	//pub, priv, err = ed25519.GenerateKey(rand.Reader)

	leaves := loadLeaves()
	addLoadedLeaves(leaves, store)

	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/verify/", verifyReq)
	http.HandleFunc("/add/", addReq)
	http.HandleFunc("/remove/", removeReq)
	//http.HandleFunc("/reset/", resetReq)
	http.Handle("/", fs)
	http.ListenAndServe(port, nil)
}

// /verify/stanford.edu/asdfasdf9as8d7f0as98df7as
func verifyReq(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path[1:], "/")
	addr := req[1]
	val := req[2]
	log.Println("Transparency.verifyReq():val=" + val)
	res := make(map[string][]string)
	if len(addr) > 0 && len(val) > 0 {
		n := merkle.FindNode(store, val)
		if n == nil || n.Deleted {
			res["error"] = []string{"Invalid"}
		} else {
			log.Println("Transparency.verifyReq():node.Val=" + n.Val)
			res["root_hash"] = []string{n.RootHash(store)}
			res["inclusion_proof"] = n.InclusionProof(store)

			/*
				res["public_key"] = []string{hex.EncodeToString(pub[:])}
				r, err := hex.DecodeString(n.RootHash(store))
				if err != nil {
					log.Fatal(err)
				}
				s := ed25519.Sign(priv, r)[:]
				fmt.Println(s)
				res["signature"] = []string{hex.EncodeToString(s)}
			*/
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
	// TODO naive authentication, redo this before production-ready
	if r.Header.Get("X-Access-Token") != authToken {
		log.Println("addReq:failed authentication check")
		http.Error(w, "Authentication Failed", http.StatusInternalServerError)
		return
	}
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	var nodes []merkle.Node
	err := json.NewDecoder(r.Body).Decode(&nodes)
	if err != nil {
		log.Fatal(err)
	}
	root := merkle.RootEntry(store)
	t := merkle.Tree{Root: root}
	for i := range nodes {
		t.AddLeaf(&nodes[i], store)
	}
}

func removeReq(w http.ResponseWriter, r *http.Request) {
	log.Print("removeReq():received request to remove leaves=")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))

	// TODO naive authentication, redo this before production-ready
	if r.Header.Get("X-Access-Token") != authToken {
		log.Println("removeReq():failed authentication check")
		http.Error(w, "Authentication Failed", http.StatusInternalServerError)
		return
	}

	log.Println("removeReq():passed authentication")
	if r.Body == nil {
		log.Println("removeReq():no body found")
		http.Error(w, "Please send a request body", 400)
		return
	}
	var nodes []merkle.Node
	err = json.NewDecoder(r.Body).Decode(&nodes)
	if err != nil {
		log.Fatal(err)
	}
	merkle.RemoveLeaves(nodes, store)
}

// TODO may not want this in the db
func resetReq(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Access-Token") != authToken {
		log.Println("resetReq():failed authentication check")
		http.Error(w, "Authentication Failed", http.StatusInternalServerError)
		return
	}
	store.DropTables()
	store.AddTables()
}

func addLoadedLeaves(leaves []*merkle.Node, s *merkle.Store) {
	tree := &merkle.Tree{}
	for _, n := range leaves {
		if merkle.FindNode(store, n.Val) == nil {
			tree.AddLeaf(n, s)
		}
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
	//nodes = append(nodes, &merkle.Node{Val: "98cfd7226e2670eafa1f06e370d97be8047c3031e3785ac9438bfdb32e1e4041"})
	return nodes
}
