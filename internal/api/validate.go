package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	model "github.com/igortoigildin/todo_app/internal/model"
	storage "github.com/igortoigildin/todo_app/internal/storage"
)

const (
	jwtSecret = "your-secret-key" // string for JWT secret
)

func idValid(id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	db, err := storage.ConnectDB("scheduler.db")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM scheduler WHERE id= :id;", sql.Named("id", id))
	if err != nil {
		return false, err
	}
	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Println(err.Error())
			return false, err
		}
		tasks = append(tasks, task)
	}
	return len(tasks) != 0, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return err
}

func checkIfTaskRequestValid(w http.ResponseWriter, task model.Task) bool {
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

func checkPass(passStruct model.PassStruct) (string, error) {
	var signedToken string
	var err error
	envPass := os.Getenv("TODO_PASSWORD")
	if passStruct.Password == envPass {
		var secretKey = []byte(jwtSecret)
		myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		signedToken, err = myToken.SignedString(secretKey)
		if err != nil {
			return "", err
		}
	}
	return signedToken, nil
}

func validateTask(w http.ResponseWriter, task model.Task) (model.Task, error) {
	check, err := idValid(task.Id)
	if err != nil && task.Id != "" {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return task, err
	}
	if !check {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return task, fmt.Errorf("not found")
	}
	// Setting current date if date received is ""
	if task.Date == "" {
		task.Date = currentDate()
	}
	if !checkIfTaskRequestValid(w, task) {
		return task, fmt.Errorf("task request not valid")
	}
	timeNow, err := time.Parse("20060102", currentDate())
	if err != nil {
		JSONError(w, "Interanal server error", http.StatusInternalServerError)
		return task, err
	}
	dateReceived, _ := time.Parse("20060102", task.Date)
	// check if date is correct
	task, err = validateDate(task, timeNow, dateReceived)
	if err != nil {
		JSONError(w, "repeat format is not valid", http.StatusBadRequest)
		return task, err
	}
	return task, nil
}

func validateDate(task model.Task, timeNow time.Time, dateReceived time.Time) (model.Task, error) {
	// check if date in request is before current date or empty
	if task.Date != "" && dateReceived.Unix() < timeNow.Unix() && task.Repeat != "" {
		nextDate, err := NextDate(timeNow, task.Date, task.Repeat)
		if err != nil {
			return task, fmt.Errorf("repeat format is not valid")
		}
		task.Date = nextDate
	} else if task.Date != "" && dateReceived.Unix() < timeNow.Unix() && task.Repeat == "" {
		task.Date = currentDate()
	} else if task.Date == "" || task.Repeat == "" {
		task.Date = currentDate()
	}
	return task, nil
}

func validateNewTask(w http.ResponseWriter, task model.Task) (model.Task, error) {
	// Setting current date if date received is ""
	if task.Date == "" {
		task.Date = currentDate()
	}
	if !checkIfTaskRequestValid(w, task) {
		return task, fmt.Errorf("task request not valid")
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
	return task, nil
}
