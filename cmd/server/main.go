package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/igortoigildin/todo_app/config"
	"github.com/igortoigildin/todo_app/internal/api"
	"github.com/igortoigildin/todo_app/internal/dbs"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.LoadConfig()
	dbs.CreateDB()
	r := chi.NewRouter()
	r.Get("/api/nextdate", api.MyRequestHandler)
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	fmt.Println("Starting the server on :7540...")
	log.Fatal(http.ListenAndServe(":"+ cfg.Port, r))
}





