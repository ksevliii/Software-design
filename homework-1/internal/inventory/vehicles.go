package inventory

// scooter (самокат): inventarizable, beginner-friendly, braked
type Scooter struct {
	inventoryNumber string
	simplicityScore int
	brakesWorking   bool
}

func NewScooter(invNum string, score int, brakesWorking bool) *Scooter {
	return &Scooter{inventoryNumber: invNum, simplicityScore: score, brakesWorking: brakesWorking}
}

func (s *Scooter) InventoryNumber() string { return s.inventoryNumber }
func (s *Scooter) Type() string            { return "Самокат" }
func (s *Scooter) SimplicityScore() int    { return s.simplicityScore }
func (s *Scooter) BrakesWorking() bool     { return s.brakesWorking }

// electric scooter (электросамокат): inventarizable, energy consumer, beginner-friendly, chargeable, braked
type ElectricScooter struct {
	inventoryNumber string
	dailyEnergy     float64
	simplicityScore int
	chargePercent   int
	brakesWorking   bool
}

func NewElectricScooter(invNum string, energy float64, score int, charge int, brakesWorking bool) *ElectricScooter {
	return &ElectricScooter{
		inventoryNumber: invNum,
		dailyEnergy:     energy,
		simplicityScore: score,
		chargePercent:   charge,
		brakesWorking:   brakesWorking,
	}
}

func (es *ElectricScooter) InventoryNumber() string { return es.inventoryNumber }
func (es *ElectricScooter) Type() string            { return "Электросамокат" }
func (es *ElectricScooter) SimplicityScore() int    { return es.simplicityScore }
func (es *ElectricScooter) DailyEnergy() float64    { return es.dailyEnergy }
func (es *ElectricScooter) ChargePercent() int      { return es.chargePercent }
func (es *ElectricScooter) BrakesWorking() bool     { return es.brakesWorking }

// bike (велосипед): inventarizable, beginner-friendly, braked
type Bike struct {
	inventoryNumber string
	simplicityScore int
	brakesWorking   bool
}

func NewBike(invNum string, score int, brakesWorking bool) *Bike {
	return &Bike{inventoryNumber: invNum, simplicityScore: score, brakesWorking: brakesWorking}
}

func (b *Bike) InventoryNumber() string { return b.inventoryNumber }
func (b *Bike) Type() string            { return "Велосипед" }
func (b *Bike) SimplicityScore() int    { return b.simplicityScore }
func (b *Bike) BrakesWorking() bool     { return b.brakesWorking }

// electric bike (электровелосипед): inventarizable, energy consumer, beginner-friendly, chargeable, braked
type ElectricBike struct {
	inventoryNumber string
	dailyEnergy     float64
	simplicityScore int
	chargePercent   int
	brakesWorking   bool
}

func NewElectricBike(invNum string, energy float64, score int, charge int, brakesWorking bool) *ElectricBike {
	return &ElectricBike{
		inventoryNumber: invNum,
		dailyEnergy:     energy,
		simplicityScore: score,
		chargePercent:   charge,
		brakesWorking:   brakesWorking,
	}
}

func (eb *ElectricBike) InventoryNumber() string { return eb.inventoryNumber }
func (eb *ElectricBike) Type() string            { return "Электровелосипед" }
func (eb *ElectricBike) SimplicityScore() int    { return eb.simplicityScore }
func (eb *ElectricBike) DailyEnergy() float64    { return eb.dailyEnergy }
func (eb *ElectricBike) ChargePercent() int      { return eb.chargePercent }
func (eb *ElectricBike) BrakesWorking() bool     { return eb.brakesWorking }
