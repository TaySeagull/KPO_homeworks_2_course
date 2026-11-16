package application

import "hw_2/internal/domain"

// StudentRepository defines operations for storing Student aggregates
type StudentRepository interface {
	GetByID(id string) (domain.Student, error)

	GetAll() ([]domain.Student, error)

	Save(student domain.Student) error

	ExistsByID(id string) (bool, error)

	ExistsEvaluationForCriterion(criterionID string) (bool, error)
}

// RubricRepository defines operations for storing the Rubric aggregate
type RubricRepository interface {
	Get() (domain.Rubric, error)

	Save(rubric domain.Rubric) error
}
