package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetAudioDuration(path string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		os.Getenv("AUDIO_PATH")+path,
	)

	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var probe struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out, &probe); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	var duration float64
	if _, err := fmt.Sscanf(probe.Format.Duration, "%f", &duration); err != nil {
		return 0, fmt.Errorf("invalid duration value: %w", err)
	}

	return duration, nil
}
