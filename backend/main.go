package main

import (
	"github.com/joho/godotenv"
	"log"
	"music-player/internal/web"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	web.StartServer()
}
