package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/igortoigildin/todo_app/config"
	"github.com/igortoigildin/todo_app/internal/model"
	storage "github.com/igortoigildin/todo_app/internal/storage"
)

type TodosHandler struct {
	repo storage.Storage
}

func NewTodosHandler(repo storage.Storage) TodosHandler {
	return TodosHandler{
		repo: repo,
	}
}

func (h TodosHandler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	var passStruct model.PassStruct
	err := json.NewDecoder(r.Body).Decode(&passStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signedToken, err := checkPass(passStruct)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusBadRequest)
		return
	}
	if signedToken == "" {
		JSONError(w, "Неверный пароль", http.StatusInternalServerError)
		return
	}
	var token model.Token
	token.Token = signedToken
	resp, err := json.Marshal(token)
	if err != nil {
		JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func (h TodosHandler) RequestHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")
	parsedNow, err := time.Parse(yymmdd, now)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprint("%w", err)))
		return
	}
	result, err := NextDate(parsedNow, date, repeat)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprint("%w", err)))
		return
	}
	_, _ = w.Write([]byte(result))
}

func (h TodosHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if task.Id == "" {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	err = h.repo.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blank := make(map[string]interface{}, 0)
	resp, err := json.Marshal(blank)
	if err != nil {
		JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func (h TodosHandler) TaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if task.Id == "" {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	switch task.Repeat {
	case "":
		err = h.repo.DeleteTask(id)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		blank := make(map[string]interface{}, 0)
		resp, err := json.Marshal(blank)
		if err != nil {
			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	default:
		timeNow, err := time.Parse(yymmdd, currentDate())
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
		err = h.repo.UpdateTask(task)
		if err != nil {
			JSONError(w, "Задача не найдена", http.StatusInternalServerError)
			return
		}
		blank := make(map[string]interface{}, 0)
		resp, err := json.Marshal(blank)
		if err != nil {

			JSONError(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	}
}

func (h TodosHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check if task received is valid
	task, err = validateTask(w, task)
	if err != nil {
		return
	}
	// check if received task exists in DB
	tempTask, err := h.repo.GetTaskByID(task.Id)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if tempTask.Id == "" {
		JSONError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	// sending task to db
	err = h.repo.UpdateTask(task)
	if err != nil {
		JSONError(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}
	var emptyTask model.Task
	resp, err := json.Marshal(emptyTask)
	if err != nil {
		JSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func (h TodosHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		JSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		JSONError(w, "Задача не найдена", http.StatusInternalServerError)
		return
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
	_, _ = w.Write(resp)
}

func (h TodosHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	searchValue := r.URL.Query().Get("search")
	result := performSearch(w, searchValue, h)
	resp, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (h TodosHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err = validateTask(w, task)
	if err != nil {
		return
	}
	// sending task to db and get id
	var taskId model.IdStrusct
	id, err := h.repo.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	taskId.Id = id
	// preparing json with task
	resp, err := json.Marshal(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func auth(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if pass is present
		pass := cfg.Pass
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
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
