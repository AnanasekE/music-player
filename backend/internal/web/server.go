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
)

func StartServer() {
	mux := http.NewServeMux()
	registerRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(mux),
	}

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func registerRoutes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir(os.Getenv("AUDIO_PATH")))
	mux.Handle("/music/", http.StripPrefix("/music/", fs))

	mux.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTracks(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/audio-paths", func(w http.ResponseWriter, r *http.Request) {
		paths := db.GetNotAddedSongPaths()
		data, err := json.Marshal(paths)
		utils.Check(err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(data))
	})

	mux.HandleFunc("/add-track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAddTrack(w, r)
			return
		}
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/upload-track", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/upload-tracks", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleGetTracks(w http.ResponseWriter) {
	songs := db.GetAllSongs()
	data, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, "Failed song marshal", http.StatusInternalServerError)
	}
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
