package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Node comment
type Node struct {
	Parent, L, R          *Node
	LVal, RVal, Name, Val string
	Level, Epoch          uint
	Path                  string
	// For DB purposes
	ID, PID, LID, RID, TID int
}

// HashVal comment
func (n *Node) HashVal() string {
	h := sha256.New()
	if n.IsLeaf() {
		return n.Val
	}
	if len(n.LVal) == 0 && n.L != nil {
		n.LVal = n.L.HashVal()
	}
	if len(n.RVal) == 0 && n.R != nil {
		n.RVal = n.R.HashVal()
	}
	io.WriteString(h, hashEmpty(n.LVal))
	io.WriteString(h, hashEmpty(n.RVal))
	return hex.EncodeToString(h.Sum(nil))
}

// Reset comment
func (n *Node) Reset() {
	n.LVal = ""
	n.RVal = ""
}

// subHash comment
func hashEmpty(subHash string) string {
	if len(subHash) > 0 {
		return subHash
	}
	h := sha256.New()
	h.Write([]byte("EMPTY NODE"))
	return hex.EncodeToString(h.Sum(nil))
}

// InclusionProof comment
func (n *Node) InclusionProof() []string {
	var p []string
	c := n
	fmt.Println("--INCLUSION PROOF--")
	fmt.Println(c.Val)
	for c.Parent != nil {
		if c.Parent.LVal == c.HashVal() {
			p = append(p, c.Parent.RVal+"_R")
		} else {
			p = append(p, c.Parent.LVal+"_L")
		}
		c = c.Parent
	}
	return p
}

// RootHash comment
func (n *Node) RootHash() string {
	c := n
	for c.Parent != nil {
		c = c.Parent
	}
	return c.HashVal()
}

// Sib comment
func (n *Node) Sib() *Node {
	if n.Parent != nil {
		if n.IsR() {
			return n.Parent.L
		}
		return n.Parent.R
	}
	return nil
}

// Prefix comment
func (n *Node) Prefix(buf bytes.Buffer) string {
	for n.Parent != nil {
		if n.IsR() {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
	}
	return buf.String()
}

// IsR comment
func (n *Node) IsR() bool {
	return n.Parent.R == n
}

// Dir comment
func (n *Node) Dir() string {
	if n.IsR() {
		return "R"
	}
	return "L"
}

// IsLeaf comment
func (n *Node) IsLeaf() bool {
	return n.L == nil && n.R == nil
}

var leaves []*Node

// Leaves comment
func (n *Node) Leaves() []*Node {
	leaves = []*Node{}
	return n.leaves()
}

// Leaves comment
func (n *Node) leaves() []*Node {
	if n.IsLeaf() {
		leaves = append(leaves, n)
	} else {
		if n.L != nil {
			n.L.leaves()
		}
		if n.R != nil {
			n.R.leaves()
		}
	}
	return leaves
}

// FindLeaf comment
func FindLeaf(s *Store, val string) *Node {
	q := "select * from nodes where val = $1"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	n := MapToNodes(rows)
	if len(n) > 0 {
		return n[0]
	}
	return nil
}

// GetLeaves comment
func GetLeaves(s *Store) []*Node {
	q := "select * from nodes where (lval = '') is not false and (rval = '') is not false"
	rows, err := s.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	nodes := MapToNodes(rows)
	return nodes
}

// Addr comment
func (n *Node) Addr() string {
	if n.Parent != nil {
		return n.Parent.Addr() + "/" + strconv.Itoa(n.Parent.ID)
	}
	return ""
}

// FindPath comment
func (n *Node) FindPath(s *Store) []*Node {
	ids := strings.Split(n.Path, "/")
	ids = ids[1:]
	q := "select * from nodes where id = any($1::integer[])"
	rows, err := s.DB.Query(q, fmt.Sprintf("{%s}", strings.Join(ids, ", ")))
	if err != nil {
		log.Fatal(err)
	}
	path := MapToNodes(rows)
	path = append(path, n)
	AssocNodes(path)
	OrderPath(path)
	return path
}

// OrderPath comment
func OrderPath(n []*Node) []*Node {
	var leaf int
	for p, v := range n {
		if len(v.Val) > 0 {
			leaf = p
		}
	}
	var o []*Node
	c := n[leaf]
	for c != nil {
		o = append(o, c)
		c = c.Parent
	}
	return o
}
