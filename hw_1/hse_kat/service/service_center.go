package service

import "hse_kat/internal/domain"

// A structure that represents a service center
type ServiceCenter struct {
	inspector Inspector
}

// Creates a new service center with the given inspector
// DOESNT KNOW ABOUNT ITS REALIZATION
func NewServiceCenter(inspector Inspector) *ServiceCenter {
	return &ServiceCenter{
		inspector: inspector,
	}
}

// Inspects the given vehicle using the service center's inspector
func (s *ServiceCenter) Inspect(vehicle domain.Vehicle) bool {
	return s.inspector.Inspect(vehicle)
}
