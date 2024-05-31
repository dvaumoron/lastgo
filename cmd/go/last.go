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
	"strings"

	"github.com/dvaumoron/lastgo/pkg/datefile"
	"github.com/dvaumoron/lastgo/pkg/goversion"
)

const (
	binGo      = "go/bin/go"
	goName     = "go"
	versionCmd = "version"
)

func main() {
	conf := InitConfigFromEnv()

	installedVersion, versionMessage := getInstalledVersion(conf)
	notInstalled := installedVersion == ""
	if notInstalled || datefile.OutsideInterval(conf.dateFilePath, conf.checkInterval) {
		if lastVersionDesc := getLastVersion(conf); installedVersion != lastVersionDesc.version {
			doUpdate := true
			if notInstalled {
				fmt.Println("Install", lastVersionDesc.version)
			} else if conf.noConfirm {
				fmt.Println("Update to", lastVersionDesc.version)
			} else {
				fmt.Print("Update to ", lastVersionDesc.version, " ? [y/N]:")

				buffer := make([]byte, 1)
				os.Stdin.Read(buffer)
				readed := buffer[0]

				doUpdate = readed == 'y' || readed == 'Y'
			}

			if doUpdate {
				versionMessage = ""
				if notInstalled {
					if err := install(conf.rootPath, lastVersionDesc); err != nil {
						fmt.Println("Unable to install", lastVersionDesc.version, ":", err)
						os.Exit(1)
					}
				} else {
					var builder strings.Builder
					builder.WriteString("GOTOOLCHAIN=")
					builder.WriteString(lastVersionDesc.version)
					builder.WriteString("+auto")

					runGo(conf.rootPath, []string{"env", "-w", builder.String()}, stdoutSetter)
				}
			}
		}
	}

	cmdArgs := os.Args[1:]
	if versionMessage != "" && len(cmdArgs) > 0 && cmdArgs[0] == versionCmd {
		fmt.Print(versionMessage)
		return
	}

	runGo(conf.rootPath, cmdArgs, stdoutSetter)
}

func getInstalledVersion(conf config) (string, string) {
	if _, err := os.Stat(filepath.Join(conf.rootPath, goName)); err != nil {
		if os.IsNotExist(err) {
			return "", ""
		}

		fmt.Println("Unable to read installation directory :", err)
		os.Exit(1)
	}

	var outBuilder strings.Builder
	runGo(conf.rootPath, []string{versionCmd}, func(cmd *exec.Cmd) {
		cmd.Stdout = &outBuilder
	})
	versionMessage := outBuilder.String()

	return goversion.Find(versionMessage), versionMessage
}

func runGo(installPath string, cmdArgs []string, outSetter func(*exec.Cmd)) {
	cmd := exec.Command(filepath.Join(installPath, binGo), cmdArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	outSetter(cmd)

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		fmt.Println("Failure during go call :", err)
	}
}

func stdoutSetter(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
}
