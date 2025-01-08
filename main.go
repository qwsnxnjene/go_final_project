package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/qwsnxnjene/go_final_project/db"
	"github.com/qwsnxnjene/go_final_project/tests"
)

func startServer() {
	ports := ":" + strconv.Itoa(tests.Port)
	webDir := "./web"

	if p := os.Getenv("PORT"); p != "" {
		ports = ":" + p
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(ports, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("can't start a server: %v", err))
		return
	}
}

func createDB() {
	dbPath := tests.DBFile
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		dbPath = path
	}
	_, err := os.Stat(dbPath)

	var install bool
	if err != nil {
		install = true
	}
	if install {
		db.CreateTable(dbPath)
	}
}

func main() {
	startServer()
	createDB()
}
