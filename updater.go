package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Updater handles application updates
type Updater struct {
	RepoOwner string
	RepoName  string
	Version   string
}

// NewUpdater creates a new Updater
func NewUpdater(owner, repo, version string) *Updater {
	return &Updater{
		RepoOwner: owner,
		RepoName:  repo,
		Version:   version,
	}
}

// CheckAndUpdate checks for updates and performs an update if available
func (u *Updater) CheckAndUpdate() {
	latestVersion, downloadURL, err := u.fetchLatestVersion()
	if err != nil {
		log.Printf("Failed to fetch latest version: %v\n", err)
		return
	}

	if latestVersion != u.Version {
		log.Printf("New version available: %s (current: %s)\n", latestVersion, u.Version)
		err := u.downloadAndUpdate(downloadURL)
		if err != nil {
			log.Printf("Failed to update to latest version: %v\n", err)
		} else {
			log.Println("Application updated successfully. Restarting...")
			u.restartApplication()
		}
	} else {
		log.Println("Already running the latest version.")
	}
}

// fetchLatestVersion fetches the latest release information
func (u *Updater) fetchLatestVersion() (string, string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.RepoOwner, u.RepoName)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch release: %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return "", "", err
	}

	if len(release.Assets) == 0 {
		return "", "", fmt.Errorf("no assets found in the latest release")
	}

	return release.TagName, release.Assets[0].BrowserDownloadURL, nil
}

// downloadAndUpdate downloads and updates the application
func (u *Updater) downloadAndUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: %s", resp.Status)
	}

	tempFile, err := os.CreateTemp("", "fs-monitor-update-*")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Rename(tempFile.Name(), execPath)
	if err != nil {
		return err
	}

	return nil
}

// restartApplication restarts the application
func (u *Updater) restartApplication() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v\n", err)
	}

	err = exec.Command(execPath).Start()
	if err != nil {
		log.Fatalf("Failed to restart application: %v\n", err)
	}

	os.Exit(0)
}
