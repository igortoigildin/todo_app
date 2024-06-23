package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/igortoigildin/todo_app/internal/dbs"
)
type Task struct {
	Id                  string  `json:"id"`
	Date 				string	`json:"date"`
	Title 				string	`json:"title"`
	Comment				string	`json:"comment"`
	Repeat 				string	`json:"repeat"`
}
type IdStrusct struct {
	Id 					int64	`json:"id"`
}
type passStruct struct {
	Password 			string  `json:"password"`
}
type token struct {
	Token 				string `json:"token"`
}

func TaskRouter() chi.Router {
	dbs.CreateDB()
	r := chi.NewRouter()
	r.Route("/api/task", func(r chi.Router) {
		r.Post("/", auth(AddNewTask))
		r.Get("/", auth(GetTask))
		r.Put("/", auth(ChangeTask))
		r.Delete("/", auth(DeleteTask))
	})
	r.Get("/api/nextdate", MyRequestHandler)
	r.Get("/api/tasks", auth(GetTasksHandler))
	r.Post("/api/task/done", auth(TaskDone))
	r.Post("/api/sign", SigninHandler)
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	fmt.Println("Starting the server on :7540...")
	return r
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	var passStruct passStruct
	err := json.NewDecoder(r.Body).Decode(&passStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	envPass := os.Getenv("TODO_PASSWORD")
	if passStruct.Password == envPass {
		var secretKey = []byte("your-secret-key")
		myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		signedToken, err := myToken.SignedString(secretKey)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusBadRequest)
		}
		var token token
		token.Token = signedToken
		resp, err := json.Marshal(token)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	} else {
		JSONError(w, "Неверный пароль", http.StatusInternalServerError)
		return
	}
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

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	// open and check db connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()
	check, err := idValid(id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !check {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var blank []int
	resp, err := json.Marshal(blank)
	if err != nil {
		JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	// open and check db connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()
	check, err := idValid(id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !check {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	rows, err := db.Query("SELECT * FROM scheduler WHERE id = :id;",
	sql.Named("id", id))
	if err != nil {
		JSONError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var task Task
	for rows.Next() {	
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat ); err != nil {
			log.Println(err.Error())
			return
		}
	}
	switch task.Repeat {
	case "":
		_, err = db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var blank []int
		resp, err := json.Marshal(blank)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	default:	
		timeNow, err := time.Parse("20060102", currentDate())
		if err != nil {
			JSONError(w, "Interanal server error", http.StatusInternalServerError)
			return
		}
		date, err := NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			JSONError(w, "Interanal server error", http.StatusInternalServerError)
			return
		}
		task.Date = date
		err = replaceTaskDB(w, task)
		if err != nil {
			JSONError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		}
		var blank []int
		resp, err := json.Marshal(blank)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func ChangeTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	check, err := idValid(task.Id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !check {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	// Setting current date if date received is ""
	if task.Date == "" {
		task.Date = currentDate()
	}
	if !checkIfTaskRequestValid(w, task) {
		return
	}
	dateReceived, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "Дата представлена в некорректном формате", http.StatusBadRequest)
		return
	}
	timeNow, err := time.Parse("20060102", currentDate())
	if err != nil {
		JSONError(w, "Interanal server error", http.StatusInternalServerError)
		return
	}
	// check if date in request is before current date or empty
	if task.Date != "" && dateReceived.Unix() < timeNow.Unix() && task.Repeat != "" {
		nextDate, err := NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			JSONError(w, "Repeat format is not valid", http.StatusBadRequest)
			return
		}
		task.Date = nextDate
	} else if task.Date == "" || task.Repeat == "" {
		task.Date = currentDate()
	}
	// sending task to db 
	err = replaceTaskDB(w, task)
	if err != nil {
		JSONError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}
	var emptyTask Task
	resp, err := json.Marshal(emptyTask)
	if err != nil {
		JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	// open and check DB connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM scheduler WHERE id = :id;",
	sql.Named("id", id))
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var task Task
	for rows.Next() {	
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat ); err != nil {
			log.Println(err.Error())
			return
		}
	}
	if task.Id == "" {
		JSONError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
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
	defer db.Close()
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
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
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
	dateReceived, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "Дата представлена в некорректном формате", http.StatusBadRequest)
	}
	timeNow, err := time.Parse("20060102", currentDate())
	if err != nil {
		JSONError(w, "Interanal server error", http.StatusInternalServerError)
	}
	// check if date in request is before current date or empty
	if task.Date != "" && dateReceived.Unix() < timeNow.Unix() && task.Repeat != "" {
		nextDate, err := NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			JSONError(w, "Repeat format is not valid", http.StatusBadRequest)
		}
		task.Date = nextDate
	} else if task.Date != "" && dateReceived.Unix() < timeNow.Unix() && task.Repeat == "" {
		task.Date = currentDate()
	} else if task.Date == "" && task.Repeat == "" {
		task.Date = currentDate()
	}
	// sending task to db and get id
	taskId, err := sendTaskToDB(w, task)
	if err != nil {
		log.Println(err.Error())
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

func sendTaskToDB(w http.ResponseWriter, task Task) (IdStrusct, error) {
	var taskId IdStrusct
	// open and check db connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()
	// sending received task to db
	res, err := db.Exec("INSERT INTO scheduler (date, comment, title, repeat) VALUES (:date, :comment, :title, :repeat)",
	sql.Named("date", task.Date),
	sql.Named("comment", task.Comment),
	sql.Named("title", task.Title),
	sql.Named("repeat", task.Repeat))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return taskId, err
	}
	// getting the last inserted task
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return taskId, err
	}
	taskId.Id = id
	return taskId, nil
}

func replaceTaskDB(w http.ResponseWriter, task Task) (error) {
	// open and check db connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()
	// sending received task to db
	_, err = db.Exec("REPLACE INTO scheduler (id, date, comment, title, repeat) VALUES (:id, :date, :comment, :title, :repeat)",
	sql.Named("id", task.Id),
	sql.Named("date", task.Date),
	sql.Named("comment", task.Comment),
	sql.Named("title", task.Title),
	sql.Named("repeat", task.Repeat))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

func checkIfTaskRequestValid(w http.ResponseWriter, task Task) bool {
	// check if title line is empty
	if task.Title == "" {
		JSONError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return false
	}
	// check if time format is valid
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "Дата представлена в некорректном формате", http.StatusBadRequest)
		return false
	}
	// check repeat format
	if task.Repeat != "" {
		timeNow, err := time.Parse("20060102", currentDate())
		if err != nil {
			JSONError(w, "Interanal server error", http.StatusInternalServerError)
			return false
		}
		_, err = NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			JSONError(w, "Repeat format is not valid", http.StatusBadRequest)
			return false
		}
	}
	return true
}



























