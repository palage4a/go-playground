package generator_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/palage4a/go-playground/cmd/distrubed-jsonl/generator"
)

// Record соответствует структуре, используемой при генерации JSONL файла.
type Record struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

func TestUniformDistribution(t *testing.T) {
	// Задайте параметры, которые использовались при генерации файла.
	// Например, если диапазон ключей от 1 до 200:
	k := 1
	n := 200
	// Общее число записей, которое генерировалось.
	totalRecords := int(1e6)
	outputFile := "./testdata/output.jsonl" // Убедитесь, что путь к файлу правильный.

	// Количество возможных ключей:
	numKeys := n - k + 1

	// Проверяем, что общее число записей делится равномерно на число ключей.
	if totalRecords%numKeys != 0 {
		t.Fatalf("Общее число записей (%d) не делится на число ключей (%d)", totalRecords, numKeys)
	}

	err := generator.Run(k, n, totalRecords, outputFile)
	if err != nil {
		t.Fatalf("Ошибка при генерации файла: %v", err)
	}

	// Ожидаемое количество записей для каждого ключа.
	expectedCount := totalRecords / numKeys

	// Открываем файл, который нужно проверить.
	file, err := os.Open(outputFile)
	if err != nil {
		t.Fatalf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Считаем количество вхождений каждого ключа.
	counts := make(map[int]int)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		var rec Record
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("Ошибка декодирования JSON в строке %d: %v", lineNumber, err)
		}
		counts[rec.Key]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Проверяем, что для каждого ключа в диапазоне [k, n] количество записей равно expectedCount.
	for key := k; key <= n; key++ {
		if count, ok := counts[key]; !ok {
			t.Errorf("Ключ %d отсутствует в файле", key)
		} else if count != expectedCount {
			t.Errorf("Для ключа %d ожидается %d записей, получено %d", key, expectedCount, count)
		}
	}
}
