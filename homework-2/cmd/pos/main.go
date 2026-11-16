package main

import (
	"flag"
	"fmt"
	"os"

	"pos/internal/application"
	"pos/internal/di"
)

func main() {
	directory := flag.String("data", "data", "Directory for JSON course files")
	mode := flag.String("mode", string(application.ModeStrict), "Calculation mode: strict or missing_as_zero")
	flag.Parse()
	app, err := di.BuildApp(*directory, application.CalculationMode(*mode), os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка запуска:", err)
		os.Exit(1)
	}
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка выполнения:", err)
		os.Exit(1)
	}
}
