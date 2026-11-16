package application

import (
	"fmt"
	"hw_2/internal/domain"
)

// EvaluationService handles application operations related to evaluations
type EvaluationService struct {
	studentRepository StudentRepository
	rubricRepository  RubricRepository
}

func NewEvaluationService(
	studentRepository StudentRepository,
	rubricRepository RubricRepository,
) *EvaluationService {
	return &EvaluationService{
		studentRepository: studentRepository,
		rubricRepository:  rubricRepository,
	}
}

// SetEvaluation creates or updates an evaluation for a student
func (service *EvaluationService) SetEvaluation(
	studentID string,
	criterionID string,
	score int,
	comment string,
) error {
	student, err := service.studentRepository.GetByID(studentID)

	if err != nil {
		return fmt.Errorf("get student: %w", err)
	}

	rubric, err := service.rubricRepository.Get()

	if err != nil {
		return fmt.Errorf("get rubric: %w", err)
	}

	criterion, exists := rubric.GetCriterion(criterionID)

	if !exists {
		return fmt.Errorf("criterion with ID %s does not exist", criterionID)
	}

	scoreValue, err := domain.NewScoreValue(score)

	if err != nil {
		return fmt.Errorf("create score: %w", err)
	}

	if err := student.SetEvaluation(
		criterionID,
		scoreValue,
		criterion.MaxScore(),
		comment,
	); err != nil {
		return fmt.Errorf("set evaluation: %w", err)
	}

	if err := service.studentRepository.Save(student); err != nil {
		return fmt.Errorf("save student: %w", err)
	}

	return nil
}
