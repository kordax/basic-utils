/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package ufile

import (
	"io"
	"os"
	"regexp"
)

// MustRead panics in case of any error.
// It is a convenience function that simplifies error handling for cases where
// errors are unexpected or should halt program execution.
// Use MustRead when you're confident the operation should not fail under normal conditions,
// such as reading embedded resources or files that are guaranteed to exist.
func MustRead(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content
}

// CreateFile creates a new file at the specified path and optionally writes content to it.
// If the file already exists, it will be truncated before writing the content.
// Returns an error if any operation fails.
func CreateFile(path string, content *string) error {
	// Create or truncate the file at the specified path.
	file, err := os.Create(path)
	if err != nil {
		// Return the error to the caller if file creation fails.
		return err
	}
	defer file.Close() // Ensure the file is closed when the function exits.

	// If content is provided (not nil), write it to the file.
	if content != nil {
		_, err = io.WriteString(file, *content)
		if err != nil {
			// Return the error to the caller if writing fails.
			return err
		}
	}

	// Return nil on success, indicating no error occurred.
	return nil
}

// Exists checks if a file or directory at the specified path exists.
// It returns true if the path exists, false otherwise.
func Exists(path string) bool {
	// Attempt to get the status of the path. If an error is returned,
	// use os.IsNotExist to determine if the error is due to the path not existing.
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListFiles lists the names of all files (not directories) in the specified directory.
// If a regex pattern is provided, only files matching the pattern are included.
// Returns a slice of file names and any error encountered.
func ListFiles(dir string, regex *string) ([]string, error) {
	var compiled *regexp.Regexp // Holds the compiled regular expression, if provided.

	// If a regex pattern is provided, compile it into a Regexp object.
	// Return any errors encountered during compilation.
	if regex != nil {
		var err error
		compiled, err = regexp.Compile(*regex)
		if err != nil {
			return nil, err
		}
	}

	// Read the directory specified by 'dir' and return an error if unsuccessful.
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileList []string // Slice to hold the names of files to be returned.
	for _, file := range files {
		// Skip directories.
		if file.IsDir() {
			continue
		}

		// If a regex pattern was compiled and the file name matches the pattern,
		// or if no pattern was provided, add the file name to the fileList.
		if compiled == nil || compiled.MatchString(file.Name()) {
			fileList = append(fileList, file.Name())
		}
	}

	// Return the list of file names and nil for the error if everything was successful.
	return fileList, nil
}
