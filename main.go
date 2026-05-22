package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	err := fo.initLog()

	if err != nil {
		return fmt.Errorf("Ошибка инициализации логирования: %w", err)
	}

	_, err = os.Stat(sourcePath)

	if err != nil {
		return fmt.Errorf("ошибка отсутствия исходного файла: %w", err)
	}
	fo.logSuccess("Исходный файл: " + sourcePath)
	fullPath := filepath.Join(fo.sourceDir, targetDir)
	err = os.MkdirAll(fullPath, 0644)

	if err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}
	_, err = os.Stat(sourcePath)

	if err != nil {
		return fmt.Errorf("исходный файл отсутствует ошибка: %w", err)
	}
	name := filepath.Base(sourcePath)
	ext := filepath.Ext(sourcePath)
	newPath := filepath.Join(fullPath, name)
	_, err = os.Stat(newPath)

	if err == nil {
		fo.logSuccess("Существующий: " + newPath)
		newPath = strings.TrimSuffix(newPath, ext) + "_" + time.Now().Format("2006-01-02_15-04-05") + ext
	} else {
		fo.logSuccess("Целевая папка: " + targetDir)
	}

	fmt.Println(newPath)
	err = os.Rename(sourcePath, newPath)

	if err != nil {
		return fmt.Errorf("ошибка перемещения файла: %w", err)
	}
	fo.logSuccess("Результат: " + newPath)
	return nil
}

func main() {
	organizer, err := NewFileOrganizer("/home/neo")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Println("FileOrganizer создан для директории:", organizer.sourceDir)

	err = organizer.moveFile("/home/neo/greet01.s", "go_tst")
	if err != nil {
		fmt.Println(err)
		return
	}
}
