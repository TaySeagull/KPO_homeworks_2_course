package main

import "fmt"

// A structure that represents a park of hse_kat
type Park struct {
	vehicles      []Vehicle
	equipment     []Equipment
	serviceCenter *ServiceCenter
}

// Creates a new park with the given service center
func NewPark(serviceCenter *ServiceCenter) *Park {
	return &Park{
		serviceCenter: serviceCenter,
	}
}

// Adds a vehicle to the park if it passes inspection
func (park *Park) AddVehicle(vehicle Vehicle) bool {
	if park == nil || vehicle == nil || park.serviceCenter == nil {
		return false
	}

	if !park.serviceCenter.Inspect(vehicle) {
		return false
	}

	park.vehicles = append(park.vehicles, vehicle)
	return true
}

// Adds equipment to the park
func (park *Park) AddEquipment(equipment Equipment) bool {
	if park == nil || equipment == nil {
		return false
	}

	park.equipment = append(park.equipment, equipment)
	return true
}

// Calculates the Energy
func (park *Park) CalculateDailyEnergy() float64 {
	if park == nil {
		return 0
	}

	var total float64

	for _, vehicle := range park.vehicles {
		if consumer, ok := vehicle.(EnergyConsumer); ok {
			total += consumer.EnergyConsumption()
		}
	}

	for _, equipment := range park.equipment {
		if consumer, ok := equipment.(EnergyConsumer); ok {
			total += consumer.EnergyConsumption()
		}
	}

	return total
}

// Gets the list of the beginner friendly bikes in the park
func (park *Park) GetBeginnerVehicles() []Vehicle {
	if park == nil {
		return nil
	}

	var result []Vehicle

	for _, vehicle := range park.vehicles {
		if vehicle.Simplicity() >= 6 {
			result = append(result, vehicle)
		}
	}

	return result
}

// Prints the inventory of the park
func (park *Park) PrintInventory() {
	if park == nil {
		return
	}

	fmt.Println("=== Park Inventory ===")

	fmt.Println("\nVehicles:")
	for _, vehicle := range park.vehicles {
		fmt.Printf("%s | Inventory number: %s\n",
			vehicle.Name(),
			vehicle.InventoryNumber(),
		)
	}

	fmt.Println("\nEquipment:")
	for _, equipment := range park.equipment {
		fmt.Printf("%s | Inventory number: %s\n",
			equipment.Name(),
			equipment.InventoryNumber(),
		)
	}
}
