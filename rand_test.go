package main_test

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Skip("experiments with uniform distribution of random numbers")
	rand.NewSource(time.Now().UnixNano())

	numbers := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		numbers = append(numbers, rand.Intn(200))
	}

	// Определяем границы отрезков
	buckets := []struct {
		min int
		max int
	}{
		{1, 20},
		{21, 40},
		{41, 80},
		{81, 120},
		{121, 160},
		{161, 200},
	}

	// Создаем срез для хранения количества чисел в каждом отрезке
	counts := make([]int, len(buckets))

	// Считаем количество чисел в каждом отрезке
	for _, num := range numbers {
		for i, bucket := range buckets {
			if num >= bucket.min && num < bucket.max {
				counts[i]++
				break
			}
		}
	}

	// Ожидаемое количество чисел в каждом отрезке
	expectedCount := len(numbers) / len(buckets)

	// Проверяем, что количество чисел в каждом отрезке близко к ожидаемому
	for i, count := range counts {
		if count < expectedCount-10 || count > expectedCount+10 {
			t.Errorf("Количество чисел в отрезке %d-%d: %d, ожидаемое: %d±10",
				buckets[i].min, buckets[i].max, count, expectedCount)
		}
	}

}

func TestRandom2(t *testing.T) {
	t.Skip("experiments with uniform distribution of random numbers")
	rand.NewSource(time.Now().UnixNano())

	numbers := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		numbers = append(numbers, rand.Intn(4))
	}

	// Определяем границы отрезков
	buckets := []struct {
		min int
		max int
	}{
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
	}

	// Создаем срез для хранения количества чисел в каждом отрезке
	counts := make([]int, len(buckets))

	// Считаем количество чисел в каждом отрезке
	for _, num := range numbers {
		for i, bucket := range buckets {
			if num >= bucket.min && num < bucket.max {
				counts[i]++
				break
			}
		}
	}

	// Ожидаемое количество чисел в каждом отрезке
	expectedCount := len(numbers) / len(buckets)

	// Проверяем, что количество чисел в каждом отрезке близко к ожидаемому
	for i, count := range counts {
		t.Logf("OK: Количество чисел в отрезке %d-%d: %d, ожидаемое: %d±10", buckets[i].min, buckets[i].max, count, expectedCount)
		if count < expectedCount-10 || count > expectedCount+10 {
			t.Errorf("WRONG: Количество чисел в отрезке %d-%d: %d, ожидаемое: %d±10",
				buckets[i].min, buckets[i].max, count, expectedCount)
		}
	}

}

type KeyConfig struct {
	Min int64
	Max int64
}

func ParseKeyFlag(flag *string) KeyConfig {
	if flag == nil {
		return KeyConfig{}
	}

	// -sk <number>
	max, err := strconv.ParseInt(*flag, 10, 64)
	if err == nil {
		return KeyConfig{
			Min: 1,
			Max: max,
		}
	}

	// -sk <min>-<max>
	minMax := strings.Split(*flag, "-")
	if len(minMax) == 2 {
		min, err1 := strconv.ParseInt(minMax[0], 10, 64)
		max, err2 := strconv.ParseInt(minMax[1], 10, 64)
		if err1 == nil && err2 == nil {
			return KeyConfig{
				Min: min,
				Max: max,
			}
		}
	}

	return KeyConfig{}
}

func TestParseKeyFlag(t *testing.T) {
	tests := []struct {
		flag   *string
		config KeyConfig
	}{
		{nil, KeyConfig{}},
		{new(string), KeyConfig{}},
		{func() *string { i := "10"; return &i }(), KeyConfig{Min: 1, Max: 10}},
		{func() *string { i := "2-20"; return &i }(), KeyConfig{Min: 2, Max: 20}},
	}
	for _, test := range tests {
		result := ParseKeyFlag(test.flag)
		assert.Equal(t, test.config, result, "ParseKeyFlag(%v) = %v, want %v", test.flag, result, test.config)
	}
}
