package utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Saves an uploaded file to disk and returns (path, size, contentType, error)
func SaveUploadedFile(fh *multipart.FileHeader, saveLocation string) (string, int64, string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", 0, "", err
	}
	defer file.Close()

	if err := os.MkdirAll(saveLocation, os.ModePerm); err != nil {
		return "", 0, "", err
	}

	dstPath := filepath.Join(saveLocation, fh.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", 0, "", err
	}
	defer dst.Close()

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	file.Seek(0, 0)

	size, err := io.Copy(dst, file)
	if err != nil {
		return "", 0, "", err
	}

	return dstPath, size, contentType, nil
}
