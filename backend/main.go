package main

import (
	"log"
	"music-player/internal/db"
	"music-player/internal/web"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db.LoadTracksMetadata()
	// db.AddSong(db.Song{Title: "Bluest Flame", Author: "CerberVT", LengthSec: 147, FilePath: "'CerberVT - Bluest Flame by Benny Blanco and Selena Gomez.m4a'"})
	web.StartServer()
}
