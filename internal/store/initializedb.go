package store

import (
	"database/sql"
	"newforum/config"

	"golang.org/x/crypto/bcrypt"
)

func InitializeDB(config *config.Config) (*sql.DB, error){
	db, err := sql.Open(config.DB.Driver, config.DB.Dsn)
	if err != nil {
		return nil, err
	} 
	
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil{
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	

	if err = CreateTables(db); err != nil {
		return nil, err
	}
	
	return db, nil

}

func CreateTables(db *sql.DB) error {

	stmt := 
	`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		hashed_password CHAR(60) NOT NULL,
		created DATETIME NOT NULL,
		token TEXT,
		role INTEGER DEFAULT 3 NOT NULL,
    requested INTEGER DEFAULT 0 NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS snippets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		image VARCHAR(100) NOT NULL,
		created DATETIME NOT NULL,
		reported INTEGER DEFAULT 0 NOT NULL,
    approved INTEGER DEFAULT 0 NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) 
	);
	
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		snippet_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created DATETIME NOT NULL,
		reported INTEGER DEFAULT 0 NOT NULL,
		FOREIGN KEY (snippet_id) REFERENCES snippets(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS reactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		snippet_id INTEGER,
		comment_id INTEGER,
		user_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		UNIQUE (snippet_id, comment_id, user_id, type),
		FOREIGN KEY (snippet_id) REFERENCES snippets(id) ON DELETE CASCADE,
		FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY,
		receiver_id INTEGER NOT NULL,
		sender_id INTEGER NOT NULL,
		author_name TEXT NOT NULL,
		action_type TEXT NOT NULL,
		snippet_id INTEGER NOT NULL,
		comment_id INTEGER DEFAULT 0,
		timestamp DATETIME NOT NULL,
		is_read INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		snippet_id INTEGER NOT NULL,
		category TEXT NOT NULL,
		FOREIGN KEY (snippet_id) REFERENCES snippets(id) ON DELETE CASCADE
	);
	`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Cheburek"), 12)
	if err != nil {
		return err
	}

	stmt = `
	INSERT OR IGNORE INTO users (name, email, hashed_password, created, role, requested)
	VALUES (?, ?, ?, datetime('now','+6 hours'), 1, 0);
	`
	
	_, err = db.Exec(stmt, "Meduza", "meduza@gmail.com", string(hashedPassword))
	
	return err
}
