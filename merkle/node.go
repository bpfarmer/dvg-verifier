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

// SibVal val
func (n *Node) SibVal() string {
	if n.Sib() != nil {
		return n.Sib().HashVal()
	}
	return hashEmpty("")
}

// InclusionProof comment
func (n *Node) InclusionProof() []string {
	var p []string
	curNode := n
	for curNode.Parent != nil {
		p = append(p, curNode.SibVal()+"_"+curNode.Dir())
		curNode = curNode.Parent
	}
	return p
}

// Sib comment
func (n *Node) Sib() *Node {
	if n.Parent != nil {
		if n.IsR() {
			return n.Parent.R
		}
		return n.Parent.L
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

// FindNode comment
func FindNode(s *Store, val string) *Node {
	q := "select * from nodes where val = $1"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	n := MapToNodes(rows)
	return n[0]
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
	fmt.Println(ids)
	var v []string
	for i := 1; i <= len(ids); i++ {
		v = append(v, "$"+strconv.Itoa(i))
	}
	q := "select * from nodes where id in (" + strings.Join(v, ",") + ")"
	fmt.Println(q)
	rows, err := s.DB.Query(q, ids)
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
