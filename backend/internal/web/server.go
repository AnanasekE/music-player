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

	mux.HandleFunc("/tracks", trackHandler)
	mux.HandleFunc("/tracks/paths", getAudioPathsHandler)
	mux.HandleFunc("/tracks/{id}", handleTracksByID)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
