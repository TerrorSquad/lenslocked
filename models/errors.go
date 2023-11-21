package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address is already taken")
	ErrNotFound   = errors.New("models: resource not found")
)

type FileError struct {
	Message string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("Invalid file: %v", fe.Message)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seeking file: %w", err)
	}
	contentType := http.DetectContentType(testBytes)
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return nil
		}
	}
	return FileError{Message: fmt.Sprintf("invalid content type: %v", contentType)}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{Message: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename))}
	}
	return nil
}
