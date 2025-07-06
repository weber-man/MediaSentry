package main

import (
	"log"

	"github.com/farmergreg/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

func watch(path string, onCreate func(string), onChange func(string), onDelete func(string)) {
    watcher, err := rfsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
				if event.Op == fsnotify.Create {
					onCreate(event.Name)
				}
                if event.Op == fsnotify.Write {
                    onChange(event.Name)
                }
				if event.Op == fsnotify.Remove {
                    onDelete(event.Name)
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()

    err = watcher.AddRecursive(path)
    if err != nil {
        log.Fatal(err)
    }
    <-make(chan struct{})
}