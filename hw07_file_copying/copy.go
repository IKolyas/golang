package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNoSearchFileError     = errors.New("no search file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return errors.New("необходимо указать путь к исходному и целевому файлам")
	}

	// Получаем абсолютные пути
	absFromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return errors.New("не удалось получить доступ к исходному файлу")
	}

	absToPath, err := filepath.Abs(toPath)
	if err != nil {
		return errors.New("не удалось получить доступ к целевому файлу")
	}

	if absFromPath == absToPath {
		return errors.New("исходный файл и целевой не должны совпадать")
	}

	// Открытие исходного файла
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return ErrNoSearchFileError
	}
	defer srcFile.Close()

	// Проверяем, является ли файл специальным
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return errors.New("ошибка при получении информации о файле")
	}

	if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return errors.New("нельзя копировать специальный файл /dev/urandom")
	}

	// Проверка размера исходного файла
	fileInfo, err = srcFile.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	// Устанавливаем позицию для чтения
	if _, err := srcFile.Seek(offset, io.SeekStart); err != nil {
		fmt.Printf("Ошибка установки курсора: %v\n", err)
		return err
	}

	// Устанавливаем лимит копируемых байт
	if limit <= 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	// Создание целевого файла
	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Создаем прогресс-бар
	bar := pb.StartNew(int(limit))
	defer bar.Finish()

	// Используем io.Copy с ограничением по лимиту
	reader := io.LimitReader(srcFile, limit)
	writer := io.MultiWriter(dstFile, bar)

	if _, err := io.Copy(writer, reader); err != nil {
		return errors.New("ошибка копирования данных")
	}

	return nil
}
