package main

import (
	"os"
	"strings"
)

func mediaFolder() (string) {
	var mediaFolder string = "./tmp" // TODO: Change to /media for the Docker container
	if os.Getenv("MEDIA_FOLDER") != "" {
		mediaFolder = os.Getenv("MEDIA_FOLDER")
	}
	return mediaFolder
}

func fileExtensions() []string {
	var extensions []string = []string{".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv"}
	if os.Getenv("FILE_EXTENSIONS") != "" {
		extensions = strings.Split(os.Getenv("FILE_EXTENSIONS"), ",")
	}
	return extensions
}