package fileprocessor

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileProcessor struct {
	directory string
}

func NewFileProcessor(directory string) *FileProcessor {
	return &FileProcessor{directory: directory}
}

func (fp *FileProcessor) GetLatestCSVFile() ([][]string, error) {
	files, err := os.ReadDir(fp.directory)
	if err != nil {
		return nil, fmt.Errorf("could not read directory: %v", err)
	}

	var latestFile os.DirEntry
	var latestModTime time.Time

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".csv" {
			continue
		}

		lockFilePath := filepath.Join(fp.directory, file.Name()+".lock")
		if _, err := os.Stat(lockFilePath); err == nil {
			// Lock file exists, skip this CSV file
			continue
		}

		filePath := filepath.Join(fp.directory, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			return nil, fmt.Errorf("could not stat file: %v", err)
		}

		if info.ModTime().After(latestModTime) {
			latestModTime = info.ModTime()
			latestFile = file
		}
	}

	if latestFile == nil {
		return nil, fmt.Errorf("no CSV files found in directory")
	}

	latestFilePath := filepath.Join(fp.directory, latestFile.Name())
	lockFilePath := latestFilePath + ".lock"

	if err := fp.createLockFile(lockFilePath); err != nil {
		return nil, err
	}

	records, err := fp.readCSVFile(latestFilePath)
	if err != nil {
		return nil, err
	}

	// Delete files after successful processing
	if err := fp.deleteFiles(latestFilePath, lockFilePath); err != nil {
		return nil, err
	}

	return records, nil
}

func (fp *FileProcessor) createLockFile(lockFilePath string) error {
	lockFile, err := os.Create(lockFilePath)
	if err != nil {
		return fmt.Errorf("could not create lock file: %v", err)
	}
	defer lockFile.Close()

	return nil
}

func (fp *FileProcessor) readCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read CSV file: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	return records, nil
}

func (fp *FileProcessor) deleteFiles(filePath, lockFilePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete CSV file: %v", err)
	}

	if err := os.Remove(lockFilePath); err != nil {
		return fmt.Errorf("could not delete lock file: %v", err)
	}

	return nil
}
