package config

import (
	"os"
	"strings"
)

func MediaFolder() (string) {
	var mediaFolder string = "/media" // TODO: Change to /media for the Docker container
	if os.Getenv("MEDIA_FOLDER") != "" {
		mediaFolder = os.Getenv("MEDIA_FOLDER")
	}
	return mediaFolder
}

func FileExtensions() []string {
	var extensions []string = []string{".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv"}
	if os.Getenv("FILE_EXTENSIONS") != "" {
		extensions = strings.Split(os.Getenv("FILE_EXTENSIONS"), ",")
	}
	return extensions
}