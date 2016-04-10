package merkle

import (
	"bytes"
	"crypto/sha256"
	"log"
	"strings"
)

// Node comment
type Node struct {
	Parent, L, R *Node
	LVal, RVal   []byte
	Name, Val    []byte
	Level, Epoch uint
	Path         string
	// For DB purposes
	ID, PID, LID, RID, TID int64
}

// HashVal comment
func (n *Node) HashVal() []byte {
	sha256 := sha256.New()
	if n.Name != nil {
		sha256.Write(n.Name)
	}
	if n.Val != nil {
		sha256.Write(n.Val)
	}
	if n.LVal == nil && n.L != nil {
		n.LVal = n.L.HashVal()
		sha256.Write(hashEmpty(n.LVal))
	}
	if n.RVal == nil && n.R != nil {
		n.RVal = n.R.HashVal()
		sha256.Write(hashEmpty(n.RVal))
	}
	return sha256.Sum(nil)
}

// Reset comment
func (n *Node) Reset() {
	n.LVal = nil
	n.RVal = nil
}

// subHash comment
func hashEmpty(subHash []byte) []byte {
	if subHash != nil {
		return subHash
	}
	sha256 := sha256.New()
	sha256.Write([]byte("EMPTY NODE"))
	return sha256.Sum(nil)
}

// SibVal val
func (n *Node) SibVal() []byte {
	if n.Sib() != nil {
		return n.Sib().HashVal()
	}
	return hashEmpty(nil)
}

// InclusionProof comment
func (n *Node) InclusionProof() [][]byte {
	var p [][]byte
	curNode := n
	for curNode.Parent != nil {
		p = append(p, curNode.SibVal())
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

// Find comment
func (n *Node) Find(s *Store, val string) *Node {
	q := "select * from nodes where val = ?"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	return MapToNodes(rows)[0]
}

// SetParent comment
func (n *Node) SetParent(p *Node) {
	n.Parent = p
	n.Path = p.Path + "/" + string(p.ID)
}

// FindPath comment
func (n *Node) FindPath(s *Store, val string) []*Node {
	ids := strings.Split(n.Path, "/")
	q := "select * from nodes where id in ?"
	rows, err := s.DB.Query(q, ids)
	if err != nil {
		log.Fatal(err)
	}
	nodes := MapToNodes(rows)
	nodes = append(nodes, n)
	AssocNodes(nodes)
	OrderNodes(nodes)
	return nodes
}

// OrderNodes comment
func OrderNodes(nodes []*Node) []*Node {
	for _, n in range nodes {
		
	}
}
