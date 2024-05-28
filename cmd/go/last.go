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
	"os/exec"
	"path/filepath"

	"github.com/dvaumoron/lastgo/pkg/datefile"
	"github.com/dvaumoron/lastgo/pkg/goversion"
)

const binGo = "go/bin/go"

func main() {
	conf := InitConfigFromEnv()

	installedVersion := getInstalledVersion(conf)
	if datefile.OutsideInterval(conf.dateFilePath, conf.checkInterval) {
		if lastVersionDesc := getLastVersion(conf); installedVersion != lastVersionDesc.version {
			fmt.Print("Update to ", lastVersionDesc.version)
			doUpdate := true
			if conf.askConfirm {
				fmt.Print(" ? [y/N]:")

				buffer := make([]byte, 1)
				os.Stdin.Read(buffer)
				readed := buffer[0]

				doUpdate = readed == 'y' || readed == 'Y'
			} else {
				fmt.Println()
			}

			if doUpdate {
				if installedVersion != "" {
					if err := os.RemoveAll(filepath.Join(conf.rootPath, installedVersion)); err != nil {
						fmt.Println("Fail to remove old version :", err)
					}
				}

				if err := install(conf.rootPath, lastVersionDesc); err != nil {
					fmt.Println("Unable to install", lastVersionDesc.version, ":", err)
					os.Exit(1)
				}
				installedVersion = lastVersionDesc.version
			}
		}
	}

	runGo(conf.rootPath, installedVersion)
}

func getInstalledVersion(conf config) string {
	entries, err := os.ReadDir(conf.rootPath)
	if err != nil {
		fmt.Println("Unable to read installation directory :", err)
		os.Exit(1)
	}

	dirNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			dirNames = append(dirNames, entry.Name())
		}
	}

	return goversion.Last(dirNames)
}

func runGo(installPath string, installedVersion string) {
	cmdArgs := os.Args[1:]
	cmd := exec.Command(filepath.Join(installPath, installedVersion, binGo), cmdArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Failure during go call :", err)
	}
}
