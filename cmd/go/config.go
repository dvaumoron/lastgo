package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	checkName            = "last-check.txt"
	defaultCheckInterval = 24 * time.Hour
	defaultDownloadURL   = "https://go.dev/dl"
)

type config struct {
	checkInterval time.Duration
	downloadURL   string
	rootPath      string
	dateFilePath  string
}

func InitConfigFromEnv() config {
	checkInterval := defaultCheckInterval
	if checkIntervalEnv := os.Getenv("LASTGO_CHECK_INTERVAL"); checkIntervalEnv != "" {
		parsedInterval, err := time.ParseDuration(checkIntervalEnv)
		if err == nil {
			checkInterval = parsedInterval
		} else {
			fmt.Println("Failed to parse LASTGO_CHECK_INTERVAL, keep using default (24h) :", err)
		}
	}

	downloadURL := os.Getenv("LASTGO_DOWNLOAD_URL")
	if downloadURL == "" {
		downloadURL = defaultDownloadURL
	}

	rootPath := os.Getenv("LASTGO_ROOT")
	if rootPath == "" {
		userPath, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		rootPath = filepath.Join(userPath, ".lastgo")
	}

	return config{
		checkInterval: checkInterval,
		downloadURL:   downloadURL,
		rootPath:      rootPath,
		dateFilePath:  filepath.Join(rootPath, checkName),
	}
}
