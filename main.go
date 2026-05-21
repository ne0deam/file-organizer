package main

import (
	"fmt"
	"os"
)

var DefaultRules = map[string]string{
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".pdf":  "Documents",
	".doc":  "Documents",
	".docx": "Documents",
	".txt":  "Documents",
	".mp3":  "Music",
	".wav":  "Music",
	".mp4":  "Video",
	".avi":  "Video",
	".zip":  "Archives",
	".rar":  "Archives",
}

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if sourceDir == "" {
		return nil, fmt.Errorf("sourceDir не может быть пустым")
	}

	fileInfo, err := os.Stat(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о %q: %w", sourceDir, err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%q не является директорией", sourceDir)
	}

	return &FileOrganizer{
		sourceDir: sourceDir,
		rulesMap:  DefaultRules,
	}, nil
}

func main() {
	organizer, err := NewFileOrganizer("/home/neo/stratus")
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	fmt.Println("FileOrganizer создан для директории:", organizer.sourceDir)
}
