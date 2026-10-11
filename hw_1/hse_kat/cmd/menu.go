package main

import (
	"fmt"

	"hse_kat/internal/domain"
	"hse_kat/service"
)

// Adds a vehicle selected from the registry
func addVehicle(park *service.Park) {
	types := getVehicleTypes()

	printVehicleMenu(types)

	choice, ok := readInt("Ваш выбор: ")
	if !ok || choice == 0 {
		return
	}

	if choice < 1 || choice > len(types) {
		fmt.Println("Неизвестный тип транспорта.")
		return
	}

	selected := types[choice-1]

	name, ok := readText("Название: ")
	if !ok {
		return
	}

	inventoryNumber, ok := readText("Инвентарный номер: ")
	if !ok {
		return
	}

	simplicity, ok := readSimplicity()
	if !ok {
		return
	}

	status, ok := readStatus()
	if !ok {
		return
	}

	var extra float64

	if selected.readExtra != nil {
		extra, ok = selected.readExtra()
		if !ok {
			return
		}
	}

	var vehicle domain.Vehicle

	vehicle = selected.create(
		name,
		inventoryNumber,
		simplicity,
		status,
		extra,
	)

	if park.AddVehicle(vehicle) {
		fmt.Printf("Транспорт %s успешно добавлен в парк!\n", name)
	} else {
		fmt.Printf(
			"Транспорт %s не прошёл техническую проверку.\n",
			name,
		)
	}
}

// Adds equipment selected from the registry
func addEquipment(park *service.Park) {
	types := getEquipmentTypes()

	printEquipmentMenu(types)

	choice, ok := readInt("Ваш выбор: ")
	if !ok || choice == 0 {
		return
	}

	if choice < 1 || choice > len(types) {
		fmt.Println("Неизвестный тип оборудования.")
		return
	}

	selected := types[choice-1]

	name, ok := readText("Название: ")
	if !ok {
		return
	}

	inventoryNumber, ok := readText("Инвентарный номер: ")
	if !ok {
		return
	}

	equipment := selected.create(name, inventoryNumber)

	if park.AddEquipment(equipment) {
		fmt.Printf("Оборудование %s добавлено в парк!\n", name)
	} else {
		fmt.Println("Не удалось добавить оборудование.")
	}
}
