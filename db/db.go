package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func CreateTable(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("[db.CreateTable]: can't open database: %v", err)
	}
	createQuery := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "19700101",
		title VARCHAR(32) NOT NULL DEFAULT "",
		comment VARCHAR(256) NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX date_scheduler on scheduler (date);
	`

	_, err = db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("[db.CreateTable]: can't create table: %v", err)
	}
	return nil
}
