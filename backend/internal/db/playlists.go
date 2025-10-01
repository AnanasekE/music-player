package db

import (
	"database/sql"
	"errors"
)

type Playlist struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	Songs       []Song  `json:"songs,omitempty"` // filled when loading
}

func CreatePlaylist(name string, description *string) (int64, error) {
	stmt, err := SqliteDB.Prepare(`INSERT INTO playlists (name, description) VALUES (?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func AddSongToPlaylist(playlistID, songID int64, position int) error {
	_, err := SqliteDB.Exec(
		`INSERT INTO playlist_songs (playlist_id, song_id, position) VALUES (?, ?, ?)`,
		playlistID, songID, position,
	)
	return err
}

func GetPlaylistWithSongs(playlistID int64) (*Playlist, error) {
	var p Playlist
	err := SqliteDB.QueryRow(
		`SELECT id, name, description, created_at FROM playlists WHERE id = ?`,
		playlistID,
	).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := SqliteDB.Query(`
        SELECT s.id, s.title, s.author, s.lengthSec, s.filePath, s.coverPath
        FROM songs s
        INNER JOIN playlist_songs ps ON s.id = ps.song_id
        WHERE ps.playlist_id = ?
        ORDER BY ps.position ASC`, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Song
		err := rows.Scan(&s.ID, &s.Title, &s.Author, &s.LengthSec, &s.FilePath, &s.CoverPath)
		if err != nil {
			return nil, err
		}
		p.Songs = append(p.Songs, s)
	}

	return &p, nil
}
func GetPlaylistByID(id int64) (*Playlist, error) {
	var p Playlist
	err := SqliteDB.QueryRow(`
		SELECT id, name, description, created_at
		FROM playlists
		WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Playlist not found") // playlist not found
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}
