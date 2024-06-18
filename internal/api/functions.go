package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/igortoigildin/todo_app/internal/dbs"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat cannot be empty")
	}
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("parsing error")
	}
	parsedRepeat := strings.Split(repeat, " ")
	if len(parsedRepeat) > 2 {
		return "", errors.New("incorrect repeat format")
	}
	switch parsedRepeat[0] {
	case "d":
		if len(parsedRepeat) < 2 {
			return "", errors.New("incorrect repeat format")
		}
		days, err := strconv.Atoi(parsedRepeat[1])
		if err != nil || days > 400 || days < 0 {
			return "", errors.New("incorrect repeat format")
		}
		for  {
			parsedDate = parsedDate.AddDate(0, 0, days)
			if parsedDate.Unix() > now.Unix() {
				break
			}
		}
		return parsedDate.Format("20060102"), nil
	case "y":
		if len(parsedRepeat) > 1 {
			return "", errors.New("incorrect repeat format")
		}
		for  {
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.Unix() > now.Unix() {
				break
			}
		}
		return parsedDate.Format("20060102"), nil
	default: 
		return "", errors.New("incorrect repeat format")
	}
}

func currentDate() string {
	return time.Now().Format("20060102")
}

func isDateValue(stringDate string) bool {
	stringDate = strings.Replace(stringDate, ".", "", -1)
	_, err := time.Parse("02012006", stringDate)
	return err == nil
}

func formatDate(date string) (string, error) {
	stringDate := strings.Replace(date, ".", "", -1)
	result, err := time.Parse("02012006", stringDate)
	if err != nil {
		log.Println("unable to parse date")
		return "", err
	}
	dateFormatted := result.Format("20060102")
	return dateFormatted, nil
}

func sendTaskToDB(w http.ResponseWriter, task Task) (IdStrusct, error) {
	var taskId IdStrusct
	// open and check db connection 
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
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
// creating method for error handling in JSON
func JSONError(w http.ResponseWriter, err interface{}, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
	result := make(map[string]interface{})
	result["error"] = "Ошибка"
    json.NewEncoder(w).Encode(result)
}

func checkIfTaskRequestValid(w http.ResponseWriter, task Task) bool {
	// check if title line is empty
	if task.Title == "" {
		JSONError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return false
	}
	// check if time format is valid
	dateReceived, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "Дата представлена в некорректном формате", http.StatusBadRequest)
		return false
	}
	// check repeat format
	timeNow, err := time.Parse("20060102", currentDate())
	if err != nil {
		JSONError(w, "Interanal server error", http.StatusInternalServerError)
		return false
	}
	// check if repeat format is valid and if date received before current date
	if task.Repeat != "" && dateReceived.Unix() < timeNow.Unix() {
		nextDate, err := NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			JSONError(w, "Repeat format is not valid", http.StatusBadRequest)
			return false
		}
		task.Date = nextDate
	} else if task.Repeat == "" && dateReceived.Unix() < timeNow.Unix() {
		task.Date = currentDate()
	}
	return true
}