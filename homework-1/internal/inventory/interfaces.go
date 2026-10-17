package inventory

// for any inventory (not for employees!!!)
type Inventarizable interface {
	InventoryNumber() string
	Type() string
}

// only for items that consume energy
type EnergyConsuming interface {
	DailyEnergy() float64
}

// to evaluate the ease of use
type SimplicityEvaluated interface {
	SimplicityScore() int // must be from 1 (complex) to 10 (easy)
}

// for devices with a battery (charge level 0-100%)
// separated from EnergyConsuming because not all electric items
// have a battery (example: ChargingStation)
type Chargeable interface {
	ChargePercent() int
}

// for wheeled devices with a braking system
type Braked interface {
	BrakesWorking() bool
}
