package domain

// An interface that defines the methods for a vehicle
type Vehicle interface {
	InventoryItem
	Simplicity() int
	Status() Status
}

// An interface that defines the methods for equipment
type Equipment interface {
	InventoryItem
	Status() Status
}
