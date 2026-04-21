package tool

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ListFiles(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fmt.Sprintf("%v", fileNames), nil
}

func EditFile(path string, old, new string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	newData := strings.ReplaceAll(string(data), old, new)

	return os.WriteFile(path, []byte(newData), 0644)
}

func CreateFile(path, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func RunCommand(command, workdir string, args []string) (string, string, int) {
	cmd := exec.Command(command, args...)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return stdout.String(), stderr.String(), exitCode
}
