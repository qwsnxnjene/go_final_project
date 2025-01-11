package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/qwsnxnjene/go_final_project/db"
	"github.com/qwsnxnjene/go_final_project/handlers"
	"github.com/qwsnxnjene/go_final_project/tests"
)

func startServer() {
	ports := ":" + strconv.Itoa(tests.Port)
	webDir := "./web"

	if p := os.Getenv("PORT"); p != "" {
		ports = ":" + p
	}

	db.CreateDB()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/task", handlers.AddTaskHandler)
	http.HandleFunc("/api/tasks", handlers.TasksHandler)

	err := http.ListenAndServe(ports, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("can't start a server: %v", err))
		return
	}
}

func main() {
	startServer()
}
