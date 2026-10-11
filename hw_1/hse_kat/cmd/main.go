package main

import (
	"fmt"
	"hse_kat/container"
)

func main() {
	app := container.NewContainer()
	park := app.Park

	for {
		printMainMenu()

		choice, ok := readInt("Выберите действие: ")
		if !ok {
			fmt.Println("\nЗавершение программы.")
			return
		}

		switch choice {
		case 1:
			printInventory(park.GetInventory())

		case 2:
			addVehicle(park)

		case 3:
			addEquipment(park)

		case 4:
			fmt.Printf(
				"\nОбщее ежедневное энергопотребление: %.2f\n",
				park.CalculateDailyEnergy(),
			)

		case 5:
			fmt.Println("\nТранспорт для новичков:")

			vehicles := park.GetBeginnerVehicles()

			if len(vehicles) == 0 {
				fmt.Println("Подходящий транспорт не найден.")
				continue
			}

			for _, vehicle := range vehicles {
				fmt.Printf(
					"%s (простота: %d)\n",
					vehicle.Name(),
					vehicle.Simplicity(),
				)
			}

		case 0:
			fmt.Println("Спасибо за использование! До свидания!")
			return

		default:
			fmt.Println("Неизвестная команда. Попробуйте ещё раз.")
		}
	}
}
