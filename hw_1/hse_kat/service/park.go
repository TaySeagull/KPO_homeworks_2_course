package service

import (
	"hse_kat/internal/domain"
)

// A structure that represents a park of hse_kat
type Park struct {
	vehicles      []domain.Vehicle
	equipment     []domain.Equipment
	serviceCenter *ServiceCenter
}

// Creates a new park with the given service center
func NewPark(serviceCenter *ServiceCenter) *Park {
	return &Park{
		serviceCenter: serviceCenter,
	}
}

// Adds a vehicle to the park if it passes inspection
func (park *Park) AddVehicle(vehicle domain.Vehicle) bool {
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
func (park *Park) AddEquipment(equipment domain.Equipment) bool {
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
		if consumer, ok := vehicle.(domain.EnergyConsumer); ok {
			total += consumer.EnergyConsumption()
		}
	}

	for _, equipment := range park.equipment {
		if consumer, ok := equipment.(domain.EnergyConsumer); ok {
			total += consumer.EnergyConsumption()
		}
	}

	return total
}

// Gets the list of the beginner friendly bikes in the park
func (park *Park) GetBeginnerVehicles() []domain.Vehicle {
	if park == nil {
		return nil
	}

	var result []domain.Vehicle

	for _, vehicle := range park.vehicles {
		if vehicle.Simplicity() >= 6 {
			result = append(result, vehicle)
		}
	}

	return result
}

// Returns all inventory items in the park
func (park *Park) GetInventory() []domain.InventoryItem {
	if park == nil {
		return nil
	}

	items := make(
		[]domain.InventoryItem,
		0,
		len(park.vehicles)+len(park.equipment),
	)

	for _, vehicle := range park.vehicles {
		items = append(items, vehicle)
	}

	for _, equipment := range park.equipment {
		items = append(items, equipment)
	}

	return items
}
