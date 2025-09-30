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
	"strings"
)

// http track handler
func handleTracks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTracks(w)
	case http.MethodPost:
		handleUploadTrack(w, r)
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

func handleGetAudioPaths(w http.ResponseWriter) {
	paths := db.GetNotAddedSongPaths()
	data, err := json.Marshal(paths)
	utils.Check(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func handleAddTrack(w http.ResponseWriter, r *http.Request) {
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

	fh := fileHeaders[0]

	savedPath, size, contentType, err := utils.SaveUploadedFile(fh, os.Getenv("AUDIO_PATH")+"uploads/")
	if err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	modifiedString, _ := strings.CutPrefix(savedPath, "music/")
	audioLengthSec, err := utils.GetAudioDuration(modifiedString)

	db.AddSong(db.Song{Title: title, Author: author, LengthSec: int(audioLengthSec), FilePath: modifiedString})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Uploaded: %s (%d bytes, %s)\nTitle: %s, Author: %s\n",
		savedPath, size, contentType, title, author)
}
