package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dvaumoron/lastgo/config"
)

func main() {
	conf := config.InitFromEnv()

	// TODO

	runGo("", "")
}

func runGo(installPath string, installedVersion string) {
	cmdArgs := os.Args[1:]
	cmd := exec.Command(filepath.Join(installPath, installedVersion, config.GoName), cmdArgs...)
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
