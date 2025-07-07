package main

import (
	"log"
    "time"
    "sync"

	"github.com/farmergreg/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

func watch(path string, onCreate func(string), onChange func(string), onDelete func(string), onReady func(string)) {
    watcher, err := rfsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    debouncer := NewDebouncedWatcher(10*time.Second, onReady)

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                debouncer.handleEvent(event)
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

type DebouncedWatcher struct {
	debounceDuration time.Duration
	timers map[string]*time.Timer
	mu sync.Mutex
    OnFileReady func(string)
}

func NewDebouncedWatcher(debounceDuration time.Duration, onFileReady func(string)) *DebouncedWatcher {
	return &DebouncedWatcher{
		debounceDuration: debounceDuration,
		timers:           make(map[string]*time.Timer),
		OnFileReady:      onFileReady,
	}
}

func (d *DebouncedWatcher) handleEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
		path := event.Name

		d.mu.Lock()
		defer d.mu.Unlock()

		if timer, found := d.timers[path]; found {
			timer.Stop()
		}

		d.timers[path] = time.AfterFunc(d.debounceDuration, func() {
			if d.OnFileReady != nil {
				d.OnFileReady(path)
			}
			d.mu.Lock()
			delete(d.timers, path)
			d.mu.Unlock()
		})
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		path := event.Name
		log.Printf("Datei gelöscht: %s", path)
		d.mu.Lock()
		if timer, found := d.timers[path]; found {
			timer.Stop()
			delete(d.timers, path)
		}
		d.mu.Unlock()
	}
}