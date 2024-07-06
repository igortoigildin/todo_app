package storage

import task "github.com/igortoigildin/todo_app/internal/model"

type Storage interface {
	GetTaskByID(id string) (task.Task, error)
	GetTasksByDate(date string) ([]task.Task, error)
	GetTasksByPhrase(phrase string) ([]task.Task, error)
	UpdateTask(task task.Task) error
	CreateTask(task task.Task) (int64, error) ///// done
	GetAllTasks() ([]task.Task, error)
	DeleteTask(id string) error
}
