package di

import (
	"fmt"
	"io"

	"github.com/samber/do/v2"

	"pos/internal/adapters/console"
	"pos/internal/application"
	"pos/internal/domain"
)

func BuildService(
	repository application.CourseRepository,
	exporter application.FinalScoreExporter,
	mode application.CalculationMode) (*application.Service, error) {
	if repository == nil {
		return nil, fmt.Errorf("course repository must not be nil")
	}
	if exporter == nil {
		return nil, fmt.Errorf("final score exporter must not be nil")
	}
	if !mode.Valid() {
		return nil, fmt.Errorf("unknown calculation mode: %q", mode)
	}
	injector := do.New()
	do.ProvideValue(injector, repository)
	do.ProvideValue(injector, exporter)
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
	service, err := do.Invoke[*application.Service](injector)
	if err != nil {
		return nil, fmt.Errorf("build application service: %w", err)
	}
	return service, nil
}

func BuildApp(
	repository application.CourseRepository,
	exporter application.FinalScoreExporter,
	mode application.CalculationMode,
	reader io.Reader,
	output io.Writer) (*console.App, error) {
	service, err := BuildService(repository, exporter, mode)
	if err != nil {
		return nil, err
	}
	app, err := console.NewApp(service, reader, output)
	if err != nil {
		return nil, fmt.Errorf("build console application: %w", err)
	}
	return app, nil
}
