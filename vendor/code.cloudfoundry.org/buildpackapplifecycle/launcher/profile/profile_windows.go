// +build windows

package profile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ProfileEnv(appDir, tempDir string, stdout io.Writer, stderr io.Writer) ([]string, error) {
	fi, err := os.Stat(tempDir)
	if err != nil {
		return nil, fmt.Errorf("invalid temp dir: %s", err.Error())
	} else if !fi.IsDir() {
		return nil, errors.New("temp dir must be a directory")
	}

	envOutputFile := filepath.Join(tempDir, "launcher.env")
	defer os.Remove(envOutputFile)

	batchFileLines := []string{
		"@echo off",
		fmt.Sprintf("cd %s", appDir),
		`(for /r %i in (..\profile.d\*) do %i)`,
		`(for /r %i in (.profile.d\*) do %i)`,
		`(if exist .profile.bat ( .profile.bat ))`,
		fmt.Sprintf("set > %s", envOutputFile),
	}

	cmd := exec.Command("cmd", "/c", strings.Join(batchFileLines, " & "))
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return []string{}, fmt.Errorf("running profile scripts failed: %s", err.Error())
	}
	out, err := ioutil.ReadFile(envOutputFile)
	if err != nil {
		return []string{}, err
	}

	cleanedVars := []string{}
	vars := strings.Split(string(out), "\n")
	for _, v := range vars {
		if v != "" {
			cleanedVars = append(cleanedVars, strings.TrimSpace(v))
		}
	}

	return cleanedVars, nil
}