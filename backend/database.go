package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error

	db, err = sql.Open("sqlite3", "./calendar.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = createTables(); err != nil {
		log.Fatal("migrate:", err)
	}
	log.Println("Database is ready")
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id),
		date TEXT NOT NULL,
		title TEXT,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id),
		date TEXT NOT NULL,
		start_time TEXT,
		end_time TEXT,
		title TEXT, 
		description TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id),
		name TEXT NOT NULL,
		color TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS habit_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id REFERENCES habits(id),
		date TEXT,
		done INTEGER DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS goals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER REFERENCES users(id),
		period TEXT,
		date TEXT,
		title TEXT,
		done INTEGER
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

type Note struct {
	ID         int
	Date       string
	Title      string
	Content    string
	Created_at string
}

func AddNote(date, title, content string) error {
	stmt, err := db.Prepare("INSERT INTO notes (date,title,content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(date, title, content)
	return err

}

func GetsNotesByDate(date string) ([]Note, error) {
	rows, err := db.Query("SELECT id, date, title, content, created_at FROM notes WHERE date = ?", date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.Date, &n.Title, &n.Content, &n.Created_at)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil

}
