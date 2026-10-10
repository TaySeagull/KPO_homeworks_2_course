package main

// A structure that represents a helmet
type Helmet struct {
	name            string
	inventoryNumber string
	status          Status
}

// Returns name of the helmet
func (helmet Helmet) Name() string {
	return helmet.name
}

// Returns inventory number of the helmet
func (helmet Helmet) InventoryNumber() string {
	return helmet.inventoryNumber
}

// Returns status of the helmet
func (helmet Helmet) Status() Status {
	return helmet.status
}
