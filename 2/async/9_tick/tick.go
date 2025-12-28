package main

import (
	"fmt"
	"math/rand"
	"time"
)

// accumulateAndFlush аккумулирует значения из inChan и каждые 3 секунды отправляет сумму в outChan, обнуляя аккумулятор.
func accumulateAndFlush(inChan <-chan int, outChan chan<- int) {
	var sum int
	ticker := time.NewTicker(3 * time.Second) // Таймер на 3 секунды
	defer ticker.Stop()

	for {
		select {
		case val, ok := <-inChan:
			if !ok {
				// Если inChan закрыт, отправляем последнюю сумму и выходим.
				outChan <- sum
				close(outChan)
				return
			}
			sum += val // Аккумулируем значение
		case <-ticker.C:
			// Прошло 3 секунды: отправляем сумму и обнуляем
			outChan <- sum
			sum = 0
		}
	}
}

// generator генерирует случайные числа и отправляет в inChan.
func generator(inChan chan<- int) {
	defer close(inChan)       // Закрываем канал после завершения
	for i := 0; i < 20; i++ { // Генерируем 20 значений для примера (можно изменить)
		val := rand.Intn(10) + 1 // Случайное от 1 до 10
		inChan <- val
		time.Sleep(500 * time.Millisecond) // Задержка для симуляции
	}
}

func main() {
	inChan := make(chan int)
	outChan := make(chan int)

	// Запускаем аккумулятор
	go accumulateAndFlush(inChan, outChan)

	// Запускаем генератор
	go generator(inChan)

	// Читаем из outChan и выводим суммы (для демонстрации)
	for sum := range outChan {
		fmt.Printf("Накопленная сумма за 3 секунды: %d\n", sum)
	}
}
