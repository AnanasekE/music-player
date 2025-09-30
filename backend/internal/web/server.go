package web

import (
	"fmt"
	"log"
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
		switch r.Method {
		case http.MethodGet:
			handleGetAudioPaths(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/add-track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleAddTrack(w, r)
			return
		}
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	})

	mux.HandleFunc("/upload-track", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handleUploadTrack(w, r)
			return
		}
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
