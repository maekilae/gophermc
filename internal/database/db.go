package db

import (
	"database/sql"
	"github.com/cockroachdb/pebble"
	_ "github.com/mattn/go-sqlite3"
)

type McDB struct {
	wdb *pebble.DB
	udb *sql.DB
}

func OpenDB(wfile, ufile string) (*McDB, error) {
	var db McDB
	err := db.initPebble(wfile)
	if err != nil {
		return &db, err
	}
	err = db.initSql(ufile)

	return &db, err
}
func (db *McDB) initPebble(file string) error {
	var err error
	db.wdb, err = pebble.Open("data/"+file, &pebble.Options{})
	if err != nil {
		return err
	}
	return nil
}

func (db *McDB) initSql(file string) error {
	var err error
	db.udb, err = sql.Open("sqlite3", "data/"+file+".sqlite")
	if err != nil {
		return err
	}
	err = db.initSqlTables()
	return err
}

func (db *McDB) initSqlTables() error {
	tx, err := db.udb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
	    uuid INTEGER NOT NULL PRIMARY KEY,
	    username TEXT
	);`
	tx.Exec(sqlStmt)
	sqlStmt = `
	CREATE TABLE IF NOT EXISTS inventories (
	    player_uuid INTEGER NOT NULL PRIMARY KEY,
	    data BLOB
	    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    	    FOREIGN KEY(player_uuid) REFERENCES users(uuid)
	);`
	tx.Exec(sqlStmt)

	tx.Commit()

	return nil
}
