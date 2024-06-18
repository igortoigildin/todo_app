package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
// creating method for error handling in JSON
func JSONError(w http.ResponseWriter, err interface{}, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
	result := make(map[string]interface{})
	result["error"] = "Ошибка"
    json.NewEncoder(w).Encode(result)
}

