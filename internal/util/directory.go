package util

import (
	"os"
	"path/filepath"
)

// ListFilesDirectory returns a list of files in a directory
func ListFilesDirectory(pathname string) []string {
	var r []string
	files, err := os.ReadDir(pathname)
	CheckErr(err)

	for _, f := range files {
		r = append(r, f.Name())
	}
	return r
}

// GetAbsolutePath find absolute path of from relative path
func GetAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	CheckErr(err)
	return filepath.Clean(absPath)
}

// GetWorkingDirectory get application working directory
func GetWorkingDirectory() string {
	dir, err := os.Getwd()
	CheckErr(err)
	return filepath.Clean(dir)
}

func CheckIfDirExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateDirIfNotExists(dir string) error {
	if !CheckIfDirExists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func GetFileName(path string) string {
	return filepath.Base(path)
}
