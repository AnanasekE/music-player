package songs

import (
	"log"
	"os"
	"path/filepath"
)

type FileInfo struct {
	FileName string
	Path     string
}

func GetAllSongs() []FileInfo {
	audioDir := os.Getenv("AUDIO_PATH")

	var songs []FileInfo

	err := filepath.WalkDir(audioDir, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !dirEntry.IsDir() {
			songs = append(songs, FileInfo{FileName: dirEntry.Name(), Path: path})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return songs
}
