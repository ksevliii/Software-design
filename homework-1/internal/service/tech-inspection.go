package service

import "vyshkat/internal/inventory"

type TechInspector interface {
	IsAccept(item inventory.Inventarizable) bool
}

type ServiceCenter struct{}

func NewServiceCenter() *ServiceCenter {
	return &ServiceCenter{}
}

// first checks data format validity, then checks serviceability
func (s *ServiceCenter) IsAccept(item inventory.Inventarizable) bool {
	if !s.isValidFormat(item) {
		return false
	}
	return s.serviceability(item)
}

// users can do strange things so we need to check input data
func (s *ServiceCenter) isValidFormat(item inventory.Inventarizable) bool {
	invNum := item.InventoryNumber()
	if len(invNum) < 3 {
		return false
	}
	if simple, ok := item.(inventory.SimplicityEvaluated); ok {
		score := simple.SimplicityScore()
		if score < 1 || score > 10 {
			return false
		}
	}
	return true
}

// do you want to ride a bike without brakes? I don’t, so i check this case
func (s *ServiceCenter) serviceability(item inventory.Inventarizable) bool {
	if chargeable, ok := item.(inventory.Chargeable); ok {
		if chargeable.ChargePercent() < 70 || chargeable.ChargePercent() > 100 {
			return false
		}
	}
	if braked, ok := item.(inventory.Braked); ok {
		if !braked.BrakesWorking() {
			return false
		}
	}
	return true
}
