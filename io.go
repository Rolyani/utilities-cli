package main

import (
	"fmt"
	"os"
	"strings"
)

func readFile(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("The path is a directory")
	}
	return os.ReadFile(path)
}

func writeNewFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func backupFile(path string) (string, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	backupPath := path + ".bak"
	err = os.WriteFile(backupPath, original, 0644)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

func overwriteFile(path string, data []byte) error {
	dir := "."
	if idx := strings.LastIndex(path, string(os.PathSeparator)); idx != -1 {
		dir = path[:idx]
		if dir == "" {
			dir = string(os.PathSeparator)
		}
	}

	tmp, err := os.CreateTemp(dir, "utilities-cli-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, path)

}
