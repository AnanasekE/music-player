package web

import (
	"encoding/json"
	"fmt"
	"log"
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
		contentType := r.Header.Get("Content-type")
		switch contentType {
		case "application/json":
			handleAddTrack(w, r)
		case "multipart/form-data":
			handleUploadTrack(w, r)
		}
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
		data, err := json.Marshal(song)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
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
	data, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, "Failed song marshal", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func getAudioPathsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		paths := db.GetNotAddedSongPaths()
		data, err := json.Marshal(paths)
		utils.Check(err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAddTrack(w http.ResponseWriter, r *http.Request) {
	// TODO: send added song back to the server
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	filePath := r.FormValue("filePath")
	cover := r.FormValue("cover")

	notAddedSongs := db.GetNotAddedSongPaths()

	if !slices.Contains(notAddedSongs, filePath) {
		http.Error(w, "Song alread added", http.StatusInternalServerError)
	}

	audioLengthSec, err := utils.GetAudioDuration(filePath)

	if err != nil {
		log.Println(err)
		http.Error(w, "Audio length check failed", http.StatusInternalServerError)
		return
	}

	if cover != "" {
		db.AddSong(db.Song{Title: title, Author: author, FilePath: filePath, LengthSec: int(audioLengthSec), CoverPath: &cover})
	} else {
		db.AddSong(db.Song{Title: title, Author: author, FilePath: filePath, LengthSec: int(audioLengthSec)})
	}

	w.Write([]byte("Song added successfully"))
}

func handleUploadTrack(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Receiving files")
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "File too big or invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")

	fileHeaders := r.MultipartForm.File["file"]
	if len(fileHeaders) == 0 {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	// TODO: update to allow more than 1 header
	fh := fileHeaders[0]

	savedPath, size, contentType, err := utils.SaveUploadedFile(fh, os.Getenv("AUDIO_PATH")+"uploads/")
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedString, _ := strings.CutPrefix(savedPath, "music/")
	audioLengthSec, err := utils.GetAudioDuration(modifiedString)

	err = db.AddSong(db.Song{Title: title, Author: author, LengthSec: int(audioLengthSec), FilePath: modifiedString})
	if err != nil {
		http.Error(w, "Failed to add song", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Uploaded: %s (%d bytes, %s)\nTitle: %s, Author: %s\n",
		savedPath, size, contentType, title, author)
}
