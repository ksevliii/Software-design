package main

import (
	"fmt"
	"vyshkat/internal/service"
)

func runInteractive(fleet *service.Fleet, inspector service.TechInspector) {
	fmt.Println("\n~~~ ИНТЕРАКТИВНЫЙ РЕЖИМ ВЫШКАТ ~~~")
	for {
		printOptions()
		var option int
		fmt.Scan(&option)
		switch option {
		case 1:
			tryAdd(fleet, inspector)
		case 2:
			printFleetStats(fleet)
		case 3:
			printForBeginners(fleet)
		case 4:
			printInventory(fleet)
		case 5:
			fmt.Println("До встречи!")
			return
		default:
			fmt.Println("Неизвестная опция")
		}
		fmt.Println()
	}
}

func printOptions() {
	fmt.Println("Выберите опцию ВышКат:")
	fmt.Println("1. Добавить технику")
	fmt.Println("2. Отчёт по парку")
	fmt.Println("3. Список для новичков")
	fmt.Println("4. Весь инвентарь")
	fmt.Println("5. Выход")
	fmt.Print("\n> ")
}
