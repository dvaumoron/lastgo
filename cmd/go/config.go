/*
 *
 * Copyright 2024 lastgo authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	askConfirm    bool
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
			fmt.Println("Failed to access user directory :", err)
			os.Exit(1)
		}

		rootPath = filepath.Join(userPath, ".lastgo")
	}

	if err := os.MkdirAll(rootPath, 0755); err != nil {
		fmt.Println("Failed to create directory :", err)
		os.Exit(1)
	}

	return config{
		checkInterval: checkInterval,
		downloadURL:   downloadURL,
		rootPath:      rootPath,
		dateFilePath:  filepath.Join(rootPath, checkName),
		askConfirm:    strings.TrimSpace(os.Getenv("LASTGO_ASK")) != "",
	}
}
