package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat cannot be empty")
	}
	parsedDate, err := time.Parse(yymmdd, date)
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
		for {
			parsedDate = parsedDate.AddDate(0, 0, days)
			if parsedDate.Unix() > now.Unix() {
				break
			}
		}
		return parsedDate.Format(yymmdd), nil
	case "y":
		if len(parsedRepeat) > 1 {
			return "", errors.New("incorrect repeat format")
		}
		for {
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.Unix() > now.Unix() {
				break
			}
		}
		return parsedDate.Format(yymmdd), nil
	default:
		return "", errors.New("incorrect repeat format")
	}
}
