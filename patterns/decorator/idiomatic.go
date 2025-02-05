package decorator

import (
	"fmt"
	"time"
)

// Структура с методом, который мы будем декорировать
type MyStruct struct {
	Value int
}

// Метод, который мы будем декорировать
func (m *MyStruct) Multiply(factor int) int {
	return m.Value * factor
}

// Тип функции, представляющей метод Multiply
type MultiplyFuncType func(*MyStruct, int) int

// Декоратор, добавляющий логирование времени выполнения
func timingDecorator(fn MultiplyFuncType) MultiplyFuncType {
	return func(m *MyStruct, factor int) int {
		start := time.Now()
		result := fn(m, factor)
		fmt.Printf("Execution time: %s\n", time.Since(start))
		return result
	}
}

func main() {
	// Создаем экземпляр структуры
	myStruct := &MyStruct{Value: 10}

	// Создаем декорированную версию метода Multiply
	decoratedMultiply := timingDecorator((*MyStruct).Multiply)

	// Вызываем декорированный метод
	result := decoratedMultiply(myStruct, 5)
	fmt.Printf("Result: %d\n", result)
}
