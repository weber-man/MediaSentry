package main

import (
	"log"
	"os"
	"strings"

	"github.com/weberman/MediaSentry/src/config"
	"github.com/weberman/MediaSentry/src/database"
	"github.com/weberman/MediaSentry/src/media"
	"github.com/weberman/MediaSentry/src/watcher"
)

var db *database.Database

func main() {
	log.Println("MediaSentry is starting...")
	log.Println("Initiating SQLite Database...")
	db = database.InitDB()
	defer db.Close()
	log.Println("Database initialized successfully.")

	log.Println("Starting watching Files...")
	watcher.Watch(config.MediaFolder(), onCreate, onChange, onDelete, onReady)

}

func onCreate(path string) {
	// Handle file creation
	log.Println("File created:", path)
	// if checkShouldFileBeUsed(path) {
	// 	checks(path)
	// }
}

func onChange(path string) {
	// Handle file modification
	// log.Println("File changed:", path)
}

func onDelete(path string) {
	// Handle file deletion
	log.Println("File deleted:", path)
}

func onReady(path string) {
	// Handle file ready for processing
	log.Println("File is ready for processing:", path)
	if checkShouldFileBeUsed(path) {
		media.Checks(path)
	}
}

func checkShouldFileBeUsed(path string) bool {
	if !isFile(path) {
		return false
	}
	var allowedfileExtensions []string = config.FileExtensions()
	ext := getFileExtension(path)

	var isAllowedExt bool = false
	for _, allowedExt := range allowedfileExtensions {
		if allowedExt == ext {
			isAllowedExt = true
			break
		}
	}
	return isAllowedExt
}

func getFileExtension(path string) string {
	if len(path) > 0 && path[len(path)-1] == '/' {
		return ""
	}
	parts := strings.Split(path, ".")
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[len(parts)-1]
}

func isFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}
