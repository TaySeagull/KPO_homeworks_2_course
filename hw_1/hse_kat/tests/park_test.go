package tests

import (
	"hse_kat/internal/domain"
	"hse_kat/internal/entities"
	"hse_kat/service"
	"math"
	"testing"
)

// Creates a park for testing
func createTestPark() *service.Park {
	inspector := &service.BasicInspector{}
	serviceCenter := service.NewServiceCenter(inspector)
	return service.NewPark(serviceCenter)
}

// Tests adding a vehicle that passes inspection
func TestAddVehicleAvailable(test *testing.T) {
	park := createTestPark()

	scooter := entities.NewScooter(
		"Scooter 1",
		"SCOOT-001",
		8,
		2.5,
		domain.Available,
	)

	if !park.AddVehicle(scooter) {
		test.Fatal("expected available scooter to be added")
	}

	if got := len(park.GetInventory()); got != 1 {
		test.Fatalf("expected 1 inventory item, got %d", got)
	}
}

// Tests rejecting a vehicle that requires repair
func TestAddVehicleNeedsRepair(test *testing.T) {
	park := createTestPark()

	scooter := entities.NewScooter(
		"Scooter 2",
		"SCOOT-002",
		8,
		2.5,
		domain.NeedsRepair,
	)

	if park.AddVehicle(scooter) {
		test.Fatal("expected vehicle requiring repair to be rejected")
	}

	if got := len(park.GetInventory()); got != 0 {
		test.Fatalf("expected 0 inventory items, got %d", got)
	}
}

// Tests adding equipment to the park
func TestAddEquipment(test *testing.T) {
	park := createTestPark()

	helmet := entities.NewHelmet(
		"Helmet 1",
		"HELMET-001",
		domain.Available,
	)

	if !park.AddEquipment(helmet) {
		test.Fatal("expected helmet to be added")
	}

	if got := len(park.GetInventory()); got != 1 {
		test.Fatalf("expected 1 inventory item, got %d", got)
	}
}

// Tests calculating total daily energy consumption
func TestCalculateDailyEnergy(test *testing.T) {
	park := createTestPark()

	scooter1 := entities.NewScooter(
		"Scooter 1",
		"SCOOT-001",
		8,
		2.5,
		domain.Available,
	)

	scooter2 := entities.NewScooter(
		"Scooter 2",
		"SCOOT-002",
		7,
		3.0,
		domain.Available,
	)

	bike := entities.NewBike(
		"Bike 1",
		"BIKE-001",
		8,
		domain.Available,
	)

	if !park.AddVehicle(scooter1) {
		test.Fatal("expected first scooter to be added")
	}

	if !park.AddVehicle(scooter2) {
		test.Fatal("expected second scooter to be added")
	}

	if !park.AddVehicle(bike) {
		test.Fatal("expected bike to be added")
	}

	got := park.CalculateDailyEnergy()
	want := 5.5

	if math.Abs(got-want) > 1e-9 {
		test.Errorf("expected energy consumption %.2f, got %.2f", want, got)
	}
}

// Tests selecting vehicles suitable for beginners
func TestGetBeginnerVehicles(test *testing.T) {
	park := createTestPark()

	easyBike := entities.NewBike(
		"Easy Bike",
		"BIKE-001",
		8,
		domain.Available,
	)

	difficultBike := entities.NewBike(
		"Difficult Bike",
		"BIKE-002",
		4,
		domain.Available,
	)

	if !park.AddVehicle(easyBike) {
		test.Fatal("expected easy bike to be added")
	}

	if !park.AddVehicle(difficultBike) {
		test.Fatal("expected difficult bike to be added")
	}

	got := park.GetBeginnerVehicles()

	if len(got) != 1 {
		test.Fatalf("expected 1 beginner vehicle, got %d", len(got))
	}

	if got[0].Name() != "Easy Bike" {
		test.Errorf("expected Easy Bike, got %s", got[0].Name())
	}
}

// Tests returning all inventory items
func TestGetInventory(test *testing.T) {
	park := createTestPark()

	scooter := entities.NewScooter(
		"Scooter 1",
		"SCOOT-001",
		8,
		2.5,
		domain.Available,
	)

	helmet := entities.NewHelmet(
		"Helmet 1",
		"HELMET-001",
		domain.Available,
	)

	if !park.AddVehicle(scooter) {
		test.Fatal("expected scooter to be added")
	}

	if !park.AddEquipment(helmet) {
		test.Fatal("expected helmet to be added")
	}

	items := park.GetInventory()

	if len(items) != 2 {
		test.Fatalf("expected 2 inventory items, got %d", len(items))
	}
}

// Tests behavior when calculating energy for a nil park
func TestCalculateDailyEnergyNilPark(test *testing.T) {
	var park *service.Park

	got := park.CalculateDailyEnergy()

	if got != 0 {
		test.Errorf("expected 0 for nil park, got %.2f", got)
	}
}

// Tests that a vehicle with simplicity exactly 6 is suitable for beginners
func TestGetBeginnerVehiclesBoundary(test *testing.T) {
	park := createTestPark()

	bike := entities.NewBike(
		"Boundary Bike",
		"BIKE-006",
		6,
		domain.Available,
	)

	if !park.AddVehicle(bike) {
		test.Fatal("expected bike to be added")
	}

	got := park.GetBeginnerVehicles()

	if len(got) != 1 {
		test.Fatalf("expected 1 beginner vehicle, got %d", len(got))
	}

	if got[0].Name() != "Boundary Bike" {
		test.Errorf("expected Boundary Bike, got %s", got[0].Name())
	}
}

// Tests that a vehicle with simplicity 5 is not suitable for beginners
func TestGetBeginnerVehiclesBelowBoundary(test *testing.T) {
	park := createTestPark()

	bike := entities.NewBike(
		"Hard Bike",
		"BIKE-005",
		5,
		domain.Available,
	)

	if !park.AddVehicle(bike) {
		test.Fatal("expected bike to be added")
	}

	got := park.GetBeginnerVehicles()

	if len(got) != 0 {
		test.Fatalf("expected 0 beginner vehicles, got %d", len(got))
	}
}

// Tests that an empty park has zero energy consumption
func TestCalculateDailyEnergyEmptyPark(test *testing.T) {
	park := createTestPark()

	got := park.CalculateDailyEnergy()

	if got != 0 {
		test.Errorf("expected 0 energy consumption, got %.2f", got)
	}
}

// Tests that an empty park has no beginner vehicles
func TestGetBeginnerVehiclesEmptyPark(test *testing.T) {
	park := createTestPark()

	got := park.GetBeginnerVehicles()

	if len(got) != 0 {
		test.Fatalf("expected 0 beginner vehicles, got %d", len(got))
	}
}

// Tests that an empty park has an empty inventory
func TestGetInventoryEmptyPark(test *testing.T) {
	park := createTestPark()

	got := park.GetInventory()

	if len(got) != 0 {
		test.Fatalf("expected 0 inventory items, got %d", len(got))
	}
}

// Tests that a rejected vehicle does not affect energy consumption
func TestRejectedVehicleDoesNotAffectEnergy(test *testing.T) {
	park := createTestPark()

	scooter := entities.NewScooter(
		"Broken Scooter",
		"SCOOT-BROKEN",
		8,
		10.0,
		domain.NeedsRepair,
	)

	if park.AddVehicle(scooter) {
		test.Fatal("expected broken scooter to be rejected")
	}

	if got := park.CalculateDailyEnergy(); got != 0 {
		test.Errorf("expected 0 energy consumption, got %.2f", got)
	}
}

// Tests that a bike without EnergyConsumer does not affect energy consumption
func TestBikeDoesNotAffectEnergy(test *testing.T) {
	park := createTestPark()

	bike := entities.NewBike(
		"Bike 1",
		"BIKE-001",
		8,
		domain.Available,
	)

	if !park.AddVehicle(bike) {
		test.Fatal("expected bike to be added")
	}

	if got := park.CalculateDailyEnergy(); got != 0 {
		test.Errorf("expected 0 energy consumption, got %.2f", got)
	}
}

// Tests adding several beginner vehicles
func TestGetSeveralBeginnerVehicles(test *testing.T) {
	park := createTestPark()

	bike := entities.NewBike(
		"Easy Bike",
		"BIKE-001",
		7,
		domain.Available,
	)

	scooter := entities.NewScooter(
		"Easy Scooter",
		"SCOOT-001",
		10,
		2.5,
		domain.Available,
	)

	hardBike := entities.NewBike(
		"Hard Bike",
		"BIKE-002",
		3,
		domain.Available,
	)

	if !park.AddVehicle(bike) {
		test.Fatal("expected bike to be added")
	}

	if !park.AddVehicle(scooter) {
		test.Fatal("expected scooter to be added")
	}

	if !park.AddVehicle(hardBike) {
		test.Fatal("expected hard bike to be added")
	}

	got := park.GetBeginnerVehicles()

	if len(got) != 2 {
		test.Fatalf("expected 2 beginner vehicles, got %d", len(got))
	}
}

// Tests returning nil for inventory of a nil park
func TestGetInventoryNilPark(test *testing.T) {
	var park *service.Park

	got := park.GetInventory()

	if got != nil {
		test.Errorf("expected nil inventory for nil park, got %v", got)
	}
}

// Tests returning nil for beginner vehicles of a nil park
func TestGetBeginnerVehiclesNilPark(test *testing.T) {
	var park *service.Park

	got := park.GetBeginnerVehicles()

	if got != nil {
		test.Errorf("expected nil beginner vehicles for nil park, got %v", got)
	}
}

// Injection of fakeInspector

type fakeInspector struct {
	result bool
}

func (inspector *fakeInspector) Inspect(vehicle domain.Vehicle) bool {
	return inspector.result
}

func TestAddVehicleWithFakeInspector(test *testing.T) {
	fake := &fakeInspector{result: false}
	serviceCenter := service.NewServiceCenter(fake)
	park := service.NewPark(serviceCenter)

	scooter := entities.NewScooter(
		"Scooter",
		"SC-001",
		8,
		2.5,
		domain.Available,
	)

	added := park.AddVehicle(scooter)

	if added {
		test.Error("expected vehicle not to be added")
	}

	if len(park.GetInventory()) != 0 {
		test.Error("expected park to remain empty")
	}
}
