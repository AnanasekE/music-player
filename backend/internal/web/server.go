package web

import (
	"encoding/json"
	"fmt"
	"log"
	"music-player/internal/songs"
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
		allSongs := songs.GetAllSongs()
		data, err := json.Marshal(allSongs)
		if err != nil {
			log.Fatal("Borked")
		}
		w.Write(data)
	})

	fs := http.FileServer(http.Dir(os.Getenv("AUDIO_PATH")))
	mux.Handle("/music/", http.StripPrefix("/music/", fs))
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
