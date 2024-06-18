package dbs

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)


func CreateDB() {
	// create DB file if not yet created
	if dbCheck("scheduler.db") {
		file, err := os.Create("scheduler.db")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()
		log.Println("db created")
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	defer db.Close()
	// check db connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	// creating table for todo tasks
	CreateTable(db)
}

func ConnectDB(name string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", name)
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
		return nil, err
	}
	// caller ConnectDB should close DB
	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
		return nil, err
	}
	return db, nil
}


func dbCheck(DBname string) bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
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
	date TEXT,
	title TEXT,
	comment TEXT,
	repeat TEXT
	);`)
	if err != nil {
		log.Fatal("failed to create table", err)
	}
	_, err = db.Exec(`CREATE INDEX index_date
	ON SCHEDULER(date)
	;`)
	if err != nil {
		log.Println("failed to create index", err)
	}
	log.Println("table created")
}

