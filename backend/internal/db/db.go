package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var SqliteDB *sql.DB

func InitDB() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "songs.db" // default
	}

	var err error
	SqliteDB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS songs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		lengthSec INTEGER NOT NULL,
		filePath TEXT NOT NULL UNIQUE,
		coverPath TEXT
	);`

	createPlaylists := `
CREATE TABLE IF NOT EXISTS playlists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	createPlaylistSongs := `
CREATE TABLE IF NOT EXISTS playlist_songs (
    playlist_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    position INTEGER,
    PRIMARY KEY (playlist_id, song_id),
    FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);`

	_, err = SqliteDB.Exec(createPlaylists)
	if err != nil {
		log.Fatal(err)
	}
	_, err = SqliteDB.Exec(createPlaylistSongs)
	if err != nil {
		log.Fatal(err)
	}
	_, err = SqliteDB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
