package main

import (
	"fmt"
	"hse_kat/internal/domain"
	"hse_kat/internal/entities"
)

// A type of vehicle available for creation
type vehicleType struct {
	name      string
	readExtra func() (float64, bool)
	create    func(string, string, int, domain.Status, float64) domain.Vehicle
}

// A type of equipment available for creation
type equipmentType struct {
	name   string
	create func(string, string) domain.Equipment
}

// Returns all registered vehicle types
func getVehicleTypes() []vehicleType {
	return []vehicleType{
		{
			name: "Самокат",

			readExtra: func() (float64, bool) {
				for {
					energy, ok := readFloat(
						"Ежедневное энергопотребление: ",
					)
					if !ok {
						return 0, false
					}

					if energy < 0 {
						fmt.Println(
							"Энергопотребление не может быть отрицательным.",
						)
						continue
					}

					return energy, true
				}
			},

			create: func(
				name string,
				inventoryNumber string,
				simplicity int,
				status domain.Status,
				energy float64,
			) domain.Vehicle {
				return entities.NewScooter(
					name,
					inventoryNumber,
					simplicity,
					energy,
					status,
				)
			},
		},
		{
			name: "Велосипед",

			readExtra: nil,

			create: func(
				name string,
				inventoryNumber string,
				simplicity int,
				status domain.Status,
				_ float64,
			) domain.Vehicle {
				return entities.NewBike(
					name,
					inventoryNumber,
					simplicity,
					status,
				)
			},
		},
	}
}

// getEquipmentTypes returns all registered equipment types.
func getEquipmentTypes() []equipmentType {
	return []equipmentType{
		{
			name: "Шлем",

			create: func(
				name string,
				inventoryNumber string,
			) domain.Equipment {
				return entities.NewHelmet(
					name,
					inventoryNumber,
					domain.Available,
				)
			},
		},
	}
}
