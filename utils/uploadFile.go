package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateDirectoryIfNotExists creates a directory if it doesn't exist
func CreateDirectoryIfNotExists(directory string) error {
	if err := os.MkdirAll(directory, 0775); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// SaveFiles saves multiple files to a destination directory
func SaveFiles(files map[string]io.Reader, directory string, allowedExtensions []string) ([]string, error) {
	var savedPaths []string

	for filename, fileReader := range files {
		// Validate file type
		extension := filepath.Ext(filename)
		if err := ValidateFileType(extension, allowedExtensions); err != nil {
			return nil, fmt.Errorf("invalid file type for %s: %w", filename, err)
		}

		// Ensure the directory path length is within bounds
		if len(directory) > 460 {
			return nil, fmt.Errorf("directory path is too long: %s", directory)
		}

		// Sanitize the original filename
		sanitizedFilename := sanitizeFilename(filename, extension)

		// Generate new file name with the original filename, UUID, and timestamp included
		newFileName := fmt.Sprintf("%s-%s-%s-%s%s",
			time.Now().Format("060102"), // Date (YYMMDD)
			time.Now().Format("150405"), // Time (HHMMSS)
			uuid.New().String(),         // UUID
			sanitizedFilename,           // Original filename (sanitized)
			extension,                   // Original file extension
		)

		dstPath := filepath.Join(directory, newFileName)
		if len(dstPath) > 512 {
			return nil, fmt.Errorf("file path is too long: %s", dstPath)
		}

		if err := SaveFile(fileReader, dstPath); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %w", filename, err)
		}

		normalizedPath := filepath.ToSlash(dstPath)
		savedPaths = append(savedPaths, fmt.Sprintf("/%s", normalizedPath))
	}

	return savedPaths, nil
}
func SaveFile(src io.Reader, dst string) error {
	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, src); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// ValidateFileType checks the type of the file based on its extension
func ValidateFileType(filename string, allowedExtensions []string) error {
	ext := filepath.Ext(filename)
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return nil // File type is valid
		}
	}
	return errors.New("invalid file type")
}

// sanitizeFilename removes unwanted characters from the original filename
func sanitizeFilename(filename, extension string) string {
	// Remove the extension from the filename for sanitization
	baseFilename := strings.TrimSuffix(filename, extension)

	// Replace spaces with underscores
	baseFilename = strings.ReplaceAll(baseFilename, " ", "_")

	// For example, allow alphanumeric characters, hyphens, underscores, and periods
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	baseFilename = reg.ReplaceAllString(baseFilename, "")

	return baseFilename
}
