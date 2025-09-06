package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"music-player/internal/db"
	"music-player/internal/utils"
	"net/http"
	"os"
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
	mux.HandleFunc("/tracks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTracks(w)
		case http.MethodPost:
			handleAddTrack(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fs := http.FileServer(http.Dir(os.Getenv("AUDIO_PATH")))
	mux.Handle("/music/", http.StripPrefix("/music/", fs))

	mux.HandleFunc("/audio-paths", func(w http.ResponseWriter, r *http.Request) {
		paths := db.GetAllSongPaths()
		data, err := json.Marshal(paths)
		utils.Check(err)
		w.Write([]byte(data))
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

func handleAddTrack(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var newTrack db.Song
	if err := json.Unmarshal(body, &newTrack); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
	}
}

func handleGetTracks(w http.ResponseWriter) {
	songs := db.GetAllSongs()
	data, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, "Failed song marshal", http.StatusInternalServerError)
	}
	w.Write([]byte(data))
}
