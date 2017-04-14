package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"math"
	"strconv"
)

// Tree comment
type Tree struct {
	Root  *Node
	Nodes [][]*Node
}

// CountLeaves comment
func (t *Tree) CountLeaves(s *Store) int {
	var count int
	q := "select count(*) from nodes where (l_val = '') is not false and (r_val = '') is not false"
	rows, err := s.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	return count
}

// RInsert comment
func (n *Node) RInsert(node *Node, s *Store) {
	log.Println("Node.RInsert():inserting node")
	o := n.P
	p := &Node{L: n, R: node}
	p.P = o
	n.P = p
	node.P = p
	// TODO: remove this
	p.LVal = n.Val
	p.RVal = node.Val
	p.Val = n.Val + "_" + node.Val
	if o == nil {
		return
	}
	log.Println("Node.RInsert():current parent=" + o.Val)
	if o.L == n {
		o.L = p
	} else {
		o.R = p
	}
}

// LInsert comment
func (n *Node) LInsert(node *Node, s *Store) {
	log.Println("Node.LInsert():inserting node")
	o := n.P
	p := &Node{L: node, R: n, P: o}
	n.P = p
	node.P = p
	// TODO: remove this
	p.LVal = node.Val
	p.RVal = n.Val
	p.Val = node.Val + "_" + n.Val
	if o == nil {
		return
	}
	log.Println("Node.RInsert():current parent=" + o.Val)
	if o.L == n {
		o.L = p
	} else {
		o.R = p
	}
}

// LEntry comment
func (n *Node) LEntry(s *Store) *Node {
	if n.L != nil {
		return n.L
	}
	n.L = FindNode(s, n.LVal)
	if n.L == nil {
		return nil
	}
	n.L.P = n
	return n.L
}

// REntry comment
func (n *Node) REntry(s *Store) *Node {
	if n.R != nil {
		return n.R
	}
	n.R = FindNode(s, n.RVal)
	if n.R == nil {
		return nil
	}
	n.R.P = n
	return n.R
}

// PEntry comment
func (n *Node) PEntry(s *Store) *Node {
	if n.P != nil {
		return n.P
	}
	n.P = FindParent(s, n.Val)
	if n.P == nil {
		return nil
	}
	if n.P.LVal == n.Val {
		n.P.L = n
		return n.P
	}
	n.P.R = n
	return n.P
}

// FindParent comment
func FindParent(s *Store, val string) *Node {
	q := "select * from nodes where l_val = $1 or r_val = $1"
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

// FindNode comment
func FindNode(s *Store, val string) *Node {
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

/*
// GetLeaves comment
func GetLeaves(s *Store) []*Node {
	q := "select * from nodes where (l_val = '') is not false and (r_val = '') is not false"
	rows, err := s.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	nodes := MapToNodes(rows)
	return nodes
}
*/

// AddLeaf comment
func (t *Tree) AddLeaf(n *Node, s *Store) {
	node := t.Root

	if node == nil {
		log.Println("Tree.AddLeaf():handle the case of a fresh tree")
		t.Root = n
		n.Save(s)
		return
	}

	leafCount := t.CountLeaves(s)
	log.Print("Tree.AddLeaf():leafCount=" + strconv.Itoa(leafCount))

	log.Println("Tree.AddLeaf():check whether the tree is currently balanced")
	if math.Log2(float64(leafCount)) == math.Floor(math.Log2(float64(leafCount))) {
		log.Println("Tree.AddLeaf():tree is currently balanced")
		node.RInsert(n, s)
	} else {
		path := strconv.FormatInt(int64(leafCount), 2)
		log.Println("Tree.AddLeaf():binary prefix of path: " + path)

		log.Println("Tree.AddLeaf():iterate through the path to find the insertion point")
		for _, d := range path {
			log.Println("Tree.AddLeaf():current path direction=" + strconv.QuoteRune(d))
			switch d {
			case '0':
				log.Println("Tree.AddLeaf():traversing left")
				if node.LEntry(s) != nil {
					log.Println("Tree.AddLeaf():left child exists, reassigning node")
					node = node.LEntry(s)
				} else {
					log.Println("Tree.AddLeaf():no left child exists, inserting left")
					node.LInsert(n, s)
				}
			case '1':
				log.Println("Tree.AddLeaf():traversing right")
				if node.REntry(s) != nil {
					log.Println("Tree.AddLeaf():right child exists, reassigning node")
					node = node.REntry(s)
				} else {
					log.Println("Tree.AddLeaf():no right child exists, inserting right")
					node.RInsert(n, s)
				}
			}
		}
	}

	if t.Root.P != nil {
		t.Root = t.Root.P
	}

	// TODO: remove this
	node.Save(s)
	if node.P != nil {
		node.P.SaveChildren(s)
	}
	// Recursively rehash beginning with new leaf
	//walkHash(n, s)

	// Recursively save nodes affected by update
	//walkSave(n, s)
}

// walkSave comment
func walkSave(n *Node, s *Store) {
	// Save the current node
	n.Save(s)
	// Look for a parent node in memory
	if n.P != nil {
		// Look for a sibling node in memory and save
		if n.P.L == n {
			if n.P.R != nil {
				n.P.R.Save(s)
			}
		} else if n.P.L != nil {
			n.P.L.Save(s)
		}
		// Traverse the path to save nodes in memory
		walkSave(n.P, s)
	}
}

// walkHash comment
func walkHash(n *Node, s *Store) {
	// Look for a parent node
	if n.PEntry(s) != nil {
		// Assuming parent exists, find sibling in mem/db, handle empty value
		n = n.PEntry(s)
		if n.LEntry(s) != nil {
			n.LVal = n.LEntry(s).Val
		} else {
			n.LVal = ""
		}
		if n.REntry(s) != nil {
			n.RVal = n.REntry(s).Val
		} else {
			n.RVal = ""
		}
		// Hash left and write values for parent
		h := sha256.New()
		io.WriteString(h, hashEmpty(n.LVal))
		io.WriteString(h, hashEmpty(n.RVal))
		n.Val = hex.EncodeToString(h.Sum(nil))

		// Recursively traverse the path of the current node
		walkHash(n, s)
	}
}
