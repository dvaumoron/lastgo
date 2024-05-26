package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dvaumoron/lastgo/pkg/datefile"
	"github.com/dvaumoron/lastgo/pkg/goversion"
)

const GoName = "go"

func main() {
	conf := InitConfigFromEnv()

	installedVersion := getInstalledVersion(conf)
	if datefile.OutsideInterval(conf.dateFilePath, conf.checkInterval) {
		lastVersionDesc := getLastVersion(conf)

		if installedVersion != lastVersionDesc.version {
			fmt.Println("Update to", lastVersionDesc.version)

			err := os.RemoveAll(filepath.Join(conf.rootPath, installedVersion))
			if err != nil {
				fmt.Println("Fail to remove old version :", err)
			}

			if err = install(conf.rootPath, lastVersionDesc); err != nil {
				fmt.Println("Unable to install", lastVersionDesc.version, ":", err)
				os.Exit(1)
			}
			installedVersion = lastVersionDesc.version
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
	cmd := exec.Command(filepath.Join(installPath, installedVersion, GoName), cmdArgs...)
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
