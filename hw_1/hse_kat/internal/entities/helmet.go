package entities

import "hse_kat/internal/domain"

// A structure that represents a helmet
type Helmet struct {
	name            string
	inventoryNumber string
	status          domain.Status
}

// Constructor for helmet
func NewHelmet(
	name string,
	inventoryNumber string,
	status domain.Status,
) Helmet {
	return Helmet{
		name:            name,
		inventoryNumber: inventoryNumber,
		status:          status,
	}
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
func (helmet Helmet) Status() domain.Status {
	return helmet.status
}
