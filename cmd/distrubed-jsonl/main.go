package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/palage4a/go-playground/cmd/distrubed-jsonl/generator"
)

var (
	min     = flag.Int("min", 1, "")
	max     = flag.Int("max", 200, "")
	records = flag.Int("records", 1e6, "number of records")
	output  = flag.String("output", "output.jsonl", "output file name")
)

func main() {
	flag.Parse()
	// Параметры генерации

	k := *min
	n := *max
	totalRecords := *records

	err := generator.Run(*min, *max, *records, *output)
	if err != nil {
		fmt.Printf("Ошибка генерации: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Сгенерировано %d записей в диапазоне от %d до %d\n", totalRecords, k, n)
}

// package main

// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"os"
// 	"os/exec"
// 	"time"
// )

// // Record описывает структуру объекта в JSONL файле.
// type Record struct {
// 	Key   int    `json:"key"`
// 	Value string `json:"value"`
// }

// func main() {
// 	// Параметры генерации.
// 	k := 1
// 	n := 200
// 	totalRecords := int64(1e6) // Пример: 1 миллиард записей

// 	numKeys := int64(n - k + 1)
// 	if totalRecords%numKeys != 0 {
// 		fmt.Printf("totalRecords (%d) должно делиться на количество ключей (%d) без остатка.\n", totalRecords, numKeys)
// 		return
// 	}
// 	repeats := totalRecords / numKeys

// 	// Файл для записи данных с случайным префиксом.
// 	tempFileName := "temp.txt"
// 	tempFile, err := os.Create(tempFileName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer tempFile.Close()

// 	// Инициализация генератора случайных чисел.
// 	rand.Seed(time.Now().UnixNano())

// 	// Генерируем данные и сразу записываем их в файл.
// 	// Для каждого ключа от k до n повторяем repeats раз.
// 	// encoder := json.NewEncoder(tempFile)
// 	for key := k; key <= n; key++ {
// 		for i := int64(0); i < repeats; i++ {
// 			// Генерируем случайный префикс.
// 			prefix := rand.Float64()
// 			// Формируем JSON объект.
// 			rec := Record{
// 				Key:   key,
// 				Value: "dummy", // значение value не имеет значения
// 			}
// 			// Кодируем объект в JSON.
// 			// Вместо того, чтобы сразу записывать в файл, получим JSON-строку.
// 			jsonBytes, err := json.Marshal(rec)
// 			if err != nil {
// 				panic(err)
// 			}
// 			// Формируем строку с префиксом и JSON через табуляцию.
// 			line := fmt.Sprintf("%f\t%s\n", prefix, jsonBytes)
// 			if _, err := tempFile.WriteString(line); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	// Закрываем файл, чтобы убедиться, что все данные записаны.
// 	tempFile.Close()
// 	fmt.Println("Генерация временного файла завершена. Начинается внешняя сортировка...")

// 	// Выполняем внешнюю сортировку файла по первому столбцу (случайному префиксу).
// 	// Параметр "-k1,1n" указывает на числовую сортировку по первому столбцу.
// 	// Результат запишем в итоговый файл "output.jsonl".
// 	outputFileName := "output.jsonl"
// 	sortCmd := exec.Command("sort", "-k1,1n", tempFileName, "-o", outputFileName)
// 	sortCmd.Stdout = os.Stdout
// 	sortCmd.Stderr = os.Stderr

// 	if err := sortCmd.Run(); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Внешняя сортировка завершена. Итоговый файл:", outputFileName)

// 	// Если требуется убрать случайный префикс, можно выполнить дополнительную обработку.
// 	// Например, читаем отсортированный файл и записываем его в новый файл, удаляя префикс.
// 	finalFileName := "final_output.jsonl"
// 	inFile, err := os.Open(outputFileName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer inFile.Close()

// 	outFile, err := os.Create(finalFileName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer outFile.Close()

// 	// Построчная обработка: удаляем часть до первого символа табуляции.
// 	// Можно использовать bufio.Scanner для эффективного чтения.
// 	scanner := NewLineScanner(inFile)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		// Ищем первую табуляцию.
// 		var jsonPart string
// 		for i, ch := range line {
// 			if ch == '\t' {
// 				jsonPart = line[i+1:]
// 				break
// 			}
// 		}
// 		// Если табуляция не найдена, можно записать строку как есть.
// 		if jsonPart == "" {
// 			jsonPart = line
// 		}
// 		_, err := outFile.WriteString(jsonPart + "\n")
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Файл без случайного префикса сохранён как:", finalFileName)

// 	// Опционально: удаляем временные файлы.
// 	os.Remove(tempFileName)
// 	os.Remove(outputFileName)
// }

// // NewLineScanner возвращает сканер, читающий строки из файла.
// // Здесь можно использовать bufio.NewScanner, но приведём обёртку для удобства.
// func NewLineScanner(file *os.File) *bufio.Scanner {
// 	return bufio.NewScanner(file)
// }
