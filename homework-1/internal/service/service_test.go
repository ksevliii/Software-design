package service_test

import (
	"testing"
	"vyshkat/internal/di"
	"vyshkat/internal/inventory"
	"vyshkat/internal/service"
)

// TestInspector - stub for Fleet tests
type FakeInspector struct {
	AcceptAll bool
}

func (f *FakeInspector) IsAccept(_ inventory.Inventarizable) bool {
	return f.AcceptAll
}

// Fleet tests

// Test 1: Dependency substitution through DI
func TestFleet_AddItem_AcceptsWhenInspectorAllows(t *testing.T) {
	fleet := service.NewFleet()
	container := di.NewContainer(&FakeInspector{AcceptAll: true})
	inspector := container.TechInspector
	item := inventory.NewScooter("S01", 8, true)
	result := fleet.AddItem(item, inspector)
	if !result {
		t.Errorf("Expected the correct scooter to be accepted, but AddItem returned false")
	}
}

// Test 2: Duplicate inventory number rejection
func TestFleet_AddItem_RejectsDuplicateInventoryNumber(t *testing.T) {
	fleet := service.NewFleet()
	container := di.NewContainer(&FakeInspector{AcceptAll: true})
	inspector := container.TechInspector
	item1 := inventory.NewScooter("S01", 8, true)
	item2 := inventory.NewScooter("S01", 9, true)
	fleet.AddItem(item1, inspector)
	result := fleet.AddItem(item2, inspector)
	if result {
		t.Errorf("Expected a rejection due to a duplicate inventory number, but AddItem returned true")
	}
}

// Test 3: Total daily energy calculation
func TestFleet_TotalEnergyConsuming_SumsOnlyEnergyConsuming(t *testing.T) {
	fleet := service.NewFleet()
	container := di.NewContainer(&FakeInspector{AcceptAll: true})
	inspector := container.TechInspector
	fleet.AddItem(inventory.NewElectricScooter("E01", 2.0, 8, 80, true), inspector)
	fleet.AddItem(inventory.NewScooter("S01", 8, true), inspector)
	fleet.AddItem(inventory.NewHelmet("H01"), inspector)
	energy := fleet.TotalEnergyConsuming()
	if energy != 2.0 {
		t.Errorf("Expected total energy 2.0, got %f", energy)
	}
}

// Test 4: Beginner-friendly equipment filter
func TestFleet_BeginnerFriendly_FiltersByScore(t *testing.T) {
	fleet := service.NewFleet()
	container := di.NewContainer(&FakeInspector{AcceptAll: true})
	inspector := container.TechInspector
	fleet.AddItem(inventory.NewScooter("S01", 8, true), inspector)
	fleet.AddItem(inventory.NewScooter("S02", 4, true), inspector)
	fleet.AddItem(inventory.NewHelmet("H01"), inspector)
	beginners := fleet.BeginnerFriendly()
	if len(beginners) != 1 {
		t.Errorf("Expected 1 beginner-friendly item, got %d", len(beginners))
	}
}

// Test 5: Rejection when inspector denies
func TestFleet_AddItem_RejectsWhenInspectorDenies(t *testing.T) {
	fleet := service.NewFleet()
	container := di.NewContainer(&FakeInspector{AcceptAll: false})
	inspector := container.TechInspector
	item := inventory.NewScooter("S01", 8, true)
	result := fleet.AddItem(item, inspector)
	if result {
		t.Errorf("Expected rejection when inspector denies, but AddItem returned true")
	}
}

// Test 6: Empty fleet energy calculation
func TestFleet_TotalEnergyConsuming_EmptyFleet(t *testing.T) {
	fleet := service.NewFleet()
	if fleet.TotalEnergyConsuming() != 0 {
		t.Errorf("Expected 0 for empty fleet, got %f", fleet.TotalEnergyConsuming())
	}
}

// Test 7: Empty fleet beginner filter
func TestFleet_BeginnerFriendly_EmptyFleet(t *testing.T) {
	fleet := service.NewFleet()
	beginners := fleet.BeginnerFriendly()
	if len(beginners) != 0 {
		t.Errorf("Expected 0 beginner-friendly items for empty fleet, got %d", len(beginners))
	}
}

// ServiceCenter tests (table-driven tests)

func TestServiceCenter_IsAccept(t *testing.T) {
	tests := []struct {
		name     string
		item     inventory.Inventarizable
		expected bool
	}{
		// Negative cases
		{
			name:     "Rejects short inventory number",
			item:     inventory.NewScooter("AB", 8, true),
			expected: false,
		},
		{
			name:     "Rejects invalid simplicity score",
			item:     inventory.NewScooter("S01", 15, true),
			expected: false,
		},
		{
			name:     "Rejects low battery charge",
			item:     inventory.NewElectricScooter("E01", 1.5, 8, 50, true),
			expected: false,
		},
		{
			name:     "Rejects broken brakes",
			item:     inventory.NewScooter("S02", 8, false),
			expected: false,
		},

		// Positive cases
		{
			name:     "Accepts valid item",
			item:     inventory.NewScooter("S01", 8, true),
			expected: true,
		},
		{
			name:     "Accepts helmet (no charge/brakes required)",
			item:     inventory.NewHelmet("H01"),
			expected: true,
		},
		{
			name:     "Accepts sufficient battery charge",
			item:     inventory.NewElectricScooter("E01", 1.5, 8, 80, true),
			expected: true,
		},

		// Boundary cases
		{
			name:     "Accepts simplicity score at exact minimum (1)",
			item:     inventory.NewScooter("S01", 1, true),
			expected: true,
		},
		{
			name:     "Accepts simplicity score at exact maximum (10)",
			item:     inventory.NewScooter("S02", 10, true),
			expected: true,
		},
		{
			name:     "Accepts charge at exact threshold (70%)",
			item:     inventory.NewElectricScooter("E01", 1.5, 8, 70, true),
			expected: true,
		},
		{
			name:     "Rejects charge just below threshold (69%)",
			item:     inventory.NewElectricScooter("E02", 1.5, 8, 69, true),
			expected: false,
		},
		{
			name:     "Rejects simplicity score below minimum (0)",
			item:     inventory.NewScooter("S03", 0, true),
			expected: false,
		},
		{
			name:     "Rejects simplicity score above maximum (11)",
			item:     inventory.NewScooter("S04", 11, true),
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inspector := service.NewServiceCenter()
			result := inspector.IsAccept(tt.item)
			if result != tt.expected {
				t.Errorf("Failed %v: got IsAccept(%v) = %v, want %v", tt.name, tt.item.InventoryNumber(), result, tt.expected)
			}
		})
	}
}
