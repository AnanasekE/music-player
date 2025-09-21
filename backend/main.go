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

	db.InitDB()
	web.StartServer()
}
