package merkle

import (
	"database/sql"
	"log"

	// Using the postgres driver
	_ "github.com/lib/pq"
)

var schema = `
	CREATE TABLE IF NOT EXISTS nodes (
		id serial primary key,
		val  varchar(64),
		l_val varchar(64),
		r_val varchar(64),
		deleted boolean
	);
`

// Store comment
type Store struct {
	DB *sql.DB
}

// Storable comment
type Storable interface{}

// AddTables comment
func (s Store) AddTables() {
	s.DB.Exec(schema)
}

// DropTables comment
func (s Store) DropTables() {
	s.DB.Exec("DROP TABLE nodes;")
}

// Exec comment
func (s Store) Exec(op func(tx *sql.Tx)) {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	op(tx)
	tx.Commit()
}
