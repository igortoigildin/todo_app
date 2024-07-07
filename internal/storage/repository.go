package storage

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/igortoigildin/todo_app/config"
	task "github.com/igortoigildin/todo_app/internal/model"
)

const Limit = 30 // limit rows for db queries results

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func InitPostgresDB(cfg *config.Config) *sql.DB {
	dbName := cfg.DBname
	// create DB file if not yet created
	if dbCheck(dbName) {
		file, err := os.Create(dbName)
		if err != nil {
			log.Println(err.Error())
		}
		defer file.Close()
		log.Println("db created")
	}
	// open db connection
	db, err := ConnectDB(dbName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	// create table for todo tasks
	CreateTable(db)
	return db
}

func (s *Repository) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	return err
}

func (s *Repository) GetTaskByID(id string) (task.Task, error) {
	var task task.Task
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE id = :id;", sql.Named("id", id))
	if err != nil {
		return task, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Println(err.Error())
			return task, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Println(err.Error())
	}
	return task, nil
}

func (s *Repository) GetTasksByDate(date string) ([]task.Task, error) {
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit;", sql.Named("date", date),
		sql.Named("limit", Limit))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return iterateRows(rows)
}

func (s *Repository) GetTasksByPhrase(phrase string) ([]task.Task, error) {
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE title LIKE :name OR comment LIKE :name ORDER BY date LIMIT :limit;",
		sql.Named("name", phrase),
		sql.Named("limit", Limit))
	if err != nil {
		return nil, err
	}
	return iterateRows(rows)
}

func (s *Repository) GetAllTasks() ([]task.Task, error) {
	rows, err := s.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit;", sql.Named("limit", Limit))
	if err != nil {
		return nil, err
	}
	return iterateRows(rows)
}

func iterateRows(rows *sql.Rows) ([]task.Task, error) {
	var tasks []task.Task
	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Println(err.Error())
			_ = rows.Close()
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return tasks, nil
}

func (s *Repository) CreateTask(task task.Task) (int64, error) {
	// sending received task to db
	res, err := s.db.Exec("INSERT INTO scheduler (date, comment, title, repeat) VALUES (:date, :comment, :title, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("title", task.Title),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	// getting the last inserted task
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Repository) UpdateTask(task task.Task) error {
	_, err := s.db.Exec("REPLACE INTO scheduler (id, date, comment, title, repeat) VALUES (:id, :date, :comment, :title, :repeat)",
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("comment", task.Comment),
		sql.Named("title", task.Title),
		sql.Named("repeat", task.Repeat))
	return err
}

func ConnectDB(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", name)
	if err != nil {
		log.Printf("unable to open database: %v", err)
		return nil, err
	}
	// caller ConnectDB should close DB
	err = db.Ping()
	if err != nil {
		log.Printf("unable to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}

// check if db with DBname already created
func dbCheck(DBname string) bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), DBname)
	_, err = os.Stat(dbFile)
	var install bool = false
	if err != nil {
		install = true
	}
	return install
}

func CreateTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
	id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	date char(8),
	title TEXT,
	comment TEXT,
	repeat TEXT(128)
	);`)
	if err != nil {
		log.Println("failed to create table", err)
		return
	}
	_, err = db.Exec(`CREATE INDEX index_date
	ON SCHEDULER(date)
	;`)
	if err != nil {
		log.Println("failed to create index", err)
		return
	}
	log.Println("table created")
}
