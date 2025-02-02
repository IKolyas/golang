package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

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
	// Открытие исходного файла
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return ErrNoSearchFileError
	}
	defer srcFile.Close()

	// Проверка размера исходного файла
	fileInfo, err := srcFile.Stat()
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

	bar := pb.StartNew(int(limit))

	// Чтение данных из исходного файла и запись в целевой файл
	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(dstFile)
	for {
		buf := make([]byte, 1024*8) // Размер буфера для чтения
		n, err := reader.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if limit < 0 {
			break
		}

		_, err = writer.Write(buf[:limit])
		if err != nil {
			return err
		}

		limit -= int64(n)
		bar.Add(int(limit) + n)
	}

	// Закрытие целевого файла и освобождение ресурсов
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Ошибка закрытия целевого файла: %v\n", err)
		return err
	}

	bar.Finish()

	return nil
}
