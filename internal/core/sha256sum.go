package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"

	"github.com/ismdeep/sha256sum-go/pkg/sha256util"
)

// SHA256Sum model
type SHA256Sum struct {
	directory string
	filename  string
}

// NewSHA256Sum new a sha256sum instance
func NewSHA256Sum(directory string, filename string) *SHA256Sum {
	return &SHA256Sum{
		directory: directory,
		filename:  filename,
	}
}

// Generate generate sha256sum check file
func (receiver *SHA256Sum) Generate() error {
	directory, err := filepath.Abs(receiver.directory)
	if err != nil {
		return err
	}

	// get all file under directory
	var fileList []string
	if err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip directory
		if info.IsDir() {
			return nil
		}

		path = path[len(directory)+1:]
		if path != receiver.filename {
			fileList = append(fileList, path)
		}

		return nil
	}); err != nil {
		return err
	}

	// create progress bar
	bar := pb.StartNew(len(fileList))

	// open sha256sum.txt
	outputFile, err := os.Create(receiver.filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// walk file list and calculate sha256sum
	for _, filePath := range fileList {
		// calculate SHA-256
		checksum, err := sha256util.ByFilepath(filepath.Join(directory, filePath))
		if err != nil {
			bar.Finish()
			return err
		}

		// write sha256sum and filepath to sha256sum.txt
		if _, err := fmt.Fprintf(outputFile, "%s  %s\n", checksum, filePath); err != nil {
			bar.Finish()
			return err
		}

		// update progress bar
		bar.Increment()
	}

	// finish progress bar
	bar.Finish()

	return nil
}

// Verify sha256sum check file
func (receiver *SHA256Sum) Verify() error {
	// read sha256sum.txt
	filePath := receiver.filename
	checksums, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(checksums), "\n")

	// create progress bar
	bar := pb.StartNew(len(lines))

	// get lines
	for _, line := range lines {
		// ignore empty line
		if line == "" {
			continue
		}

		// get sha256sum and file path
		idx := strings.Index(line, "  ")
		if idx != 64 {
			bar.Finish()
			return fmt.Errorf("invalid format in %v", line)
		}
		expectedChecksum := line[:idx]
		path := line[idx+2:]

		// get actual SHA-256
		actualChecksum, err := sha256util.ByFilepath(filepath.Join(receiver.directory, path))
		if err != nil {
			bar.Finish()
			return err
		}

		// check actual sha256sum and expected sha256sum
		if actualChecksum != expectedChecksum {
			bar.Finish()
			return fmt.Errorf("verification failed for %v", line)
		}

		// update progress bar
		bar.Increment()
	}

	// finish progress bar
	bar.Finish()
	return nil
}
