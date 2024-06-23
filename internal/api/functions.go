package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
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

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if pass is present
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwt string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			valid := true
			err = verifyToken(jwt)
			if err != nil {
				valid = false
			}
			if !valid{
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func verifyToken(tokenSting string) (error) {
	token, err := jwt.Parse(tokenSting, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return  err
	}
	if !token.Valid {
		return  fmt.Errorf("invalid token")
	}
	return err
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
// method for error handling in JSON
func JSONError(w http.ResponseWriter, err interface{}, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
	result := make(map[string]interface{})
	result["error"] = "Ошибка"
    json.NewEncoder(w).Encode(result)
}

func idValid(id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	db, err := dbs.ConnectDB("scheduler.db")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM scheduler WHERE id= :id;",
	sql.Named("id", id))
	if err != nil {
		return false, err
	}
	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Println(err.Error())
			return false, err
		}
		tasks = append(tasks, task)	
	}
	return len(tasks) != 0, nil
}