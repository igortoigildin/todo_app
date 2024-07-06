package api

import (
	"fmt"
	"net/http"

	"github.com/igortoigildin/todo_app/internal/model"
)


func performSearch(w http.ResponseWriter, searchValue string, h TodosHandler) map[string]interface{} {
	var tasks []model.Task
	var err error
	// if search request for specific date
	if searchValue != "" && isDateValue(searchValue) {
		date, err := formatDate(searchValue)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return nil
		}
		tasks, err = h.repo.GetTasksByDate(date)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return nil
		}
	} else if searchValue != "" && !isDateValue(searchValue) {
		// if search reqest for specific phrase
		name := fmt.Sprintf("%%%s%%", searchValue)
		tasks, err = h.repo.GetTasksByPhrase(name)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return nil
		}
	} else {
		// general reqest for all tasks	
		tasks, err = h.repo.GetAllTasks()
		if err != nil {
			JSONError(w, "Bad request", http.StatusInternalServerError)
			return nil
		}
	}
	temp := make([]model.Task, 0)
	temp = append(temp, tasks...)
	result := make(map[string]interface{})
	result["tasks"] = temp
	return result
}