package main

import (
	"fmt"
	"log"
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

func (fo *FileOrganizer) initLog() error {
	file, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("открытие лога: %w", err)
	}
	fo.logFile = file

	log.SetOutput(fo.logFile)
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	log.Println("[SUCCESS]", message)
}

func (fo *FileOrganizer) logError(message string) {
	log.Println("[ERROR]", message)
}

func (fo *FileOrganizer) Close() error {
	if fo.logFile == nil {
		return nil
	}
	err := fo.logFile.Close()
	fo.logFile = nil
	return err
}

func main() {
	organizer, err := NewFileOrganizer("/home/neo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("FileOrganizer создан для директории:", organizer.sourceDir)
	err = organizer.initLog()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer organizer.Close()

	organizer.logSuccess("lol kek piza cheburek")
}
