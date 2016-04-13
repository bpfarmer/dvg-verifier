package merkle

import (
	"database/sql"
	"log"
)

// MapToNodes comment
func MapToNodes(rows *sql.Rows) []*Node {
	var ID, PID, LID, RID, TID int
	var Level, Epoch uint
	var Name, Val, LVal, RVal []byte
	var Path string
	var nodes []*Node
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ID, &Name, &Val, &LVal, &RVal, &PID, &LID, &RID, &TID, &Level, &Epoch, &Path)
		if err != nil {
			log.Fatal(err)
		}
		nodes = append(nodes, &Node{
			ID:    ID,
			Name:  Name,
			Val:   Val,
			LVal:  LVal,
			RVal:  RVal,
			PID:   PID,
			LID:   LID,
			RID:   RID,
			TID:   TID,
			Level: Level,
			Epoch: Epoch,
		})
	}
	return nodes
}

// Save comment
func (n *Node) Save(s *Store) {
	var q string
	var id int
	if n.ID == 0 {
		q = `INSERT INTO nodes (name, val, l_val, r_val, p_id, l_id, r_id, t_id, level, epoch, path)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id;`
	} else {
		q = `UPDATE nodes SET name=$1, val=$2, l_val=$3, r_val=$4, p_id=$5,
        l_id=$6, r_id=$7, t_id=$8, level=$9, epoch=$10, path=$11 WHERE id = $12;`
	}
	s.Save(func(tx *sql.Tx) {
		stmt, err := tx.Prepare(q)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if n.ID == 0 {
			err = stmt.QueryRow(n.Name, n.Val, n.LVal, n.RVal, NullID(n.Parent), NullID(n.L), NullID(n.R), 0, n.Level, n.Epoch, n.Path()).Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
			n.ID = id
		} else {
			_, err := stmt.Exec(n.Name, n.Val, n.LVal, n.RVal, NullID(n.Parent), NullID(n.L), NullID(n.R), 0, n.Level, n.Epoch, n.Path(), n.ID)
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
		if n.Parent == nil {
			k = SliceIndex(len(r), func(i int) bool { return r[i].ID == n.PID })
			if k != -1 {
				n.Parent = r[k]
				if r[k].LID == n.ID {
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
