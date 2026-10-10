package main

// An interface that defines the methods for an inspector
type Inspector interface {
	Inspect(vehicle Vehicle) bool
}

// A basic implementation of the Inspector interface
type BasicInspector struct{}

// Checks if the vehicle needs inspection and returns true if it does, false otherwise
func (i *BasicInspector) Inspect(vehicle Vehicle) bool {
	return vehicle != nil && vehicle.Status() == Available
}
