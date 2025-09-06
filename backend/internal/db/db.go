package db

import (
	"encoding/json"
	"io/fs"
	"music-player/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	tracks []Song
	mu     sync.Mutex
)

type Song struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	LengthSec int    `json:"lengthSec"`
	FilePath  string `json:"filePath"`
}

func LoadTracksMetadata() {
	path := os.Getenv("TRACK_METADATA_JSON_PATH")
	if path == "" {
		panic("TRACK_METADATA_JSON_PATH not set")
	}

	data, err := os.ReadFile(path)
	utils.Check(err)

	err = json.Unmarshal(data, &tracks)
	utils.Check(err)
}

func SaveTrackMetadata() {
	path := os.Getenv("TRACK_METADATA_JSON_PATH")
	data, err := json.MarshalIndent(tracks, "", "  ")
	utils.Check(err)

	err = os.WriteFile(path, data, 0644)
	utils.Check(err)
}

func AddSong(song Song) {
	mu.Lock()
	defer mu.Unlock()

	tracks = append(tracks, song)
	SaveTrackMetadata()
}

func RemoveSong(songId int) {
	mu.Lock()
	defer mu.Unlock()

	if songId < 0 || songId >= len(tracks) {
		return
	}

	tracks = append(tracks[:songId], tracks[songId+1:]...)
	SaveTrackMetadata()
}

func GetSongData(title string) *Song {
	mu.Lock()
	defer mu.Unlock()

	for _, track := range tracks {
		if track.Title == title {
			return &track
		}
	}
	return nil
}

func GetAllSongs() []Song {
	mu.Lock()
	defer mu.Unlock()

	copyTracks := make([]Song, len(tracks))
	copy(copyTracks, tracks)
	return copyTracks
}

func GetAllSongPaths() []string {
	var paths []string

	err := filepath.WalkDir(os.Getenv("AUDIO_PATH"), func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dirEntry.IsDir() {
			paths = append(paths, strings.TrimPrefix(path, "music/"))
		}
		return nil
	})
	utils.Check(err)
	return paths
}
