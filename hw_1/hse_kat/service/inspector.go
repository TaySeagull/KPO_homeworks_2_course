package service

import "hse_kat/internal/domain"

// An interface that defines the methods for an inspector
type Inspector interface {
	Inspect(vehicle domain.Vehicle) bool
}

// A basic implementation of the Inspector interface
type BasicInspector struct{}

// Checks if the vehicle needs inspection and returns true if it does, false otherwise
func (i *BasicInspector) Inspect(vehicle domain.Vehicle) bool {
	return vehicle != nil && vehicle.Status() == domain.Available
}
