package merkle

import (
	"database/sql"
	"log"

	// Using the postgres driver
	_ "github.com/lib/pq"
)

var schema = `
	CREATE TABLE nodes (
		id serial primary key,
		name varchar(64),
		val  varchar(64),
		l_val varchar(64),
		r_val varchar(64),
		p_id integer,
		l_id integer,
		r_id integer,
		t_id integer,
		level integer,
		epoch integer,
		path text
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

// Save comment
func (s Store) Save(op func(tx *sql.Tx)) {
	tx, err := s.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	op(tx)
	tx.Commit()
}
