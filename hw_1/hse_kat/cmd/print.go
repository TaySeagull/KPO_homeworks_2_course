package main

import (
	"fmt"
	"hse_kat/internal/domain"
)

// Prints the main menu
func printMainMenu() {
	fmt.Println("\n=== ВышКат: управление парком ===")
	fmt.Println("1. Показать инвентарь")
	fmt.Println("2. Добавить транспорт")
	fmt.Println("3. Добавить оборудование")
	fmt.Println("4. Рассчитать энергопотребление")
	fmt.Println("5. Показать транспорт для новичков")
	fmt.Println("0. Выход")
}

// Prints vehicle types
func printVehicleMenu(types []vehicleType) {
	fmt.Println("\nКакой транспорт добавить?")

	for i, vehicleType := range types {
		fmt.Printf("%d. %s\n", i+1, vehicleType.name)
	}

	fmt.Println("0. Назад")
}

// Prints equipment types
func printEquipmentMenu(types []equipmentType) {
	fmt.Println("\nКакое оборудование добавить?")

	for i, equipmentType := range types {
		fmt.Printf("%d. %s\n", i+1, equipmentType.name)
	}

	fmt.Println("0. Назад")
}

// Prints technical condition options
func printStatusMenu() {
	fmt.Println("\nТехническое состояние:")
	fmt.Println("1. Исправен")
	fmt.Println("2. Требует ремонта")
}

// Prints all inventory items
func printInventory(items []domain.InventoryItem) {
	fmt.Println("\n=== Инвентарь парка ===")

	if len(items) == 0 {
		fmt.Println("Инвентарь пуст.")
		return
	}

	for _, item := range items {
		fmt.Printf(
			"%s | Инвентарный номер: %s\n",
			item.Name(),
			item.InventoryNumber(),
		)
	}
}
