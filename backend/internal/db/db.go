package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"music-player/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

var SqliteDB *sql.DB

type Song struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	LengthSec int     `json:"lengthSec"`
	FilePath  string  `json:"filePath"`
	CoverPath *string `json:"coverPath,omitempty"`
}

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

	_, err = SqliteDB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func AddSong(song Song) error {
	stmt, err := SqliteDB.Prepare(`
	INSERT INTO songs (title, author, lengthSec, filePath, coverPath)
	VALUES (?, ?, ?, ?, ?) `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(song.Title, song.Author, song.LengthSec, song.FilePath, song.CoverPath)
	if err != nil {
		return err
	}

	song.ID, _ = res.LastInsertId()
	return nil
}

func RemoveSong(id int64) error {
	_, err := SqliteDB.Exec(`DELETE FROM songs WHERE id = ?`, id)
	return err
}

func GetSongByID(id int64) (*Song, error) {
	var s Song
	err := SqliteDB.QueryRow(`SELECT id, title, author, lengthSec, filePath, coverPath FROM songs WHERE id = ?`, id).
		Scan(&s.ID, &s.Title, &s.Author, &s.LengthSec, &s.FilePath, &s.CoverPath)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSongByTitle(title string) (*Song, error) {
	var s Song
	err := SqliteDB.QueryRow(`SELECT id, title, author, lengthSec, filePath, coverPath FROM songs WHERE title = ?`, title).
		Scan(&s.ID, &s.Title, &s.Author, &s.LengthSec, &s.FilePath, &s.CoverPath)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetAllSongs() ([]Song, error) {
	rows, err := SqliteDB.Query(`SELECT id, title, author, lengthSec, filePath, coverPath FROM songs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []Song
	for rows.Next() {
		var s Song
		err := rows.Scan(&s.ID, &s.Title, &s.Author, &s.LengthSec, &s.FilePath, &s.CoverPath)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, s)
	}

	return tracks, nil
}

func GetAllSongPaths() []string {
	var paths []string

	err := filepath.WalkDir(os.Getenv("AUDIO_PATH"), func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dirEntry.IsDir() {
			paths = append(paths, strings.TrimPrefix(path, "music/"))
		}
		return nil
	})
	utils.Check(err)
	return paths
}

func GetNotAddedSongPaths() []string {
	paths := GetAllSongPaths()
	allSongs, err := GetAllSongs()
	if err != nil {
		log.Println("Failed to get not added song paths")
		return make([]string, 0)
	}

	addedPaths := make(map[string]struct{})
	for _, song := range allSongs {
		addedPaths[song.FilePath] = struct{}{}
	}

	var notAdded []string
	for _, path := range paths {
		if _, exists := addedPaths[path]; !exists {
			notAdded = append(notAdded, path)
		}
	}

	return notAdded
}

func SaveFile(fileName string, fileData []byte) error {
	audioPath := os.Getenv("AUDIO_PATH")
	_, err := os.Stat(audioPath + fileName)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("file already exists")
	}
	os.WriteFile(audioPath+fileName, fileData, 0644)
	return nil
}
