package entities

import "hse_kat/internal/domain"

// A structure that represents a bike
type Bike struct {
	name            string
	inventoryNumber string
	simplicity      int
	status          domain.Status
}

// Constructor for bike
func NewBike(
	name string,
	inventoryNumber string,
	simplicity int,
	status domain.Status,
) Bike {
	return Bike{
		name:            name,
		inventoryNumber: inventoryNumber,
		simplicity:      simplicity,
		status:          status,
	}
}

// Returns name of the bike
func (bike Bike) Name() string {
	return bike.name
}

// Returns inventory number of the bike
func (bike Bike) InventoryNumber() string {
	return bike.inventoryNumber
}

// Returns simplicity of the bike
func (bike Bike) Simplicity() int {
	return bike.simplicity
}

// Returns status of the bike
func (bike Bike) Status() domain.Status {
	return bike.status
}
