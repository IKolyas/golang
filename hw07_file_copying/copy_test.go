package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Создание временного исходного файла
	srcFile, err := os.CreateTemp("", "src")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(srcFile.Name())

	// Запись данных в исходный файл
	srcData := []byte("Hello, World!")
	if _, err := srcFile.Write(srcData); err != nil {
		t.Fatalf("Ошибка записи в исходный файл: %v", err)
	}

	// Создание временного целевого файла
	dstFile, err := os.CreateTemp("", "dst")
	if err != nil {
		t.Fatalf("Ошибка создания временного файла: %v", err)
	}
	defer os.Remove(dstFile.Name())

	// Тестирование копирования файла без ограничений и смещения
	if err := Copy(srcFile.Name(), dstFile.Name(), 5, 2); err != nil {
		t.Fatalf("Ошибка копирования файла: %v", err)
	}

	// Проверка содержимого целевого файла
	dstData, err := os.ReadFile(dstFile.Name())
	if err != nil {
		t.Fatalf("Ошибка чтения целевого файла: %v", err)
	}

	if string(dstData) != string(srcData[5:7]) {
		t.Errorf("Неверное содержимое целевого файла: ожидается %q, получено %q", srcData[2:7], dstData)
	}
}

func TestCopyFileInvalidSource(t *testing.T) {
	t.Run("invalid source", func(t *testing.T) {
		from := "invalid/path/to/source.txt"
		to := "some/destination.txt"
		limit := 0
		offset := 0

		err := Copy(from, to, int64(limit), int64(offset))
		require.Equal(t, err, ErrNoSearchFileError)
	})
}
