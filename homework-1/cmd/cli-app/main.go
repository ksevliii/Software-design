package main

import (
	"vyshkat/internal/di"
	"vyshkat/internal/service"
)

func main() {
	realInspector := service.NewServiceCenter()
	container := di.NewContainer(realInspector)
	fleet := service.NewFleet()
	runDemo(container.TechInspector)
	runInteractive(fleet, container.TechInspector)
}
