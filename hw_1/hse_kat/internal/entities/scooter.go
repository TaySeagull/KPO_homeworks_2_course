package entities

import "hse_kat/internal/domain"

// An struct that represents a scooter
type Scooter struct {
	name            string
	inventoryNumber string
	simplicity      int
	Energy          float64
	status          domain.Status
}

// Constructor for scooter
func NewScooter(
	name string,
	inventoryNumber string,
	simplicity int,
	Energy float64,
	status domain.Status,
) Scooter {
	return Scooter{
		name:            name,
		inventoryNumber: inventoryNumber,
		simplicity:      simplicity,
		Energy:          Energy,
		status:          status,
	}
}

// Returns name of the scooter
func (scooter Scooter) Name() string {
	return scooter.name
}

// Returns inventory number of the scooter
func (scooter Scooter) InventoryNumber() string {
	return scooter.inventoryNumber
}

// Returns simplicity of the scooter
func (scooter Scooter) Simplicity() int {
	return scooter.simplicity
}

// Returns status of the scooter
func (scooter Scooter) Status() domain.Status {
	return scooter.status
}

// Returns daily energy consumption of the scooter
func (scooter Scooter) EnergyConsumption() float64 {
	return scooter.Energy
}
