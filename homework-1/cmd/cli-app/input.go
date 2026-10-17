package main

import (
	"fmt"
	"vyshkat/internal/inventory"
	"vyshkat/internal/service"
)

func tryAdd(fleet *service.Fleet, inspector service.TechInspector) {
	item := readItem()
	if item == nil {
		return
	}

	if fleet.AddItem(item, inspector) {
		fmt.Printf("\n%s (номер: %s) принят в парк\n", item.Type(), item.InventoryNumber())
	} else {
		fmt.Printf("\n%s (номер: %s) отклонён техосмотром\n", item.Type(), item.InventoryNumber())
	}
}

func readItem() inventory.Inventarizable {
	fmt.Println("Введите тип: 1-самокат, 2-электросамокат, 3-велосипед, 4-электровелосипед, 5-шлем, 6-зарядная станция, 7-стойка для парковки")
	var itemType string
	fmt.Scan(&itemType)
	var num string
	fmt.Print("Инвентарный номер: ")
	fmt.Scan(&num)
	switch itemType {
	case "1": // Scooter: inventarizable, beginner-friendly, braked
		return inventory.NewScooter(num, askScore(), askBrakes())
	case "2": // ElectricScooter: inventarizable, energy-consuming, beginner-friendly, chargeable, braked
		return inventory.NewElectricScooter(num, askEnergy(), askScore(), askCharge(), askBrakes())
	case "3": // Bike: inventarizable, beginner-friendly, braked
		return inventory.NewBike(num, askScore(), askBrakes())
	case "4": // ElectricBike: inventarizable, energy-consuming, beginner-friendly, chargeable, braked
		return inventory.NewElectricBike(num, askEnergy(), askScore(), askCharge(), askBrakes())
	case "5": // Helmet: inventarizable
		return inventory.NewHelmet(num)
	case "6": // ChargingStation: inventarizable, energy-consuming
		return inventory.NewChargingStation(num, askEnergy())
	case "7": // VehicleStand: inventarizable
		return inventory.NewVehicleStand(num)
	default:
		fmt.Println("Неизвестный тип")
		return nil
	}
}

// helpers
func askScore() int {
	var score int
	fmt.Print("Простота для новичка (1-10): ")
	fmt.Scan(&score)
	return score
}

func askBrakes() bool {
	var brakes int
	fmt.Print("Тормоза исправны? (1 - да, 0 - нет): ")
	fmt.Scan(&brakes)
	return brakes == 1
}

func askCharge() int {
	var charge int
	fmt.Print("Уровень заряда батареи (0-100%): ")
	fmt.Scan(&charge)
	return charge
}

func askEnergy() float64 {
	var energy float64
	fmt.Print("Энергопотребление в сутки (кВт·ч): ")
	fmt.Scan(&energy)
	return energy
}


