package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/igortoigildin/todo_app/internal/dbs"
)

type Task struct {
	Date 				string	`json:"date"`
	Title 				string	`json:"title"`
	Comment				string	`json:"comment"`
	Repeat 				string	`json:"repeat"`
	Id                  string  `json:"id"`
}
type IdStrusct struct {
	Id 					int64	`json:"id"`
}

func MyRequestHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat") 
	parsedNow, err := time.Parse("20060102", now) 
	if err != nil {
		w.Write([]byte(fmt.Sprint("%w", err)))
		return
	}
	result, err := NextDate(parsedNow, date, repeat)
	if err != nil {
		w.Write([]byte(fmt.Sprint("%w", err)))
		return
	}
	w.Write([]byte(result))
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddNewTask(w, r)
	}
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := make([]Task, 0)	
	var limit int = 30
	var rows *sql.Rows
	// open and check DB connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	search := r.URL.Query().Get("search")
	// if search request for specific date
	if search != "" && isDateValue(search) {
		date, err := formatDate(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, err = db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit;",
		sql.Named("date", date),
		sql.Named("limit", limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if search != "" && !isDateValue(search) {
	// if search reqest for specific phrase
		name := fmt.Sprintf("%%%s%%", search)
		rows, err = db.Query("SELECT * FROM scheduler WHERE title LIKE :name OR comment LIKE :name ORDER BY date LIMIT :limit;",
		sql.Named("name", name),
		sql.Named("limit", limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
	// general reqest for all tasks
		rows, err = db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit;",
		sql.Named("limit", limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Title, &task.Comment, &task.Date, &task.Repeat, &task.Id); err != nil {
			log.Println(err.Error())
			return
		}
		tasks = append(tasks, task)	
	}
	result := make(map[string]interface{})
	result["tasks"] = tasks
	resp, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Setting current date if date received is ""
	if task.Date == "" {
		task.Date = currentDate()
	}
	if !checkIfTaskRequestValid(w, task) {
		return
	}
	// sending task to db and get id
	taskId, err := sendTaskToDB(w, task)
	if err != nil {
		return 
	}
	// preparing json with task
	resp, err := json.Marshal(taskId)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}



























