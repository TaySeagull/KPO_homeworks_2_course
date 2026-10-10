package main

// A structure that represents a bike
type Bike struct {
	name            string
	inventoryNumber string
	simplicity      int
	status          Status
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
func (bike Bike) Status() Status {
	return bike.status
}
