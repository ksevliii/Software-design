package di

import "vyshkat/internal/service"

type Container struct {
	TechInspector service.TechInspector
	// future possible dependencies
}

func NewContainer(inspector service.TechInspector) *Container {
	return &Container{
		TechInspector: inspector,
	}
}
