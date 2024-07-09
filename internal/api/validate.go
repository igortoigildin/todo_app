package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	model "github.com/igortoigildin/todo_app/internal/model"
)

const (
	jwtSecret = "your-secret-key" // string for JWT secret
)

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
	// Setting current date if date received is ""
	if task.Date == "" {
		task.Date = currentDate()
	}
	if !validateTaskRequest(w, task) {
		return task, fmt.Errorf("task request not valid")
	}
	timeNow, err := time.Parse(yymmdd, currentDate())
	if err != nil {
		JSONError(w, "Interanal server error", http.StatusInternalServerError)
		return task, err
	}
	dateReceived, _ := time.Parse(yymmdd, task.Date)
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
	} else if task.Date == "" && task.Repeat == "" {
		task.Date = currentDate()
	}
	return task, nil
}

func validateTaskRequest(w http.ResponseWriter, task model.Task) bool {
	// check if title line is empty
	if task.Title == "" {
		JSONError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return false
	}
	// check if time format is valid
	_, err := time.Parse(yymmdd, task.Date)
	if err != nil {
		JSONError(w, "Дата представлена в некорректном формате", http.StatusBadRequest)
		return false
	}
	// check repeat format
	if task.Repeat != "" {
		timeNow, err := time.Parse(yymmdd, currentDate())
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
