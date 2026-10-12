package main

import (
	"bufio"
	"fmt"
	"hse_kat/internal/domain"
	"os"
	"strconv"
	"strings"
)

// new scaner
var scanner = bufio.NewScanner(os.Stdin)

// Reads a string from the console
func readText(prompt string) (string, bool) {
	fmt.Print(prompt)

	if !scanner.Scan() {
		return "", false
	}

	return strings.TrimSpace(scanner.Text()), true
}

// Reads an integer from the console
func readInt(prompt string) (int, bool) {
	for {
		text, ok := readText(prompt)
		if !ok {
			return 0, false
		}

		value, err := strconv.Atoi(text)
		if err == nil {
			return value, true
		}

		fmt.Println("Введите целое число.")
	}
}

// Reads a float number from the console
func readFloat(prompt string) (float64, bool) {
	for {
		text, ok := readText(prompt)
		if !ok {
			return 0, false
		}

		text = strings.ReplaceAll(text, ",", ".")

		value, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return value, true
		}

		fmt.Println("Введите число, например 2.5.")
	}
}

// Reads a valid simplicity value from 1 to 10
func readSimplicity() (int, bool) {
	for {
		value, ok := readInt("Простота использования (1–10): ")
		if !ok {
			return 0, false
		}

		if value >= 1 && value <= 10 {
			return value, true
		}

		fmt.Println("Простота должна быть от 1 до 10.")
	}
}

// Reads the technical condition of a vehicle
func readStatus() (domain.Status, bool) {
	printStatusMenu()

	for {
		choice, ok := readInt("Выберите состояние: ")
		if !ok {
			return domain.Unknown, false
		}

		switch choice {
		case 1:
			return domain.Available, true
		case 2:
			return domain.NeedsRepair, true
		default:
			fmt.Println("Выберите 1 или 2.")
		}
	}
}
