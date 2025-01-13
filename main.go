package main

import (
	"fmt"
	"log"
	"time"
)

const (
	repoOwner       = "linuxmant"  // Replace with your GitHub username
	repoName        = "go-monitor" // Replace with your repository name
	currentVersion  = "v0.1.0"     // Update this with your current version
	checkInterval   = 1 * time.Hour
	queueBufferSize = 2
)

func main() {
	// Initialize and start the file monitor
	fileMonitor, err := NewFileMonitor("./files", queueBufferSize)
	if err != nil {
		log.Fatalf("Failed to initialize file monitor: %v\n", err)
	}

	go fileMonitor.StartMonitoring()
	go fileMonitor.ProcessQueue()

	// Initialize and start the updater
	updater := NewUpdater(repoOwner, repoName, currentVersion)
	go func() {
		fmt.Println("Startedd checking for updates...")
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("\tchecking for update")
			updater.CheckAndUpdate()
		}
	}()

	// Block the main goroutine
	select {}
}
