package main

import (
	"log"
	"net/http"

	"github.com/igortoigildin/todo_app/config"
	"github.com/igortoigildin/todo_app/internal/api"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.LoadConfig()
	log.Fatal(http.ListenAndServe(":"+ cfg.Port, api.TaskRouter()))
}





