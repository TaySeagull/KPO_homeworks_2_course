package main

import "fmt"

func main() {
	container := NewContainer()
	Park := container.Park

	// Create vehicles
	scooter := Scooter{
		name:            "City Scooter",
		inventoryNumber: "SCOOT-001",
		simplicity:      8,
		Energy:          2.5,
		status:          Available,
	}

	bike := Bike{
		name:            "Mountain Bike",
		inventoryNumber: "BIKE-001",
		simplicity:      7,
		status:          Available,
	}

	// Add vehicles to the Park
	fmt.Println("Adding vehicles:")

	fmt.Println(scooter.Name(), Park.AddVehicle(scooter))
	fmt.Println(bike.Name(), Park.AddVehicle(bike))

	// Create equipment
	helmet := Helmet{
		name:            "Safety Helmet",
		inventoryNumber: "HELM-001",
		status:          Available,
	}

	// Add equipment to the Park
	fmt.Println("\nAdding equipment:")

	fmt.Println(helmet.Name(), Park.AddEquipment(helmet))

	// Calculate the energy consumption
	fmt.Printf(
		"\nTotal energy consumption: %.2f\n",
		Park.CalculateDailyEnergy(),
	)

	// Display beginner-friendly vehicles
	fmt.Println("\nBeginner-friendly vehicles:")

	for _, vehicle := range Park.GetBeginnerVehicles() {
		fmt.Printf(
			"%s (simplicity: %d)\n",
			vehicle.Name(),
			vehicle.Simplicity(),
		)
	}

	fmt.Println()
	Park.PrintInventory()
}
