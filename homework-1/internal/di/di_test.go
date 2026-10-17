package di_test

import (
	"testing"
	"vyshkat/internal/di"
	"vyshkat/internal/inventory"
)

// TestInspector - stub for DI container tests
type FakeInspector struct {
	id        int
	AcceptAll bool
}

func (f *FakeInspector) IsAccept(_ inventory.Inventarizable) bool {
	return f.AcceptAll
}

// Test 1: container creation with dependency
func TestDIContainer_RegisterAndGetDependency(t *testing.T) {
	inspector := &FakeInspector{id: 1}
	container := di.NewContainer(inspector)
	if container.TechInspector == nil {
		t.Fatalf("Expected TechInspector to be set, got nil")
	}
	if container.TechInspector != inspector {
		t.Errorf("Expected container to hold the passed inspector instance, got %v, want %v", container.TechInspector, inspector)
	}
}

// Test 2: dependency substitution with fake
func TestDIContainer_GetDependency_NotRegistered(t *testing.T) {
	container := di.NewContainer(nil)
	if container.TechInspector != nil {
		t.Errorf("Expected container TechInspector to be nil when not registered, got %v", container.TechInspector)
	}
}

// Test 3: different containers are isolated
func TestDIContainer_RegisterOverwritesPrevious(t *testing.T) {
	inspector1 := &FakeInspector{id: 1}
	inspector2 := &FakeInspector{id: 2}
	container1 := di.NewContainer(inspector1)
	container2 := di.NewContainer(inspector2)
	if container1.TechInspector == container2.TechInspector {
		t.Error("Expected different containers to hold different inspector instances")
	}
	if container1.TechInspector != inspector1 {
		t.Errorf("Expected container1 to hold inspector1, got %v, want %v", container1.TechInspector, inspector1)
	}
	if container2.TechInspector != inspector2 {
		t.Errorf("Expected container2 to hold inspector2, got %v, want %v", container2.TechInspector, inspector2)
	}
}
