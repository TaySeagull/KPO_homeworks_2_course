package container

import "hse_kat/service"

// Composition root for the hse_kat
type Container struct {
	Inspector     service.Inspector
	ServiceCenter *service.ServiceCenter
	Park          *service.Park
}

// Creates a container with inspector, service and park
func NewContainer() *Container {
	inspector := &service.BasicInspector{}

	serviceCenter := service.NewServiceCenter(inspector)

	park := service.NewPark(serviceCenter)

	return &Container{
		Inspector:     inspector,
		ServiceCenter: serviceCenter,
		Park:          park,
	}
}
