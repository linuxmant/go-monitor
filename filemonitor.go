package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileMonitor handles file monitoring
type FileMonitor struct {
	Watcher     *fsnotify.Watcher
	FileQueue   chan string
	Done        chan bool
	QueueSize   int
	PathToWatch string
}

// NewFileMonitor creates a new FileMonitor
func NewFileMonitor(path string, queueSize int) (*FileMonitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileMonitor{
		Watcher:     watcher,
		FileQueue:   make(chan string, queueSize),
		Done:        make(chan bool),
		QueueSize:   queueSize,
		PathToWatch: path,
	}, nil
}

// StartMonitoring starts monitoring the file system
func (fm *FileMonitor) StartMonitoring() {
	go func() {
		for {
			select {
			case event, ok := <-fm.Watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("Event detected: %v\n", event)
				fm.FileQueue <- event.Name
			case err, ok := <-fm.Watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	if err := fm.Watcher.Add(fm.PathToWatch); err != nil {
		log.Fatalf("Error adding path to watcher: %v", err)
	}

	fmt.Printf("Monitoring changes in: %s\n", fm.PathToWatch)
	<-fm.Done
}

// StopMonitoring stops the file monitoring process
func (fm *FileMonitor) StopMonitoring() {
	close(fm.FileQueue)
	close(fm.Done)
	fm.Watcher.Close()
}

// ProcessQueue processes files in the queue
func (fm *FileMonitor) ProcessQueue() {
	for file := range fm.FileQueue {
		fmt.Printf("Processing file: %s\n", file)
		time.Sleep(1 * time.Second) // Simulate file processing
	}
}
