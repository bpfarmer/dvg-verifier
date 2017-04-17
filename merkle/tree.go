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

// CountNodesAtLevel just for debugging purposes
/*
func (n *Node) CountNodesAtLevel(c int, level int, s *Store) int {
	count := 0
	if c < level {
		l := n.LEntry(s)
		r := n.REntry(s)

		if l != nil {
			count += l.CountNodesAtLevel(c+1, level, s)
		}

		if r != nil {
			count += r.CountNodesAtLevel(c+1, level, s)
		}
	} else {
		return 1
	}
	return count
}*/

// ShiftInsert comment
func (n *Node) ShiftInsert(node *Node, leafCount int, s *Store) {
	l := 0
	o := n
	if leafCount%2 == 0 {
		d := targetShiftDepth(leafCount)
		for l < d && o.PEntry(s) != nil {
			o = o.PEntry(s)
			l++
		}
	}
	p := &Node{L: o, R: node, P: o.P}
	o.P = p
	node.P = p

	if p.P == nil {
		return
	}
	if p.P.L == o {
		p.P.L = p
	} else {
		p.P.R = p
	}
}

func targetShiftDepth(leafCount int) int {
	n := int(math.Floor(math.Log2(float64(leafCount))))
	// Find the complete right subtree
	for leafCount != int(math.Exp2(float64(n))) {
		if leafCount-int(math.Exp2(float64(n))) >= 0 {
			leafCount -= int(math.Exp2(float64(n)))
		}
		n--
	}
	return n
}

// RMostEntry returns the right-most descendant
func (n *Node) RMostEntry(s *Store) *Node {
	node := n
	for node.REntry(s) != nil {
		node = node.REntry(s)
	}
	return node
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
	log.Println("Tree.PEntry():looking for parent of val=" + n.Val)
	if n.P != nil {
		log.Println("Tree.PEntry():found parent in memory=" + n.P.Val)
		return n.P
	}
	n.P = FindParent(s, n.Val)
	if n.P == nil {
		log.Println("Tree.PEntry():couldn't find parent in db, returning nil")
		return nil
	}
	if n.P.LVal == n.Val {
		log.Println("Tree.PEntry():setting node as left child of parent")
		n.P.L = n
		return n.P
	}
	log.Println("Tree.PEntry():setting node as right child of parent")
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

// AddLeaf comment
func (t *Tree) AddLeaf(n *Node, s *Store) {
	if t.Root == nil {
		log.Println("Tree.AddLeaf():handle the case of a fresh tree")
		t.Root = n
		n.Save(s)
		return
	}

	leafCount := t.CountLeaves(s)
	log.Print("Tree.AddLeaf():leafCount=" + strconv.Itoa(leafCount))
	node := t.Root.RMostEntry(s)
	node.ShiftInsert(n, leafCount, s)

	if t.Root.P != nil {
		t.Root = t.Root.P
	}

	// Recursively rehash beginning with new leaf
	walkHash(n, s)

	// Recursively save nodes affected by update
	walkSave(n, s)
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
