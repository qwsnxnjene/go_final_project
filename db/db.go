package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/qwsnxnjene/go_final_project/tests"
	_ "modernc.org/sqlite"
)

var ActualDbPath string

// CreateDB() получает путь к файлу с базой данных и, если необходимо, создаёт таблицу
func CreateDB() {
	dbPath := tests.DBFile
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		dbPath = path
	}
	_, err := os.Stat(dbPath)

	ActualDbPath = dbPath

	// if there is no DB, create one
	var install bool
	if err != nil {
		install = true
	}
	if install {
		CreateTable(dbPath)
	}
}

// CreateTable(path) создаёт таблицу в файле, указанном в path
// возвращает nil при успешном создании таблицы, иначе ошибку
func CreateTable(path string) error {
	db, err := sql.Open("sqlite", path)

	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("[db.CreateTable]: can't open database: %v", err)
	}
	createQuery := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "19700101",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment VARCHAR(256) NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX date_scheduler on scheduler (date);
	`

	_, err = db.Exec(createQuery)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("[db.CreateTable]: can't create table: %v", err)
	}
	return nil
}

// OpenSql() возвращает таблицу для работы с ней
func OpenSql() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ActualDbPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("[db.CreateTable]: can't open database: %v", err)
	}

	return db, nil
}
