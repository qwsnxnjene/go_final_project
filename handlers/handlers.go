package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/qwsnxnjene/go_final_project/db"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, idToAdd)))
	}
}

func TasksHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks := []Task{}
	dB, err := db.OpenSql()

	if err != nil {
		rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
		return
	}
	defer dB.Close()
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`
	rows, err := dB.Query(query)
	if err != nil {
		rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t Task

		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка работы с БД %v", err).Error())))
			return
		}
		tasks = append(tasks, t)
	}

	var jsonTask struct {
		Tasks []Task `json:"tasks"`
	}

	jsonTask.Tasks = tasks
	data, err := json.Marshal(jsonTask)
	if err != nil {
		fmt.Println("goal")
		rw.Write([]byte(fmt.Sprintf(`{"error":"%v"}`, fmt.Errorf("ошибка сериализации %v", err).Error())))
		return
	}

	fmt.Println(string(data))

	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}
