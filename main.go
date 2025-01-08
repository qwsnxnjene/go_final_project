package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	http.HandleFunc("/api/nextdate", NextDateHandler)

	err := http.ListenAndServe(ports, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("can't start a server: %v", err))
		return
	}
}

func NextDateHandler(rw http.ResponseWriter, r *http.Request) {
	date, repeat := r.FormValue("date"), r.FormValue("repeat")
	now, err := time.Parse("20060102", r.FormValue("now"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}
	newDate, err := NextDate(now, date, repeat)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	} else {
		rw.Write([]byte(newDate))
	}
}

func createDB() {
	dbPath := tests.DBFile
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		dbPath = path
	}
	_, err := os.Stat(dbPath)

	// if there is no DB, create one
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
