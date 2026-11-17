package main

import (
	"flag"
	"fmt"
	"os"

	"pos/internal/adapters/csvexport"
	"pos/internal/adapters/jsonstore"
	"pos/internal/application"
	"pos/internal/di"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	directory := flag.String("data", "data", "Directory for JSON course files")
	modeFlag := flag.String("mode", string(application.ModeStrict), "Calculation mode: strict or missing_as_zero")
	flag.Parse()
	mode := application.CalculationMode(*modeFlag)
	if !mode.Valid() {
		return fmt.Errorf("unknown calculation mode: %q", mode)
	}
	repository, err := jsonstore.NewRepository(*directory)
	if err != nil {
		return fmt.Errorf("create repository: %w", err)
	}
	app, err := di.BuildApp(repository, csvexport.Exporter{}, mode, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	return app.Run()
}
