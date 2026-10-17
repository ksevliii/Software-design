package inventory_test

import (
	"testing"

	"vyshkat/internal/inventory"
)

// TestInterfaceContracts - interface correctness
func TestInterfaceContracts(t *testing.T) {
	tests := []struct {
		name                 string
		item                 inventory.Inventarizable
		implementsEnergy     bool
		implementsSimplicity bool
		implementsCharge     bool
		implementsBrakes     bool
	}{
		{
			name:                 "ElectricScooter implements all 4 optional interfaces",
			item:                 inventory.NewElectricScooter("E01", 1.5, 8, 80, true),
			implementsEnergy:     true,
			implementsSimplicity: true,
			implementsCharge:     true,
			implementsBrakes:     true,
		},
		{
			name:                 "Scooter implements Beginner and Brakes",
			item:                 inventory.NewScooter("S01", 8, true),
			implementsEnergy:     false,
			implementsSimplicity: true,
			implementsCharge:     false,
			implementsBrakes:     true,
		},
		{
			name:                 "Helmet implements only Inventarizable",
			item:                 inventory.NewHelmet("H01"),
			implementsEnergy:     false,
			implementsSimplicity: false,
			implementsCharge:     false,
			implementsBrakes:     false,
		},
		{
			name:                 "ChargingStation implements Energy",
			item:                 inventory.NewChargingStation("CS01", 5.0),
			implementsEnergy:     true,
			implementsSimplicity: false,
			implementsCharge:     false,
			implementsBrakes:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, isEnergy := tt.item.(inventory.EnergyConsuming)
			_, isBeginner := tt.item.(inventory.SimplicityEvaluated)
			_, isCharge := tt.item.(inventory.Chargeable)
			_, isBrakes := tt.item.(inventory.Braked)
			if isEnergy != tt.implementsEnergy {
				t.Errorf("EnergyConsuming: got %v, want %v", isEnergy, tt.implementsEnergy)
			}
			if isBeginner != tt.implementsSimplicity {
				t.Errorf("SimplicityEvaluated: got %v, want %v", isBeginner, tt.implementsSimplicity)
			}
			if isCharge != tt.implementsCharge {
				t.Errorf("Chargeable: got %v, want %v", isCharge, tt.implementsCharge)
			}
			if isBrakes != tt.implementsBrakes {
				t.Errorf("Braked: got %v, want %v", isBrakes, tt.implementsBrakes)
			}
		})
	}
}
