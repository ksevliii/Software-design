package inventory

// helmet (шлем): inventarizable
type Helmet struct {
	inventoryNumber string
}

func NewHelmet(invNum string) *Helmet {
	return &Helmet{inventoryNumber: invNum}
}

func (h *Helmet) InventoryNumber() string { return h.inventoryNumber }
func (h *Helmet) Type() string            { return "Шлем" }

// charging station (зарядная станция): inventarizable and energy comsumer
type ChargingStation struct {
	inventoryNumber string
	dailyEnergy     float64
}

func NewChargingStation(invNum string, energy float64) *ChargingStation {
	return &ChargingStation{inventoryNumber: invNum, dailyEnergy: energy}
}

func (c *ChargingStation) InventoryNumber() string { return c.inventoryNumber }
func (c *ChargingStation) Type() string            { return "Зарядная станция" }
func (c *ChargingStation) DailyEnergy() float64    { return c.dailyEnergy }

// VehicleStand (стойка для парковки): inventarizable
type VehicleStand struct {
	inventoryNumber string
}

func NewVehicleStand(invNum string) *VehicleStand {
	return &VehicleStand{inventoryNumber: invNum}
}

func (v *VehicleStand) InventoryNumber() string { return v.inventoryNumber }
func (v *VehicleStand) Type() string            { return "Стойка для парковки" }
