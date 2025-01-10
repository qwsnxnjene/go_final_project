package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/qwsnxnjene/go_final_project/db"
	"github.com/qwsnxnjene/go_final_project/tests"

	_ "modernc.org/sqlite"
)

func startServer() {
	ports := ":" + strconv.Itoa(tests.Port)
	webDir := "./web"

	if p := os.Getenv("PORT"); p != "" {
		ports = ":" + p
	}

	db.CreateDB()

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", AddTaskHandler)

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

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTaskHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

		decoder := json.NewDecoder(r.Body)
		var t Task
		err := decoder.Decode(&t)
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`,
				fmt.Errorf("ошибка десериализации %v", err).Error())))
			return
		}
		defer r.Body.Close()

		date, title, comment, repeat := t.Date, t.Title, t.Comment, t.Repeat

		if len(title) == 0 {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error":"заголовок не может быть пустым"}`))
			return
		}

		if len(date) == 0 {
			date = time.Now().Format("20060102")
		}

		dateTo, err := time.Parse("20060102", date)
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(`{"error":"некорректная дата"}`))
			return
		}

		newDate := ""
		if len(repeat) > 0 {
			newDate, err = NextDate(time.Now(), date, repeat)
			if err != nil {
				//rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, err.Error())))
				return
			}
			if timeDiff(time.Now(), dateTo) {
				date = newDate
			}
		}

		//если дата меньше сегодняшнего числа
		if timeDiff(time.Now(), dateTo) {
			fmt.Println(date)
			date = time.Now().Format("20060102")
		}

		database, err := db.OpenSql()
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
			return
		}

		defer database.Close()

		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
		res, err := database.Exec(query,
			sql.Named("date", date),
			sql.Named("title", title),
			sql.Named("comment", comment),
			sql.Named("repeat", repeat),
		)
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
			return
		}
		idToAdd, err := res.LastInsertId()
		if err != nil {
			//rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
			return
		}

		fmt.Println("okay")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, idToAdd)))
	}
}

func main() {
	startServer()
}
