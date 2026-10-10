package main

// Composition root for the hse_kat
type Container struct {
	Inspector     Inspector
	ServiceCenter *ServiceCenter
	Park          *Park
}

func NewContainer() *Container {
	inspector := &BasicInspector{}

	serviceCenter := NewServiceCenter(inspector)

	park := NewPark(serviceCenter)

	return &Container{
		Inspector:     inspector,
		ServiceCenter: serviceCenter,
		Park:          park,
	}
}
