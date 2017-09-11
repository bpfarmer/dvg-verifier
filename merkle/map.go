package merkle

import (
	"database/sql"
	"log"
)

// MapToNodes comment
func MapToNodes(rows *sql.Rows) []*Node {
	var ID int
	var Val, LVal, RVal string
	var Deleted bool
	var nodes []*Node
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ID, &Val, &LVal, &RVal, &Deleted)
		if err != nil {
			log.Fatal(err)
		}
		nodes = append(nodes, &Node{
			ID:      ID,
			Val:     Val,
			LVal:    LVal,
			RVal:    RVal,
			Deleted: Deleted,
		})
	}
	return nodes
}

// Save comment
func (n *Node) Save(s *Store) {
	var q string
	var id int
	if n.ID == 0 {
		q = `INSERT INTO nodes (val, l_val, r_val, deleted)
        VALUES($1, $2, $3, $4)
        RETURNING id;`
	} else {
		q = `UPDATE nodes SET val=$2, l_val=$3, r_val=$4, deleted=$5 WHERE id = $1;`
	}
	s.Save(func(tx *sql.Tx) {
		stmt, err := tx.Prepare(q)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if n.ID == 0 {
			err = stmt.QueryRow(n.Val, n.LVal, n.RVal, n.Deleted).Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
			n.ID = id
		} else {
			_, err := stmt.Exec(n.ID, n.Val, n.LVal, n.RVal, n.Deleted)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

// AllNodes comment
func AllNodes(s *Store) []*Node {
	rows, err := s.DB.Query("SELECT * FROM nodes;")
	if err != nil {
		log.Fatal(err)
	}
	return AssocNodes(MapToNodes(rows))
}

// AssocNodes comment
func AssocNodes(l []*Node) []*Node {
	r := l
	for _, n := range r {
		var k int
		if n.P == nil {
			k = SliceIndex(len(r), func(i int) bool { return r[i].Val == n.LVal || r[i].Val == n.RVal })
			if k != -1 {
				n.P = r[k]
				if r[k].LVal == n.Val {
					r[k].L = n
				} else {
					r[k].R = n
				}
			}
		}
	}
	return l
}

// SliceIndex comment
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// NullID comment
func NullID(n *Node) int {
	if n != nil {
		return n.ID
	}
	return 0
}
