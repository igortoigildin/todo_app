package api

import (
	"fmt"
	"net/http"
	"time"
)

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