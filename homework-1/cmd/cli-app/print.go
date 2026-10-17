package main

import (
	"fmt"
	"vyshkat/internal/inventory"
	"vyshkat/internal/service"
)

func printFleetStats(fleet *service.Fleet) {
	fmt.Printf("Всего в парке: %d ед.\n", fleet.InventoryCount())
	fmt.Printf("Энергопотребление в сутки: %.2f кВт·ч\n", fleet.TotalEnergyConsuming())
	fmt.Printf("Для новичков: %d ед.\n", len(fleet.BeginnerFriendly()))
}

func printForBeginners(fleet *service.Fleet) {
	fmt.Println("\nИнвентарь для новичков:")
	items := fleet.BeginnerFriendly()
	if len(items) == 0 {
		fmt.Println("Нет техники для новичков")
		return
	}
	for _, item := range items {
		if variable, ok := item.(inventory.SimplicityEvaluated); ok {
			fmt.Printf("- %s (номер: %s, оценка: %d)\n", item.Type(),
				item.InventoryNumber(), variable.SimplicityScore())
		}
	}
}

func printInventory(fleet *service.Fleet) {
	fmt.Println("\nВесь доступный инвентарь:")
	items := fleet.AllInventory()
	if fleet.InventoryCount() == 0 {
		fmt.Println("Парк пуст")
		return
	}
	for _, item := range items {
		fmt.Printf("[%s] %s\n", item.InventoryNumber(), item.Type())
	}
}
