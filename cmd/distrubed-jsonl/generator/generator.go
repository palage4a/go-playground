package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"math/rand"
)

// Record описывает структуру объекта в JSONL файле.
type Record struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

func Run(k, n, totalRecords int, output string) error {
	// Проверяем, что totalRecords делится равномерно на количество ключей
	rangeCount := n - k + 1
	if totalRecords%rangeCount != 0 {
		return fmt.Errorf("ошибка: Общее число записей (%d) должно быть кратно количеству ключей (%d) для равномерного распределения", totalRecords, rangeCount)
	}
	repeats := totalRecords / rangeCount

	// Генерируем срез ключей: каждый ключ из диапазона [k, n] повторяется `repeats` раз.
	keys := make([]int, 0, totalRecords)
	for i := k; i <= n; i++ {
		for j := 0; j < repeats; j++ {
			keys = append(keys, i)
		}
	}

	// Перемешиваем срез ключей случайным образом
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// Создаём (или открываем) выходной файл
	outFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer outFile.Close()

	encoder := json.NewEncoder(outFile)

	// Генерируем записи и записываем их в файл по одной строке
	// Каждая запись кодируется в формате JSON и записывается с переводом строки.
	for _, key := range keys {
		rec := Record{
			Key:   key,
			Value: "dummy", // Значение value можно заменить на любое другое
		}
		if err := encoder.Encode(rec); err != nil {
			return fmt.Errorf("encode: %w", err)
		}
	}

	return nil
}
