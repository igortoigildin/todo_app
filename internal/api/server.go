package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/igortoigildin/todo_app/config"
	"github.com/igortoigildin/todo_app/internal/storage"
)

func RunServer() {
	cfg := config.LoadConfig()
	db := storage.InitPostgresDB(cfg)
	Repository := storage.NewRepository(db)
	handler := NewTodosHandler(Repository)
	r := chi.NewRouter()
	r.Route("/api/task", func(r chi.Router) {
		r.Post("/", auth(handler.CreateTask, cfg))
		r.Get("/", auth(handler.GetTaskByID, cfg))
		r.Put("/", auth(handler.UpdateTask, cfg))
		r.Delete("/", auth(handler.DeleteTask, cfg))
	})
	r.Get("/api/nextdate", handler.RequestHandler)
	r.Get("/api/tasks", auth(handler.GetTasksHandler, cfg))
	r.Post("/api/task/done", auth(handler.TaskDone, cfg))
	r.Post("/api/sign", handler.SigninHandler)
	r.Handle("/*", http.FileServer(http.Dir("./web")))

	log.Printf("Starting server on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
