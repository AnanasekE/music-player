package web

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"music-player/internal/db"
	"music-player/internal/utils"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
)

func trackHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTracks(w)
	case http.MethodPost:
		handleCreateTrack(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTracksByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Track ID required", http.StatusBadRequest)
		return
	}
	idStrInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid song id", http.StatusBadRequest)
	}
	switch r.Method {
	case http.MethodGet:
		song, err := db.GetSongByID(int64(idStrInt))
		if err != nil {
			http.Error(w, "Song not found", http.StatusInternalServerError)
		}
		writeJSON(w, song)
	case http.MethodDelete:
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetTracks(w http.ResponseWriter) {
	songs, err := db.GetAllSongs()
	if err != nil {
		http.Error(w, "Failed to get songs from db", http.StatusInternalServerError)
	}
	writeJSON(w, songs)
}

func getAudioPathsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		paths := db.GetNotAddedSongPaths()
		writeJSON(w, paths)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Can handle uploaded files or existing filePaths
func handleCreateTrack(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20) // allow up to 50MB
	if err != nil {
		http.Error(w, "Invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	filePath := r.FormValue("filePath") // for existing files
	cover := r.FormValue("cover")

	var created []db.Song

	// Case 1: user uploads new file(s)
	if files := r.MultipartForm.File["file"]; len(files) > 0 {
		for _, fh := range files {
			song, err := processUploadedFile(fh, title, author, cover)
			if err != nil {
				http.Error(w, "Failed to process upload: "+err.Error(), http.StatusInternalServerError)
				return
			}
			created = append(created, song)
		}
		writeJSON(w, created)
		return
	}

	// Case 2: use existing filePath
	if filePath != "" {
		song, err := processExistingFile(filePath, title, author, cover)
		if err != nil {
			http.Error(w, "Failed to add existing track: "+err.Error(), http.StatusInternalServerError)
			return
		}
		created = append(created, song)
		writeJSON(w, created)
		return
	}

	http.Error(w, "No file uploaded or filePath provided", http.StatusBadRequest)
}

func processUploadedFile(fh *multipart.FileHeader, title, author, cover string) (db.Song, error) {
	savedPath, _, _, err := utils.SaveUploadedFile(fh, os.Getenv("AUDIO_PATH")+"uploads/")
	if err != nil {
		return db.Song{}, err
	}

	modifiedPath, _ := strings.CutPrefix(savedPath, "music/")
	audioLengthSec, err := utils.GetAudioDuration(modifiedPath)
	if err != nil {
		return db.Song{}, err
	}

	song := db.Song{
		Title:     title,
		Author:    author,
		FilePath:  modifiedPath,
		LengthSec: int(audioLengthSec),
	}
	if cover != "" {
		song.CoverPath = &cover
	}

	if err := db.AddSong(song); err != nil {
		return db.Song{}, err
	}

	return song, nil
}

func processExistingFile(filePath, title, author, cover string) (db.Song, error) {
	notAddedSongs := db.GetNotAddedSongPaths()
	if !slices.Contains(notAddedSongs, filePath) {
		return db.Song{}, fmt.Errorf("song already added")
	}

	audioLengthSec, err := utils.GetAudioDuration(filePath)
	if err != nil {
		return db.Song{}, err
	}

	song := db.Song{
		Title:     title,
		Author:    author,
		FilePath:  filePath,
		LengthSec: int(audioLengthSec),
	}
	if cover != "" {
		song.CoverPath = &cover
	}

	if err := db.AddSong(song); err != nil {
		return db.Song{}, err
	}

	return song, nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
