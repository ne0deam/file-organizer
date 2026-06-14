package main

import (
	"fmt"
	"io/fs"
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
		sourceDir = "/home/neo/srcd"
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
	log.SetOutput(os.Stdout)
	fo.logFile = nil
	return err
}

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	_, err := os.Stat(sourcePath)
	if err != nil {
		fo.logError(err.Error())
		return fmt.Errorf("ошибка существования файла %w", err)
	}

	fo.logSuccess("Исходный файл:" + sourcePath)
	fullTargetDir := filepath.Join(fo.sourceDir, targetDir)
	err = os.MkdirAll(fullTargetDir, os.ModePerm)
	if err != nil {
		fo.logError(err.Error())
		return fmt.Errorf("ошибка создания директории: %w", err)
	}

	fileName := filepath.Base(sourcePath)
	ext := filepath.Ext(fileName)
	fullPath := filepath.Join(fullTargetDir, fileName)
	_, err = os.Stat(fullPath)
	if err == nil {
		fo.logSuccess("Существующий:" + fullPath)
		fullPath = strings.TrimSuffix(fullPath, ext) + "_" + time.Now().Format("2006-01-02_15-04-05") + ext
	} else {
		fo.logSuccess("Целевая папка: " + targetDir)
	}
	err = os.Rename(sourcePath, fullPath)
	if err != nil {
		fo.logError(err.Error())
		return fmt.Errorf("ошибка перемещения файла: %w", err)
	}

	fo.logSuccess("Результат: " + fullPath)
	return nil
}

func (fo *FileOrganizer) initLog() (*os.File, error) {
	logFile, err := os.OpenFile("organizer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть лог файл: %w", err)
	}
	log.SetOutput(logFile)
	return logFile, err
}

func (fo *FileOrganizer) Organize() error {
	fo.initLog()
	err := filepath.WalkDir(fo.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // пропускаем директории
		}

		// обработка файла
		fmt.Println(path)
		if filepath.Dir(path) != fo.sourceDir {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if dir, ok := fo.rulesMap[ext]; ok {
			err := fo.moveFile(path, dir)
			if err != nil {
				fo.logError(err.Error())
				return err
			}

			fo.processedFiles += 1
		}

		return nil
	})
	if err != nil {
		fo.logError(err.Error())
		return err
	}

	return nil
}

func main() {
	organizer, err := NewFileOrganizer("/home/neo/srcd")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer organizer.Close()

	organizer.initLog()
	err = organizer.Organize()
	if err != nil {
		fmt.Println(err)
		return
	}
}
