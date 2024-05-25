package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	GoName = "go"

	defaultCheckInterval = 24 * time.Hour
	defaultDownloadURL   = "https://go.dev/dl"
)

type Config struct {
	CheckInterval time.Duration
	DownloadURL   string
	RootPath      string
}

func InitFromEnv() Config {
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

	return Config{
		CheckInterval: checkInterval,
		DownloadURL:   downloadURL,
		RootPath:      rootPath,
	}
}
