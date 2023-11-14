package util

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type FileMetadata struct {
	File       *os.File
	Filename   string
	FileInfo   os.FileInfo
	Size       int64
	IsNotExist bool
	Empty      bool
	Error      error
	Filepath   string
	MimeType   string
	Hash       string
}

// OpenFileMetadata opens a file and returns a FileMetadata struct
func OpenFileMetadata(path string, hash bool) *FileMetadata {
	m := FileMetadata{
		Filename:   filepath.Base(path),
		Size:       0,
		File:       nil,
		FileInfo:   nil,
		Error:      nil,
		IsNotExist: false,
		Empty:      false,
		Filepath:   filepath.Clean(path),
		MimeType:   mime.TypeByExtension(filepath.Ext(path)),
		Hash:       "",
	}

	if m.File, m.Error = os.Open(path); m.Error != nil {
		return &m
	}

	if m.File == nil {
		m.Error = errors.New("file is nil")
		return &m
	}

	m.FileInfo, m.Error = m.File.Stat()
	m.Size = m.FileInfo.Size()

	if os.IsNotExist(m.Error) {
		m.Error = errors.New("file not found")
		m.IsNotExist = true
		return &m
	}

	if m.Size == 0 {
		m.Empty = true
		m.Error = errors.New("file is empty")
		return &m
	}

	if hash {
		h := sha256.New()
		if _, m.Error = io.Copy(h, m.File); m.Error != nil {
			CheckErr(m.Error)
		}

		// calculate hash and convert to hex
		m.Hash = fmt.Sprintf("%x", h.Sum(nil))

		// reset file pointer
		if _, m.Error = m.File.Seek(0, 0); m.Error != nil {
			CheckErr(m.Error)
		}
	}
	return &m
}

func getFilenameWithoutExtension(fileName string) string {
	i := strings.LastIndex(fileName, ".")
	if i > 0 {
		return fileName[:i]
	}
	return ""
}

func ChangeExtension(fileName string, newExtension string) string {
	return getFilenameWithoutExtension(fileName) + "." + newExtension
}
