package service

import "vyshkat/internal/inventory"

// transport fleet 
type Fleet struct {
	inventory []inventory.Inventarizable
}

func NewFleet() *Fleet {
	return &Fleet{inventory: make([]inventory.Inventarizable, 0)}
}

func (f *Fleet) AddItem(item inventory.Inventarizable, inspector TechInspector) bool {
	if !inspector.IsAccept(item) {
		return false
	}
	// i want inventoryNumber to be unique so i should check it
	for _, accepted := range f.inventory {
		if accepted.InventoryNumber() == item.InventoryNumber() {
			return false
		}
	}
	f.inventory = append(f.inventory, item) // yeah we finally did it
	return true
}

// consuming per day
func (f *Fleet) TotalEnergyConsuming() float64 {
	totalEnergy := 0.0
	for _, item := range f.inventory {
		if energyConsumer, ok := item.(inventory.EnergyConsuming); ok {
			totalEnergy += energyConsumer.DailyEnergy()
		}
	}
	return totalEnergy
}

// help for beginners
// *он просто новичок в самокатах*
func (f *Fleet) BeginnerFriendly() []inventory.Inventarizable {
	forBeginners := make([]inventory.Inventarizable, 0)
	for _, item := range f.inventory {
		if variable, ok := item.(inventory.SimplicityEvaluated); ok {
			if variable.SimplicityScore() >= 6 {
				forBeginners = append(forBeginners, item)
			}
		}
	}
	return forBeginners
}

// swag items!!!
func (f *Fleet) AllInventory() []inventory.Inventarizable {
	return f.inventory
}

// how many swag items!!!
func (f *Fleet) InventoryCount() int {
	return len(f.inventory)
}
