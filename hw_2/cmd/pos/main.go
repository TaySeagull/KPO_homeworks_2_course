package main

import (
	"fmt"
	"os"

	"hw_2/internal/container"
	"hw_2/internal/presentation/console"
)

func main() {
	applicationContainer, err := container.New(container.Config{
		StudentsFile: "data/students.json",
		RubricFile:   "data/rubric.json",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Не удалось запустить приложение:", err)
		os.Exit(1)
	}

	defer applicationContainer.Shutdown()

	consoleApp := console.New(
		applicationContainer.Facade,
		os.Stdin,
		os.Stdout,
	)

	consoleApp.Run()
}
