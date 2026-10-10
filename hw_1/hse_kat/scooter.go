package main

// An struct that represents a scooter
type Scooter struct {
	name            string
	inventoryNumber string
	simplicity      int
	Energy          float64
	status          Status
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
func (scooter Scooter) Status() Status {
	return scooter.status
}

// Returns daily energy consumption of the scooter
func (scooter Scooter) EnergyConsumption() float64 {
	return scooter.Energy
}
