package di

import (
	"fmt"
	"io"

	"github.com/samber/do/v2"

	"pos/internal/adapters/console"
	"pos/internal/adapters/csvexport"
	"pos/internal/adapters/jsonstore"
	"pos/internal/application"
	"pos/internal/domain"
)

func BuildApp(directory string, mode application.CalculationMode, reader io.Reader, output io.Writer) (*console.App, error) {
	injector := do.New()
	do.Provide(injector, func(_ do.Injector) (domain.FinalScoreStrategy, error) {
		switch mode {
		case application.ModeStrict:
			return domain.StrictStrategy{}, nil
		case application.ModeMissingAsZero:
			return domain.MissingAsZeroStrategy{}, nil
		default:
			return nil, fmt.Errorf("unknown calculation mode: %q", mode)
		}
	})
	do.Provide(injector, func(_ do.Injector) (application.CourseRepository, error) {
		return jsonstore.NewRepository(directory)
	})
	do.Provide(injector, func(_ do.Injector) (application.FinalScoreExporter, error) {
		return csvexport.Exporter{}, nil
	})
	do.Provide(injector, func(i do.Injector) (*domain.FinalScoreService, error) {
		strategy, err := do.Invoke[domain.FinalScoreStrategy](i)
		if err != nil {
			return nil, err
		}
		return domain.NewFinalScoreService(strategy)
	})
	do.Provide(injector, func(i do.Injector) (*application.Service, error) {
		repository, err := do.Invoke[application.CourseRepository](i)
		if err != nil {
			return nil, err
		}
		scoring, err := do.Invoke[*domain.FinalScoreService](i)
		if err != nil {
			return nil, err
		}
		exporter, err := do.Invoke[application.FinalScoreExporter](i)
		if err != nil {
			return nil, err
		}
		return application.NewService(repository, scoring, exporter, mode)
	})
	do.Provide(injector, func(i do.Injector) (*console.App, error) {
		service, err := do.Invoke[*application.Service](i)
		if err != nil {
			return nil, err
		}
		return console.NewApp(service, reader, output)
	})
	app, err := do.Invoke[*console.App](injector)
	if err != nil {
		return nil, fmt.Errorf("build application: %w", err)
	}
	return app, nil
}
