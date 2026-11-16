package container

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/do/v2"

	"hw_2/internal/application"
	"hw_2/internal/infrastructure/csvexport"
	"hw_2/internal/infrastructure/jsonstore"
)

// Config contains paths to application data files
type Config struct {
	StudentsFile string
	RubricFile   string
}

// Container holds the dependency injection container and application facade
type Container struct {
	injector *do.RootScope
	Facade   *application.ApplicationFacade
}

// New creates and configures the application dependency injection container
func New(config Config) (*Container, error) {
	if config.StudentsFile == "" {
		config.StudentsFile = "data/students.json"
	}

	if config.RubricFile == "" {
		config.RubricFile = "data/rubric.json"
	}

	if err := os.MkdirAll(filepath.Dir(config.StudentsFile), 0755); err != nil {
		return nil, fmt.Errorf("create students data directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(config.RubricFile), 0755); err != nil {
		return nil, fmt.Errorf("create rubric data directory: %w", err)
	}

	injector := do.New()

	// Register rubric repository
	do.Provide[application.RubricRepository](
		injector,
		func(i do.Injector) (application.RubricRepository, error) {
			return jsonstore.NewJSONRubricRepository(config.RubricFile), nil
		},
	)

	// Register student repository
	do.Provide[application.StudentRepository](
		injector,
		func(i do.Injector) (application.StudentRepository, error) {
			rubricRepository, err := do.Invoke[application.RubricRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve rubric repository: %w", err)
			}

			return jsonstore.NewJSONStudentRepository(
				config.StudentsFile,
				rubricRepository,
			), nil
		},
	)

	// Register CSV exporter
	do.Provide[application.FinalScoreExporter](
		injector,
		func(i do.Injector) (application.FinalScoreExporter, error) {
			return csvexport.NewCSVExporter(), nil
		},
	)

	// Register student service
	do.Provide[*application.StudentService](
		injector,
		func(i do.Injector) (*application.StudentService, error) {
			repository, err := do.Invoke[application.StudentRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student repository: %w", err)
			}

			return application.NewStudentService(repository), nil
		},
	)

	// Register rubric service
	do.Provide[*application.RubricService](
		injector,
		func(i do.Injector) (*application.RubricService, error) {
			rubricRepository, err := do.Invoke[application.RubricRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve rubric repository: %w", err)
			}

			studentRepository, err := do.Invoke[application.StudentRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student repository: %w", err)
			}

			return application.NewRubricService(
				rubricRepository,
				studentRepository,
			), nil
		},
	)

	// Register evaluation service
	do.Provide[*application.EvaluationService](
		injector,
		func(i do.Injector) (*application.EvaluationService, error) {
			studentRepository, err := do.Invoke[application.StudentRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student repository: %w", err)
			}

			rubricRepository, err := do.Invoke[application.RubricRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve rubric repository: %w", err)
			}

			return application.NewEvaluationService(
				studentRepository,
				rubricRepository,
			), nil
		},
	)

	// Register score service
	do.Provide[*application.ScoreService](
		injector,
		func(i do.Injector) (*application.ScoreService, error) {
			studentRepository, err := do.Invoke[application.StudentRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student repository: %w", err)
			}

			rubricRepository, err := do.Invoke[application.RubricRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve rubric repository: %w", err)
			}

			return application.NewScoreService(
				studentRepository,
				rubricRepository,
			), nil
		},
	)

	// Register export service
	do.Provide[*application.ExportService](
		injector,
		func(i do.Injector) (*application.ExportService, error) {
			studentRepository, err := do.Invoke[application.StudentRepository](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student repository: %w", err)
			}

			scoreService, err := do.Invoke[*application.ScoreService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve score service: %w", err)
			}

			exporter, err := do.Invoke[application.FinalScoreExporter](i)
			if err != nil {
				return nil, fmt.Errorf("resolve score exporter: %w", err)
			}

			return application.NewExportService(
				studentRepository,
				scoreService,
				exporter,
			), nil
		},
	)

	// Register application facade
	do.Provide[*application.ApplicationFacade](
		injector,
		func(i do.Injector) (*application.ApplicationFacade, error) {
			studentService, err := do.Invoke[*application.StudentService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve student service: %w", err)
			}

			rubricService, err := do.Invoke[*application.RubricService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve rubric service: %w", err)
			}

			evaluationService, err := do.Invoke[*application.EvaluationService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve evaluation service: %w", err)
			}

			scoreService, err := do.Invoke[*application.ScoreService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve score service: %w", err)
			}

			exportService, err := do.Invoke[*application.ExportService](i)
			if err != nil {
				return nil, fmt.Errorf("resolve export service: %w", err)
			}

			return application.NewApplicationFacade(
				studentService,
				rubricService,
				evaluationService,
				scoreService,
				exportService,
			), nil
		},
	)

	facade, err := do.Invoke[*application.ApplicationFacade](injector)
	if err != nil {
		injector.Shutdown()
		return nil, fmt.Errorf("initialize application facade: %w", err)
	}

	return &Container{
		injector: injector,
		Facade:   facade,
	}, nil
}

// Shutdown releases resources managed by the dependency injection container
func (container *Container) Shutdown() {
	if container == nil || container.injector == nil {
		return
	}

	container.injector.Shutdown()
}
