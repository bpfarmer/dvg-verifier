package merkle

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
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

// ShiftInsert comment
func (n *Node) ShiftInsert(node *Node, leafCount int, s *Store) {
	l := 0
	if leafCount%2 == 0 {
		d := targetShiftDepth(leafCount)
		for l < d && n.PEntry(s) != nil {
			n = n.PEntry(s)
			l++
		}
	}

	p := &Node{L: n, R: node, P: n.P}
	n.P = p
	node.P = p

	log.Println("ShiftInsert - Inserting a new node with children:")
	log.Println("Left Node: " + n.Val)
	log.Println("Right Node: " + node.Val)

	if p.P == nil {
		return
	}
	if p.P.L == n {
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
	log.Println("targetShiftDepth - " + strconv.Itoa(n))
	return n
}

// RemoveLeaves comment
func RemoveLeaves(nodes []Node, s *Store) {
	log.Println("RemoveLeaves: Removing nodes")
	var vals []string
	for i := range nodes {
		vals = append(vals, nodes[i].Val)
	}
	log.Println("RemoveLeaves: vals = " + strings.Join(vals, ","))
	q := `UPDATE nodes SET deleted=true WHERE VAL IN ($1);`
	s.Exec(func(tx *sql.Tx) {
		stmt, err := tx.Prepare(q)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(strings.Join(vals, ","))
		if err != nil {
			log.Fatal(err)
		}
	})
}

// AddLeaf comment
func (t *Tree) AddLeaf(n *Node, s *Store) {
	// Handle the case of a fresh tree
	if t.Root == nil {
		t.Root = n
		n.Save(s)
		return
	}

	log.Println("AddLeaf: Dealing with existing tree with root: " + t.Root.Val)
	var o = FindNode(s, n.Val)
	if o != nil {
		o.Deleted = false
		n = o
	} else {
		leafCount := t.CountLeaves(s)
		node := t.Root.RMostEntry(s)
		log.Println("AddLeaf: Below node: " + t.Root.Val)
		log.Println("AddLeaf: Rightmost Entry: " + node.Val)
		node.ShiftInsert(n, leafCount, s)
	}

	// Recursively rehash beginning with new leaf
	walkHash(n, s)

	// Recursively save nodes affected by update
	walkSave(n, s)

	if t.Root.PEntry(s) != nil {
		t.Root = t.Root.PEntry(s)
		log.Println("Shifting Root to " + t.Root.Val)
	}
}

// walkSave comment
func walkSave(n *Node, s *Store) {
	// Save the current node
	n.Save(s)
	// Look for a parent node in memory
	if n.PEntry(s) != nil {
		// Look for a sibling node in memory and save
		if n.P.LEntry(s) == n {
			if n.P.REntry(s) != nil {
				n.P.R.Save(s)
			}
		} else if n.P.LEntry(s) != nil {
			n.P.LEntry(s).Save(s)
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
		// Hash left and right values for parent
		h := sha256.New()
		io.WriteString(h, hashEmpty(n.LVal))
		io.WriteString(h, hashEmpty(n.RVal))
		n.Val = hex.EncodeToString(h.Sum(nil))

		// Recursively traverse the path of the current node
		walkHash(n, s)
	}
}

// RootEntry comment
func RootEntry(s *Store) *Node {
	q := "select * from nodes limit 1"
	rows, err := s.DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	n := MapToNodes(rows)
	if len(n) > 0 {
		p := n[0]
		for p.PEntry(s) != nil {
			p = p.PEntry(s)
		}
		return p
	}
	return nil
}
