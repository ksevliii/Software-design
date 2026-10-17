package main

import (
	"fmt"
	"vyshkat/internal/inventory"
	"vyshkat/internal/service"
)

func runDemo(inspector service.TechInspector) {
	fmt.Println("~~~ ДЕМОНСТРАЦИЯ СЦЕНАРИЕВ ВЫШКАТ ~~~")
	fleet := service.NewFleet()

	// case 1: correct items
	fmt.Println("--- Сценарий 1: приём корректной техники ---")
	demoAdd(fleet, inspector, inventory.NewElectricScooter("ES01", 1.5, 8, 85, true), "электросамокат, заряд 85%")
	demoAdd(fleet, inspector, inventory.NewScooter("S02", 7, true), "самокат, тормоза исправны")
	demoAdd(fleet, inspector, inventory.NewElectricBike("EB067", 2.0, 7, 90, true), "электровелосипед, заряд 90%")
	demoAdd(fleet, inspector, inventory.NewBike("B08", 9, true), "велосипед")
	demoAdd(fleet, inspector, inventory.NewHelmet("H228"), "шлем")
	demoAdd(fleet, inspector, inventory.NewChargingStation("CS45", 5.0), "зарядная станция")
	demoAdd(fleet, inspector, inventory.NewVehicleStand("VS09"), "стойка для парковки")

	// case 2: troubles
	fmt.Println("\n--- Сценарий 2: отклонения техосмотром ---")
	demoAdd(fleet, inspector, inventory.NewScooter("AB", 8, true), "короткий номер (< 3 символов)")
	demoAdd(fleet, inspector, inventory.NewScooter("S99", 15, true), "невалидная оценка простоты (15 > 10)")
	demoAdd(fleet, inspector, inventory.NewElectricScooter("ES92", 1.5, 8, 50, true), "низкий заряд батареи (50% < 70%)")
	demoAdd(fleet, inspector, inventory.NewScooter("S13", 7, false), "неисправные тормоза")

	// case 3: duplicate
	fmt.Println("\n--- Сценарий 3: попытка добавить дубликат (по серийному номеру) ---")
	demoAdd(fleet, inspector, inventory.NewElectricScooter("ES01", 1.5, 8, 80, true), "дубликат номера ES01")

	// result report
	fmt.Println("\n~~~ ИТОГОВЫЕ ОТЧЁТЫ ПО ДЕМО-ПАРКУ ~~~")
	printFleetStats(fleet)
	printForBeginners(fleet)
	printInventory(fleet)

	fmt.Println("\n... Переход в интерактивный режим ...")
}

func demoAdd(fleet *service.Fleet, inspector service.TechInspector, item inventory.Inventarizable, reason string) {
	accepted := fleet.AddItem(item, inspector)
	status := "Принято!"
	if !accepted {
		status = "Отклонено:("
	}
	fmt.Printf("%s %s [%s]: %s\n", status, item.Type(), item.InventoryNumber(), reason)
}
