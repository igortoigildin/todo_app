package api

import (
	"log"
	"strings"
	"time"
)

const (
	yymmdd = "20060102" // date format constant
	ddmmyy = "02012006" // date format constant
)

func currentDate() string {
	return time.Now().Format(yymmdd)
}

func isDateValue(stringDate string) bool {
	stringDate = strings.Replace(stringDate, ".", "", -1)
	_, err := time.Parse(ddmmyy, stringDate)
	return err == nil
}

func formatDate(date string) (string, error) {
	stringDate := strings.Replace(date, ".", "", -1)
	result, err := time.Parse(ddmmyy, stringDate)
	if err != nil {
		log.Println("unable to parse date")
		return "", err
	}
	dateFormatted := result.Format(yymmdd)
	return dateFormatted, nil
}
