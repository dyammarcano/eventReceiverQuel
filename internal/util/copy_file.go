package util

import (
	"io"
	"os"
)

func CopyFile(sourceFile, targetFile string) error {
	// Open the file for reading
	file, err := os.Open(sourceFile)
	if err != nil {
		return err
	}

	defer DeferCloseFatal(file)

	// Create the target file
	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	defer DeferCloseFatal(target)

	// Copy the contents from file to target
	if _, err = io.Copy(target, file); err != nil {
		return err
	}

	// Sync the file to ensure the contents are written to disk
	if err = target.Sync(); err != nil {
		return err
	}

	return nil
}
