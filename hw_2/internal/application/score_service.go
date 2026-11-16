package application

import (
	"fmt"

	"hw_2/internal/domain"
)

// ScoreService executes the use case of calculating a student's final score
type ScoreService struct {
	studentRepository StudentRepository
	rubricRepository  RubricRepository
	strategy          ScoreCalculationStrategy
}

// NewScoreService creates a score service with the default weighted strategy
func NewScoreService(
	studentRepository StudentRepository,
	rubricRepository RubricRepository,
) *ScoreService {
	return NewScoreServiceWithStrategy(
		studentRepository,
		rubricRepository,
		WeightedScoreStrategy{},
	)
}

// NewScoreServiceWithStrategy creates a score service with a selected calculation strategy
func NewScoreServiceWithStrategy(
	studentRepository StudentRepository,
	rubricRepository RubricRepository,
	strategy ScoreCalculationStrategy,
) *ScoreService {
	if strategy == nil {
		strategy = WeightedScoreStrategy{}
	}

	return &ScoreService{
		studentRepository: studentRepository,
		rubricRepository:  rubricRepository,
		strategy:          strategy,
	}
}

// CalculateStudentScore calculates the final score for a specific student
func (service *ScoreService) CalculateStudentScore(
	studentID string,
) (domain.FinalScore, error) {
	student, err := service.studentRepository.GetByID(studentID)
	if err != nil {
		return domain.FinalScore{}, fmt.Errorf("get student: %w", err)
	}

	rubric, err := service.rubricRepository.Get()
	if err != nil {
		return domain.FinalScore{}, fmt.Errorf("get rubric: %w", err)
	}

	finalScore, err := service.strategy.Calculate(
		student.GetEvaluations(),
		rubric.GetCriteria(),
	)
	if err != nil {
		return domain.FinalScore{}, fmt.Errorf("calculate final score: %w", err)
	}

	return finalScore, nil
}
